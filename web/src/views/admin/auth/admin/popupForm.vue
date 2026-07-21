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
                    <el-form-item :label="t('auth.admin.username')" prop="username">
                        <el-input
                            type="string"
                            v-model="formItems.username"
                            :placeholder="t('common.pleaseEnter', { field: t('auth.admin.username') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('auth.admin.nickname')" prop="nickname">
                        <el-input
                            type="string"
                            v-model="formItems.nickname"
                            :placeholder="t('common.pleaseEnter', { field: t('auth.admin.nickname') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('auth.admin.avatar')">
                        <AgUpload type="image" v-model="formItems.avatar" />
                    </el-form-item>

                    <el-form-item :label="t('common.email')" prop="email">
                        <el-input
                            type="string"
                            v-model="formItems.email"
                            :placeholder="t('common.pleaseEnter', { field: t('common.email') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('common.mobile')" prop="mobile">
                        <el-input
                            v-model="formItems.mobile"
                            type="string"
                            :placeholder="t('common.pleaseEnter', { field: t('common.mobile') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('common.password')" prop="password">
                        <el-input
                            v-model="formItems.password"
                            type="password"
                            autocomplete="new-password"
                            :placeholder="
                                manager.form.operate == 'create'
                                    ? t('common.pleaseEnter', { field: t('common.password') })
                                    : t('auth.admin.leaveBlankIfUnchanged')
                            "
                        ></el-input>
                    </el-form-item>

                    <el-form-item prop="bio" :label="t('auth.admin.bio')">
                        <el-input
                            @keyup.enter.stop=""
                            @keyup.ctrl.enter="manager.submitForm(formRef)"
                            v-model="formItems.bio"
                            type="textarea"
                            :placeholder="t('common.pleaseEnter', { field: t('auth.admin.bio') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('common.status')">
                        <el-radio-group v-model="formItems.status">
                            <el-radio value="enable" :border="true">{{ t('common.enable') }}</el-radio>
                            <el-radio value="disable" :border="true">{{ t('common.disable') }}</el-radio>
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
import { reactive, watch, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { regularPassword, buildValidatorRule } from '/@/utils/validate'
import type { FormItemRule } from 'element-plus'
import { useConfig } from '/@/stores/config'
import AgUpload from '/@/components/agInput/components/agUpload.vue'

interface Props {
    manager: TableManagerInstance
}

const props = defineProps<Props>()
const formItems = defineModel<AnyObj>('formItems', { required: true })

const config = useConfig()
const formRef = useTemplateRef('formRef')

const { t } = useI18n()

const rules: Partial<Record<string, FormItemRule[]>> = reactive({
    username: [buildValidatorRule({ name: 'required', title: t('auth.admin.username') }), buildValidatorRule({ name: 'account' })],
    nickname: [buildValidatorRule({ name: 'required', title: t('auth.admin.nickname') })],
    email: [buildValidatorRule({ name: 'email', message: t('common.invalidEntry', { field: t('common.email') }) })],
    mobile: [buildValidatorRule({ name: 'mobile', message: t('common.invalidEntry', { field: t('common.mobile') }) })],
    password: [
        {
            validator: (rule: any, val: string, callback: Function) => {
                if (props.manager.form.operate == 'create') {
                    if (!val) {
                        return callback(new Error(t('common.pleaseEnter', { field: t('common.password') })))
                    }
                } else {
                    if (!val) {
                        return callback()
                    }
                }
                if (!regularPassword(val)) {
                    return callback(new Error(t('common.invalidEntry', { field: t('common.password') })))
                }
                return callback()
            },
            trigger: 'blur',
        },
    ],
})

watch(
    () => props.manager.form.operate,
    (newVal) => {
        // 创建密码字段必填，编辑非必填
        rules.password![0].required = newVal == 'create'
    }
)
</script>

<style scoped lang="scss"></style>
