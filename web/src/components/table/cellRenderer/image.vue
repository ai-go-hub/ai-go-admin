<template>
    <div>
        <el-image
            v-if="cellValue"
            :hide-on-click-modal="true"
            :preview-teleported="true"
            :preview-src-list="[fullURL(cellValue)]"
            :src="fullURL(cellValue)"
            class="ag-table-render-image"
            v-bind="invokeTableContextDataFun(props.columnConfig.customRenderAttr?.image, { row, columnConfig, column, cellValue, index })"
        ></el-image>
    </div>
</template>

<script setup lang="ts">
import { invokeTableContextDataFun } from '/@/components/table/index'
import { CellRendererProps } from '/@/components/table/types'
import { fullURL } from '/@/utils/common'

const props = defineProps<CellRendererProps>()
</script>

<style scoped lang="scss">
.ag-table-render-image {
    height: 36px;
    width: 36px;
}
</style>
