<script setup lang="ts">
// Reusable confirm dialog wrapping shadcn-vue AlertDialog.
// Use case: const confirmed = await useConfirm('Delete?')
import {
  AlertDialog,
  AlertDialogContent,
  AlertDialogTitle,
  AlertDialogDescription,
  AlertDialogCancel,
  AlertDialogAction,
} from '@/components/ui/alert-dialog'
import { buttonVariants } from '@/components/ui/button'
import { cn } from '@/lib/utils'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ 'update:open': [v: boolean]; confirm: [] }>()
</script>

<template>
  <AlertDialog :open="open" @update:open="(v: boolean) => $emit('update:open', v)">
    <AlertDialogContent>
      <AlertDialogTitle><slot name="title">Confirm</slot></AlertDialogTitle>
      <AlertDialogDescription>
        <slot name="description">Are you sure?</slot>
      </AlertDialogDescription>
      <div class="flex justify-end gap-2">
        <AlertDialogCancel :class="cn(buttonVariants({ variant: 'outline', size: 'sm' }), 'text-xs')" @click="$emit('update:open', false)">Cancel</AlertDialogCancel>
        <AlertDialogAction :class="cn(buttonVariants({ variant: 'destructive', size: 'sm' }), 'text-xs')" @click="$emit('confirm')">
          <slot name="confirm">Delete</slot>
        </AlertDialogAction>
      </div>
    </AlertDialogContent>
  </AlertDialog>
</template>
