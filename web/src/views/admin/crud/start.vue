<template>
    <div class="default-main">
        <div class="crud-title">{{ t('crud.index.start') }}</div>
        <div class="start-opt">
            <el-row :gutter="20">
                <el-col :xs="24" :span="6">
                    <div @click="changeStep('create')" class="start-item suspension">
                        <div class="start-item-title">{{ t('crud.index.create') }}</div>
                        <div class="start-item-remark">{{ t('crud.index.createAnew') }}</div>
                    </div>
                </el-col>
                <el-col @click="onShowDialog('ai')" :xs="24" :span="6">
                    <div class="start-item suspension">
                        <div class="start-item-title">{{ t('crud.index.ai') }}</div>
                        <div class="start-item-remark">{{ t('crud.index.aiSubtitle') }}</div>
                    </div>
                </el-col>
                <el-col @click="onShowDialog('table')" :xs="24" :span="6">
                    <div class="start-item suspension">
                        <div class="start-item-title">{{ t('crud.index.selectTable') }}</div>
                        <div class="start-item-remark">{{ t('crud.index.selectTableSubtitle') }}</div>
                    </div>
                </el-col>
                <el-col @click="onShowDialog('log')" :xs="24" :span="6">
                    <div class="start-item suspension">
                        <div class="start-item-title">{{ t('crud.index.record') }}</div>
                        <div class="start-item-remark">{{ t('crud.index.recordContinue') }}</div>
                    </div>
                </el-col>
            </el-row>
            <el-row justify="center">
                <el-col :span="20" class="crud-tips suspension">
                    <b>{{ t('crud.index.quickExperience') }}</b>
                    <ol>
                        <li>
                            {{ t('crud.index.experience11') }}
                            <code>{{ t('crud.index.experience12') }}</code>
                            （{{ t('crud.index.experience13') }}）
                        </li>
                        <li>
                            {{ t('crud.index.experience21') }}
                            <code>{{ t('crud.index.create') }}</code>
                            {{ t('crud.index.or') }}
                            <code> {{ t('crud.index.experience22') }}{{ t('crud.index.experience23') }} </code>
                        </li>
                        <li>
                            {{ t('crud.index.experience31') }} <code>{{ t('crud.index.experience32') }}</code>
                            {{ t('crud.index.experience33') }}
                            <code>{{ t('crud.index.experience34') }}</code>
                        </li>
                    </ol>
                    <el-alert v-if="!isDev()" class="no-dev" type="warning" :show-icon="true" :closable="false">
                        <template #title>
                            <span>{{ t('crud.index.experience41') }}{{ t('crud.index.experience42') }}</span>
                            <span>
                                {{ t('crud.index.experience43') }}
                            </span>
                        </template>
                    </el-alert>
                </el-col>
            </el-row>

            <el-dialog
                class="ag-operate-dialog select-table-dialog"
                v-model="state.showSelectTableDialog"
                :title="t('crud.index.selectTable')"
                :destroy-on-close="true"
            >
                <el-form
                    :label-width="140"
                    @keyup.enter="onSubmitSelectTableForm()"
                    class="select-table-form"
                    ref="formRef"
                    :model="state"
                    :rules="rules"
                >
                    <el-form-item :label-width="140" :label="t('crud.index.table')" prop="table">
                        <RemoteSelect
                            v-model="state.table"
                            pk="table"
                            field="comment"
                            :remote-url="tableListUrl"
                            @row="onTableStartChange"
                            :remote-params="{
                                exclusions: [
                                    'schema_migrations',
                                    'areas',
                                    'tokens',
                                    'captchas',
                                    'admin_group_access',
                                    'configs',
                                    'admin_logs',
                                    'crud_logs',
                                ],
                            }"
                        />
                    </el-form-item>
                    <el-alert
                        v-if="state.successRecord"
                        class="success-record-alert"
                        :title="t('crud.index.generatedTip')"
                        :show-icon="true"
                        :closable="false"
                        type="warning"
                    />
                </el-form>
                <template #footer>
                    <div :style="{ width: 'calc(100% * 0.9)' }">
                        <el-button @click="state.showSelectTableDialog = false">{{ $t('common.cancel') }}</el-button>
                        <el-button :loading="state.loading" @click="onSubmitSelectTableForm()" type="primary">{{ t('common.confirm') }}</el-button>
                        <el-button v-if="state.successRecord" @click="onLogStart" type="success">
                            {{ t('crud.index.recordStart') }}
                        </el-button>
                    </div>
                </template>
            </el-dialog>

            <CrudLog v-model="state.showLogDialog" />

            <AIDialog v-model="state.showAIDialog" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { checkLog, tableListUrl } from '@/api/admin/crud'
import RemoteSelect from '@/components/agInput/components/remoteSelect.vue'
import { buildValidatorRule } from '@/utils/validate'
import AIDialog from '@/views/admin/crud/ai.vue'
import { changeStep } from '@/views/admin/crud/index'
import CrudLog from '@/views/admin/crud/log.vue'
import type { FormItemRule } from 'element-plus'
import { reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const formRef = useTemplateRef('formRef')

const state = reactive({
    table: '',
    loading: false,
    successRecord: 0,
    showAIDialog: false,
    showLogDialog: false,
    showSelectTableDialog: false,
})

const onShowDialog = (type: 'table' | 'log' | 'ai') => {
    if (type == 'table') {
        state.table = ''
        state.successRecord = 0
        state.showSelectTableDialog = true
    } else if (type == 'log') {
        state.showLogDialog = true
    } else if (type == 'ai') {
        state.showAIDialog = true
    }
}

const rules: Partial<Record<string, FormItemRule[]>> = {
    table: [buildValidatorRule({ name: 'required', message: t('common.pleaseSelect', { field: t('crud.index.table') }) })],
}

const onSubmitSelectTableForm = () => {
    formRef.value?.validate((valid) => {
        if (valid) {
            changeStep('table', { table: state.table })
        }
    })
}

const onTableStartChange = () => {
    if (state.table) {
        // 检查是否有CRUD记录
        state.loading = true
        checkLog(state.table)
            .then((res) => {
                state.successRecord = res.data.data.id
            })
            .finally(() => {
                state.loading = false
            })
    }
}

const onLogStart = () => {
    if (state.successRecord) {
        changeStep('log', { id: state.successRecord })
    }
}

const isDev = () => {
    return import.meta.env.DEV
}
</script>

<style scoped lang="scss">
:deep(.select-table-dialog) .el-dialog__body {
    height: unset;
    .select-table-form {
        width: 88%;
        padding: 40px 0;
    }
    .success-record-alert {
        width: calc(100% - 140px);
        margin-left: 140px;
        margin-bottom: 30px;
        margin-top: -10px;
    }
}
.suspension {
    transition: all 0.3s ease;
}
.suspension:hover {
    -webkit-transform: translateY(-4px) scale(1.02);
    -moz-transform: translateY(-4px) scale(1.02);
    -ms-transform: translateY(-4px) scale(1.02);
    -o-transform: translateY(-4px) scale(1.02);
    transform: translateY(-4px) scale(1.02);
    -webkit-box-shadow: 0 14px 24px rgba(0, 0, 0, 0.2);
    box-shadow: 0 14px 24px rgba(0, 0, 0, 0.2);
    z-index: 2147483600;
    border-radius: 6px;
}
.crud-title {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: var(--el-font-size-extra-large);
    font-weight: bold;
    padding-top: 120px;
}
.start-opt {
    display: block;
    width: 70%;
    margin: 40px auto;
}
.start-item {
    background-color: #e1eaf9;
    border-radius: var(--el-border-radius-base);
    padding: 25px;
    margin-bottom: 20px;
    cursor: pointer;
}
.start-item-title {
    font-size: var(--el-font-size-large);
    color: var(--ag-color-primary-light);
}
.start-item-remark {
    display: block;
    line-height: 18px;
    min-height: 48px;
    padding-top: 12px;
    color: #92969a;
}
.crud-tips {
    margin-top: 60px;
    padding: 20px;
    background-color: rgba($color: #ffffff, $alpha: 0.6);
    border-radius: var(--el-border-radius-base);
    color: var(--el-color-info);
    b {
        font-size: 15px;
        padding-left: 10px;
    }
    .no-dev {
        margin-top: 10px;
    }
    code {
        color: #3594f7;
        background-color: #3baafa1a;
        display: inline-block;
        padding: 0 4px;
        border-radius: 2px;
        line-height: 22px;
    }
    ol {
        margin: 0.6em 0;
        padding-left: 1.6em;
        li {
            margin: 0.5em 0;
            line-height: 1.6;
        }
    }
}
@at-root .dark {
    .start-item {
        background-color: #1d1e1f;
    }
    .crud-tips {
        background-color: rgba($color: #1d1e1f, $alpha: 0.4);
    }
}
</style>
