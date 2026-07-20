<template>
    <div>
        <template v-if="isArray(cellValue) && cellValue.length">
            <el-image
                v-for="(item, idx) in cellValue"
                :key="idx"
                :initial-index="idx"
                :preview-teleported="true"
                :preview-src-list="arrayFullURL(cellValue)"
                class="ag-table-render-images-item"
                :src="fullURL(item)"
                :hide-on-click-modal="true"
                v-bind="invokeTableContextDataFun(props.columnConfig.customRenderAttr?.image, { row, columnConfig, column, cellValue, index })"
            ></el-image>
        </template>
    </div>
</template>

<script setup lang="ts">
import { isArray } from 'lodash-es'
import { invokeTableContextDataFun } from '/@/components/table/index'
import { CellRendererProps } from '/@/components/table/types'
import { arrayFullURL, fullURL } from '/@/utils/common'

const props = defineProps<CellRendererProps>()
</script>

<style scoped lang="scss">
.ag-table-render-images-item {
    height: 36px;
    width: 36px;
    margin: 0 5px;
}
</style>
