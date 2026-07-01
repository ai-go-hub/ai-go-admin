<template>
    <div>
        <div
            v-loading="true"
            element-loading-background="var(--ag-bg-color-overlay)"
            :element-loading-text="$t('common.loading')"
            class="default-main ag-main-loading"
        ></div>
        <div v-if="state.showReload" class="loading-footer">
            <el-button @click="refresh" type="warning">{{ $t('common.reload') }}</el-button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, reactive } from 'vue'
import router from '/@/router/index'
import { useMenu } from '/@/stores/menu'
import { getFirstMenu } from '/@/utils/router'

let timer: number
const menu = useMenu()
const state = reactive({
    maximumWait: 1000 * 3,
    showReload: false,
})

const refresh = () => {
    router.go(0)
}

timer = window.setTimeout(() => {
    state.showReload = true
}, state.maximumWait)

onMounted(() => {
    if (menu.rawData) {
        let firstRoute = getFirstMenu(menu.rawData)
        if (firstRoute) {
            router.push(firstRoute.path)
        }
    }
})

onUnmounted(() => {
    clearTimeout(timer)
})
</script>

<style scoped lang="scss">
.ag-main-loading {
    height: 300px;
    display: flex;
    align-items: center;
    justify-content: center;
}
.loading-footer {
    display: flex;
    align-items: center;
    justify-content: center;
}
</style>
