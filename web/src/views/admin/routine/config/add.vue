<template>
    <el-dialog class="ag-operate-dialog" :close-on-click-modal="false" :model-value="props.modelValue" @close="closeForm">
        <template #header>
            <div class="title">{{ t('routine.config.addConfigurationItem') }}</div>
        </template>
        <el-scrollbar class="ag-table-form-scrollbar">
            <div class="ag-operate-form ag-add-form" :style="config.layout.shrink ? '' : 'width: calc(100% - ' + state.labelWidth / 2 + 'px)'">
                <el-form
                    ref="formRef"
                    @keyup.enter="onAddSubmit()"
                    :rules="rules"
                    :model="state.addConfig"
                    :label-position="config.layout.shrink ? 'top' : 'right'"
                    :label-width="state.labelWidth"
                >
                    <el-form-item :label="t('routine.config.variableGroup')" prop="group">
                        <el-select
                            class="w100"
                            v-model="state.addConfig.group"
                            :placeholder="t('common.pleaseSelect', { field: t('routine.config.variableGroup') })"
                        >
                            <el-option v-for="(title, name) in configGroup" :key="name" :label="title" :value="name"></el-option>
                        </el-select>
                    </el-form-item>

                    <el-form-item :label="t('routine.config.variableName')" prop="name">
                        <el-input
                            v-model="state.addConfig.name"
                            :placeholder="t('common.pleaseEnter', { field: t('routine.config.variableName') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.config.variableTitle')" prop="title">
                        <el-input
                            v-model="state.addConfig.title"
                            :placeholder="t('common.pleaseEnter', { field: t('routine.config.variableTitle') })"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.config.variableType')" prop="type">
                        <el-select
                            class="w100"
                            v-model="state.addConfig.type"
                            :placeholder="t('common.pleaseSelect', { field: t('routine.config.variableType') })"
                        >
                            <el-option v-for="item in inputTypes" :key="item" :label="item" :value="item"></el-option>
                        </el-select>
                    </el-form-item>

                    <el-form-item :label="t('routine.config.tip')" prop="tip">
                        <el-input v-model="state.addConfig.tip" :placeholder="t('common.pleaseEnter', { field: t('routine.config.tip') })"></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.config.validationRules')" prop="rule">
                        <el-select
                            class="w100"
                            v-model="state.addConfig.rule"
                            :placeholder="t('common.pleaseSelect', { field: t('routine.config.validationRules') })"
                            :multiple="true"
                        >
                            <el-option v-for="(title, key) in validatorType" :key="key" :label="title" :value="key"></el-option>
                        </el-select>
                    </el-form-item>

                    <el-form-item v-if="dictRequiredTypes.includes(state.addConfig.type)" :label="t('routine.config.dictionaryData')" prop="dict">
                        <el-input
                            type="textarea"
                            v-model="state.addConfig.dict"
                            :placeholder="t('routine.config.strAttrTip')"
                            :rows="3"
                            @keyup.enter.stop=""
                            @keyup.ctrl.enter="onAddSubmit()"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('routine.config.inputExtend')" prop="input_extend">
                        <el-input
                            type="textarea"
                            v-model="state.addConfig.input_extend"
                            :placeholder="t('routine.config.strAttrTip')"
                            :rows="3"
                            @keyup.enter.stop=""
                            @keyup.ctrl.enter="onAddSubmit()"
                        ></el-input>
                    </el-form-item>

                    <el-form-item :label="t('common.weigh')" prop="weigh">
                        <el-input-number v-model="state.addConfig.weigh" controls-position="right" class="w100" />
                    </el-form-item>
                </el-form>
            </div>
        </el-scrollbar>
        <template #footer>
            <div :style="'width: calc(100% - ' + state.labelWidth / 1.8 + 'px)'">
                <el-button @click="closeForm">{{ t('common.cancel') }}</el-button>
                <el-button :loading="state.submitLoading" @click="onAddSubmit()" type="primary"> {{ t('common.add') }} </el-button>
            </div>
        </template>
    </el-dialog>
</template>

<script setup lang="ts">
import type { FormRules } from 'element-plus'
import { cloneDeep } from 'lodash-es'
import { reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { TableManagerAPI } from '/@/api/table'
import { inputTypes, dictRequiredTypes } from '/@/components/agInput/index'
import { useConfig } from '/@/stores/config'
import { parseStrAttr } from '/@/utils/common'
import { buildValidatorRule, validatorType } from '/@/utils/validate'

const config = useConfig()
const api = new TableManagerAPI('/admin/routine/config/')

interface Props {
    modelValue: boolean
    configGroup: AnyObj
}

const props = withDefaults(defineProps<Props>(), {
    modelValue: false,
    configGroup: () => {
        return {}
    },
})

const emits = defineEmits<{
    (e: 'update:modelValue', value: boolean): void
}>()

const closeForm = () => {
    emits('update:modelValue', false)
}

const { t } = useI18n()
const formRef = useTemplateRef('formRef')

const state: {
    inputTypes: AnyObj
    labelWidth: number
    submitLoading: boolean
    addConfig: {
        group: string

        name: string
        title: string
        type: string
        tip: string | null
        rule: string | string[] | null
        dict: string | null
        input_extend: string | null
        value: string | null

        weigh: number
    }
} = reactive({
    inputTypes: {},
    labelWidth: 160,
    submitLoading: false,
    addConfig: {
        group: '',

        name: '',
        title: '',
        type: '',
        tip: null,
        rule: [],
        dict: `key1=value1
key2=value2`,
        input_extend: '',
        value: null,

        weigh: 0,
    },
})

const rules = reactive<FormRules>({
    group: [
        buildValidatorRule({ name: 'required', trigger: 'change', message: t('common.pleaseSelect', { field: t('routine.config.variableGroup') }) }),
    ],
    name: [buildValidatorRule({ name: 'required', title: t('routine.config.variableName') }), buildValidatorRule({ name: 'varName' })],
    title: [buildValidatorRule({ name: 'required', title: t('routine.config.variableTitle') })],
    type: [
        buildValidatorRule({ name: 'required', trigger: 'change', message: t('common.pleaseSelect', { field: t('routine.config.variableType') }) }),
    ],
    weigh: [buildValidatorRule({ name: 'integer', title: t('common.weigh') })],
})

const onAddSubmit = () => {
    formRef.value?.validate((valid) => {
        if (valid) {
            const config = cloneDeep(state.addConfig)

            // 默认值需要设为字符串的字段类型
            if (['files', 'images'].includes(config.type)) {
                config.value = ''
            }

            // rule 使用 , 号分割入库
            if (Array.isArray(config.rule)) {
                config.rule = config.rule.length ? config.rule.join(',') : null
            }

            // 字典和扩展数据格式化为键值对并转为 JSON 入库
            config.input_extend = config.input_extend ? JSON.stringify(parseStrAttr(config.input_extend)) : null
            config.dict = dictRequiredTypes.includes(config.type) && config.dict ? JSON.stringify(parseStrAttr(config.dict)) : null

            api.post('create', '', config).then(() => {
                closeForm()
            })
        }
    })
}
</script>

<style scoped lang="scss"></style>
