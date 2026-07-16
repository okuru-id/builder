<script setup lang="ts">
// Right inspector: edits the selected node's props + classes. Delete/duplicate actions.
// Always fills its parent (RightPanel tab). Context header + centered empty state.
import { computed, inject, ref } from 'vue'
import {
  IconCopy, IconTrash, IconArrowUp, IconArrowDown, IconUnlink, IconCode, IconSettings,
  IconTypography, IconPhoto, IconLink,
  IconSquare, IconSection, IconGridDots, IconH1, IconCircuitPushbutton, IconMinus,
  IconStar, IconForms, IconInputSearch, IconComponents, IconPointer,
} from '@tabler/icons-vue'
import { Button, buttonVariants } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { BUILDER_KEY } from '@/components/builder/injection'
import type { Breakpoint } from '@/types/page-builder'
import { findNode } from '@/components/builder/tree-utils'
import { ICONS, ICON_LIST } from '@/lib/icon-map'
import LayoutSection from './inspector/LayoutSection.vue'
import TypographySection from './inspector/TypographySection.vue'
import SpacingSection from './inspector/SpacingSection.vue'
import BackgroundSection from './inspector/BackgroundSection.vue'
import BorderSection from './inspector/BorderSection.vue'
import AppearanceSection from './inspector/AppearanceSection.vue'
import SizeSection from './inspector/SizeSection.vue'
import InspectorSection from './inspector/InspectorSection.vue'
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogCancel,
  AlertDialogAction,
} from '@/components/ui/alert-dialog'

const store = inject(BUILDER_KEY, null)!

const node = computed(() => store.selectedNode.value)

// ponytail: node-type icon/color map duplicated in TreeRow + NodeTreePanel.
// Extract to shared node-meta module when a 4th consumer appears.
const NODE_ICON: Record<string, any> = {
  frame: IconSquare,
  section: IconSection,
  grid: IconGridDots,
  text: IconTypography,
  heading: IconH1,
  image: IconPhoto,
  button: IconCircuitPushbutton,
  link: IconLink,
  divider: IconMinus,
  icon: IconStar,
  form: IconForms,
  input: IconInputSearch,
  component: IconComponents,
}
const NODE_COLOR: Record<string, string> = {
  frame: 'text-violet-500',
  section: 'text-sky-500',
  grid: 'text-cyan-500',
  text: 'text-muted-foreground',
  heading: 'text-amber-500',
  image: 'text-emerald-500',
  button: 'text-rose-500',
  link: 'text-blue-500',
  divider: 'text-neutral-400',
  icon: 'text-amber-500',
  form: 'text-emerald-500',
  input: 'text-sky-500',
  component: 'text-purple-500',
}

function setProp(key: string, value: unknown) {
  if (!node.value) return
  store.patchNode(node.value.id, { props: { ...node.value.props, [key]: value } })
}
function setName(value: string | number) {
  if (!node.value) return
  store.patchNode(node.value.id, { name: String(value) })
}

// Per-breakpoint visibility. Toggling M/T/D adds/removes the breakpoint in
// node.hiddenOn, which codegen maps to Tailwind responsive utilities.
const BP_OPTIONS: { key: Breakpoint; label: string }[] = [
  { key: 'mobile', label: 'M' },
  { key: 'tablet', label: 'T' },
  { key: 'desktop', label: 'D' },
]
function bpHidden(bp: Breakpoint): boolean {
  return !!node.value?.hiddenOn?.includes(bp)
}
function toggleBp(bp: Breakpoint) {
  if (!node.value) return
  const cur = new Set(node.value.hiddenOn ?? [])
  if (cur.has(bp)) cur.delete(bp)
  else cur.add(bp)
  store.patchNode(node.value.id, { hiddenOn: [...cur] })
}

const showDeleteConfirm = ref(false)
function delConfirm() {
  showDeleteConfirm.value = true
}
function confirmDelete() {
  if (node.value) {
    store.removeNode(node.value.id)
    showDeleteConfirm.value = false
  }
}
function dup() {
  if (node.value) store.duplicateNode(node.value.id)
}
// Sibling position for up/down disabled state.
const atTop = computed(() => {
  if (!node.value) return true
  const f = findNode(store.tree.value.root, node.value.id)
  return !f || !f.parent || f.index === 0
})
const atBottom = computed(() => {
  if (!node.value) return true
  const f = findNode(store.tree.value.root, node.value.id)
  return !f || !f.parent || f.index >= f.parent.children.length - 1
})
function moveUp() {
  if (node.value) store.moveSiblingNode(node.value.id, -1)
}
function moveDown() {
  if (node.value) store.moveSiblingNode(node.value.id, 1)
}

function addClass(e: Event) {
  const input = e.target as HTMLInputElement
  const cls = input.value.trim()
  if (!cls || !node.value) return
  // Support space-separated classes
  const newClasses = cls.split(/\s+/).filter((c) => c && !node.value!.classes.includes(c))
  if (newClasses.length) {
    store.patchNode(node.value.id, { classes: [...node.value.classes, ...newClasses] })
  }
  input.value = ''
}

function removeClass(idx: number) {
  if (!node.value) return
  const classes = [...node.value.classes]
  classes.splice(idx, 1)
  store.patchNode(node.value.id, { classes })
}
</script>

<template>
  <div class="flex h-full min-h-0 flex-col">
    <!-- Compact context header: stays visible while scrolling the sections. -->
    <div v-if="node" class="flex shrink-0 items-center gap-1.5 border-b border-border bg-card px-3 py-2">
      <component :is="NODE_ICON[node.type] ?? IconSquare" class="size-3.5 shrink-0" :class="NODE_COLOR[node.type] ?? 'text-muted-foreground'" />
      <span class="truncate text-xs font-medium">{{ node.name || node.type }}</span>
      <span class="ml-auto shrink-0 rounded bg-muted px-1.5 py-0.5 text-[10px] uppercase tracking-wide text-muted-foreground">{{ node.type }}</span>
    </div>

    <!-- Centered empty state -->
    <div v-if="!node" class="flex flex-1 flex-col items-center justify-center gap-2 p-6 text-center">
      <IconPointer class="size-8 text-muted-foreground/40" />
      <p class="text-sm font-medium text-muted-foreground">No selection</p>
      <p class="max-w-[16rem] text-[11px] leading-relaxed text-muted-foreground/70">
        Click an element on the canvas to edit its properties.
      </p>
    </div>

    <div v-else class="flex-1 overflow-auto">
      <!-- Umum: identity + actions, same shell as style sections -->
      <InspectorSection title="General" :icon="IconSettings">
        <div v-if="node.id !== 'root'" class="flex items-center justify-between">
          <div class="flex items-center gap-0.5">
            <Button variant="ghost" size="icon-sm" :disabled="atTop" title="Move up" @click="moveUp">
              <IconArrowUp class="size-4" />
            </Button>
            <Button variant="ghost" size="icon-sm" :disabled="atBottom" title="Move down" @click="moveDown">
              <IconArrowDown class="size-4" />
            </Button>
            <Button variant="ghost" size="icon-sm" title="Duplicate" @click="dup">
              <IconCopy class="size-4" />
            </Button>
          </div>
          <Button variant="ghost" size="icon-sm" class="text-red-500 hover:bg-red-500/10 hover:text-red-600 dark:text-red-400" title="Delete" @click="delConfirm">
            <IconTrash class="size-4" />
          </Button>
        </div>

        <div class="flex items-center gap-2">
          <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Name</Label>
          <Input :model-value="node.name" class="h-7 flex-1 text-xs" @update:model-value="setName" />
        </div>

        <div class="flex items-center gap-2">
          <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Type</Label>
          <div class="flex-1 text-xs text-muted-foreground">{{ node.type }}</div>
        </div>

        <div
          v-if="node.componentId"
          class="flex items-center justify-between rounded-md border border-primary/20 bg-primary/10 px-2.5 py-1.5"
        >
          <div class="text-[11px] text-primary">
            Instance #{{ node.componentId }}
          </div>
          <Button variant="outline" size="sm" class="h-7 px-2 text-xs" @click="store.breakInstance(node.id)">
            <IconUnlink class="size-3.5" /> Detach
          </Button>
        </div>

        <!-- Per-breakpoint visibility: hide the node on Mobile / Tablet / Desktop.
             Maps to Tailwind hidden md:block / md:hidden lg:block / lg:hidden. -->
        <div v-if="node.id !== 'root'" class="flex items-center justify-between">
          <span class="text-[11px] text-muted-foreground">Hide on</span>
          <div class="flex gap-1">
            <button
              v-for="bp in BP_OPTIONS"
              :key="bp.key"
              class="h-6 w-7 rounded border text-[11px] font-medium transition-colors"
              :class="bpHidden(bp.key)
                ? 'border-primary bg-primary/10 text-primary'
                : 'border-border text-muted-foreground hover:bg-muted'"
              :title="`Hide on ${bp.key}`"
              @click="toggleBp(bp.key)"
            >{{ bp.label }}</button>
          </div>
        </div>
      </InspectorSection>

      <!-- Text-like: editable text + heading level -->
      <template v-if="node.type === 'text' || node.type === 'heading' || node.type === 'button' || node.type === 'link'">
        <InspectorSection title="Text" :icon="IconTypography">
          <div class="flex items-start gap-2">
            <Label class="mt-1.5 w-12 shrink-0 text-[11px] text-muted-foreground">Text</Label>
            <Textarea
              :model-value="String(node.props.text ?? '')"
              rows="2"
              class="flex-1 text-xs"
              @update:model-value="(v) => setProp('text', v)"
            />
          </div>
          <div v-if="node.type === 'heading'" class="flex items-center gap-2">
            <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Level</Label>
            <Select :model-value="String(node.props.level ?? 2)" @update:model-value="(v) => setProp('level', Number(v))">
              <SelectTrigger class="h-7 flex-1 text-xs"><SelectValue /></SelectTrigger>
              <SelectContent>
                <SelectItem v-for="l in [1, 2, 3, 4, 5, 6]" :key="l" :value="String(l)">H{{ l }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </InspectorSection>
      </template>

      <template v-if="node.type === 'image'">
        <InspectorSection title="Image" :icon="IconPhoto">
          <div class="flex items-center gap-2">
            <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">URL</Label>
            <Input :model-value="String(node.props.src ?? '')" class="h-7 flex-1 text-xs" @update:model-value="(v) => setProp('src', v)" />
          </div>
          <div class="flex items-center gap-2">
            <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Alt</Label>
            <Input :model-value="String(node.props.alt ?? '')" class="h-7 flex-1 text-xs" @update:model-value="(v) => setProp('alt', v)" />
          </div>
        </InspectorSection>
      </template>

      <template v-if="node.type === 'link'">
        <InspectorSection title="Link" :icon="IconLink">
          <div class="flex items-center gap-2">
            <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Href</Label>
            <Input :model-value="String(node.props.href ?? '#')" class="h-7 flex-1 text-xs" @update:model-value="(v) => setProp('href', v)" />
          </div>
        </InspectorSection>
      </template>

      <template v-if="node.type === 'icon'">
        <InspectorSection title="Icon" :icon="IconPhoto">
          <div class="grid grid-cols-6 gap-1">
            <button
              v-for="name in ICON_LIST"
              :key="name"
              class="flex h-8 w-8 items-center justify-center rounded-md border text-muted-foreground transition-colors"
              :class="node.props.icon === name ? 'border-primary bg-primary/10 text-primary' : 'border-border hover:bg-muted'"
              :title="name"
              @click="setProp('icon', name)"
            >
              <svg v-if="ICONS[name]" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round" class="size-4">
                <path v-for="seg in ICONS[name]" :key="seg[1].key" :d="seg[1].d" />
              </svg>
            </button>
          </div>
        </InspectorSection>
      </template>

      <template v-if="node.type === 'input'">
        <InspectorSection title="Input" :icon="IconPhoto">
          <div class="flex items-center gap-2">
            <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Label</Label>
            <Input :model-value="node.props.label ?? ''" class="h-7 flex-1 text-xs" @update:model-value="(v) => setProp('label', v)" />
          </div>
          <div class="flex items-center gap-2">
            <Label class="w-12 shrink-0 text-[11px] text-muted-foreground">Placeholder</Label>
            <Input :model-value="node.props.placeholder ?? ''" class="h-7 flex-1 text-xs" @update:model-value="(v) => setProp('placeholder', v)" />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-[11px] text-muted-foreground">Type</span>
            <select
              :value="node.props.inputType ?? 'text'"
              class="h-7 rounded-md border border-input bg-background px-2 text-xs"
              @change="(e) => setProp('inputType', (e.target as HTMLSelectElement).value)"
            >
              <option value="text">text</option>
              <option value="email">email</option>
              <option value="tel">tel</option>
              <option value="number">number</option>
              <option value="password">password</option>
              <option value="url">url</option>
            </select>
          </div>
          <div class="flex items-center gap-2">
            <input
              type="checkbox"
              :checked="node.props.required ?? false"
              class="size-3.5"
              @change="(e) => setProp('required', (e.target as HTMLInputElement).checked)"
            />
            <span class="text-[11px] text-muted-foreground">Required</span>
          </div>
        </InspectorSection>
      </template>

      <!-- Style sections: only for non-instance nodes -->
      <template v-if="!node.componentId">
        <LayoutSection />
        <TypographySection />
        <SpacingSection />
        <BackgroundSection />
        <BorderSection />
        <AppearanceSection />
        <SizeSection />

        <!-- Custom classes: arbitrary Tailwind class input -->
        <InspectorSection title="Custom Class" :icon="IconCode" :default-open="false">
          <div class="flex flex-wrap gap-1">
            <span
              v-for="(cls, idx) in node.classes"
              :key="idx"
              class="group inline-flex items-center gap-0.5 rounded bg-muted px-1.5 py-0.5 text-[11px] font-mono text-foreground"
            >
              {{ cls }}
              <button
                class="ml-0.5 rounded-sm text-muted-foreground opacity-0 transition-opacity hover:text-red-500 group-hover:opacity-100"
                @click="removeClass(idx)"
                title="Remove class"
              >&times;</button>
            </span>
          </div>
          <Input
            :model-value="''"
            class="h-7 font-mono text-xs"
            placeholder="Add class, press Enter"
            @keydown.enter.prevent="addClass"
          />
        </InspectorSection>
      </template>
      <div v-else class="border-b border-border/50 px-3 py-2 text-[11px] text-muted-foreground">
        Edit the master component to change instances. Click Detach above to make an independent copy.
      </div>
    </div>

    <AlertDialog v-model:open="showDeleteConfirm">
      <AlertDialogContent>
        <AlertDialogTitle>Delete element</AlertDialogTitle>
        <AlertDialogDescription>
          Are you sure you want to delete <strong>{{ node?.name }}</strong>? This action cannot be undone.
        </AlertDialogDescription>
        <div class="flex justify-end gap-2">
          <AlertDialogCancel :class="cn(buttonVariants({ variant: 'outline', size: 'sm' }), 'text-xs')">Cancel</AlertDialogCancel>
          <AlertDialogAction :class="cn(buttonVariants({ variant: 'destructive', size: 'sm' }), 'text-xs')" @click="confirmDelete">
            Delete
          </AlertDialogAction>
        </div>
      </AlertDialogContent>
    </AlertDialog>
  </div>
</template>
