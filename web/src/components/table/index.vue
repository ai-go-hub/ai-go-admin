<template>
    <div>
        <slot name="header"></slot>
        <el-table
            ref="tableRef"
            class="ag-data-table"
            header-cell-class-name="table-header-cell"
            :default-expand-all="manager.table.expandAll"
            :data="manager.table.data"
            :row-key="manager.table.pk"
            :border="true"
            :default-sort="defaultSort"
            v-loading="manager.table.loading"
            stripe
            @select="onSelect"
            @select-all="onSelectAll"
            @selection-change="onSelectionChange"
            @sort-change="onSortChange"
            @row-dblclick="onTableDblclick"
            v-bind="$attrs"
        >
            <slot name="columnPrepend"></slot>
            <template v-for="(item, key) in manager.table.column">
                <template v-if="item.show !== false">
                    <!-- 渲染为 slot -->
                    <slot v-if="item.render == 'slot'" :name="item.slotName"></slot>

                    <el-table-column
                        v-else
                        :key="`${key}-column`"
                        v-bind="item"
                        :column-key="item['columnKey'] || (item.prop ? `table-column-${item.prop}` : '') || uuid()"
                    >
                        <!-- ./cellRenderer/ 文件夹内的每个组件为一种字段渲染器，组件名称为渲染器名称 -->
                        <template v-if="item.render" #default="scope">
                            <component
                                :row="scope.row"
                                :columnConfig="item"
                                :column="scope.column"
                                :cellValue="getCellValue(scope.row, item, scope.column, scope.$index)"
                                :index="scope.$index"
                                :manager="props.manager"
                                :is="cellRenderer[item.render] ?? cellRenderer['default']"
                                :key="getRenderKey(key, item, scope)"
                            />
                        </template>
                    </el-table-column>
                </template>
            </template>
            <slot name="columnAppend"></slot>
        </el-table>

        <div v-if="props.pagination && manager.table.total" class="ag-table-pagination">
            <el-pagination
                :page-size="manager.table.filter!.limit"
                :page-sizes="pageSizes"
                :current-page="manager.table.filter!.page"
                background
                :layout="config.layout.shrink ? 'prev, next, jumper' : 'sizes,total, ->, prev, pager, next, jumper'"
                :total="manager.table.total"
                @update:page-size="onTableSizeChange"
                @update:current-page="onTableCurrentChange"
            ></el-pagination>
        </div>

        <slot name="footer"></slot>
    </div>
</template>

<script setup lang="ts">
import { getCellValue } from '@/components/table/index'
import { useConfig } from '@/stores/config'
import { uuid } from '@/utils/random'
import type { ElTable, Sort, TableColumnCtx } from 'element-plus'
import type { Component } from 'vue'
import { computed, nextTick, useTemplateRef } from 'vue'

const config = useConfig()
const tableRef = useTemplateRef('tableRef')
type ElTableProps = Partial<InstanceType<typeof ElTable>['$props']>

interface Props extends /* @vue-ignore */ ElTableProps {
    manager: TableManagerInstance
    pagination?: boolean
}
const props = withDefaults(defineProps<Props>(), {
    pagination: true,
})

const cellRenderer: Record<string, Component> = {}
const cellRendererComponents: Record<string, any> = import.meta.glob('./cellRenderer/**.vue', { eager: true })
for (const key in cellRendererComponents) {
    const fileName = key.replace('./cellRenderer/', '').replace('.vue', '')
    cellRenderer[fileName] = cellRendererComponents[key].default
}

const getRenderKey = (key: number, item: TableColumn, scope: any) => {
    if (item.getRenderKey && typeof item.getRenderKey == 'function') {
        return item.getRenderKey(scope.row, item, scope.column, scope.$index)
    }
    if (item.render == 'switch') {
        return item.render + item.prop
    }
    return key + scope.$index + '-' + item.render + '-' + (item.prop ? '-' + item.prop + '-' + scope.row[item.prop] : '')
}

const onTableSizeChange = (val: number) => {
    props.manager.handleEvent('page-size-change', { size: val })
}

const onTableCurrentChange = (val: number) => {
    props.manager.handleEvent('current-page-change', { page: val })
}

const onSortChange = ({ order, prop }: { order: string; prop: string }) => {
    props.manager.handleEvent('sort-change', { prop: prop, order: order ? (order == 'descending' ? 'desc' : 'asc') : '' })
}

const onTableDblclick = (row: TableRow, column: TableColumnCtx<TableRow>) => {
    props.manager.handleEvent('column-dblclick', { row, column })
}

const onSelectionChange = (selections: TableRow[]) => {
    props.manager.handleEvent('selection-change', selections)
}

const pageSizes = computed(() => {
    let defaultSizes = [10, 20, 50, 100]
    if (props.manager.table.filter!.limit) {
        if (!defaultSizes.includes(props.manager.table.filter!.limit)) {
            defaultSizes.push(props.manager.table.filter!.limit)
        }
    }
    return defaultSizes
})

const defaultSort = computed(() => {
    if (props.manager.table.filter?.sort) {
        return {
            prop: props.manager.table.filter?.sort,
            order: props.manager.table.filter?.order == 'desc' ? 'descending' : 'ascending',
        } as Sort
    }
    return undefined
})

/*
 * 全选和取消全选
 * 实现子级同时选择和取消选中
 */
const onSelectAll = (selections: TableRow[]) => {
    if (isSelectAll(selections.map((row: TableRow) => row[props.manager.table.pk!].toString()))) {
        selections.map((row: TableRow) => {
            if (row.children) {
                selectChildren(row.children, true)
            }
        })
    } else {
        tableRef.value?.clearSelection()
    }
}

/*
 * 是否是全选操作
 * 只检查第一个元素是否被选择
 * 全选时: selectIDs为所有元素的id
 * 取消全选时: selectIDs为所有子元素的id
 */
const isSelectAll = (selectIDs: string[]) => {
    let data = props.manager.table.data as TableRow[]
    for (const key in data) {
        return selectIDs.includes(data[key][props.manager.table.pk!].toString())
    }
    return false
}

/*
 * 选择子项-递归
 */
const selectChildren = (children: TableRow[], type: boolean) => {
    children.map((j: TableRow) => {
        toggleSelection(j, type)
        if (j.children) {
            selectChildren(j.children, type)
        }
    })
}

/*
 * 执行选择操作
 */
const toggleSelection = (row: TableRow, type: boolean) => {
    if (row) {
        nextTick(() => {
            tableRef.value?.toggleRowSelection(row, type)
        })
    }
}

/*
 * 手动选择时，同时选择子级
 */
const onSelect = (selections: TableRow[], row: TableRow) => {
    if (
        selections.some((item: TableRow) => {
            return row[props.manager.table.pk!] === item[props.manager.table.pk!]
        })
    ) {
        if (row.children) {
            selectChildren(row.children, true)
        }
    } else {
        if (row.children) {
            selectChildren(row.children, false)
        }
    }
}

/*
 * 设置折叠所有-递归
 */
const toggleExpansion = (children: TableRow[], expanded: boolean) => {
    for (const key in children) {
        tableRef.value?.toggleRowExpansion(children[key], expanded)
        if (children[key].children) {
            toggleExpansion(children[key].children!, expanded)
        }
    }
}

/**
 * 获取 el-table 的 ref
 */
const getElTableRef = () => {
    return tableRef.value
}

/*
 * 折叠所有
 */
const toggleExpansionAll = (expanded: boolean) => {
    toggleExpansion(props.manager.table.data!, expanded)
}

defineExpose({
    getElTableRef,
    toggleExpansionAll,
})
</script>

<style scoped lang="scss">
.ag-data-table {
    width: 100%;
    :deep(.table-header-cell) .cell {
        color: var(--el-text-color-primary);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
}
.ag-table-pagination {
    box-sizing: border-box;
    width: 100%;
    max-width: 100%;
    background-color: var(--ag-bg-color-overlay);
    padding: 13px 15px;
}
</style>
