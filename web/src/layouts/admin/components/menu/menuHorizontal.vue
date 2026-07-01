<template>
    <div class="layouts-menu-horizontal">
        <div class="menu-horizontal-logo" v-if="config.layout.menuShowTopBar">
            <Logo />
        </div>
        <el-scrollbar ref="layoutMenuScrollbarRef" class="horizontal-menus-scrollbar">
            <el-menu
                :popper-style="{
                    '--el-menu-bg-color': config.getColorValue('headerBarBackground'),
                    '--el-menu-text-color': config.getColorValue('headerBarTabColor'),
                    '--el-menu-active-color': config.getColorValue('headerBarTabActiveColor'),
                    '--el-menu-hover-bg-color': config.getColorValue('headerBarHoverBackground'),
                    '--el-menu-active-bg-color': config.getColorValue('headerBarTabActiveBackground'),
                    '--el-menu-hover-text-color': config.getColorValue('headerBarTabColor'),
                }"
                popper-class="menu-horizontal-streamline-popper"
                ref="layoutMenuRef"
                class="menu-horizontal"
                mode="horizontal"
                :default-active="state.defaultActive"
            >
                <MenuTree :extends="{ position: 'horizontal', level: 1 }" :menus="menu.rawData" />
            </el-menu>
        </el-scrollbar>
        <NavMenu />
    </div>
</template>

<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { onBeforeRouteUpdate, useRoute, type RouteLocationNormalizedLoaded } from 'vue-router'
import Logo from '/@/layouts/admin/components/logo.vue'
import MenuTree from '/@/layouts/admin/components/menu/menuTree.vue'
import NavMenu from '/@/layouts/admin/components/navMenu.vue'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { layoutMenuRef, layoutMenuScrollbarRef } from '/@/stores/ref'
import { getMenuDataByRoute } from '/@/utils/common'
import horizontalScroll from '/@/utils/horizontalScroll'
import { getMenuKey } from '/@/utils/router'

const menu = useMenu()
const route = useRoute()
const config = useConfig()

const state = reactive({
    defaultActive: '',
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
 * 滚动条横向滚动到激活菜单所在位置
 */
const verticalMenusScroll = () => {
    setTimeout(() => {
        let activeMenu: HTMLElement | null = document.querySelector('.el-menu.menu-horizontal li.is-active')
        if (activeMenu) {
            layoutMenuScrollbarRef.value?.setScrollLeft(activeMenu.offsetLeft)
        }
    }, 500)
}

onMounted(() => {
    currentRouteActive(route)
    verticalMenusScroll()

    new horizontalScroll(layoutMenuScrollbarRef.value!.wrapRef!)
})

onBeforeRouteUpdate((to) => {
    currentRouteActive(to)
})
</script>

<style lang="scss">
.menu-horizontal-streamline-popper {
    .el-menu--horizontal {
        .el-menu-item:not(.is-disabled):hover,
        .el-menu-item:not(.is-disabled):focus {
            color: var(--el-menu-text-color);
        }
        .el-menu-item.is-active,
        .el-menu-item.is-active:hover,
        .el-sub-menu.is-active > .el-sub-menu__title,
        .el-sub-menu.is-active > .el-sub-menu__title:hover {
            color: var(--el-menu-active-color);
        }
    }
}
</style>

<style scoped lang="scss">
.layouts-menu-horizontal {
    display: flex;
    align-items: center;
    width: 100vw;
    height: var(--el-header-height);
    background-color: v-bind('config.getColorValue("headerBarBackground")');
    border-bottom: 1px solid var(--el-color-info-light-8);
}
.menu-horizontal-logo {
    width: 180px;
    background-color: v-bind('config.getColorValue("headerBarBackground")');
}
.horizontal-menus-scrollbar {
    flex: 1;
    height: var(--el-header-height);
}
.menu-horizontal {
    border: none;
    --el-menu-bg-color: v-bind('config.getColorValue("headerBarBackground")');
    --el-menu-text-color: v-bind('config.getColorValue("headerBarTabColor")');
    --el-menu-active-color: v-bind('config.getColorValue("headerBarTabActiveColor")');
    --el-menu-hover-bg-color: v-bind('config.getColorValue("headerBarHoverBackground")');
    --el-menu-active-bg-color: v-bind('config.getColorValue("headerBarTabActiveBackground")');
    --el-menu-hover-text-color: v-bind('config.getColorValue("headerBarTabColor")');
}

:deep(.el-sub-menu),
:deep(.el-menu-item) {
    .el-sub-menu__title {
        background-color: var(--el-menu-bg-color);
        &:hover {
            background-color: var(--el-menu-hover-bg-color);
        }
    }
    &:hover {
        color: var(--el-menu-hover-text-color) !important;
        background-color: var(--el-menu-hover-bg-color);
    }
    &.is-active {
        background-color: var(--el-menu-active-bg-color);
        .el-sub-menu__title {
            background-color: var(--el-menu-active-bg-color);
        }
        &:hover {
            color: var(--el-menu-active-color) !important;
        }
    }
}
</style>
