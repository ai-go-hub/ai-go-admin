<template>
    <div>
        {{ datetime }}
    </div>
</template>

<script setup lang="ts">
import { dayjs } from 'element-plus'
import { computed } from 'vue'
import { CellRendererProps } from '/@/components/table/types'

const props = defineProps<CellRendererProps>()

const datetime = computed(() => {
    if (!props.cellValue || props.cellValue == '0001-01-01T00:00:00Z') {
        return '-'
    }
    return dayjs(props.cellValue).format(props.columnConfig.custom?.format ? props.columnConfig.custom.format : 'YYYY-MM-DD HH:mm:ss')
})
</script>
