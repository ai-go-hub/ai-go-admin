<template>
    <div class="default-main">
        <TableHeader
            :manager="tableManager"
            v-model:com-search="tableManager.comSearch"
            :buttons="['refresh', 'edit', 'delete', 'comSearch', 'quickSearch', 'columnDisplay']"
        />
        <Table :manager="tableManager" />

        <DialogForm :manager="tableManager" v-model:form-items="tableManager.form.items!" />
    </div>
</template>

<script setup lang="ts">
import { TableManagerAPI } from '@/api/table'
import TableHeader from '@/components/table/header/index.vue'
import { getDefaultOptButtons } from '@/components/table/index'
import Table from '@/components/table/index.vue'
import { useTableManager } from '@/hooks/useTableManager'
import { TableColumnCtx } from 'element-plus'
import { useI18n } from 'vue-i18n'
import DialogForm from './dialogForm.vue'
import { formatFileSize, previewRenderFormatter } from './index'

const { t } = useI18n()
const optButtons = getDefaultOptButtons(['edit', 'delete'])

const tableManager = useTableManager({
    api: new TableManagerAPI('/admin/routine/attachment/'),
    table: {
        column: [
            { type: 'selection', align: 'center', operator: false },
            { label: 'ID', prop: 'id', align: 'center', operator: 'BETWEEN', width: 70 },
            { label: t('routine.attachment.userID'), prop: 'user_id', align: 'center', operator: 'eq', width: 100 },
            { label: t('routine.attachment.userType'), prop: 'user_type', align: 'center', operator: 'ILIKE', width: 120 },
            { label: t('routine.attachment.topic'), prop: 'topic', align: 'center', operator: 'ILIKE', quickSearch: true },
            {
                label: t('routine.attachment.preview'),
                prop: 'url',
                align: 'center',
                operator: false,
                render: 'image',
                width: 80,
                formatter: previewRenderFormatter,
            },
            {
                label: t('routine.attachment.name'),
                prop: 'name',
                width: 160,
                align: 'center',
                operator: 'ILIKE',
                quickSearch: true,
                showOverflowTooltip: true,
            },
            {
                label: t('routine.attachment.size'),
                prop: 'size',
                align: 'center',
                operator: 'BETWEEN',
                comSearchPlaceholder: [t('routine.attachment.minByte'), t('routine.attachment.maxByte')],
                width: 100,
                formatter: (_row: TableRow, _col: TableColumnCtx<TableRow>, val: number) => {
                    return formatFileSize(val)
                },
            },
            { label: t('routine.attachment.mimetype'), prop: 'mimetype', align: 'center', operator: 'ILIKE' },
            { label: t('routine.attachment.quote'), prop: 'quote', align: 'center', operator: 'BETWEEN', width: 100 },
            { label: t('routine.attachment.driver'), prop: 'driver', align: 'center', operator: 'ILIKE' },
            { show: false, label: t('routine.attachment.sha1'), prop: 'sha1', align: 'center', operator: 'ILIKE' },
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
                label: t('routine.attachment.lastUploadAt'),
                prop: 'last_upload_at',
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
                width: 100,
                render: 'buttons',
                buttons: optButtons,
                operator: false,
            },
        ],
        filter: {
            sort: 'last_upload_at',
            order: 'desc',
        },
    },
})

tableManager.initCtx()
tableManager.getData()
</script>

<style scoped lang="scss"></style>
