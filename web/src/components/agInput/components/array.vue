<template>
    <div>
        <el-row :gutter="10">
            <el-col :span="10" class="ag-array-key">{{ state.keyTitle }}</el-col>
            <el-col :span="10" class="ag-array-value">{{ state.valueTitle }}</el-col>
        </el-row>
        <el-row class="ag-array-item" v-for="(item, idx) in state.value" :gutter="10" :key="idx">
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
import { reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'

type agInputArray = { key: string; value: string }
interface Props {
    modelValue: agInputArray[]
    keyTitle?: string
    valueTitle?: string
}

const { t } = useI18n()

const props = withDefaults(defineProps<Props>(), {
    modelValue: () => [],
    keyTitle: '',
    valueTitle: '',
})

const state = reactive({
    value: props.modelValue,
    keyTitle: props.keyTitle ? props.keyTitle : t('common.arrayKey'),
    valueTitle: props.valueTitle ? props.valueTitle : t('common.arrayValue'),
})

const onAddArrayItem = () => {
    state.value.push({
        key: '',
        value: '',
    })
}

const onDelArrayItem = (idx: number) => {
    state.value.splice(idx, 1)
}

watch(
    () => props.modelValue,
    (newVal) => {
        state.value = newVal
    }
)
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
