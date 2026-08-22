<template>
    <div>
        <div v-html="columnConfig.customTemplate ? columnConfig.customTemplate(row, props.columnConfig, column, cellValue, index) : ''"></div>
    </div>
</template>

<script setup lang="ts">
import { CellRendererProps } from '@/components/table/types'

const props = defineProps<CellRendererProps>()
</script>
