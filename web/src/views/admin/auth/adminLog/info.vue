<template>
    <!-- 查看详情 -->
    <el-dialog class="ag-operate-dialog" :model-value="manager.form.operate == 'info'" @close="manager.toggleForm">
        <template #header>
            <div class="title">{{ t('common.info') }}</div>
        </template>
        <el-scrollbar v-loading="manager.form.loading" class="ag-table-form-scrollbar">
            <div class="ag-operate-form" :class="'ag-' + manager.form.operate + '-form'">
                <el-descriptions :column="2" border>
                    <el-descriptions-item :label="t('auth.adminLog.adminId')">
                        {{ manager.form.extend!.info.id }}
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('auth.adminLog.username')">
                        {{ manager.form.extend!.info.username }}
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('auth.adminLog.title')">
                        {{ manager.form.extend!.info.title }}
                    </el-descriptions-item>
                    <el-descriptions-item :label="t('auth.adminLog.ip')">
                        {{ manager.form.extend!.info.ip }}
                    </el-descriptions-item>
                    <el-descriptions-item :width="120" :span="2" :label="t('auth.adminLog.url')">
                        {{ manager.form.extend!.info.url }}
                    </el-descriptions-item>
                    <el-descriptions-item :width="120" :span="2" :label="t('auth.adminLog.userAgent')">
                        {{ manager.form.extend!.info.user_agent }}
                    </el-descriptions-item>
                    <el-descriptions-item :width="120" :span="2" :label="t('common.createdAt')">
                        {{ dayjs(manager.form.extend!.info.create_at).format('YYYY-MM-DD HH:mm:ss') }}
                    </el-descriptions-item>
                    <el-descriptions-item :width="120" :span="2" :label="t('auth.adminLog.data')">
                        <el-tree class="table-el-tree" :data="manager.form.extend!.info.data" :props="{ label: 'label', children: 'children' }" />
                    </el-descriptions-item>
                </el-descriptions>
            </div>
        </el-scrollbar>
    </el-dialog>
</template>

<script setup lang="ts">
import { dayjs } from 'element-plus'
import { useI18n } from 'vue-i18n'

interface Props {
    manager: TableManagerInstance
}

const { t } = useI18n()
defineProps<Props>()
</script>

<style scoped lang="scss"></style>
