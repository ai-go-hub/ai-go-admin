<template>
    <div class="default-main">
        <el-row v-loading="state.loading" :gutter="20">
            <el-col class="xs-mb-20" :xs="24" :sm="16">
                <el-form
                    v-if="!state.loading"
                    ref="formRef"
                    @submit.prevent=""
                    @keyup.enter="onSubmit()"
                    :model="state.form"
                    :rules="state.rules"
                    label-position="top"
                    :key="state.formKey"
                >
                    <el-tabs v-model="state.activeTab" type="border-card" :before-leave="onBeforeLeave">
                        <el-tab-pane class="config-tab-pane" v-for="(group, key) in state.config" :key="key" :name="group.name" :label="group.title">
                            <div class="config-form-item" v-for="(item, idx) in group.configs" :key="idx">
                                <template v-if="item.group == state.activeTab">
                                    <!-- 富文本在 dialog 内全屏编辑器时必须拥有很高的 z-index，此处选择单独为 editor 设定较小的 z-index -->
                                    <el-form-item v-if="item.type == 'editor'" :label="item.title" :prop="item.name" :key="'editor-' + item.id">
                                        <AgInput
                                            :type="item.type"
                                            v-model="state.form[item.name]"
                                            :attr="{
                                                style: {
                                                    zIndex: 99,
                                                },
                                                ...item.input_extend,
                                            }"
                                            :placeholder="item.tip"
                                            @keyup.enter.stop=""
                                            @keyup.ctrl.enter="onSubmit()"
                                        />
                                    </el-form-item>

                                    <!-- textarea -->
                                    <el-form-item
                                        v-else-if="item.type == 'textarea'"
                                        :label="item.title"
                                        :prop="item.name"
                                        :key="'textarea-' + item.id"
                                    >
                                        <AgInput
                                            :type="item.type"
                                            v-model="state.form[item.name]"
                                            :attr="{ rows: 3, ...item.input_extend }"
                                            :placeholder="item.tip"
                                            @keyup.enter.stop=""
                                            @keyup.ctrl.enter="onSubmit()"
                                        />
                                    </el-form-item>

                                    <!-- 其他 -->
                                    <el-form-item v-else :label="item.title" :prop="item.name" :key="'other-' + item.id">
                                        <AgInput
                                            :type="item.type"
                                            v-model="state.form[item.name]"
                                            :placeholder="item.tip"
                                            :attr="{
                                                dict: item.dict,
                                                ...item.input_extend,
                                            }"
                                        />
                                    </el-form-item>

                                    <div class="config-form-item-name">${{ item.name }}</div>
                                    <div class="del-config-form-item">
                                        <el-popconfirm
                                            @confirm="onDelConfig(item)"
                                            v-if="item.allow_del"
                                            :confirmButtonText="t('common.delete')"
                                            :title="t('routine.config.areYouSureToDelete')"
                                        >
                                            <template #reference>
                                                <Icon class="close-icon" size="15" name="el-close" />
                                            </template>
                                        </el-popconfirm>
                                    </div>
                                </template>
                            </div>
                            <div v-if="group.name == 'mail'" class="send-test-mail">
                                <el-button @click="onTestSendMail()">{{ t('routine.config.testMailSending') }}</el-button>
                            </div>
                            <el-button type="primary" @click="onSubmit()">{{ t('common.save') }}</el-button>
                        </el-tab-pane>
                        <el-tab-pane
                            name="add_config"
                            class="config-tab-pane config-tab-pane-add"
                            :label="t('routine.config.addConfigurationItem')"
                        ></el-tab-pane>
                    </el-tabs>
                </el-form>
            </el-col>
            <el-col :xs="24" :sm="8">
                <el-card :header="t('routine.config.quickConfigurationEntry')">
                    <el-button v-for="(item, idx) in state.quickEntrance" class="quick_entrance" :key="idx">
                        <div @click="router.push({ name: item['value'] })">{{ item['key'] }}</div>
                    </el-button>
                </el-card>
            </el-col>
        </el-row>

        <AddFrom v-if="!state.loading" v-model="state.showAddForm" :config-group="state.configGroup" />
    </div>
</template>

<script setup lang="ts">
import type { FormItemRule } from 'element-plus'
import { ElMessageBox, ElNotification } from 'element-plus'
import { snakeCase } from 'lodash-es'
import { onActivated, onDeactivated, onMounted, onUnmounted, reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import AddFrom from './add.vue'
import { sendTestMail } from '/@/api/admin/routine'
import AgInput from '/@/components/agInput/index.vue'
import { adminBaseRoutePath } from '/@/router/static/adminBase'
import type { Site } from '/@/stores/interface/config'
import { useConfig } from '/@/stores/config'
import { uuid } from '/@/utils/random'
import { buildValidatorRule, type BuildValidatorParams } from '/@/utils/validate'
import { closeHotUpdate, openHotUpdate } from '/@/utils/vite'
import { useRouter } from 'vue-router'
import { TableManagerAPI } from '/@/api/table'
import { inputModelValueTypes } from '/@/components/agInput/index'

defineOptions({
    name: 'routine/config',
})

const { t } = useI18n()
const config = useConfig()
const router = useRouter()
const formRef = useTemplateRef('formRef')
const api = new TableManagerAPI('/admin/routine/config/')

const state: {
    loading: boolean
    config: AnyObj
    configGroup: AnyObj
    activeTab: string
    showAddForm: boolean
    rules: Partial<Record<string, FormItemRule[]>>
    form: AnyObj
    quickEntrance: AnyObj
    formKey: string
} = reactive({
    loading: true,
    config: [],
    configGroup: {},
    activeTab: '',
    showAddForm: false,
    rules: {},
    form: {},
    quickEntrance: {},
    formKey: uuid(),
})

const getData = () => {
    api.list<AnyObj>()
        .then((res) => {
            state.config = res.data.data.list

            for (const key in state.config) {
                for (const ckey in state.config[key].configs) {
                    let { type, value, dict, input_extend } = state.config[key].configs[ckey]

                    // 数组类型的值先调用 JSON.parse
                    if (inputModelValueTypes.array.includes(type)) {
                        state.config[key].configs[ckey].value = value ? JSON.parse(value) : []
                    }

                    // 其他 JSON 数据字段
                    state.config[key].configs[ckey].dict = dict ? JSON.parse(dict) : {}
                    state.config[key].configs[ckey].input_extend = input_extend ? JSON.parse(input_extend) : {}
                }
            }

            state.configGroup = res.data.data.configGroup
            state.quickEntrance = res.data.data.quickEntrance
            if (!state.activeTab) {
                for (const key in state.configGroup) {
                    state.activeTab = key
                    break
                }
            }
            let formNames: AnyObj = {}
            let rules: Partial<Record<string, FormItemRule[]>> = {}
            for (const key in state.config) {
                for (const lKey in state.config[key].configs) {
                    if (state.config[key].configs[lKey].rule) {
                        let ruleStr = state.config[key].configs[lKey].rule.split(',')
                        let ruleArr: AnyObj = []
                        ruleStr.forEach((item: string) => {
                            ruleArr.push(
                                buildValidatorRule({ name: item as BuildValidatorParams['name'], title: state.config[key].configs[lKey].title })
                            )
                        })
                        rules = Object.assign(rules, {
                            [state.config[key].configs[lKey].name]: ruleArr,
                        })
                    }
                    formNames[state.config[key].configs[lKey].name] =
                        state.config[key].configs[lKey].type == 'number'
                            ? parseFloat(state.config[key].configs[lKey].value)
                            : state.config[key].configs[lKey].value
                }
            }

            state.form = formNames
            state.rules = rules
            state.formKey = uuid()
        })
        .finally(() => {
            state.loading = false
        })
}

const onBeforeLeave = (newTabName: string | number) => {
    if (newTabName == 'add_config') {
        state.showAddForm = true
        return false
    }
}

const onSubmit = () => {
    formRef.value?.validate((valid) => {
        if (valid) {
            // 只提交当前tab的表单数据
            const formData: AnyObj = {}
            for (const key in state.config) {
                if (state.config[key].name != state.activeTab) {
                    continue
                }
                for (const lKey in state.config[key].configs) {
                    const { type, name } = state.config[key].configs[lKey]

                    // 组装需要提交的数据
                    formData[name] = state.form[name] ?? ''

                    // array 类型的值转为 JSON 字符串提交
                    if (inputModelValueTypes.array.includes(type)) {
                        formData[name] = JSON.stringify(formData[name])
                    }

                    // number 类型的值转为字符串提交
                    if (inputModelValueTypes.number.includes(type)) {
                        if (type == 'switch') {
                            // switch
                            if (typeof formData[name] != 'string') {
                                formData[name] = formData[name] ? '1' : '0'
                            }
                        } else if (type == 'number') {
                            // number，避免 NAN
                            formData[name] = formData[name] ? String(formData[name]) : '0'
                        } else {
                            // number 和 remoteSelect
                            formData[name] = String(formData[name])
                        }
                    }
                }
            }
            api.post('update', state.activeTab, formData).then(() => {
                // 更新状态商店数据
                for (const key in config.site) {
                    const formDataKey = snakeCase(key)
                    if (formData[formDataKey] && config.site[key as keyof Site] != formData[formDataKey]) {
                        ;(config.site[key as keyof Site] as any) = formData[formDataKey]
                    }
                }

                if (formData.entrance && formData.entrance != adminBaseRoutePath) {
                    window.open(window.location.href.replace(adminBaseRoutePath, formData.entrance))
                    window.close()
                }
            })
        }
    })
}

const onDelConfig = (config: AnyObj) => {
    api.delete([config.id]).then(() => {
        getData()
    })
}

const onTestSendMail = () => {
    if (!state.form.smtp_server || !state.form.smtp_port || !state.form.smtp_user || !state.form.smtp_pass || !state.form.smtp_sender_mail) {
        ElNotification({
            type: 'error',
            message: t('routine.config.pleaseEnterCorrectMailConfig'),
        })
        return false
    }

    ElMessageBox.prompt(t('routine.config.pleaseEnterRecipientEmail'), t('routine.config.testMailSending'), {
        confirmButtonText: t('routine.config.sendOut'),
        cancelButtonText: t('common.cancel'),
        inputPattern: /[\w!#$%&'*+/=?^_`{|}~-]+(?:\.[\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\w](?:[\w-]*[\w])?\.)+[\w](?:[\w-]*[\w])?/,
        inputErrorMessage: t('routine.config.pleaseEnterCorrectEmail'),
        beforeClose: (action, instance, done) => {
            if (action === 'confirm') {
                instance.confirmButtonLoading = true
                instance.confirmButtonText = t('routine.config.sending')
                sendTestMail(state.form, instance.inputValue).finally(() => {
                    done()
                })
            } else {
                done()
            }
        },
    })
}

onMounted(() => {
    getData()
    closeHotUpdate('config')
})
onActivated(() => {
    closeHotUpdate('config')
})
onDeactivated(() => {
    openHotUpdate('config')
})
onUnmounted(() => {
    openHotUpdate('config')
})
</script>

<style scoped lang="scss">
.send-test-mail {
    padding-bottom: 20px;
}
.el-tabs--border-card {
    border: none;
    box-shadow: var(--el-box-shadow-light);
    border-radius: var(--el-border-radius-base);
}
.el-tabs--border-card :deep(.el-tabs__header) {
    background-color: var(--ag-bg-color);
    border-bottom: none;
    border-top-left-radius: var(--el-border-radius-base);
    border-top-right-radius: var(--el-border-radius-base);
}
.el-tabs--border-card :deep(.el-tabs__item.is-active) {
    border: 1px solid transparent;
}
.el-tabs--border-card :deep(.el-tabs__nav-wrap) {
    border-top-left-radius: var(--el-border-radius-base);
    border-top-right-radius: var(--el-border-radius-base);
}
.el-card :deep(.el-card__header) {
    height: 40px;
    padding: 0;
    line-height: 40px;
    border: none;
    padding-left: 20px;
    background-color: var(--ag-bg-color);
}
.config-tab-pane {
    padding: 5px;
}
.config-tab-pane-add {
    width: 80%;
}
.config-form-item {
    display: flex;
    align-items: center;
    .el-form-item {
        flex: 13;
    }
    .config-form-item-name {
        opacity: 0;
        flex: 3;
        font-size: 13px;
        color: var(--el-text-color-disabled);
        padding-left: 20px;
    }
    .del-config-form-item {
        cursor: pointer;
        flex: 1;
        padding-left: 10px;
    }
    .close-icon {
        display: none;
    }
    &:hover {
        .config-form-item-name {
            opacity: 1;
        }
        .close-icon {
            display: block;
            color: var(--el-text-color-disabled) !important;
        }
    }
}
.quick_entrance {
    margin-left: 10px;
    margin-bottom: 10px;
}
@media screen and (max-width: 768px) {
    .xs-mb-20 {
        margin-bottom: 20px;
    }
}
</style>
