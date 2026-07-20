<template>
    <div>
        <div :style="{ background: props.cellValue }" class="ag-table-cell-render-color"></div>
    </div>
</template>

<script setup lang="ts">
import { CellRendererProps } from '/@/components/table/types'

const props = defineProps<CellRendererProps>()
</script>

<style scoped lang="scss">
.ag-table-cell-render-color {
    height: 25px;
    width: 100%;
}
</style>
