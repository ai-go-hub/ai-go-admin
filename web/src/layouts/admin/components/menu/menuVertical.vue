<template>
    <el-scrollbar ref="layoutMenuScrollbarRef" class="vertical-menus-scrollbar">
        <el-menu
            class="layouts-menu-vertical"
            :collapse-transition="false"
            :unique-opened="config.layout.menuUniqueOpened"
            :default-active="state.defaultActive"
            :collapse="config.layout.menuCollapse"
            ref="layoutMenuRef"
        >
            <MenuTree :menus="menu.rawData" />
        </el-menu>
    </el-scrollbar>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { onBeforeRouteUpdate, useRoute, type RouteLocationNormalizedLoaded } from 'vue-router'
import MenuTree from '/@/layouts/admin/components/menu/menuTree.vue'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { layoutMenuRef, layoutMenuScrollbarRef } from '/@/stores/ref'
import { getMenuDataByRoute } from '/@/utils/common'
import { getMenuKey } from '/@/utils/router'

const menu = useMenu()
const route = useRoute()
const config = useConfig()

const state = reactive({
    defaultActive: '',
})

const verticalMenusScrollbarHeight = computed(() => {
    const menuTopBarHeight = config.layout.menuShowTopBar ? 50 : 0
    const asideFooterToolbarHeight = config.layout.menuCollapse ? 100 : 50
    return 'calc(100% - ' + (menuTopBarHeight + asideFooterToolbarHeight) + 'px)'
})

/**
 * 激活当前路由对应的菜单
 */
const currentRouteActive = (currentRoute: RouteLocationNormalizedLoaded) => {
    // 以路由 fullPath 匹配的菜单优先，且 fullPath 无匹配时，回退到 path 的匹配菜单
    const tabView = getMenuDataByRoute(currentRoute)
    if (tabView) {
        state.defaultActive = getMenuKey(tabView, tabView.meta!.matched as string)
    }
}

/**
 * 滚动条滚动到激活菜单所在位置
 */
const verticalMenusScroll = () => {
    setTimeout(() => {
        let activeMenu: HTMLElement | null = document.querySelector('.el-menu.layouts-menu-vertical li.is-active')
        if (activeMenu) {
            layoutMenuScrollbarRef.value?.setScrollTop(activeMenu.offsetTop)
        }
    }, 500)
}

onMounted(() => {
    currentRouteActive(route)
    verticalMenusScroll()
})

onBeforeRouteUpdate((to) => {
    currentRouteActive(to)
})
</script>

<style scoped lang="scss">
.vertical-menus-scrollbar {
    height: v-bind(verticalMenusScrollbarHeight);
    background-color: v-bind('config.getColorValue("menuBackground")');
}
.layouts-menu-vertical {
    border: 0;
    --el-menu-bg-color: v-bind('config.getColorValue("menuBackground")');
    --el-menu-text-color: v-bind('config.getColorValue("menuColor")');
    --el-menu-active-color: v-bind('config.getColorValue("menuActiveColor")');
    --el-menu-hover-bg-color: v-bind('config.getColorValue("menuHoverBackground")');
    --el-menu-active-bg-color: v-bind('config.getColorValue("menuActiveBackground")');
}
</style>
