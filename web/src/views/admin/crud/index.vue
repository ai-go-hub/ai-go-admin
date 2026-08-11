<template>
    <div>
        <component :is="state.step"></component>
    </div>
</template>

<script setup lang="ts">
import { onActivated, onDeactivated, onMounted, onUnmounted } from 'vue'
import { closeHotUpdate, openHotUpdate } from '/@/utils/vite'
import Design from '/@/views/admin/crud/design.vue'
import { state } from '/@/views/admin/crud/index'
import Start from '/@/views/admin/crud/start.vue'

defineOptions({
    name: 'crud/crud',
    components: { Start, Design },
})

onMounted(() => {
    closeHotUpdate('crud')
})

onActivated(() => {
    closeHotUpdate('crud')
})

onDeactivated(() => {
    openHotUpdate('crud')
})

onUnmounted(() => {
    openHotUpdate('crud')
})
</script>

<style scoped lang="scss"></style>
