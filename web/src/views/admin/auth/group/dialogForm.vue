<template>
    <!-- 对话框表单 -->
    <el-dialog
        class="ag-operate-dialog"
        :close-on-click-modal="false"
        :model-value="['create', 'update'].includes(manager.form.operate!)"
        @close="manager.toggleForm"
        :destroy-on-close="true"
        :draggable="true"
    >
        <template #header>
            <div class="title">
                {{ manager.form.operate == 'create' ? t('common.add') : t('common.edit') }}
            </div>
        </template>
        <el-scrollbar v-loading="manager.form.loading" class="ag-table-form-scrollbar">
            <div
                class="ag-operate-form"
                :class="'ag-' + manager.form.operate + '-form'"
                :style="config.layout.shrink ? '' : 'width: calc(100% - ' + manager.form.labelWidth! / 2 + 'px)'"
            >
                <el-form
                    ref="formRef"
                    @submit.prevent=""
                    @keyup.enter="manager.submitForm(formRef)"
                    :model="formItems"
                    :label-position="config.layout.shrink ? 'top' : 'right'"
                    :label-width="manager.form.labelWidth + 'px'"
                    :rules="rules"
                    v-if="!manager.form.loading"
                >
                    <el-form-item :label="t('auth.group.pid')" prop="pid">
                        <RemoteSelect
                            v-model="formItems.pid"
                            field="name"
                            remote-url="/admin/auth/group/list"
                            :placeholder="t('common.pleaseSelect', { field: t('auth.group.pid') })"
                            :pagination="false"
                            :empty-values="[null, 0]"
                            :value-on-clear="0"
                        />
                    </el-form-item>

                    <el-form-item :label="t('auth.group.name')" prop="name">
                        <el-input v-model="formItems.name" :placeholder="t('common.pleaseEnter', { field: t('auth.group.name') })"></el-input>
                    </el-form-item>

                    <el-form-item prop="rules" :label="t('auth.group.rules')">
                        <el-tree
                            ref="treeRef"
                            :key="manager.form.extend!.treeKey"
                            :default-checked-keys="manager.form.extend!.defaultCheckedKeys"
                            :default-expand-all="true"
                            show-checkbox
                            node-key="id"
                            :props="{ children: 'children', label: 'title', class: treeNodeClass }"
                            :data="manager.form.extend!.menuRules"
                            class="w100"
                        />
                    </el-form-item>

                    <el-form-item :label="t('common.status')" prop="status">
                        <el-radio-group v-model="formItems.status" class="ag-input-item-radio">
                            <el-radio :value="1" :border="true">{{ t('common.enable') }}</el-radio>
                            <el-radio :value="0" :border="true">{{ t('common.disable') }}</el-radio>
                        </el-radio-group>
                    </el-form-item>
                </el-form>
            </div>
        </el-scrollbar>
        <template #footer>
            <div :style="'width: calc(100% - ' + manager.form.labelWidth! / 1.8 + 'px)'">
                <el-button @click="manager.toggleForm()">{{ t('common.cancel') }}</el-button>
                <el-button :loading="manager.form.submitLoading" @click="manager.submitForm(formRef)" type="primary">
                    {{ manager.form.operatePKs && manager.form.operatePKs.length > 1 ? t('common.saveAndContinue') : t('common.save') }}
                </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import RemoteSelect from '@/components/agInput/components/remoteSelect.vue'
import { useConfig } from '@/stores/config'
import { buildValidatorRule } from '@/utils/validate'
import type { FormItemRule } from 'element-plus'
import type ElTreeNode from 'element-plus/es/components/tree/src/model/node'
import { reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

interface Props {
    manager: TableManagerInstance
}

defineProps<Props>()
const { t } = useI18n()
const config = useConfig()
const formRef = useTemplateRef('formRef')
const treeRef = useTemplateRef('treeRef')
const formItems = defineModel<AnyObj>('formItems', { required: true })

const rules: Partial<Record<string, FormItemRule[]>> = reactive({
    name: [buildValidatorRule({ name: 'required', title: t('auth.group.name') })],
    rules: [
        {
            required: true,
            validator: (rule: any, val: string, callback: Function) => {
                let ids = getCheckeds()
                if (ids.length <= 0) {
                    return callback(new Error(t('common.pleaseSelect', { field: t('auth.group.rules') })))
                }
                return callback()
            },
        },
    ],
    pid: [
        {
            validator: (rule: any, val: string, callback: Function) => {
                if (!val) {
                    return callback()
                }
                if (formItems.value.id && parseInt(val) == parseInt(formItems.value.id)) {
                    return callback(new Error(t('auth.group.parentGroupSelfError')))
                }
                return callback()
            },
            trigger: 'blur',
        },
    ],
})

const treeNodeClass = (data: AnyObj, node: ElTreeNode) => {
    if (node.isLeaf) return ''
    let addClass = true
    for (const key in node.childNodes) {
        if (!node.childNodes[key].isLeaf) {
            addClass = false
        }
    }
    return addClass ? 'penultimate-node' : ''
}

const getCheckeds = () => {
    return treeRef.value!.getCheckedKeys().concat(treeRef.value!.getHalfCheckedKeys())
}

defineExpose({
    getCheckeds,
})
</script>

<style scoped lang="scss">
:deep(.penultimate-node) {
    .el-tree-node__children {
        padding-left: 60px;
        white-space: pre-wrap;
        line-height: 12px;
        .el-tree-node {
            display: inline-block;
        }
        .el-tree-node__content {
            padding-left: 5px !important;
            padding-right: 5px;
            .el-tree-node__expand-icon {
                display: none;
            }
        }
    }
}
</style>
