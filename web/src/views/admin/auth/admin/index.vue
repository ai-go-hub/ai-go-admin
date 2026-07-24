<template>
    <div class="default-main">
        <TableHeader
            :manager="tableManager"
            v-model:com-search="tableManager.comSearch"
            :buttons="['refresh', 'add', 'edit', 'delete', 'comSearch', 'quickSearch', 'columnDisplay']"
        />
        <Table :manager="tableManager" />

        <PopupForm :manager="tableManager" v-model:form-items="tableManager.form.items!" />
    </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import PopupForm from './popupForm.vue'
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
            { label: t('auth.admin.username'), prop: 'username', align: 'center', operator: 'eq', quickSearch: true },
            { label: t('auth.admin.nickname'), prop: 'nickname', align: 'center', operator: 'LIKE', quickSearch: true },
            { label: t('auth.admin.avatar'), prop: 'avatar', align: 'center', render: 'image', operator: false },
            { label: t('common.email'), prop: 'email', width: 180, align: 'center', operator: 'LIKE' },
            { label: t('common.mobile'), prop: 'mobile', align: 'center', operator: 'LIKE' },
            {
                label: t('common.lastLoginAt'),
                prop: 'lastLoginAt',
                align: 'center',
                render: 'datetime',
                comSearchRender: 'datetime',
                sortable: 'custom',
                operator: 'BETWEEN',
                width: 160,
            },
            {
                label: t('common.updatedAt'),
                prop: 'updatedAt',
                align: 'center',
                render: 'datetime',
                comSearchRender: 'datetime',
                sortable: 'custom',
                operator: 'BETWEEN',
                width: 160,
            },
            {
                label: t('common.createdAt'),
                prop: 'createdAt',
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
        },
    },
})
tableManager.initCtx()
tableManager.getData()
</script>

<style scoped lang="scss"></style>
