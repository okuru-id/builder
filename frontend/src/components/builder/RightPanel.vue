<script setup lang="ts">
// Right panel container: tabs between Inspector (style/props) and AI Agent.
import { ref } from 'vue'
import { IconAdjustmentsHorizontal, IconRobotFace } from '@tabler/icons-vue'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import InspectorPanel from './InspectorPanel.vue'
import ChatPanel from './ChatPanel.vue'

// ponytail: default to inspector. Tab state lives here (not in store) — cheap,
// local, and avoids touching the builder contract for a UI-only concern.
const tab = ref<'inspector' | 'agent'>('inspector')
</script>

<template>
  <aside class="flex h-full w-full shrink-0 flex-col border-l border-neutral-200 bg-white">
    <Tabs v-model="tab" class="flex h-full min-h-0 flex-col">
      <div class="flex justify-center border-b border-neutral-200 py-1">
        <TabsList class="h-9">
          <TabsTrigger value="inspector" class="gap-1.5 py-1 text-xs">
            <IconAdjustmentsHorizontal class="size-3.5" /> Inspector
          </TabsTrigger>
          <TabsTrigger value="agent" class="gap-1.5 py-1 text-xs">
            <IconRobotFace class="size-3.5" /> Agent
          </TabsTrigger>
        </TabsList>
      </div>

      <TabsContent value="inspector" class="mt-0 flex-1 overflow-hidden data-[state=inactive]:hidden">
        <InspectorPanel bare />
      </TabsContent>
      <TabsContent value="agent" class="mt-0 flex-1 overflow-hidden data-[state=inactive]:hidden">
        <ChatPanel />
      </TabsContent>
    </Tabs>
  </aside>
</template>
