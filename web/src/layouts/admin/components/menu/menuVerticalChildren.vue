<template>
    <el-scrollbar ref="layoutMenuScrollbarRef" class="children-vertical-menus-scrollbar">
        <el-menu
            class="layouts-menu-vertical-children"
            :collapse-transition="false"
            :unique-opened="config.layout.menuUniqueOpened"
            :default-active="state.defaultActive"
            :collapse="config.layout.menuCollapse"
            ref="layoutMenuRef"
        >
            <MenuTree v-if="menu.children.length > 0" :menus="menu.children" />
        </el-menu>
    </el-scrollbar>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, useTemplateRef } from 'vue'
import type { RouteLocationNormalizedLoaded } from 'vue-router'
import { onBeforeRouteUpdate, useRoute } from 'vue-router'
import MenuTree from '/@/layouts/admin/components/menu/menuTree.vue'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { layoutMenuRef } from '/@/stores/ref'
import { getMenuDataByRoute } from '/@/utils/common'
import { getMenuKey } from '/@/utils/router'

const menu = useMenu()
const route = useRoute()
const config = useConfig()

const layoutMenuScrollbarRef = useTemplateRef('layoutMenuScrollbarRef')

const state = reactive({
    defaultActive: '',
})

const verticalMenusScrollbarHeight = computed(() => {
    const menuTopBarHeight = config.layout.menuShowTopBar ? 60 : 0
    const asideFooterToolbarHeight = config.layout.menuCollapse ? 100 : 50
    return 'calc(100% - ' + (menuTopBarHeight + asideFooterToolbarHeight) + 'px)'
})

/**
 * 激活当前路由的菜单
 */
const currentRouteActive = (currentRoute: RouteLocationNormalizedLoaded) => {
    // 以路由 fullPath 匹配的菜单优先，且 fullPath 无匹配时，回退到 path 的匹配菜单
    const tabView = getMenuDataByRoute(currentRoute)
    if (tabView) {
        state.defaultActive = getMenuKey(tabView, tabView.meta!.matched as string)
    }

    let routeChildren = getMenuDataByRoute(currentRoute, 'above')
    if (routeChildren) {
        if (routeChildren.children && routeChildren.children.length > 0) {
            menu.setChildren(routeChildren.children)
        } else {
            menu.setChildren([routeChildren])
        }
    } else {
        menu.setChildren([])
    }
}

/**
 * 侧栏菜单滚动条滚动到激活菜单所在位置
 */
const verticalMenusScroll = () => {
    setTimeout(() => {
        let activeMenu: HTMLElement | null = document.querySelector('.el-menu.layouts-menu-vertical-children li.is-active')
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
.children-vertical-menus-scrollbar {
    height: v-bind(verticalMenusScrollbarHeight);
    background-color: v-bind('config.getColorValue("menuBackground")');
}
.layouts-menu-vertical-children {
    border: 0;
    --el-menu-bg-color: v-bind('config.getColorValue("menuBackground")');
    --el-menu-text-color: v-bind('config.getColorValue("menuColor")');
    --el-menu-active-color: v-bind('config.getColorValue("menuActiveColor")');
    --el-menu-hover-bg-color: v-bind('config.getColorValue("menuHoverBackground")');
    --el-menu-active-bg-color: v-bind('config.getColorValue("menuActiveBackground")');
}
</style>
