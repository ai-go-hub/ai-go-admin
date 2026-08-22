<template>
    <div class="default-main">
        <TableHeader
            :manager="tableManager"
            v-model:com-search="tableManager.comSearch"
            :buttons="['refresh', 'add', 'edit', 'delete', 'rowExpansion', 'comSearch', 'quickSearch', 'columnDisplay']"
        />

        <!-- 设置合适的 max-height 实现隐藏布局主体部分本身的滚动条，这样就可以监听表格的 @scroll 了 -->
        <!-- max-height = 100vh - (当前布局顶栏高度 + 表头栏高度 + 表格上边距 + 预留的底部下边距) -->
        <Table
            ref="tableRef"
            :pagination="false"
            :manager="tableManager"
            :max-height="`calc(-${adminLayoutHeaderBarHeight[config.layout.mode as keyof typeof adminLayoutHeaderBarHeight] + 75 + 16}px + 100vh)`"
            @scroll="onScroll"
            @expand-change="onExpandChange"
        />

        <DialogForm :manager="tableManager" v-model:form-items="tableManager.form.items!" />
    </div>
</template>

<script setup lang="ts">
import { TableManagerAPI } from '@/api/table'
import TableHeader from '@/components/table/header/index.vue'
import { getDefaultOptButtons } from '@/components/table/index'
import Table from '@/components/table/index.vue'
import { useTableManager } from '@/hooks/useTableManager'
import { useConfig } from '@/stores/config'
import { adminLayoutHeaderBarHeight } from '@/utils/layout'
import { cloneDeep, debounce } from 'lodash-es'
import { nextTick, onMounted, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import DialogForm from './dialogForm.vue'

const { t } = useI18n()
const config = useConfig()
const tableRef = useTemplateRef('tableRef')
const defaultIcon = 'lucide-circle-small'

/**
 * 内存缓存表格的一些状态数据，供数据刷新后恢复
 */
const sessionStateDefault = {
    expanded: [] as any[],
    scrollTop: 0,
    scrollLeft: 0,
    expandAll: false,
}
let sessionState = sessionStateDefault

/**
 * 记录表格行展开状态
 */
const onExpandChange = (row: any, expanded: boolean) => {
    if (expanded) {
        sessionState.expanded.push(row)
    } else {
        sessionState.expanded = sessionState.expanded.filter((item: any) => item.id !== row.id)
    }
}

/**
 * 记录表格滚动条位置
 */
const onScroll = debounce(({ scrollLeft, scrollTop }: { scrollLeft: number; scrollTop: number }) => {
    sessionState.scrollTop = scrollTop
    sessionState.scrollLeft = scrollLeft
}, 500)

/**
 * 记录表格行展开折叠状态
 */
const onUnfoldAll = (state: boolean) => {
    sessionState.expandAll = state
}

/**
 * 恢复已记录的表格状态
 */
const restoreState = () => {
    nextTick(() => {
        const sessionStateTemp = sessionState

        // 重置 sessionState 为默认值，恢复缓存的记录时将自动重设
        sessionState = cloneDeep(sessionStateDefault)

        for (const key in sessionStateTemp.expanded) {
            tableRef.value?.getElTableRef()?.toggleRowExpansion(sessionStateTemp.expanded[key], true)
        }
        nextTick(() => {
            if (sessionStateTemp.scrollTop || sessionStateTemp.scrollLeft) {
                tableRef.value?.getElTableRef()?.scrollTo({ top: sessionStateTemp.scrollTop || 0, left: sessionStateTemp.scrollLeft || 0 })
            }

            /**
             * expandAll 为 “是否默认展开所有行”
             * 此处表格数据已渲染，仅做顶部按钮状态标记用，不会实际上的执行展开折叠操作
             * 展开全部行之后，再只对某一行进行折叠时，expandAll 不会改变，所以此处并不根据 expandAll 值执行折叠展开所有行的操作
             */
            tableManager.table.expandAll = sessionStateTemp.expandAll
            onUnfoldAll(sessionStateTemp.expandAll)
        })
    })
}

const tableManager = useTableManager({
    api: new TableManagerAPI('/admin/auth/rule/'),
    table: {
        filter: {
            sort: 'weigh',
            order: 'desc',
        },
        expandAll: false,
        dragSortLimitField: 'pid',
        dblClickNotEditColumn: ['keepalive', 'status'],
        column: [
            { type: 'selection', align: 'center', operator: false },
            { label: t('auth.rule.title'), prop: 'title', align: 'left', width: 260, operator: 'ILIKE', quickSearch: true },
            { label: 'ID', prop: 'id', align: 'center', operator: 'BETWEEN', width: 70 },
            { label: t('auth.rule.name'), prop: 'name', align: 'center', operator: 'ILIKE', quickSearch: true },
            {
                label: t('auth.rule.type'),
                prop: 'type',
                align: 'center',
                render: 'tag',
                custom: { dir: 'primary', menu: 'success', node: 'warning' },
                dict: {
                    dir: t('auth.rule.typeDir'),
                    menu: t('auth.rule.typeMenu'),
                    node: t('auth.rule.typeNode'),
                },
                width: 100,
            },
            { label: t('common.icon'), prop: 'icon', align: 'center', render: 'icon', operator: false, width: 70 },
            { label: t('auth.rule.path'), prop: 'path', align: 'center', operator: 'ILIKE', show: false, quickSearch: true },
            { label: t('auth.rule.component'), prop: 'component', align: 'center', operator: 'ILIKE', show: false },
            {
                label: t('auth.rule.openType'),
                prop: 'open_type',
                align: 'center',
                dict: {
                    tab: t('auth.rule.openTypeTab'),
                    link: t('auth.rule.openTypeLink'),
                    iframe: t('auth.rule.openTypeIframe'),
                },
                width: 100,
                show: false,
                comSearchRender: 'select',
            },
            { label: t('common.weigh'), prop: 'weigh', align: 'center', sortable: 'custom', operator: 'BETWEEN', width: 90 },
            {
                label: t('common.updatedAt'),
                prop: 'updated_at',
                align: 'center',
                render: 'datetime',
                comSearchRender: 'datetime',
                sortable: 'custom',
                operator: 'BETWEEN',
                width: 160,
                show: false,
            },
            {
                label: t('common.createdAt'),
                prop: 'created_at',
                align: 'center',
                render: 'datetime',
                comSearchRender: 'datetime',
                sortable: 'custom',
                operator: 'BETWEEN',
                width: 160,
            },
            { label: t('auth.rule.keepalive'), prop: 'keepalive', align: 'center', width: 120, render: 'switch' },
            {
                label: t('common.status'),
                prop: 'status',
                align: 'center',
                render: 'switch',
                width: 80,
            },
            {
                label: t('common.operate'),
                align: 'center',
                width: 120,
                render: 'buttons',
                buttons: getDefaultOptButtons(),
                operator: false,
            },
        ],
    },
    form: {
        defaultItems: {
            type: 'menu',
            open_type: 'tab',
            keepalive: 1,
            status: 1,
            icon: defaultIcon,
        },
    },
})

tableManager.opts.before!.getData = () => {
    tableManager.table.expandAll = tableManager.table.filter!.quickSearchKeywords ? true : false
}

tableManager.opts.after!.getEditData = () => {
    if (tableManager.form.items && !tableManager.form.items.icon) {
        tableManager.form.items.icon = defaultIcon
    }
}

// 表格顶部按钮事件触发后的钩子
tableManager.opts.after!.tableEvent = ({ event, data }) => {
    if (event == 'toggle-expansion') {
        onUnfoldAll(data.expanded)
    }
}

// 获取到表格数据后的钩子
tableManager.opts.after!.getData = () => {
    restoreState()
}

onMounted(() => {
    tableManager.table.ref = tableRef.value
    tableManager.initCtx()
    tableManager.getData()?.then(() => {
        tableManager.initDragSort()
    })
})
</script>

<style scoped lang="scss">
.default-main {
    margin-bottom: 0;
}
</style>
