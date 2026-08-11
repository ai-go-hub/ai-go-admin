<template>
    <div>
        <el-switch
            v-if="columnConfig.prop"
            @change="onChange"
            :loading="loading"
            :model-value="modelValue"
            :active-value="getDefaultValue('active')"
            :inactive-value="getDefaultValue('inactive')"
            v-bind="invokeTableContextDataFun(columnConfig.customRenderAttr?.switch, { row, columnConfig, column, cellValue: modelValue, index })"
        />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { invokeTableContextDataFun } from '/@/components/table/index'
import { CellRendererProps } from '/@/components/table/types'

const loading = ref(false)
const props = defineProps<CellRendererProps>()
const modelValueType = typeof props.cellValue

const getDefaultValue = (type: 'active' | 'inactive') => {
    if (modelValueType == 'number') {
        return type == 'active' ? 1 : 0
    } else if (modelValueType == 'string') {
        return type == 'active' ? '1' : '0'
    } else if (modelValueType == 'boolean') {
        return type == 'active' ? true : false
    }
}

const modelValue = ref(props.cellValue == null ? getDefaultValue('inactive') : props.cellValue)

const onChange = (value: string | number | boolean) => {
    loading.value = true
    props.manager.api
        .post('update', props.row[props.manager.table.pk!], {
            // 发送全部字段，避免服务端的数据效验报错，如某字段有必填验证
            ...props.row,
            [props.columnConfig.prop!]: value,
        })
        .then(() => {
            modelValue.value = value
            props.manager.handleEvent('cell-change', { value: value, ...props })
        })
        .finally(() => {
            loading.value = false
        })
}
</script>
