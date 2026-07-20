<template>
    <div>
        <el-switch
            v-if="columnConfig.prop"
            @change="onChange"
            :model-value="cellValueNew"
            :loading="loading"
            active-value="1"
            inactive-value="0"
            v-bind="invokeTableContextDataFun(columnConfig.customRenderAttr?.switch, { row, columnConfig, column, cellValue: cellValueNew, index })"
        />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { invokeTableContextDataFun } from '/@/components/table/index'
import { CellRendererProps } from '/@/components/table/types'

const props = defineProps<CellRendererProps>()

const loading = ref(false)
const cellValueNew = ref<string | number | boolean>(typeof props.cellValue === 'number' ? props.cellValue.toString() : props.cellValue)

const onChange = (value: string | number | boolean) => {
    loading.value = true
    props.manager.api
        .post('update', props.row[props.manager.table.pk!], {
            // 发送全部字段，避免服务端的数据效验报错，如某字段有必填验证
            ...props.row,
            [props.columnConfig.prop!]: value,
        })
        .then(() => {
            cellValueNew.value = value
            props.manager.handleEvent('cell-change', { value: value, ...props })
        })
        .finally(() => {
            loading.value = false
        })
}
</script>
