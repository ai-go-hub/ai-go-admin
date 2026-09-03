<template>
    <div class="default-main">
        <el-row :gutter="30">
            <el-col :xs="24" :sm="24" :md="24" :lg="10">
                <div class="admin-info">
                    <AgUpload
                        type="image"
                        v-model="state.adminInfo.avatar"
                        topic="avatar"
                        :show-file-list="false"
                        class="avatar-upload"
                        v-if="!isEmpty(state.adminInfo)"
                    >
                        <template #default>
                            <div class="avatar-uploader">
                                <el-image fit="cover" :src="state.adminInfo.avatar ? fullURL(state.adminInfo.avatar) : ''" class="avatar">
                                    <template #error>
                                        <div class="image-slot">
                                            <Icon size="30" color="#c0c4cc" name="el-picture" />
                                        </div>
                                    </template>
                                </el-image>
                            </div>
                        </template>
                    </AgUpload>

                    <div class="admin-info-base">
                        <div class="admin-nickname">{{ state.adminInfo.nickname }}</div>
                        <div class="admin-other">
                            {{ t('common.lastLoginAt') }} {{ dayjs(state.adminInfo.last_login_at).format('YYYY-MM-DD HH:mm:ss') }}
                        </div>
                    </div>

                    <div class="admin-info-form">
                        <el-form
                            @submit.prevent=""
                            @keyup.enter="onSubmit(formRef)"
                            :key="state.formKey"
                            label-position="top"
                            :rules="rules"
                            ref="formRef"
                            :model="state.adminInfo"
                        >
                            <el-form-item :label="t('auth.admin.username')">
                                <el-input disabled v-model="state.adminInfo.username" />
                            </el-form-item>
                            <el-form-item :label="t('auth.admin.nickname')" prop="nickname">
                                <el-input
                                    :placeholder="t('common.pleaseEnter', { field: t('auth.admin.nickname') })"
                                    v-model="state.adminInfo.nickname"
                                />
                            </el-form-item>
                            <el-form-item :label="t('common.email')" prop="email">
                                <el-input :placeholder="t('common.pleaseEnter', { field: t('common.email') })" v-model="state.adminInfo.email" />
                            </el-form-item>
                            <el-form-item :label="t('common.mobile')" prop="mobile">
                                <el-input :placeholder="t('common.pleaseEnter', { field: t('common.mobile') })" v-model="state.adminInfo.mobile" />
                            </el-form-item>
                            <el-form-item :label="t('auth.admin.bio')" prop="bio">
                                <el-input
                                    @keyup.enter.stop=""
                                    @keyup.ctrl.enter="onSubmit(formRef)"
                                    :placeholder="t('common.pleaseEnter', { field: t('auth.admin.bio') })"
                                    type="textarea"
                                    v-model="state.adminInfo.bio"
                                />
                            </el-form-item>
                            <el-form-item :label="t('common.password')" prop="password">
                                <el-input
                                    type="password"
                                    autocomplete="new-password"
                                    :placeholder="t('auth.admin.leaveBlankIfUnchanged')"
                                    v-model="state.adminInfo.password"
                                />
                            </el-form-item>
                            <el-form-item>
                                <el-button type="primary" :loading="state.buttonLoading" @click="onSubmit(formRef)">
                                    {{ t('common.save') }}
                                </el-button>
                                <el-button @click="resetForm(formRef)">{{ t('common.reset') }}</el-button>
                            </el-form-item>
                        </el-form>
                    </div>
                </div>
            </el-col>

            <el-col v-loading="state.logLoading" :xs="24" :sm="24" :md="24" :lg="14" class="lg-mt-20">
                <el-card shadow="never">
                    <template #header>{{ t('auth.adminLog.operationLog') }}</template>
                    <el-timeline>
                        <el-timeline-item
                            v-for="(item, idx) in state.log"
                            :key="idx"
                            size="large"
                            :timestamp="dayjs(item.created_at).format('YYYY-MM-DD HH:mm:ss')"
                        >
                            {{ item.title }}
                        </el-timeline-item>
                    </el-timeline>
                    <el-pagination
                        v-model:current-page="state.logCurrentPage"
                        v-model:page-size="state.logPageSize"
                        :page-sizes="[11, 20, 50, 100]"
                        background
                        layout="prev, next, jumper"
                        :total="state.logTotal"
                        @size-change="onLogSizeChange"
                        @current-change="onLogCurrentChange"
                    />
                </el-card>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { getAdminProfile, updateAdminProfile, getAdminLog } from '@/api/admin/routine/index'
import { dayjs, type FormItemRule } from 'element-plus'
import { fullURL, resetForm } from '@/utils/common'
import { uuid } from '@/utils/random'
import { buildValidatorRule, regularPassword } from '@/utils/validate'
import AgUpload from '@/components/agInput/components/agUpload.vue'
import { useAdminInfo } from '@/stores/adminInfo'
import { isEmpty } from 'lodash-es'

defineOptions({
    name: 'routine/profile',
})

const { t } = useI18n()
const adminInfoStore = useAdminInfo()
const formRef = useTemplateRef('formRef')

const state: {
    adminInfo: AnyObj
    formKey: string
    buttonLoading: boolean
    log: {
        title: string
        created_at: string
        url: string
    }[]
    logFilter: AnyObj
    logCurrentPage: number
    logPageSize: number
    logTotal: number
    logLoading: boolean
} = reactive({
    adminInfo: {},
    formKey: uuid(),
    buttonLoading: false,
    log: [],
    logFilter: {
        limit: 11,
        sort: 'created_at',
        order: 'desc',
        wheres: [],
    },
    logCurrentPage: 1,
    logPageSize: 11,
    logTotal: 0,
    logLoading: true,
})

const getLog = () => {
    state.logLoading = true

    /**
     * 仅查询当前管理员日志
     * 1. 日志接口服务端层面已做限制，非超管仅能查询到自己的操作日志，所以不存在越权的问题
     * 2. 此处的筛选是针对超管的，因为个人资料页面，超管也不应该显示其他管理员的日志
     */
    state.logFilter.wheres = [
        {
            wheres: [
                {
                    field: 'admin_id',
                    value: adminInfoStore.id,
                    operator: 'eq',
                },
            ],
            or: false,
        },
    ]

    getAdminLog(state.logFilter)
        .then((res) => {
            state.log = res.data.data.list
            state.logTotal = res.data.data.total
        })
        .finally(() => {
            state.logLoading = false
        })
}

const onLogSizeChange = (limit: number) => {
    state.logPageSize = limit
    state.logFilter.limit = limit
    getLog()
}

const onLogCurrentChange = (page: number) => {
    state.logCurrentPage = page
    state.logFilter.page = page
    getLog()
}

// 表单验证规则
const rules: Partial<Record<string, FormItemRule[]>> = {
    nickname: [buildValidatorRule({ name: 'required', title: t('auth.admin.nickname') })],
    email: [buildValidatorRule({ name: 'email', message: t('common.invalidEntry', { field: t('common.email') }) })],
    mobile: [buildValidatorRule({ name: 'mobile', message: t('common.invalidEntry', { field: t('common.mobile') }) })],
    password: [
        {
            validator: (rule: any, val: string, callback: Function) => {
                if (!val) {
                    return callback()
                }
                if (!regularPassword(val)) {
                    return callback(new Error(t('common.passwordFormatError')))
                }
                return callback()
            },
            trigger: 'blur',
        },
    ],
}

// 提交表单
const onSubmit = (formEl: any) => {
    formEl?.validate((valid: boolean) => {
        if (valid) {
            state.buttonLoading = true
            updateAdminProfile(adminInfoStore.id, {
                id: state.adminInfo.id,
                avatar: state.adminInfo.avatar,
                username: state.adminInfo.username,
                nickname: state.adminInfo.nickname,
                email: state.adminInfo.email,
                mobile: state.adminInfo.mobile,
                bio: state.adminInfo.bio,
                password: state.adminInfo.password,
                status: state.adminInfo.status,
            })
                .then(() => {
                    adminInfoStore.dataFill({ nickname: state.adminInfo.nickname })
                })
                .finally(() => {
                    state.buttonLoading = false
                })
        }
    })
}

onMounted(() => {
    // 获取最新的管理员个人信息
    getAdminProfile(adminInfoStore.id).then((res) => {
        state.adminInfo = res.data.data.row
        // 重新渲染表单以记录初始值
        state.formKey = uuid()

        getLog()
    })
})
</script>

<style scoped lang="scss">
.default-main {
    margin-bottom: 0;
}
.admin-info {
    background-color: var(--ag-bg-color-overlay);
    border-radius: var(--el-border-radius-base);
    border-top: 3px solid #409eff;
    :deep(.avatar-upload) {
        display: flex;
        justify-content: center;
        .el-upload--picture-card {
            width: 110px;
            height: 110px;
            margin: 26px auto 10px auto;
            border-radius: 50%;
            box-shadow: var(--el-box-shadow-light);
            border: 1px dashed var(--el-border-color);
        }
        .el-upload--picture-card:hover {
            border-color: var(--el-color-primary);
        }
        .avatar-uploader {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 100%;
            height: 100%;
            overflow: hidden;
            border-radius: 50%;
            .avatar {
                width: 110px;
                height: 110px;
                display: block;
            }
            .image-slot {
                display: flex;
                align-items: center;
                justify-content: center;
                height: 100%;
            }
        }
    }
    .admin-info-base {
        .admin-nickname {
            font-size: 22px;
            color: var(--el-text-color-primary);
            text-align: center;
            padding: 8px 0;
        }
        .admin-other {
            color: var(--el-text-color-regular);
            font-size: 14px;
            text-align: center;
            line-height: 20px;
        }
    }
    .admin-info-form {
        padding: 15px 30px 10px 30px;
    }
}
.el-card :deep(.el-timeline-item__icon) {
    font-size: 10px;
}
@media screen and (max-width: 1200px) {
    .lg-mt-20 {
        margin-top: 20px;
    }
}
</style>
