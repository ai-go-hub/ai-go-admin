<template>
    <div>
        <template v-if="isArray(cellValue)">
            <template v-for="(tag, idx) in cellValue" :key="idx">
                <el-tag
                    v-if="![null, undefined, ''].includes(tag)"
                    class="m-4"
                    :type="getTagType(tag, columnConfig.custom)"
                    effect="light"
                    size="default"
                    v-bind="invokeTableContextDataFun(props.columnConfig.customRenderAttr?.tag, { row, columnConfig, column, cellValue, index })"
                >
                    {{ !isEmpty(columnConfig.dict) ? (columnConfig.dict[tag] ?? tag) : tag }}
                </el-tag>
            </template>
        </template>
        <template v-else>
            <el-tag
                v-if="![null, undefined, ''].includes(cellValue)"
                :type="getTagType(cellValue, columnConfig.custom)"
                effect="light"
                size="default"
                v-bind="invokeTableContextDataFun(props.columnConfig.customRenderAttr?.tag, { row, columnConfig, column, cellValue, index })"
            >
                {{ !isEmpty(columnConfig.dict) ? (columnConfig.dict[cellValue] ?? cellValue) : cellValue }}
            </el-tag>
        </template>
    </div>
</template>

<script setup lang="ts">
import { TagProps } from 'element-plus'
import { isArray, isEmpty } from 'lodash-es'
import { invokeTableContextDataFun } from '/@/components/table/index'
import { CellRendererProps } from '/@/components/table/types'

const props = defineProps<CellRendererProps>()

const getTagType = (value: string, custom: any): TagProps['type'] => {
    return !isEmpty(custom) && custom[value] ? custom[value] : 'primary'
}
</script>

<style scoped lang="scss">
.m-4 {
    margin: 4px;
}
</style>
