<template>
    <div class="w100">
        <el-row :gutter="10">
            <el-col :span="10" class="ag-array-key">{{ state.keyTitle }}</el-col>
            <el-col :span="10" class="ag-array-value">{{ state.valueTitle }}</el-col>
        </el-row>
        <el-row class="ag-array-item" v-for="(item, idx) in list" :gutter="10" :key="idx">
            <el-col :span="10">
                <el-input v-model="item.key"></el-input>
            </el-col>
            <el-col :span="10">
                <el-input v-model="item.value"></el-input>
            </el-col>
            <el-col :span="4">
                <el-button @click="onDelArrayItem(idx)" size="small" icon="el-icon-delete" circle />
            </el-col>
        </el-row>
        <el-row :gutter="10">
            <el-col :span="10" :offset="10">
                <el-button class="ag-add-array-item" @click="onAddArrayItem" icon="el-icon-plus">{{ t('common.add') }}</el-button>
            </el-col>
        </el-row>
    </div>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'
import { useI18n } from 'vue-i18n'

interface Props {
    keyTitle?: string
    valueTitle?: string
}
type agInputArray = { key: string; value: string }

const { t } = useI18n()

const props = withDefaults(defineProps<Props>(), {
    keyTitle: '',
    valueTitle: '',
})

const model = defineModel<agInputArray[] | null>({ default: () => [] })
const list = computed<agInputArray[]>(() => model.value ?? [])

const state = reactive({
    keyTitle: props.keyTitle ? props.keyTitle : t('common.arrayKey'),
    valueTitle: props.valueTitle ? props.valueTitle : t('common.arrayValue'),
})

const onAddArrayItem = () => {
    model.value = [...(model.value ?? []), { key: '', value: '' }]
}

const onDelArrayItem = (idx: number) => {
    model.value = (model.value ?? []).filter((_, i) => i !== idx)
}
</script>

<style scoped lang="scss">
.ag-array-key,
.ag-array-value {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 5px 0;
    color: var(--el-text-color-secondary);
}
.ag-array-item {
    margin-bottom: 6px;
}
.ag-add-array-item {
    float: right;
}
</style>
