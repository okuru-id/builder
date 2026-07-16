<script setup lang="ts">
// Builder AI Agent: streams LLM completions over SSE against the
// OpenAI-compatible backend proxy (/admin/api/builder/chat). Uses shadcn-vue
// Message + MessageScroller. The agent may emit fenced ```action:*``` blocks;
// those render as "Apply" buttons that mutate the tree directly.
import { computed, inject, nextTick, onMounted, ref, watch } from 'vue'
import { IconSend, IconRobot, IconUser, IconLoader2, IconCheck, IconSparkles, IconPlayerStop, IconBolt, IconPlus } from '@tabler/icons-vue'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { Textarea } from '@/components/ui/textarea'
import {
  Message,
  MessageContent,
  MessageHeader,
} from '@/components/ui/message'
import {
  MessageScroller,
  MessageScrollerButton,
  MessageScrollerContent,
  MessageScrollerItem,
  MessageScrollerProvider,
  MessageScrollerViewport,
} from '@/components/ui/message-scroller'
import { BUILDER_KEY } from '@/components/builder/injection'
import { makeNode, findNode } from '@/components/builder/tree-utils'
import { cn } from '@/lib/utils'
import { CONTAINER_TYPES, PALETTE_TYPES } from '@/types/page-builder'
import type { Node } from '@/types/page-builder'

const store = inject(BUILDER_KEY, null)!

interface ChatMsg {
  id: string
  role: 'user' | 'assistant'
  content: string
}

interface ActionPart {
  kind: 'add' | 'classes' | 'text' | 'delete' | 'move'
  payload: any
  raw: string
}
interface MsgPart {
  type: 'text' | 'action'
  text?: string
  action?: ActionPart
  applied?: boolean
}

const messages = ref<ChatMsg[]>([
  {
    id: 'intro',
    role: 'assistant',
    content: "Hi! I'm the AI Agent. Ask for layout/color/structure changes and I'll propose actions you can apply with one click.",
  },
])
const input = ref('')
const focusedNode = ref<Node | null>(null)
const busy = ref(false)
const autoApply = ref(true)
const showAutoApplyConfirm = ref(false)
let seq = 0
const uid = () => `m_${Date.now()}_${++seq}`

// Active AbortController for the in-flight stream, if any.
let abortCtl: AbortController | null = null

// Persist conversation per-page in localStorage so it survives refresh.
// ponytail: client-side only; cross-device sync is YAGNI for a builder assistant.
const storageKey = computed(() => `builder_chat_${store.page.value?.id ?? 'draft'}`)

function loadChat() {
  try {
    const raw = localStorage.getItem(storageKey.value)
    if (raw) {
      const parsed = JSON.parse(raw)
      if (Array.isArray(parsed) && parsed.length) {
        messages.value = parsed
        return
      }
    }
  } catch {
    // corrupt entry — ignore, keep intro
  }
  // No saved chat: reset to intro.
  messages.value = [{
    id: 'intro',
    role: 'assistant',
    content: "Hi! I'm the AI Agent. Ask for layout/color/structure changes and I'll propose actions you can apply with one click.",
  }]
}

function saveChat() {
  try {
    localStorage.setItem(storageKey.value, JSON.stringify(messages.value))
  } catch {
    // quota / private mode — ignore, chat just won't persist
  }
}

// Start a fresh conversation: drop history + persist empty state.
function newSession() {
  if (busy.value) stop()
  messages.value = [{
    id: 'intro',
    role: 'assistant',
    content: "Hi! I'm the AI Agent. Ask for layout/color/structure changes and I'll propose actions you can apply with one click.",
  }]
  saveChat()
}

// Track applied state per action block by its raw signature within a message.
const appliedFlags = ref<Record<string, boolean>>({})

// Parse assistant content into text + action parts.
function partsOf(m: ChatMsg): MsgPart[] {
  if (m.role !== 'assistant') return [{ type: 'text', text: m.content }]
  const out: MsgPart[] = []
  const re = /```action:(add|classes|text|delete|move)\n([\s\S]*?)```/g
  let last = 0
  let match: RegExpExecArray | null
  while ((match = re.exec(m.content)) !== null) {
    if (match.index > last) {
      out.push({ type: 'text', text: m.content.slice(last, match.index).trim() })
    }
    const kind = match[1] as ActionPart['kind']
    let payload: any = null
    try {
      payload = JSON.parse(match[2].trim())
    } catch {
      payload = null
    }
    out.push({
      type: 'action',
      action: { kind, payload, raw: match[0] },
    })
    last = re.lastIndex
  }
  if (last < m.content.length) {
    const tail = m.content.slice(last).trim()
    if (tail) out.push({ type: 'text', text: tail })
  }
  return out
}

function partKey(msgId: string, raw: string) {
  return `${msgId}::${raw}`
}

// Normalize an agent-supplied node spec into a full Node with ids.
function normalizeNode(spec: any): Node {
  const type = spec?.type ?? 'frame'
  const base = makeNode(type, {
    name: spec?.name,
    props: spec?.props,
    classes: Array.isArray(spec?.classes) ? spec.classes : [],
  })
  const node: Node = {
    ...base,
    ...spec,
    type,
    id: spec?.id ?? base.id,
    name: spec?.name ?? base.name,
    props: spec?.props ?? base.props,
    classes: Array.isArray(spec?.classes) ? spec.classes : base.classes,
    children: Array.isArray(spec?.children) ? spec.children.map(normalizeNode) : [],
  }
  return node
}

function applyAction(msgId: string, part: MsgPart) {
  const action = part.action!
  if (!action.payload) {
    toast.error('Invalid agent action')
    return
  }
  const key = partKey(msgId, action.raw)
  try {
    if (action.kind === 'add') {
      const parentId = action.payload?.parentId && action.payload.parentId !== 'root'
        ? action.payload.parentId
        : store.tree.value.root.id
      const node = normalizeNode(action.payload?.node)
      store.appendNode(node, parentId)
      toast.success('Node added')
    } else if (action.kind === 'classes') {
      // Merge via twMerge: agent classes override conflicting prefixes,
      // unrelated classes (color/spacing/...) are preserved. Fixes data loss.
      const cur = findNode(store.tree.value.root, action.payload.nodeId)?.node.classes ?? []
      const merged = cn(cur, ...(Array.isArray(action.payload.set) ? action.payload.set : []))
        .split(/\s+/)
        .filter(Boolean)
      store.patchNode(action.payload.nodeId, { classes: merged })
      toast.success('Classes updated')
    } else if (action.kind === 'text') {
      store.patchNode(action.payload.nodeId, { props: { text: action.payload.text } })
      toast.success('Text updated')
    } else if (action.kind === 'delete') {
      store.removeNode(action.payload.nodeId)
      toast.success('Node deleted')
    } else if (action.kind === 'move') {
      const parentId = action.payload?.parentId && action.payload.parentId !== 'root'
        ? action.payload.parentId
        : null
      store.moveNode(action.payload.nodeId, parentId, Number(action.payload?.index ?? -1))
      toast.success('Node moved')
    }
    appliedFlags.value[key] = true
  } catch (e: any) {
    toast.error(`Failed to apply: ${e?.message ?? 'error'}`)
  }
}

async function send() {
  const text = input.value.trim()
  if (!text || busy.value) return
  input.value = ''
  const focusForRequest = focusedNode.value
  focusedNode.value = null

  const userMsg: ChatMsg = { id: uid(), role: 'user', content: text }
  const assistantMsg: ChatMsg = { id: uid(), role: 'assistant', content: '' }
  messages.value.push(userMsg, assistantMsg)
  busy.value = true

  // ponytail: pass prior turns (drop the empty assistant placeholder) so the
  // LLM has conversational context. Tree is included server-side via system
  // prompt — we only send the user/assistant text history here.
  const history = messages.value
    .filter((m) => m.id !== assistantMsg.id && m.id !== 'intro')
    .map((m) => ({ role: m.role, content: m.content }))

  abortCtl = new AbortController()
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch('/admin/api/builder/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        Accept: 'text/event-stream',
      },
      body: JSON.stringify({
        messages: history,
        tree: store.tree.value,
        pageName: store.page.value?.name ?? '',
        focusNode: focusForRequest,
        nodeCatalog: {
          types: PALETTE_TYPES,
          containers: [...CONTAINER_TYPES],
        },
      }),
      signal: abortCtl.signal,
    })
    if (!res.ok || !res.body) {
      throw new Error(`HTTP ${res.status}`)
    }

    const reader = res.body.getReader()
    const decoder = new TextDecoder()
    let buf = ''
    while (true) {
      const { value, done } = await reader.read()
      if (done) break
      buf += decoder.decode(value, { stream: true })
      const events = buf.split('\n\n')
      buf = events.pop() ?? ''
      for (const evt of events) {
        const line = evt.split('\n').find((l) => l.startsWith('data: '))
        if (!line) continue
        const data = line.slice(6)
        if (data === '[DONE]') continue
        try {
          const payload = JSON.parse(data)
          if (payload.content) {
            assistantMsg.content += payload.content
          }
        } catch {
          // skip malformed
        }
      }
    }
    if (!assistantMsg.content) {
      assistantMsg.content = '(no response)'
    } else if (autoApply.value) {
      applyAll(assistantMsg)
    }
  } catch (e: any) {
    if (e?.name === 'AbortError') {
      // Keep partially streamed content; only mark empty as cancelled.
      if (!assistantMsg.content) assistantMsg.content = '(cancelled)'
    } else {
      assistantMsg.content = `Failed: ${e?.message || 'network error'}`
      toast.error('Agent failed')
    }
  } finally {
    abortCtl = null
    busy.value = false
    saveChat()
    await nextTick()
  }
}

function applyAll(msg: ChatMsg) {
  for (const part of partsOf(msg)) {
    if (part.type === 'action' && part.action && !appliedFlags.value[partKey(msg.id, part.action.raw)]) {
      applyAction(msg.id, part)
    }
  }
}

function requestAutoApply() {
  if (autoApply.value) {
    autoApply.value = false
    return
  }
  showAutoApplyConfirm.value = true
}

function confirmAutoApply() {
  autoApply.value = true
  showAutoApplyConfirm.value = false
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
}

// Cancel the in-flight stream. Keep whatever has streamed so far.
function stop() {
  abortCtl?.abort()
}

onMounted(loadChat)
// Reload chat when the builder switches to a different page.
watch(() => store.page.value?.id, () => loadChat())
// Consume a one-shot focus request from the Layer Tree.
watch(() => store.agentFocus.value, (node) => {
  if (!node) return
  focusedNode.value = node
  input.value = `Focus on node "${node.name || node.type}" (ID: "${node.id}"). `
  store.clearAgentFocus()
})

// Drag a node from the Layer Tree onto the composer to reference it.
// Reads the custom 'application/x-builder-node' payload set by TreeRow.
function onDropNode(e: DragEvent) {
  const id = e.dataTransfer?.getData('application/x-builder-node')
  if (!id) return
  e.preventDefault()
  const found = findNode(store.tree.value.root, id)
  if (!found) return
  store.askAgentAbout(found.node)
}
</script>

<template>
  <div class="flex h-full flex-col">
    <!-- Transcript -->
    <div class="min-h-0 flex-1 p-2">
      <MessageScrollerProvider auto-scroll default-scroll-position="last-anchor">
        <MessageScroller class="h-full">
          <MessageScrollerViewport class="h-full">
            <MessageScrollerContent class="space-y-1 px-1 py-2">
              <MessageScrollerItem
                v-for="m in messages"
                :key="m.id"
                :message-id="m.id"
                :scroll-anchor="m.role === 'user'"
              >
                <Message :align="m.role === 'user' ? 'end' : 'start'">
                  <div
                    class="flex size-7 shrink-0 items-center justify-center rounded-full"
                    :class="m.role === 'user' ? 'bg-primary text-primary-foreground' : 'bg-muted text-muted-foreground'"
                  >
                    <component :is="m.role === 'user' ? IconUser : IconRobot" class="size-4" />
                  </div>
                  <MessageContent>
                    <MessageHeader class="px-0">
                      {{ m.role === 'user' ? 'You' : 'AI Agent' }}
                    </MessageHeader>

                    <!-- Typing indicator: only while waiting for first token -->
                    <div
                      v-if="m.role === 'assistant' && busy && !m.content && m.id === messages[messages.length - 1].id"
                      class="inline-flex items-center gap-1 rounded-lg bg-muted px-3 py-2 text-xs text-muted-foreground"
                    >
                      <IconLoader2 class="size-3 animate-spin" /> typing…
                    </div>

                    <!-- Text + action parts -->
                    <div
                      v-for="(part, i) in partsOf(m)"
                      :key="i"
                    >
                      <div
                        v-if="part.type === 'text' && part.text"
                        class="rounded-lg px-3 py-2 text-xs leading-relaxed whitespace-pre-wrap"
                        :class="m.role === 'user'
                          ? 'bg-primary text-primary-foreground'
                          : 'bg-muted text-foreground'"
                      >
                        {{ part.text }}
                      </div>

                      <!-- Action block: Apply button -->
                      <div
                        v-else-if="part.type === 'action' && part.action"
                        class="mt-1 flex items-center gap-2 rounded-lg border border-primary/30 bg-primary/5 px-2.5 py-1.5"
                      >
                        <IconSparkles class="size-3.5 shrink-0 text-primary" />
                        <span class="flex-1 text-[11px] text-muted-foreground">
                          Action: <span class="font-mono">{{ part.action.kind }}</span>
                        </span>
                        <Button
                          size="sm"
                          class="h-6 gap-1 px-2 text-[11px]"
                          :variant="appliedFlags[partKey(m.id, part.action.raw)] ? 'secondary' : 'default'"
                          :disabled="appliedFlags[partKey(m.id, part.action.raw)]"
                          @click="applyAction(m.id, part)"
                        >
                          <template v-if="appliedFlags[partKey(m.id, part.action.raw)]">
                            <IconCheck class="size-3" /> Applied
                          </template>
                          <template v-else>Apply</template>
                        </Button>
                      </div>
                    </div>
                  </MessageContent>
                </Message>
              </MessageScrollerItem>
            </MessageScrollerContent>
          </MessageScrollerViewport>
          <MessageScrollerButton direction="end" />
        </MessageScroller>
      </MessageScrollerProvider>
    </div>

    <!-- Composer -->
    <div class="border-t border-border p-2">
      <div v-if="autoApply" class="mb-2 flex items-center justify-between gap-2 rounded-md border border-destructive/40 bg-destructive/10 px-2 py-1.5 text-[11px] text-destructive">
        <span class="flex items-center gap-1"><IconBolt class="size-3.5" /> Auto-apply aktif — termasuk delete.</span>
        <Button size="sm" variant="destructive" class="h-6 px-2 text-[10px]" @click="autoApply = false">Stop</Button>
      </div>
      <Textarea
        v-model="input"
        rows="2"
        class="mb-1.5 max-h-40 resize-none overflow-y-auto text-xs"
        placeholder="Request changes to this page… (Enter to send, or drag a node from the tree)"
        :disabled="busy"
        @keydown="onKey"
        @dragover.prevent
        @drop="onDropNode"
      />
      <div class="flex items-center justify-between gap-2">
        <Button size="sm" variant="outline" class="h-7 gap-1 px-2 text-[10px]" :class="autoApply ? 'border-destructive text-destructive' : ''" @click="requestAutoApply">
          <IconBolt class="size-3" /> Auto-apply {{ autoApply ? 'on' : 'off' }}
        </Button>
        <div class="flex items-center gap-1">
          <Button size="sm" variant="ghost" class="h-7 gap-1 px-2 text-[10px] text-muted-foreground" title="New session" :disabled="busy" @click="newSession">
            <IconPlus class="size-3.5" /> New
          </Button>
          <Button v-if="busy" size="sm" variant="secondary" class="h-7 px-2 text-xs" @click="stop">
            <IconPlayerStop class="size-3.5" /> Stop
          </Button>
          <Button size="sm" class="h-7 px-2 text-xs" :disabled="!input.trim() || busy" @click="send">
            <IconSend class="size-3.5" /> Send
          </Button>
        </div>
      </div>
    </div>
    <ConfirmDialog :open="showAutoApplyConfirm" @update:open="showAutoApplyConfirm = $event" @confirm="confirmAutoApply">
      <template #title>Enable auto-apply?</template>
      <template #description>All future agent actions in this session apply automatically, including delete and move. Undo remains available.</template>
      <template #confirm>Enable auto-apply</template>
    </ConfirmDialog>
  </div>
</template>
