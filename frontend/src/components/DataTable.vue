<script setup lang="ts" generic="T extends { id: string | number }">
import type {
  ColumnDef,
  ColumnFiltersState,
  SortingState,
  VisibilityState,
} from '@tanstack/vue-table'
import { RestrictToVerticalAxis } from '@dnd-kit/abstract/modifiers'
import {
  IconChevronDown,
  IconChevronLeft,
  IconChevronRight,
  IconChevronsLeft,
  IconChevronsRight,
  IconLayoutColumns,
  IconInbox,
  IconSearch,
} from '@tabler/icons-vue'
import {
  FlexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getPaginationRowModel,
  getSortedRowModel,
  useVueTable,
} from '@tanstack/vue-table'
import { DragDropProvider } from 'dnd-kit-vue'
import { h, ref } from 'vue'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import DraggableRow from './DraggableRow.vue'
import DragHandle from './DragHandle.vue'
import {
  DropdownMenu,
  DropdownMenuCheckboxItem,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

const emit = defineEmits<{
  rowClick: [row: T]
}>()

const props = withDefaults(defineProps<{
  data: T[]
  columns: ColumnDef<T>[]
  loading?: boolean
  searchKey?: string
  searchPlaceholder?: string
  draggable?: boolean
}>(), {
  loading: false,
  searchKey: undefined,
  searchPlaceholder: 'Search…',
  draggable: true,
})

const sorting = ref<SortingState>([])
const columnFilters = ref<ColumnFiltersState>([])
const columnVisibility = ref<VisibilityState>({})
const globalFilter = ref('')

// ponytail: drag = visual reorder only, no persistence yet.
// Add: emit @reorder(orderedIds) when a reorder API endpoint exists.
const allColumns = props.draggable
  ? [
      { id: 'drag', header: () => null, cell: () => h(DragHandle), enableSorting: false, enableHiding: false },
      ...props.columns,
    ] as ColumnDef<T>[]
  : props.columns

const table = useVueTable({
  get data() { return props.data },
  get columns() { return allColumns },
  getCoreRowModel: getCoreRowModel(),
  getPaginationRowModel: getPaginationRowModel(),
  getSortedRowModel: getSortedRowModel(),
  getFilteredRowModel: getFilteredRowModel(),
  onSortingChange: (u) => { sorting.value = typeof u === 'function' ? u(sorting.value) : u },
  onColumnFiltersChange: (u) => { columnFilters.value = typeof u === 'function' ? u(columnFilters.value) : u },
  onColumnVisibilityChange: (u) => { columnVisibility.value = typeof u === 'function' ? u(columnVisibility.value) : u },
  onGlobalFilterChange: (u) => { globalFilter.value = typeof u === 'function' ? u(globalFilter.value) : u },
  state: {
    get sorting() { return sorting.value },
    get columnFilters() { return columnFilters.value },
    get columnVisibility() { return columnVisibility.value },
    get globalFilter() { return globalFilter.value },
  },
})

// route global filter to a single column if searchKey given
function onSearch(value: string) {
  if (props.searchKey) {
    table.getColumn(props.searchKey)?.setFilterValue(value)
  } else {
    globalFilter.value = value
  }
}
</script>

<template>
  <div class="flex w-full flex-col gap-4">
    <!-- Toolbar -->
    <div class="flex items-center gap-2">
      <div v-if="searchKey" class="relative w-full max-w-xs">
        <IconSearch class="text-muted-foreground pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2" />
        <Input
          :model-value="globalFilter"
          :placeholder="searchPlaceholder"
          class="h-8 w-full pl-8"
          @update:model-value="onSearch(String($event))"
        />
      </div>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <Button variant="outline" size="sm" class="ml-auto">
            <IconLayoutColumns />
            <span class="hidden lg:inline">Columns</span>
            <IconChevronDown />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" class="w-40">
          <template v-for="column in table.getAllColumns().filter((c) => c.getCanHide())" :key="column.id">
            <DropdownMenuCheckboxItem
              class="capitalize"
              :model-value="column.getIsVisible()"
              @update:model-value="(value: any) => column.toggleVisibility(!!value)"
            >
              {{ column.id }}
            </DropdownMenuCheckboxItem>
          </template>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>

    <!-- Table -->
    <div class="overflow-hidden rounded-lg border">
      <component :is="draggable ? DragDropProvider : 'div'" :modifiers="draggable ? [RestrictToVerticalAxis] : undefined">
        <Table>
          <TableHeader class="bg-muted sticky top-0 z-10">
            <TableRow v-for="headerGroup in table.getHeaderGroups()" :key="headerGroup.id">
              <TableHead v-for="header in headerGroup.headers" :key="header.id" :colspan="header.colSpan">
                <FlexRender v-if="!header.isPlaceholder" :render="header.column.columnDef.header" :props="header.getContext()" />
              </TableHead>
            </TableRow>
          </TableHeader>
          <TableBody class="**:data-[slot=table-cell]:first:w-8">
            <template v-if="loading">
              <TableRow v-for="i in 3" :key="i">
                <TableCell :colspan="allColumns.length" class="h-12">
                  <div class="bg-muted h-4 w-full animate-pulse rounded" />
                </TableCell>
              </TableRow>
            </template>
            <template v-else-if="table.getRowModel().rows.length">
              <template v-if="draggable">
                <DraggableRow v-for="row in table.getRowModel().rows" :key="row.id" :row="row" :index="row.index" />
              </template>
              <template v-else>
                <TableRow
                  v-for="row in table.getRowModel().rows"
                  :key="row.id"
                  class="cursor-pointer hover:bg-muted/50"
                  @click="emit('rowClick', row.original as T)"
                >
                  <TableCell v-for="cell in row.getVisibleCells()" :key="cell.id">
                    <FlexRender :render="cell.column.columnDef.cell" :props="cell.getContext()" />
                  </TableCell>
                </TableRow>
              </template>
            </template>
            <TableRow v-else>
              <TableCell :colspan="allColumns.length" class="h-48 text-center">
                <div class="flex flex-col items-center justify-center gap-3">
                  <div class="bg-muted/50 flex size-16 items-center justify-center rounded-full">
                    <IconInbox class="text-muted-foreground/60 size-8" />
                  </div>
                  <p class="text-muted-foreground text-sm font-medium">No data yet</p>
                  <p class="text-muted-foreground/50 text-xs max-w-xs">Data will appear here once you add some entries.</p>
                </div>
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </component>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between">
      <div class="text-muted-foreground hidden flex-1 text-sm lg:flex">
        {{ table.getFilteredRowModel().rows.length }} row(s).
      </div>
      <div class="flex w-full items-center gap-8 lg:w-fit">
        <div class="hidden items-center gap-2 lg:flex">
          <Label for="rows-per-page" class="text-sm font-medium">Rows per page</Label>
          <Select
            :model-value="String(table.getState().pagination.pageSize)"
            @update:model-value="(value: any) => table.setPageSize(Number(value))"
          >
            <SelectTrigger id="rows-per-page" size="sm" class="w-20">
              <SelectValue />
            </SelectTrigger>
            <SelectContent side="top">
              <SelectItem v-for="pageSize in [10, 20, 30, 40, 50]" :key="pageSize" :value="String(pageSize)">
                {{ pageSize }}
              </SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div class="flex w-fit items-center justify-center text-sm font-medium">
          Page {{ table.getState().pagination.pageIndex + 1 }} of {{ table.getPageCount() }}
        </div>
        <div class="ml-auto flex items-center gap-2 lg:ml-0">
          <Button variant="outline" class="hidden h-8 w-8 p-0 lg:flex" :disabled="!table.getCanPreviousPage()" @click="table.setPageIndex(0)">
            <span class="sr-only">First page</span>
            <IconChevronsLeft />
          </Button>
          <Button variant="outline" size="icon" class="size-8" :disabled="!table.getCanPreviousPage()" @click="table.previousPage()">
            <span class="sr-only">Previous</span>
            <IconChevronLeft />
          </Button>
          <Button variant="outline" size="icon" class="size-8" :disabled="!table.getCanNextPage()" @click="table.nextPage()">
            <span class="sr-only">Next</span>
            <IconChevronRight />
          </Button>
          <Button variant="outline" class="hidden size-8 lg:flex" :disabled="!table.getCanNextPage()" @click="table.setPageIndex(table.getPageCount() - 1)">
            <span class="sr-only">Last page</span>
            <IconChevronsRight />
          </Button>
        </div>
      </div>
    </div>
  </div>
</template>
