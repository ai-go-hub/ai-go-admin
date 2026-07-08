<template>
    <div class="left-split-menus">
        <div class="left-split-primary-menus-scrollbar-wrap">
            <div v-if="config.layout.menuTopBarLogo" class="logo-img">
                <img src="/@/assets/logo.png" alt="logo" />
            </div>
            <el-scrollbar ref="layoutMenuScrollbarRef" class="left-split-primary-menus-scrollbar">
                <el-menu
                    class="layouts-menu-vertical primary-menus"
                    :collapse-transition="false"
                    :unique-opened="config.layout.menuUniqueOpened"
                    :default-active="state.primaryDefaultActive"
                    :collapse="true"
                    ref="layoutMenuRef"
                >
                    <el-menu-item
                        v-for="menu in menuStore.rawData"
                        @click="onClickPrimaryMenu(menu)"
                        :index="getMenuKey(menu)"
                        :key="getMenuKey(menu)"
                    >
                        <Icon :size="20" :name="menu.meta?.icon ? menu.meta?.icon : config.layout.menuDefaultIcon" />
                        <span>{{ menu.meta?.title ? menu.meta?.title : $t('layouts.untitled') }}</span>
                    </el-menu-item>
                </el-menu>
            </el-scrollbar>
        </div>

        <div v-if="menuStore.children.length" class="left-split-secondary-menus-scrollbar-wrap">
            <el-scrollbar ref="layoutSecondaryMenuScrollbarRef" class="left-split-secondary-menus-scrollbar">
                <el-menu
                    class="layouts-menu-vertical secondary-menus"
                    :collapse-transition="false"
                    :unique-opened="config.layout.menuUniqueOpened"
                    :default-active="state.secondaryDefaultActive"
                    :collapse="config.layout.menuCollapse"
                    ref="layoutSecondaryMenuRef"
                >
                    <MenuLeftSplitTree :menus="menuStore.children" />
                </el-menu>
            </el-scrollbar>
            <AsideFooterToolbar />
        </div>
    </div>
</template>

<script setup lang="ts">
import type { MenuInstance, ScrollbarInstance } from 'element-plus'
import { computed, onMounted, reactive, ref } from 'vue'
import { onBeforeRouteUpdate, RouteLocationNormalized, RouteRecordRaw, useRoute, type RouteLocationNormalizedLoaded } from 'vue-router'
import AsideFooterToolbar from '/@/layouts/admin/components/aside/toolbar/footer.vue'
import MenuLeftSplitTree from '/@/layouts/admin/components/menu/menuLeftSplitTree.vue'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { layoutMenuRef, layoutMenuScrollbarRef } from '/@/stores/ref'
import { getMenuDataByRoute } from '/@/utils/common'
import { getMenuKey, openMenu } from '/@/utils/router'

const route = useRoute()
const config = useConfig()
const menuStore = useMenu()
const menuWidth = computed(() => config.getMenuWidth())

const state = reactive({
    primaryDefaultActive: '',
    secondaryDefaultActive: '',
})

const layoutSecondaryMenuRef = ref<MenuInstance>()
const layoutSecondaryMenuScrollbarRef = ref<ScrollbarInstance>()

const verticalPrimaryMenusScrollbarHeight = computed(() => {
    const menuTopBarHeight = config.layout.menuTopBarLogo ? 61 : 0
    return 'calc(100% - ' + menuTopBarHeight + 'px)'
})

const verticalSecondaryMenusScrollbarHeight = computed(() => {
    const asideFooterToolbarHeight = config.layout.menuCollapse ? 100 : 50
    return 'calc(100% - ' + asideFooterToolbarHeight + 'px)'
})

const findRouteChildren = (menu: RouteRecordRaw | RouteLocationNormalized) => {
    let routeChildren = getMenuDataByRoute(menu, 'above')
    if (routeChildren) {
        state.primaryDefaultActive = getMenuKey(routeChildren)
    }
    if (routeChildren && routeChildren.children) {
        menuStore.setChildren(routeChildren.children)
    } else {
        menuStore.setChildren([])
    }
}

const onClickPrimaryMenu = (menu: RouteRecordRaw) => {
    if (menu.meta?.type == 'dir') {
        return findRouteChildren(menu)
    } else if (menu.meta?.type == 'menu' && route.fullPath == menu.path) {
        return findRouteChildren(menu)
    }
    openMenu(menu)
}

/**
 * 激活当前路由对应的菜单
 */
const currentRouteActive = (currentRoute: RouteLocationNormalizedLoaded) => {
    // 以路由 fullPath 匹配的菜单优先，且 fullPath 无匹配时，回退到 path 的匹配菜单
    const tabView = getMenuDataByRoute(currentRoute)
    if (tabView) {
        state.secondaryDefaultActive = getMenuKey(tabView, tabView.meta!.matched as string)
    }

    findRouteChildren(currentRoute)
}

/**
 * 滚动条滚动到激活菜单所在位置
 */
const verticalMenusScroll = () => {
    setTimeout(() => {
        let activeMenu: HTMLElement | null = document.querySelector('.primary-menus.layouts-menu-vertical li.is-active')
        if (activeMenu) {
            layoutMenuScrollbarRef.value?.setScrollTop(activeMenu.offsetTop)
        }

        let secondaryActiveMenu: HTMLElement | null = document.querySelector('.secondary-menus.layouts-menu-vertical li.is-active')
        if (secondaryActiveMenu) {
            layoutSecondaryMenuScrollbarRef.value?.setScrollTop(secondaryActiveMenu.offsetTop)
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
.left-split-menus {
    display: flex;
    height: 100%;
    .logo-img {
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 20px 0;
        margin-bottom: 16px;
        img {
            width: 26px;
        }
    }
}
.left-split-primary-menus-scrollbar-wrap {
    width: 80px;
    background-color: v-bind('config.getColorValue("menuBackgroundPrimary")');
    .left-split-primary-menus-scrollbar {
        width: 100%;
        height: v-bind(verticalPrimaryMenusScrollbarHeight);
    }
}
.left-split-secondary-menus-scrollbar-wrap {
    width: calc(v-bind(menuWidth) - 80px);
    background-color: v-bind('config.getColorValue("menuBackground")');
    .left-split-secondary-menus-scrollbar {
        width: 100%;
        padding: 8px;
        height: v-bind(verticalSecondaryMenusScrollbarHeight);
    }
}
.layouts-menu-vertical {
    border: 0;
}
.primary-menus {
    margin: 0 8px;
    --el-menu-bg-color: v-bind('config.getColorValue("menuBackgroundPrimary")');
    --el-menu-text-color: v-bind('config.getColorValue("menuColor")');
    --el-menu-active-color: v-bind('config.getColorValue("menuActiveColor")');
    --el-menu-hover-bg-color: v-bind('config.getColorValue("menuHoverBackgroundLeftSplit")');
    --el-menu-active-bg-color: v-bind('config.getColorValue("menuActiveBackgroundPrimary")');
    .el-menu-item {
        margin: 8px 0;
        border-radius: var(--el-border-radius-base);
        color: var(--el-menu-text-color);
        .icon {
            vertical-align: middle;
            text-align: center;
            flex-shrink: 0;
        }
        &.is-active {
            background-color: var(--el-menu-active-bg-color);
        }
        &.is-active {
            color: var(--el-menu-active-color);
        }
    }
}
.secondary-menus {
    --el-menu-bg-color: v-bind('config.getColorValue("menuBackground")');
    --el-menu-text-color: v-bind('config.getColorValue("menuColor")');
    --el-menu-active-color: v-bind('config.getColorValue("menuActiveColor")');
    --el-menu-hover-bg-color: v-bind('config.getColorValue("menuHoverBackgroundLeftSplit")');
    --el-menu-active-bg-color: v-bind('config.getColorValue("menuActiveBackground")');
}
</style>
