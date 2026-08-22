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

        <DialogForm ref="formRef" :manager="tableManager" v-model:form-items="tableManager.form.items!" />
    </div>
</template>

<script setup lang="ts">
import { getAdminRules } from '@/api/admin/auth/index'
import { TableManagerAPI } from '@/api/table'
import TableHeader from '@/components/table/header/index.vue'
import { getDefaultOptButtons } from '@/components/table/index'
import Table from '@/components/table/index.vue'
import { useTableManager } from '@/hooks/useTableManager'
import { useConfig } from '@/stores/config'
import { adminLayoutHeaderBarHeight } from '@/utils/layout'
import { uuid } from '@/utils/random'
import { cloneDeep, debounce } from 'lodash-es'
import { nextTick, onMounted, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import DialogForm from './dialogForm.vue'

const { t } = useI18n()
const config = useConfig()
const formRef = useTemplateRef('formRef')
const tableRef = useTemplateRef('tableRef')

const tableManager = useTableManager({
    api: new TableManagerAPI('/admin/auth/group/'),
    table: {
        filter: {
            order: 'desc',
            sort: 'created_at',
        },
        expandAll: false,
        dblClickNotEditColumn: ['status'],
        column: [
            { type: 'selection', align: 'center', operator: false },
            { label: t('auth.group.name'), prop: 'name', align: 'left', width: 280, operator: 'ILIKE', quickSearch: true },
            { label: 'ID', prop: 'id', align: 'center', operator: 'BETWEEN', width: 70 },
            { label: t('auth.group.rules'), prop: 'rules_title', align: 'center', operator: false },
            {
                label: t('common.status'),
                prop: 'status',
                align: 'center',
                render: 'switch',
                width: 100,
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
                label: t('common.operate'),
                align: 'center',
                width: 120,
                render: 'buttons',
                buttons: getDefaultOptButtons(['edit', 'delete']),
                operator: false,
            },
        ],
    },
    form: {
        defaultItems: {
            status: 1,
        },
    },
})

// 利用提交前钩子重写提交操作
tableManager.opts.before!.submitForm = ({ formEl, operate, items }) => {
    let submitCallback = () => {
        tableManager.form.submitLoading = true
        const rules = formRef.value?.getCheckeds()
        tableManager.api
            .post(operate, items[tableManager.table.pk!], {
                ...items,
                rules: rules?.join(','),
            })
            .then((res) => {
                tableManager.refresh({ event: 'submit-form', operate, items: items })
                tableManager.form.operatePKs?.shift()
                if (tableManager.form.operatePKs!.length > 0) {
                    tableManager.toggleForm('update', tableManager.form.operatePKs)
                } else {
                    tableManager.toggleForm()
                }
                tableManager.runAfter('submitForm', { res })
            })
            .finally(() => {
                tableManager.form.submitLoading = false
            })
    }

    if (formEl) {
        tableManager.form.ref = formEl
        formEl.validate((valid) => {
            if (valid) {
                submitCallback()
            }
        })
    } else {
        submitCallback()
    }
    return false
}

// 切换表单后钩子
tableManager.opts.after!.toggleForm = ({ operate }) => {
    if (operate == 'create') {
        menuRuleTreeUpdate()
    }
}

// 编辑请求完成后钩子
tableManager.opts.after!.getEditData = () => {
    menuRuleTreeUpdate()
}

const menuRuleTreeUpdate = () => {
    getAdminRules().then((res) => {
        tableManager.form.extend!.menuRules = res.data.data.list

        if (typeof tableManager.form.items?.rules == 'string') {
            tableManager.form.items.rules = tableManager.form.items.rules.split(',')
        }

        if (tableManager.form.items?.rules && tableManager.form.items.rules.length) {
            if (tableManager.form.items.rules.includes('*')) {
                let arr: number[] = []
                for (const key in tableManager.form.extend!.menuRules) {
                    arr.push(tableManager.form.extend!.menuRules[key].id)
                }
                tableManager.form.extend!.defaultCheckedKeys = arr
            } else {
                tableManager.form.extend!.defaultCheckedKeys = tableManager.form.items.rules
            }
        } else {
            tableManager.form.extend!.defaultCheckedKeys = []
        }
        tableManager.form.extend!.treeKey = uuid()
    })
}

onMounted(() => {
    tableManager.table.ref = tableRef.value
    tableManager.initCtx()
    tableManager.getData()
})

// ========================== 以下为表格状态内存缓存功能 ==========================

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

tableManager.opts.before!.getData = () => {
    tableManager.table.expandAll = tableManager.table.filter!.quickSearchKeywords ? true : false
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
</script>

<style scoped lang="scss">
.default-main {
    margin-bottom: 0;
}
</style>
