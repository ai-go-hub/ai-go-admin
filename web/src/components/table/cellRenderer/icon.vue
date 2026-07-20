<template>
    <div>
        <Icon
            color="var(--el-text-color-primary)"
            :name="props.cellValue"
            v-bind="invokeTableContextDataFun(columnConfig.customRenderAttr?.icon, { row, columnConfig, column, cellValue, index })"
        />
    </div>
</template>

<script setup lang="ts">
import Icon from '/@/components/icon/index.vue'
import { invokeTableContextDataFun } from '/@/components/table/index'
import { CellRendererProps } from '/@/components/table/types'

const props = defineProps<CellRendererProps>()
</script>
