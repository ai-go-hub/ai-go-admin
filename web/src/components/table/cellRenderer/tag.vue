<template>
    <div>
        <el-tag
            v-if="![null, undefined, ''].includes(cellValue)"
            :type="getTagType(cellValue, props.columnConfig.custom)"
            effect="light"
            size="default"
            v-bind="invokeTableContextDataFun(columnConfig.customRenderAttr?.tag, { row, columnConfig, column, cellValue, index })"
        >
            {{ !isEmpty(columnConfig.dict) ? (columnConfig.dict[cellValue] ?? cellValue) : cellValue }}
        </el-tag>
    </div>
</template>

<script setup lang="ts">
import { invokeTableContextDataFun } from '@/components/table/index'
import { CellRendererProps } from '@/components/table/types'
import { TagProps } from 'element-plus'
import { isEmpty } from 'lodash-es'

const props = defineProps<CellRendererProps>()

const getTagType = (value: string, custom: any): TagProps['type'] => {
    return !isEmpty(custom) && custom[value] ? custom[value] : 'primary'
}
</script>
