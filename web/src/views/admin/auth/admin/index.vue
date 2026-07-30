<template>
    <div class="default-main">
        <TableHeader
            :manager="tableManager"
            v-model:com-search="tableManager.comSearch"
            :buttons="['refresh', 'add', 'edit', 'delete', 'comSearch', 'quickSearch', 'columnDisplay']"
        />
        <Table :manager="tableManager" />

        <DialogForm :manager="tableManager" v-model:form-items="tableManager.form.items!" />
    </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import DialogForm from './dialogForm.vue'
import { TableManagerAPI } from '/@/api/table'
import TableHeader from '/@/components/table/header/index.vue'
import { getDefaultOptButtons } from '/@/components/table/index'
import Table from '/@/components/table/index.vue'
import { useTableManager } from '/@/hooks/useTableManager'

const { t } = useI18n()
const optButtons = getDefaultOptButtons(['edit', 'delete'])

const tableManager = useTableManager({
    api: new TableManagerAPI('/admin/auth/admin/'),
    table: {
        column: [
            { type: 'selection', align: 'center', operator: false },
            { label: 'ID', prop: 'id', align: 'center', operator: 'eq', width: 70 },
            { label: t('auth.admin.username'), prop: 'username', align: 'center', operator: 'ILIKE', quickSearch: true },
            { label: t('auth.admin.nickname'), prop: 'nickname', align: 'center', operator: 'ILIKE', quickSearch: true },

            // 角色组字段 - 显示专用
            {
                operator: false,
                label: t('auth.admin.groups'),
                prop: 'admin_group_accesses',
                formatter(row, column, cellValue) {
                    const groupNames = []
                    for (const key in cellValue) {
                        groupNames.push(cellValue[key]['group']['name'])
                    }
                    return groupNames
                },
                render: 'tags',
                align: 'center',
                width: 180,
            },
            // 角色组字段 - 筛选专用
            {
                show: false,
                label: t('auth.admin.groups'),
                prop: 'admin_group_accesses.group_id',
                operator: 'IN',
                comSearchRender: 'remoteSelect',
                comSearchRemote: {
                    remoteURL: '/admin/auth/group/list',
                },
            },

            { label: t('auth.admin.avatar'), prop: 'avatar', align: 'center', render: 'image', operator: false },
            { label: t('common.email'), prop: 'email', width: 180, align: 'center', operator: 'ILIKE' },
            { label: t('common.mobile'), prop: 'mobile', align: 'center', operator: 'ILIKE' },
            {
                label: t('common.lastLoginAt'),
                prop: 'last_login_at',
                align: 'center',
                render: 'datetime',
                comSearchRender: 'datetime',
                sortable: 'custom',
                operator: 'BETWEEN',
                width: 160,
            },
            {
                label: t('common.updatedAt'),
                prop: 'updated_at',
                align: 'center',
                render: 'datetime',
                comSearchRender: 'datetime',
                sortable: 'custom',
                operator: 'BETWEEN',
                width: 160,
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
                label: t('common.status'),
                prop: 'status',
                align: 'center',
                render: 'tag',
                custom: { enable: 'success', disable: 'danger' },
                dict: { enable: t('common.enable'), disable: t('common.disable') },
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
        dblClickNotEditColumn: ['status'],
    },
    form: {
        defaultItems: {
            status: 'enable',
            group_ids: [],
        },
    },
})

// 提交前: 将 group_ids 转为 admin_group_accesses 以便关联创建
tableManager.opts.before!.submitForm = ({ items }: { items: AnyObj }) => {
    if (items.group_ids && Array.isArray(items.group_ids)) {
        items.admin_group_accesses = items.group_ids.map((gid: number) => ({ group_id: gid }))
    }
}

// 编辑回填: 将关联表数据的 admin_group_accesses 转为前端表单的 group_ids 字段值供 RemoteSelect 绑定
tableManager.opts.after!.getEditData = ({ res }: { res: ApiResponse<any> }) => {
    const row = res.data.data.row
    if (row.admin_group_accesses && Array.isArray(row.admin_group_accesses)) {
        row.group_ids = row.admin_group_accesses.map((a: any) => a.group_id)
    }
}

tableManager.initCtx()
tableManager.getData()
</script>

<style scoped lang="scss"></style>
