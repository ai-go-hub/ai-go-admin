<template>
    <el-dialog top="10vh" class="ag-crud-log-dialog" v-model="model" :title="t('crud.index.record')" :append-to-body="true" width="70%">
        <TableHeader
            :manager="tableManager"
            v-model:com-search="tableManager.comSearch"
            :buttons="['refresh', 'delete', 'comSearch', 'quickSearch', 'columnDisplay']"
        />
        <Table :manager="tableManager" />
    </el-dialog>
</template>

<script setup lang="ts">
import { TableManagerAPI } from '@/api/table'
import TableHeader from '@/components/table/header/index.vue'
import { getDefaultOptButtons } from '@/components/table/index'
import Table from '@/components/table/index.vue'
import { useTableManager } from '@/hooks/useTableManager'
import { changeStep, state as CRUDState } from '@/views/admin/crud/index'
import { reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = defineModel<boolean>()
const state = reactive({
    init: false,
})

const optButtons: OptButton[] = [
    {
        render: 'confirm',
        name: 'copy',
        icon: 'lucide-copy-plus',
        type: 'primary',
        title: t('crud.log.continue'),
        class: 'table-row-copy',
        popconfirm: {
            confirmButtonText: t('common.confirm'),
            cancelButtonText: t('common.cancel'),
            confirmButtonType: 'primary',
            title: t('crud.log.logStartTip'),
            width: '220px',
        },
        click: (row) => {
            CRUDState.log.id = row.id
            CRUDState.log.type = 'local'
            changeStep('log')
            model.value = false
        },
    },
    ...getDefaultOptButtons(['delete']),
]

optButtons[1].title = t('crud.log.delete')
optButtons[1].popconfirm = {
    confirmButtonText: t('common.delete'),
    cancelButtonText: t('common.cancel'),
    confirmButtonType: 'danger',
    title: t('crud.log.deleteRecords'),
    width: '248px',
}

const tableManager = useTableManager({
    api: new TableManagerAPI('/admin/crud/log/'),
    table: {
        column: [
            { type: 'selection', align: 'center', operator: false },
            { label: 'ID', prop: 'id', align: 'center', operator: 'BETWEEN', width: 70 },
            { label: t('crud.index.tableName'), prop: 'name', align: 'center', operator: 'ILIKE', quickSearch: true },
            { label: t('crud.log.comment'), prop: 'comment', align: 'center', operator: 'ILIKE', quickSearch: true },
            {
                label: t('common.status'),
                prop: 'status',
                align: 'center',
                render: 'tag',
                custom: {
                    generating: 'warning',
                    succeeded: 'success',
                    failed: 'danger',
                    deleted: 'info',
                },
                dict: {
                    generating: t('crud.log.statusGenerating'),
                    succeeded: t('crud.log.statusSucceeded'),
                    failed: t('crud.log.statusFailed'),
                    deleted: t('crud.log.statusDeleted'),
                },
                width: 100,
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
            sort: 'created_at',
            order: 'desc',
        },
        dblClickNotEditColumn: ['all'],
    },
})

tableManager.initCtx()

// 打开弹窗时加载最新记录
watch(
    () => model.value,
    (val) => {
        if (val && !state.init) {
            state.init = true
            tableManager.getData()
        }
    }
)
</script>

<style scoped lang="scss">
.ag-crud-log-dialog .el-dialog__body {
    padding: 10px 20px;
}
</style>
