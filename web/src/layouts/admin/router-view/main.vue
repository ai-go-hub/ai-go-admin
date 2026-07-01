<template>
    <el-main class="layout-main">
        <el-scrollbar class="layout-main-scrollbar" :style="layoutMainScrollbarStyle" ref="layoutMainScrollbarRef">
            <router-view v-slot="{ Component }">
                <transition :name="config.layout.mainAnimation" mode="out-in">
                    <keep-alive :include="state.keepAliveComponentNameList">
                        <component :is="Component" :key="state.componentKey" />
                    </keep-alive>
                </transition>
            </router-view>
        </el-scrollbar>
    </el-main>
</template>

<script setup lang="ts">
import { nextTick, onBeforeMount, onMounted, onUnmounted, reactive, watch } from 'vue'
import { useRoute, type RouteLocationNormalized } from 'vue-router'
import { useConfig } from '/@/stores/config'
import { useNavTab } from '/@/stores/navTab'
import { layoutMainScrollbarRef, layoutMainScrollbarStyle } from '/@/stores/ref'
import { getGlobalProperties, getMenuDataByRoute } from '/@/utils/common'

defineOptions({
    name: 'layout/main',
})

const route = useRoute()
const config = useConfig()
const navTab = useNavTab()
const globalProperties = getGlobalProperties()!

const state: {
    componentKey: string
    keepAliveComponentNameList: string[]
} = reactive({
    componentKey: route.fullPath,
    keepAliveComponentNameList: [],
})

const addKeepAliveComponentName = function (keepAliveName: string | undefined) {
    if (keepAliveName) {
        let exist = state.keepAliveComponentNameList.find((name: string) => {
            return name === keepAliveName
        })
        if (exist) return
        state.keepAliveComponentNameList.push(keepAliveName)
    }
}

const addActiveRouteKeepAlive = () => {
    if (navTab.state.activeRoute) {
        const tabView = getMenuDataByRoute(navTab.state.activeRoute)
        if (tabView && typeof tabView.meta?.keepalive == 'string') {
            addKeepAliveComponentName(tabView.meta.keepalive)
        }
    }
}

onBeforeMount(() => {
    globalProperties.eventBus.on('onTabViewRefresh', (menu: RouteLocationNormalized) => {
        state.keepAliveComponentNameList = state.keepAliveComponentNameList.filter((name: string) => menu.meta.keepalive !== name)
        state.componentKey = ''
        nextTick(() => {
            state.componentKey = menu.fullPath
            addKeepAliveComponentName(menu.meta.keepalive as string)
        })
    })
    globalProperties.eventBus.on('onTabViewClose', (menu: RouteLocationNormalized) => {
        state.keepAliveComponentNameList = state.keepAliveComponentNameList.filter((name: string) => menu.meta.keepalive !== name)
    })
})

onUnmounted(() => {
    globalProperties.eventBus.off('onTabViewRefresh')
    globalProperties.eventBus.off('onTabViewClose')
})

onMounted(() => {
    // 确保刷新页面时也能正确取得当前路由 keepalive 参数（热更新）
    addActiveRouteKeepAlive()
})

watch(
    () => route.fullPath,
    () => {
        state.componentKey = route.fullPath
        addActiveRouteKeepAlive()
    }
)
</script>

<style scoped lang="scss">
.layout-container .layout-main {
    padding: 0 !important;
    overflow: hidden;
    width: 100%;
    height: 100%;
}
.layout-main-scrollbar {
    width: 100%;
    position: relative;
    overflow: hidden;
}
</style>
