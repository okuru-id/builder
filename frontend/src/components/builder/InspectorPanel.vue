<script setup lang="ts">
// Right inspector: edits the selected node's props + classes. Delete/duplicate actions.
// `bare` mode drops the <aside> wrapper so it can be embedded inside RightPanel tabs.
import { computed, inject, ref } from 'vue'
import { IconCopy, IconTrash, IconArrowUp, IconArrowDown, IconUnlink, IconCode, IconSettings, IconTypography, IconPhoto, IconLink } from '@tabler/icons-vue'
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
import { findNode } from '@/components/builder/tree-utils'
import { ICONS, ICON_LIST } from '@/lib/icon-map'
import LayoutSection from './inspector/LayoutSection.vue'
import TypographySection from './inspector/TypographySection.vue'
import SpacingSection from './inspector/SpacingSection.vue'
import BackgroundSection from './inspector/BackgroundSection.vue'
import BorderSection from './inspector/BorderSection.vue'
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

withDefaults(defineProps<{ bare?: boolean }>(), { bare: false })

const node = computed(() => store.selectedNode.value)

function setProp(key: string, value: unknown) {
  if (!node.value) return
  store.patchNode(node.value.id, { props: { ...node.value.props, [key]: value } })
}
function setName(value: string | number) {
  if (!node.value) return
  store.patchNode(node.value.id, { name: String(value) })
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
  <component :is="bare ? 'div' : 'aside'" :class="bare ? 'flex h-full min-h-0 flex-col' : 'flex w-72 flex-col border-l border-border bg-background'">
    <div v-if="!bare" class="border-b border-border px-4 py-3">
      <h2 class="text-sm font-semibold">Inspector</h2>
      <p v-if="node" class="mt-0.5 text-[11px] text-muted-foreground">{{ node.type }} · {{ node.name }}</p>
    </div>

    <div v-if="!node" class="p-4 text-sm text-muted-foreground">
      Select a node on the canvas to edit its properties.
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
  </component>
</template>
