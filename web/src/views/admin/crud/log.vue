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
import { copy } from '@/utils/common'
import { changeStep, state as CRUDState } from '@/views/admin/crud/index'
import { ElMessage } from 'element-plus'
import { reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const model = defineModel<boolean>()
const state = reactive({
    init: false,
})

const optButtons: OptButton[] = [
    {
        render: 'tip',
        name: 'sql',
        icon: 'lucide-database-arrow-down',
        type: 'success',
        title: t('crud.log.copySql'),
        class: 'table-row-sql',
        click: (row) => {
            if (!copy(row.sql)) {
                ElMessage.error(t('common.operationFailed'))
            } else {
                ElMessage.success(t('common.operationSuccess'))
            }
        },
    },
    {
        render: 'tip',
        name: 'context',
        icon: 'lucide-message-circle-code',
        type: 'success',
        title: t('crud.log.copyContext'),
        class: 'table-row-context',
        click: (row) => {
            const context =
                `已生成的 CRUD 代码上下文数据如下: \n` +
                `数据表名: ${row.table.name}\n` +
                `数据模型名: ${row.model_basic_data.name}\n` +
                `模型文件: ${row.model_basic_data.file}\n` +
                `仓储文件: ${row.table.repositoryFile}\n` +
                `服务文件: ${row.table.serviceFile}\n` +
                `控制器文件: ${row.table.handlerFile}\n` +
                `路由注册文件: ${row.table.routerFile}\n` +
                `前端表格组件（路由入口）: ${row.views_basic_data.dir}/index.vue\n` +
                `前端表单组件: ${row.views_basic_data.dir}/dialogForm.vue\n` +
                `前端中文语言包: ${row.lang_basic_data.cn_file}\n` +
                `前端英文语言包: ${row.lang_basic_data.en_file}\n`

            if (!copy(context)) {
                ElMessage.error(t('common.operationFailed'))
            } else {
                ElMessage.success(t('common.operationSuccess'))
            }
        },
    },
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

optButtons[3].title = t('crud.log.delete')
optButtons[3].popconfirm = {
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
                width: 160,
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
