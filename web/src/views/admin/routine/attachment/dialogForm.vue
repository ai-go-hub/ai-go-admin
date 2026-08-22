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
                    <el-form-item :label="t('routine.attachment.preview')">
                        <el-image class="preview-img" :preview-src-list="[previewUrl]" :src="previewUrl"></el-image>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.topic')">
                        <el-input v-model="formItems.topic" type="string" disabled></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.url')">
                        <el-input v-model="formItems.url" type="string" disabled></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.name')" prop="name">
                        <el-input
                            v-model="formItems.name"
                            type="string"
                            :placeholder="t('common.pleaseEnter', { field: t('routine.attachment.name') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.size')">
                        <el-input :model-value="formatFileSize(formItems.size)" type="string" disabled></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.mimetype')">
                        <el-input v-model="formItems.mimetype" type="string"></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.quote')">
                        <el-input v-model="formItems.quote" type="string"></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.driver')">
                        <el-input v-model="formItems.driver" type="string"></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.attachment.sha1')">
                        <el-input v-model="formItems.sha1" type="string" disabled></el-input>
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
import { useConfig } from '@/stores/config'
import { buildValidatorRule } from '@/utils/validate'
import type { FormItemRule } from 'element-plus'
import { computed, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { formatFileSize, previewRenderFormatter } from './index'

interface Props {
    manager: TableManagerInstance
}

defineProps<Props>()
const config = useConfig()
const formRef = useTemplateRef('formRef')
const formItems = defineModel<AnyObj>('formItems', { required: true })

const { t } = useI18n()

const rules: Partial<Record<string, FormItemRule[]>> = {
    name: [buildValidatorRule({ name: 'required', title: t('routine.attachment.name') })],
}

// 预览图 URL
const previewUrl = computed(() => {
    return previewRenderFormatter(formItems.value, null, formItems.value.url || '')
})
</script>

<style scoped lang="scss">
.preview-img {
    width: 60px;
    height: 60px;
}
</style>
