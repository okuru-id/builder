<script setup lang="ts">
// Builder AI Agent: streams LLM completions over SSE against the
// OpenAI-compatible backend proxy (/admin/api/builder/chat). Uses shadcn-vue
// Message + MessageScroller. The agent may emit fenced ```action:*``` blocks;
// those render as "Apply" buttons that mutate the tree directly.
import { inject, nextTick, ref } from 'vue'
import { IconSend, IconRobot, IconUser, IconLoader2, IconCheck, IconSparkles } from '@tabler/icons-vue'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
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
import { addChild, makeNode } from '@/components/builder/tree-utils'
import type { Node } from '@/types/page-builder'

const store = inject(BUILDER_KEY, null)!

interface ChatMsg {
  id: string
  role: 'user' | 'assistant'
  content: string
}

interface ActionPart {
  kind: 'add' | 'classes' | 'text'
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
const busy = ref(false)
let seq = 0
const uid = () => `m_${Date.now()}_${++seq}`

// Track applied state per action block by its raw signature within a message.
const appliedFlags = ref<Record<string, boolean>>({})

// Parse assistant content into text + action parts.
function partsOf(m: ChatMsg): MsgPart[] {
  if (m.role !== 'assistant') return [{ type: 'text', text: m.content }]
  const out: MsgPart[] = []
  const re = /```action:(add|classes|text)\n([\s\S]*?)```/g
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
  const key = partKey(msgId, action.raw)
  try {
    if (action.kind === 'add') {
      const parentId = action.payload?.parentId === 'root'
        ? store.tree.value.root.id
        : action.payload?.parentId
      const node = normalizeNode(action.payload?.node)
      store.tree.value = { root: addChild(store.tree.value.root, parentId, node) }
      store.select(node.id)
      toast.success('Node added')
    } else if (action.kind === 'classes') {
      store.patchNode(action.payload.nodeId, { classes: action.payload.set })
      toast.success('Classes updated')
    } else if (action.kind === 'text') {
      store.patchNode(action.payload.nodeId, { props: { text: action.payload.text } })
      toast.success('Text updated')
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
      }),
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
    }
  } catch (e: any) {
    assistantMsg.content = `Failed: ${e?.message || 'network error'}`
    toast.error('Agent failed')
  } finally {
    busy.value = false
    await nextTick()
  }
}

function onKey(e: KeyboardEvent) {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    send()
  }
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

                    <!-- Text parts -->
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
                      <div
                        v-else-if="m.role === 'assistant' && busy && m.id === messages[messages.length - 1].id && !m.content && i === 0"
                        class="inline-flex items-center gap-1 rounded-lg bg-muted px-3 py-2 text-xs text-muted-foreground"
                      >
                        <IconLoader2 class="size-3 animate-spin" /> typing…
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
      <Textarea
        v-model="input"
        rows="2"
        class="mb-1.5 resize-none text-xs"
        placeholder="Request changes to this page… (Enter to send)"
        :disabled="busy"
        @keydown="onKey"
      />
      <div class="flex items-center justify-between">
        <span class="text-[10px] text-muted-foreground">AI Agent · LLM_API_KEY</span>
        <Button size="sm" class="h-7 px-2 text-xs" :disabled="!input.trim() || busy" @click="send">
          <IconSend class="size-3.5" /> Send
        </Button>
      </div>
    </div>
  </div>
</template>
