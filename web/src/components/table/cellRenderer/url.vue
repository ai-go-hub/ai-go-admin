<template>
    <div>
        <el-input :model-value="props.cellValue" :placeholder="$t('common.url')">
            <template #append>
                <el-button v-blur @click="openURL(props.cellValue)">
                    <Icon size="16" color="#606266" name="el-position" />
                </el-button>
            </template>
        </el-input>
    </div>
</template>

<script setup lang="ts">
import { invokeTableContextDataFun } from '@/components/table/index'
import { CellRendererProps } from '@/components/table/types'

const props = defineProps<CellRendererProps>()

const linkAttr = invokeTableContextDataFun(props.columnConfig.customRenderAttr?.link, {
    row: props.row,
    columnConfig: props.columnConfig,
    column: props.column,
    cellValue: props.cellValue,
    index: props.index,
})

const openURL = (url: string) => {
    if (linkAttr.target == '_blank') {
        window.open(url)
    } else {
        window.location.href = url
    }
}
</script>
