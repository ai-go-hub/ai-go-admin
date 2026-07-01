<template>
    <div>
        <template v-for="(item, idx) in props.menus" :key="idx">
            <el-option
                class="select-option"
                :label="pad(' ', props.level * 4) + item.meta?.title"
                :value="item.path"
                :disabled="item.meta?.type != 'menu'"
            />
            <template v-if="item.children">
                <SelectOption :menus="item.children" :level="props.level + 1" />
            </template>
        </template>
    </div>
</template>

<script setup lang="ts">
import { pad } from 'lodash-es'
import { RouteRecordRaw } from 'vue-router'

interface Props {
    menus: RouteRecordRaw[]
    level?: number
}

const props = withDefaults(defineProps<Props>(), {
    menus: () => [],
    level: 0,
})
</script>

<style scoped lang="scss">
.select-option {
    white-space: pre;
}
</style>
