<template>
    <div class="default-main">
        <TableHeader
            :manager="tableManager"
            v-model:com-search="tableManager.comSearch"
            :buttons="['refresh', 'comSearch', 'quickSearch', 'columnDisplay']"
        />
        <Table :manager="tableManager" />

        <Info :manager="tableManager" />
    </div>
</template>

<script setup lang="ts">
import { TableManagerAPI } from '@/api/table'
import TableHeader from '@/components/table/header/index.vue'
import Table from '@/components/table/index.vue'
import { useTableManager } from '@/hooks/useTableManager'
import { buildJsonToElTreeData } from '@/utils/common'
import { cloneDeep } from 'lodash-es'
import { useI18n } from 'vue-i18n'
import Info from './info.vue'

const { t } = useI18n()

let optButtons: OptButton[] = [
    {
        render: 'tip',
        name: 'info',
        title: t('common.info'),
        text: '',
        type: 'primary',
        icon: 'lucide-info',
        class: 'table-row-info',
        click: (row: TableRow) => {
            openInfo(row)
        },
    },
]

const tableManager = useTableManager({
    api: new TableManagerAPI('/admin/auth/log/'),
    table: {
        column: [
            { label: 'ID', prop: 'id', align: 'center', operator: 'BETWEEN', width: 70 },
            {
                label: t('auth.adminLog.adminId'),
                prop: 'admin_id',
                formatter(row, column, cellValue) {
                    return cellValue ? cellValue : '-'
                },
                align: 'center',
                operator: 'eq',
                width: 100,
            },
            { label: t('auth.adminLog.username'), prop: 'username', align: 'center', operator: 'ILIKE', quickSearch: true },
            { label: t('auth.adminLog.url'), prop: 'url', align: 'center', operator: 'ILIKE' },
            { label: t('auth.adminLog.title'), prop: 'title', align: 'center', operator: 'ILIKE', quickSearch: true },
            { label: t('auth.adminLog.data'), show: false, prop: 'data', align: 'center', operator: false, width: 200 },
            { label: t('auth.adminLog.ip'), prop: 'ip', align: 'center', operator: 'ILIKE', width: 140 },
            { label: t('auth.adminLog.userAgent'), prop: 'user_agent', showOverflowTooltip: true, align: 'center', operator: 'ILIKE', width: 200 },
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
            {
                label: t('common.operate'),
                align: 'center',
                width: '100',
                render: 'buttons',
                buttons: optButtons,
                operator: false,
            },
        ],
        filter: {
            sort: 'created_at',
            order: 'desc',
        },
    },
})

tableManager.initCtx()
tableManager.getData()

// 利用双击单元格前钩子重写双击操作
tableManager.opts.before!.columnDblclick = ({ row }) => {
    openInfo(row)
    return false
}

/**
 * 点击查看详情按钮响应
 */
const openInfo = (row: TableRow) => {
    if (!row) return
    // 数据来自表格数据，未重新请求 API，深克隆，不然可能会影响表格
    let rowClone = cloneDeep(row)

    rowClone.data = rowClone.data ? [{ label: '点击展开', children: buildJsonToElTreeData(JSON.parse(rowClone.data)) }] : []
    tableManager.form.extend!['info'] = rowClone
    tableManager.form.operate = 'info'
}
</script>

<style scoped lang="scss"></style>
