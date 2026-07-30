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
                    @keyup.enter="manager.submitForm(formRef)"
                    :model="formItems"
                    :label-position="config.layout.shrink ? 'top' : 'right'"
                    :label-width="manager.form.labelWidth + 'px'"
                    :rules="rules"
                    v-if="!manager.form.loading"
                >
                    <el-form-item :label="t('auth.rule.pid')" prop="pid">
                        <RemoteSelect
                            v-model="formItems.pid"
                            field="title"
                            remote-url="/admin/auth/rule/list"
                            :remote-search-fields="['title', 'name', 'path']"
                            :placeholder="t('common.pleaseSelect', { field: t('auth.rule.pid') })"
                            :pagination="false"
                        />
                    </el-form-item>

                    <el-form-item :label="t('auth.rule.type')" prop="type" class="ag-input-item-radio">
                        <el-radio-group v-model="formItems.type">
                            <el-radio value="dir" :border="true">{{ t('auth.rule.typeDir') }}</el-radio>
                            <el-radio value="menu" :border="true">{{ t('auth.rule.typeMenu') }}</el-radio>
                            <el-radio value="node" :border="true">{{ t('auth.rule.typeNode') }}</el-radio>
                        </el-radio-group>
                    </el-form-item>

                    <el-form-item :label="t('auth.rule.title')" prop="title">
                        <el-input v-model="formItems.title" :placeholder="t('common.pleaseEnter', { field: t('auth.rule.title') })"></el-input>
                    </el-form-item>

                    <el-form-item :label="t('auth.rule.name')" prop="name">
                        <el-input v-model="formItems.name" :placeholder="t('common.pleaseEnter', { field: t('auth.rule.name') })"></el-input>
                    </el-form-item>

                    <el-form-item v-if="formItems.type !== 'node'" :label="t('common.icon')" prop="icon">
                        <IconSelect v-model="formItems.icon" />
                    </el-form-item>

                    <template v-if="formItems.type === 'menu'">
                        <el-form-item :label="t('auth.rule.openType')" prop="open_type" class="ag-input-item-radio">
                            <el-radio-group v-model="formItems.open_type">
                                <el-radio value="tab" :border="true">{{ t('auth.rule.openTypeTab') }}</el-radio>
                                <el-radio value="link" :border="true">{{ t('auth.rule.openTypeLink') }}</el-radio>
                                <el-radio value="iframe" :border="true">{{ t('auth.rule.openTypeIframe') }}</el-radio>
                            </el-radio-group>
                        </el-form-item>

                        <el-form-item v-if="formItems.open_type === 'tab'" :label="t('auth.rule.path')" prop="path">
                            <el-input v-model="formItems.path" :placeholder="t('common.pleaseEnter', { field: t('auth.rule.path') })"></el-input>
                        </el-form-item>

                        <el-form-item v-if="formItems.open_type === 'tab'" :label="t('auth.rule.component')" prop="component">
                            <el-input
                                v-model="formItems.component"
                                :placeholder="t('common.pleaseEnter', { field: t('auth.rule.component') })"
                            ></el-input>
                        </el-form-item>

                        <el-form-item v-if="['link', 'iframe'].includes(formItems.open_type)" :label="t('auth.rule.url')" prop="url">
                            <el-input v-model="formItems.url" :placeholder="t('common.pleaseEnter', { field: t('auth.rule.url') })"></el-input>
                        </el-form-item>

                        <el-form-item :label="t('auth.rule.keepalive')" prop="keepalive">
                            <el-switch v-model="formItems.keepalive" :active-value="1" :inactive-value="0" />
                        </el-form-item>
                    </template>

                    <el-form-item :label="t('auth.rule.extend')" prop="extend">
                        <el-select v-model="formItems.extend" :placeholder="t('common.pleaseSelect', { field: t('auth.rule.extend') })" clearable>
                            <el-option :label="t('auth.rule.extendAddRouteOnly')" value="add_route_only" />
                            <el-option :label="t('auth.rule.extendAddMenuOnly')" value="add_menu_only" />
                        </el-select>
                    </el-form-item>

                    <el-form-item :label="t('common.weigh')" prop="weigh">
                        <el-input-number
                            class="w100"
                            v-model="formItems.weigh"
                            :step="1"
                            :precision="0"
                            controls-position="right"
                            :placeholder="t('common.pleaseEnter', { field: t('common.weigh') })"
                        />
                    </el-form-item>

                    <el-form-item :label="t('auth.rule.remark')" prop="remark">
                        <el-input
                            @keyup.enter.stop=""
                            @keyup.ctrl.enter="manager.submitForm(formRef)"
                            v-model="formItems.remark"
                            type="textarea"
                            :placeholder="t('common.pleaseEnter', { field: t('auth.rule.remark') })"
                        ></el-input>
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
import type { FormItemRule } from 'element-plus'
import { reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import IconSelect from '/@/components/agInput/components/iconSelect.vue'
import { useConfig } from '/@/stores/config'
import { buildValidatorRule } from '/@/utils/validate'
import RemoteSelect from '/@/components/agInput/components/remoteSelect.vue'

interface Props {
    manager: TableManagerInstance
}

defineProps<Props>()
const { t } = useI18n()
const config = useConfig()
const formRef = useTemplateRef('formRef')
const formItems = defineModel<AnyObj>('formItems', { required: true })

const rules: Partial<Record<string, FormItemRule[]>> = reactive({
    type: [buildValidatorRule({ name: 'required', title: t('auth.rule.type'), trigger: 'change' })],
    title: [buildValidatorRule({ name: 'required', title: t('auth.rule.title') })],
    name: [buildValidatorRule({ name: 'required', title: t('auth.rule.name') })],
    url: [
        buildValidatorRule({ name: 'required', title: t('auth.rule.url') }),
        buildValidatorRule({ name: 'url', message: t('common.invalidEntry', { field: t('common.url') }) }),
    ],
    path: [buildValidatorRule({ name: 'required', title: t('auth.rule.path') })],
    component: [buildValidatorRule({ name: 'required', title: t('auth.rule.component') })],
    pid: [
        {
            validator: (rule: any, val: string, callback: Function) => {
                if (!val) {
                    return callback()
                }
                if (formItems.value.id && parseInt(val) == parseInt(formItems.value.id)) {
                    return callback(new Error(t('auth.rule.parentRuleSelfError')))
                }
                return callback()
            },
            trigger: 'blur',
        },
    ],
})
</script>

<style scoped lang="scss"></style>
