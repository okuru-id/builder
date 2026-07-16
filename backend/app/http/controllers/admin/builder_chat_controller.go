package admin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
	"strings"
	"time"

	ghttp "github.com/goravel/framework/contracts/http"

	"okuru/app/facades"
)

// BuilderChatController streams an OpenAI-compatible chat completion for the
// builder assistant. Config via env: LLM_BASE_URL, LLM_API_KEY, LLM_MODEL.
// When LLM_API_KEY is empty the endpoint streams a single guidance message so
// the UI still works end-to-end and tells the operator how to enable the LLM.
//
// ponytail: dedicated client with a 180s timeout so a stalled upstream cannot
// leak goroutines/connections forever. Generous budget because some routed
// models (e.g. glm-5.2) are reasoning models that can think for a long time
// before emitting content. Upgrade path: per-request timeout knob.
var llmClient = &nethttp.Client{Timeout: 180 * time.Second}

type BuilderChatController struct{}

func NewBuilderChatController() *BuilderChatController {
	return &BuilderChatController{}
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Messages    []chatMessage `json:"messages"`
	Tree        any           `json:"tree"`
	PageName    string        `json:"pageName"`
	FocusNode   any           `json:"focusNode"`
	NodeCatalog any           `json:"nodeCatalog"`
}

const fence = "`" + "`" + "`"

const builderSystemPrompt = `You are an AI agent embedded inside a visual landing-page builder (okuru.id).
You act on the user's page, represented as a tree of nodes. Each node has:
  id, type, name, props, classes (Tailwind utility classes), and children.
All styling uses Tailwind CSS utility classes only.

The frontend sends a node catalog below. It is authoritative: only create node
type values listed in Available node types. Nodes listed in
Container types may contain children; all other nodes must have no children.
For icon nodes, use Tabler icon names in PascalCase (e.g. IconStar, IconHeart,
IconArrowRight, IconMail, IconPhone) — the frontend resolves them against the
full @tabler/icons-vue set. Use the current tree's real ids for changes; never
invent ids.

Capabilities:
- Explain and critique the current design.
- Suggest concrete improvements: exact Tailwind classes, node additions,
  text rewrites, hierarchy restructures.
- Output an action block in a fenced code region when the user asks for a
  change. The frontend parses these blocks and offers an "Apply" button.
  Use ONE of these formats per action block:

  Add a node as the last child of a parent id (use "root" for top level):
  ` + fence + `action:add
  { "parentId": "root", "node": { "type": "section", "classes": ["py-16","px-6","bg-neutral-900","text-white"], "children": [ { "type": "heading", "props": { "text": "Hello", "level": 2 }, "classes": ["text-4xl","font-bold"] } ] } }
  ` + fence + `

  Patch classes on a node (replace matched prefixes):
  ` + fence + `action:classes
  { "nodeId": "<id>", "set": ["flex","flex-col","gap-4"] }
  ` + fence + `

  Rewrite text:
  ` + fence + `action:text
  { "nodeId": "<id>", "text": "New copy" }
  ` + fence + `

  Delete a node (root cannot be deleted):
  ` + fence + `action:delete
  { "nodeId": "<id>" }
  ` + fence + `

  Move/reorder a node under a new parent (use "root" for top level; index -1 = append):
  ` + fence + `action:move
  { "nodeId": "<id>", "parentId": "root", "index": -1 }
  ` + fence + `

Rules:
- Be concise. Default to short, direct answers.
- Only emit an action block when the user clearly asks for a change.
- Never invent node ids. For new nodes omit "id"; the frontend assigns one.
- Tailwind only. No custom CSS.

Current page name: %s
Available builder node catalog JSON:
%s
Focused node JSON (null means no focused node):
%s
If a focused node exists, target it or its descendants unless the user explicitly asks for a broader change.
Current tree JSON:
%s`

// Chat streams an SSE response: one event per token chunk, final [DONE].
func (c *BuilderChatController) Chat(ctx ghttp.Context) ghttp.Response {
	var in chatRequest
	if err := ctx.Request().Bind(&in); err != nil {
		return ctx.Response().Json(ghttp.StatusBadRequest, ghttp.Json{"error": err.Error()})
	}

	baseURL := strings.TrimRight(facades.Config().GetString("LLM_BASE_URL", "https://api.openai.com/v1"), "/")
	apiKey := facades.Config().GetString("LLM_API_KEY", "")
	model := facades.Config().GetString("LLM_MODEL", "gpt-4o-mini")

	treeJSON := "null"
	if in.Tree != nil {
		if b, err := json.Marshal(in.Tree); err == nil {
			treeJSON = string(b)
		}
	}
	catalogJSON := "{}"
	if in.NodeCatalog != nil {
		if b, err := json.Marshal(in.NodeCatalog); err == nil {
			catalogJSON = string(b)
		}
	}
	focusJSON := "null"
	if in.FocusNode != nil {
		if b, err := json.Marshal(in.FocusNode); err == nil {
			focusJSON = string(b)
		}
	}

	// SSE event writer. Each event: "data: <json>\n\n". End with "data: [DONE]\n\n".
	writeEvent := func(w ghttp.StreamWriter, content string) error {
		payload := map[string]string{"content": content}
		b, _ := json.Marshal(payload)
		if _, err := w.WriteString("data: " + string(b) + "\n\n"); err != nil {
			return err
		}
		return w.Flush()
	}

	return ctx.Response().Stream(ghttp.StatusOK, func(w ghttp.StreamWriter) error {
		// Not configured: emit guidance + close.
		if apiKey == "" {
			_ = writeEvent(w, "AI Agent belum dikonfigurasi. Set LLM_API_KEY (dan opsional LLM_BASE_URL, LLM_MODEL) di backend/.env untuk mengaktifkan agen.")
			_, _ = w.WriteString("data: [DONE]\n\n")
			return w.Flush()
		}

		// Build upstream request.
		sys := fmt.Sprintf(builderSystemPrompt, in.PageName, catalogJSON, focusJSON, treeJSON)
		msgs := append([]chatMessage{{Role: "system", Content: sys}}, in.Messages...)
		body := map[string]any{
			"model":    model,
			"messages": msgs,
			"stream":   true,
			// Disable chain-of-thought reasoning when the upstream model supports
			// it (e.g. GLM/Zhipu thinking models). The builder agent must emit
			// structured action blocks, so silent reasoning that delays content
			// past the client timeout reads as "(no response)". Harmless if the
			// provider ignores the field.
			"thinking": map[string]any{"type": "disabled"},
		}
		bodyBytes, _ := json.Marshal(body)

		req, err := nethttp.NewRequestWithContext(ctx.Request().Origin().Context(), nethttp.MethodPost, baseURL+"/chat/completions", bytes.NewReader(bodyBytes))
		if err != nil {
			_ = writeEvent(w, "Gagal membangun request: "+err.Error())
			_, _ = w.WriteString("data: [DONE]\n\n")
			return w.Flush()
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := llmClient.Do(req)
		if err != nil {
			_ = writeEvent(w, "Gagal menghubungi LLM: "+err.Error())
			_, _ = w.WriteString("data: [DONE]\n\n")
			return w.Flush()
		}
		defer resp.Body.Close()

		if resp.StatusCode != nethttp.StatusOK {
			raw, _ := io.ReadAll(resp.Body)
			_ = writeEvent(w, fmt.Sprintf("LLM error (HTTP %d): %s", resp.StatusCode, string(raw)))
			_, _ = w.WriteString("data: [DONE]\n\n")
			return w.Flush()
		}

		// Parse SSE stream from upstream and forward content deltas.
		scanner := bufio.NewScanner(resp.Body)
		scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content          string `json:"content"`
						ReasoningContent string `json:"reasoning_content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}
			if len(chunk.Choices) > 0 {
				// Prefer visible content. Fall back to reasoning_content so the
				// stream is never silent on models that only emit reasoning
				// (defensive; thinking is disabled above).
				piece := chunk.Choices[0].Delta.Content
				if piece == "" {
					piece = chunk.Choices[0].Delta.ReasoningContent
				}
				if piece != "" {
					if err := writeEvent(w, piece); err != nil {
						return err
					}
				}
			}
		}

		_, _ = w.WriteString("data: [DONE]\n\n")
		return w.Flush()
	})
}
