package services

import (
	"encoding/json"

	"okuru/app/facades"
	"okuru/app/models"
)

// ResolveComponentInstances walks a page tree and replaces every component
// instance node (type=component, componentId set) with a deep clone of its
// master root. The clone keeps the instance node's id so positional identity
// is preserved, but absorbs the master's type/props/classes/children.
//
// Masters are batch-loaded in a single query to avoid N+1. If a master is
// missing (deleted), the instance is replaced with a neutral placeholder div
// so publish never fails on a dangling reference.
//
// This runs server-side at publish time so the published HTML is self-contained
// and never depends on the component masters surviving in the DB.
func ResolveComponentInstances(tree map[string]any) map[string]any {
	root, _ := tree["root"].(map[string]any)
	if root == nil {
		return tree
	}
	ids := collectComponentIds(root, nil)
	masters := loadMasters(ids)
	resolved := resolveNode(root, masters)
	return map[string]any{"root": resolved}
}

func collectComponentIds(n map[string]any, out []uint) []uint {
	if id := componentIdOf(n); id > 0 {
		out = append(out, id)
	}
	for _, c := range childrenOf(n) {
		if cm, ok := c.(map[string]any); ok {
			out = collectComponentIds(cm, out)
		}
	}
	return out
}

func loadMasters(ids []uint) map[uint]map[string]any {
	if len(ids) == 0 {
		return nil
	}
	// Dedup.
	seen := make(map[uint]struct{}, len(ids))
	uniq := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		uniq = append(uniq, id)
	}
	var comps []models.LandingComponent
	if err := facades.Orm().Query().Where("id IN ?", uniq).Get(&comps); err != nil {
		return nil
	}
	m := make(map[uint]map[string]any, len(comps))
	for _, c := range comps {
		var root map[string]any
		if err := json.Unmarshal(c.Tree, &root); err == nil {
			if r, ok := root["root"].(map[string]any); ok {
				m[c.ID] = r
			}
		}
	}
	return m
}

func resolveNode(n map[string]any, masters map[uint]map[string]any) map[string]any {
	if n == nil {
		return nil
	}
	// Instance: clone master, preserve instance id.
	if id := componentIdOf(n); id > 0 {
		if master, ok := masters[id]; ok {
			clone := deepCloneMap(master)
			if c, ok := clone.(map[string]any); ok {
				c["id"] = n["id"]
				// Preserve instance-level visibility overrides so hiding a component
				// instance in the builder also hides it in published HTML.
				if ho, ok := n["hiddenOn"]; ok {
					c["hiddenOn"] = ho
				}
				resolveChildren(c, masters)
				return c
			}
		}
		// Dangling reference → placeholder.
		return map[string]any{
			"id":      n["id"],
			"type":    "frame",
			"name":    "missing-component",
			"props":   map[string]any{},
			"classes": []string{"p-4", "border", "border-dashed", "border-neutral-300", "text-xs", "text-neutral-400"},
			"children": []any{},
		}
	}
	resolveChildren(n, masters)
	return n
}

func resolveChildren(n map[string]any, masters map[uint]map[string]any) {
	kids := childrenOf(n)
	if len(kids) == 0 {
		return
	}
	out := make([]any, 0, len(kids))
	for _, c := range kids {
		if cm, ok := c.(map[string]any); ok {
			out = append(out, resolveNode(cm, masters))
		} else {
			out = append(out, c)
		}
	}
	n["children"] = out
}

func componentIdOf(n map[string]any) uint {
	if n == nil {
		return 0
	}
	v, ok := n["componentId"]
	if !ok {
		return 0
	}
	switch t := v.(type) {
	case float64:
		return uint(t)
	case uint:
		return t
	}
	return 0
}

func childrenOf(n map[string]any) []any {
	if n == nil {
		return nil
	}
	c, _ := n["children"].([]any)
	return c
}

func deepCloneMap(m map[string]any) any {
	b, err := json.Marshal(m)
	if err != nil {
		return m
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return m
	}
	return out
}
