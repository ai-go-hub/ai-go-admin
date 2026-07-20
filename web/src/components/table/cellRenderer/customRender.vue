<template>
    <div>
        <component
            :is="props.columnConfig.customRender"
            :renderRow="props.row"
            :renderColumnConfig="props.columnConfig"
            :renderValue="props.cellValue"
            :renderColumn="props.column"
            :renderIndex="props.index"
            :manager="props.manager"
        />
    </div>
</template>

<script setup lang="ts">
import { CellRendererProps } from '/@/components/table/types'

const props = defineProps<CellRendererProps>()
</script>
