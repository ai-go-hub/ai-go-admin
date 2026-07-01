<template>
    <div class="nav-tabs" ref="tabScrollbarRef">
        <div
            v-for="(item, idx) in navTab.state.list"
            @click="onTab(item)"
            @contextmenu.prevent="onContextmenu(item, $event)"
            class="ai-go-nav-tab"
            :class="navTab.state.activeIndex == idx ? 'active' : ''"
            :ref="tabsRefs.set"
            :key="idx"
        >
            {{ item.meta.title }}
            <transition @after-leave="selectNavTab(tabsRefs[navTab.state.activeIndex])" name="el-fade-in">
                <Icon v-show="navTab.state.list.length > 1" class="close-icon" @click.stop="closeTab(item)" size="15" name="el-close" />
            </transition>
        </div>
        <div :style="activeBoxStyle" class="nav-tabs-active-box"></div>
        <Contextmenu ref="contextmenuRef" :items="state.contextmenuItems" @menuClick="onContextMenuClick" />
    </div>
</template>

<script setup lang="ts">
import { useTemplateRefsList } from '@vueuse/core'
import { nextTick, onMounted, reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteUpdate, useRoute, useRouter, type RouteLocationNormalized } from 'vue-router'
import Contextmenu from '/@/components/contextmenu/index.vue'
import type { ContextMenuItem, ContextMenuItemClickEmitArg } from '/@/components/contextmenu/interface'
import { adminBaseRoutePath } from '/@/router/static/adminBase'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { useNavTab } from '/@/stores/navTab'
import { getGlobalProperties } from '/@/utils/common'
import horizontalScroll from '/@/utils/horizontalScroll'
import { getFirstMenu } from '/@/utils/router'

const { t } = useI18n()
const route = useRoute()
const config = useConfig()
const navTab = useNavTab()
const router = useRouter()
const menuStore = useMenu()

const globalProperties = getGlobalProperties()!
const tabsRefs = useTemplateRefsList<HTMLDivElement>()
const contextmenuRef = useTemplateRef('contextmenuRef')
const tabScrollbarRef = useTemplateRef('tabScrollbarRef')

const state: {
    contextmenuItems: ContextMenuItem[]
} = reactive({
    contextmenuItems: [
        { name: 'refresh', label: t('common.reload'), icon: 'lucide-refresh-cw' },
        { name: 'close', label: t('layouts.closeTab'), icon: 'lucide-x' },
        { name: 'fullScreen', label: t('layouts.tabFullscreen'), icon: 'lucide-maximize' },
        { name: 'closeOther', label: t('layouts.closeOtherTabs'), icon: 'lucide-minus' },
        { name: 'closeAll', label: t('layouts.closeAllTabs'), icon: 'lucide-octagon-minus' },
    ],
})

const activeBoxStyle = reactive({
    width: '0',
    transform: 'translateX(0px)',
})

const onTab = (menu: RouteLocationNormalized) => {
    router.push(menu.fullPath)
}

// tab 激活状态切换
const selectNavTab = function (dom: HTMLDivElement) {
    if (!dom) {
        return false
    }
    activeBoxStyle.width = dom.clientWidth + 'px'
    activeBoxStyle.transform = `translateX(${dom.offsetLeft}px)`

    if (tabScrollbarRef.value) {
        let scrollLeft = dom.offsetLeft + dom.clientWidth - tabScrollbarRef.value.clientWidth
        if (dom.offsetLeft < tabScrollbarRef.value.scrollLeft) {
            tabScrollbarRef.value.scrollTo(dom.offsetLeft, 0)
        } else if (scrollLeft > tabScrollbarRef.value.scrollLeft) {
            tabScrollbarRef.value.scrollTo(scrollLeft, 0)
        }
    }
}

const toLastTab = () => {
    const lastTab = navTab.state.list.slice(-1)[0]
    if (lastTab) {
        router.push(lastTab.fullPath)
    } else {
        router.push(adminBaseRoutePath)
    }
}

const closeTab = (route: RouteLocationNormalized) => {
    navTab._closeTab(route)
    globalProperties.eventBus.emit('onTabViewClose', route)
    if (navTab.state.activeRoute?.fullPath === route.fullPath) {
        toLastTab()
    } else {
        navTab._setActiveRoute(navTab.state.activeRoute!)
        nextTick(() => {
            selectNavTab(tabsRefs.value[navTab.state.activeIndex])
        })
    }

    contextmenuRef.value?.onHideContextmenu()
}

const closeOtherTab = (menu: RouteLocationNormalized) => {
    navTab._closeTabs(menu)
    navTab._setActiveRoute(menu)
    if (navTab.state.activeRoute?.fullPath !== route.fullPath) {
        router.push(menu!.fullPath)
    }
}

/**
 * 关闭所有tab（等同于 navTabs.closeAllTab）
 * @param menu 需要保留的标签，否则关闭全部标签
 */
const closeAllTab = (menu?: RouteLocationNormalized) => {
    let firstRoute = getFirstMenu(menuStore.rawData)
    if (menu && firstRoute && firstRoute.path == menu.fullPath) {
        return closeOtherTab(menu)
    }
    if (firstRoute && firstRoute.path == navTab.state.activeRoute?.fullPath) {
        return closeOtherTab(navTab.state.activeRoute)
    }
    navTab._closeTabs(false)
    if (firstRoute) {
        router.push(firstRoute.path)
    }
}

const onContextmenu = (menu: RouteLocationNormalized, el: MouseEvent) => {
    // 禁用刷新
    state.contextmenuItems[0].disabled = route.fullPath !== menu.fullPath
    // 禁用关闭其他和关闭全部
    state.contextmenuItems[4].disabled = state.contextmenuItems[3].disabled = navTab.state.list.length == 1 ? true : false

    const { clientX, clientY } = el
    contextmenuRef.value?.onShowContextmenu(menu, {
        x: clientX,
        y: clientY,
    })
}

const onContextMenuClick = (item: ContextMenuItemClickEmitArg<RouteLocationNormalized>) => {
    const { name, sourceData } = item
    if (!sourceData) return
    switch (name) {
        case 'refresh':
            globalProperties.eventBus.emit('onTabViewRefresh', sourceData)
            break
        case 'close':
            closeTab(sourceData)
            break
        case 'closeOther':
            closeOtherTab(sourceData)
            break
        case 'closeAll':
            closeAllTab(sourceData)
            break
        case 'fullScreen':
            if (route.fullPath !== sourceData.fullPath) {
                router.push(sourceData.fullPath as string)
            }
            navTab.setActiveFullScreen(true)
            break
    }
}

const updateTab = function (newRoute: RouteLocationNormalized) {
    // 添加tab
    navTab._addTab(newRoute)
    // 激活当前tab
    navTab._setActiveRoute(newRoute)

    nextTick(() => {
        selectNavTab(tabsRefs.value[navTab.state.activeIndex])
    })
}

onBeforeRouteUpdate(async (to) => {
    updateTab(to)
})

onMounted(() => {
    updateTab(router.currentRoute.value)
    if (tabScrollbarRef.value) {
        new horizontalScroll(tabScrollbarRef.value)
    }
})

/**
 * 通过路由路径关闭 tab（等同于 navTabs.closeTabByPath）
 * @param fullPath 需要关闭的 tab 的路径
 */
const closeTabByPath = (fullPath: string) => {
    for (const key in navTab.state.list) {
        if (navTab.state.list[key].fullPath == fullPath) {
            closeTab(navTab.state.list[key])
            break
        }
    }
}

/**
 * 修改 tab 标题（等同于 navTabs.updateTabTitle）
 * @param fullPath 需要修改标题的 tab 的路径
 * @param title 新的标题
 */
const updateTabTitle = (fullPath: string, title: string) => {
    navTab._updateTabTitle(fullPath, title)
    nextTick(() => {
        selectNavTab(tabsRefs.value[navTab.state.activeIndex])
    })
}

defineExpose({
    closeAllTab,
    closeTabByPath,
    updateTabTitle,
})
</script>

<style scoped lang="scss">
.dark {
    .close-icon {
        color: v-bind('config.getColorValue("headerBarTabColor")') !important;
    }
    .ai-go-nav-tab.active {
        .close-icon {
            color: v-bind('config.getColorValue("headerBarTabActiveColor")') !important;
        }
    }
}
.nav-tabs {
    overflow-x: auto;
    overflow-y: hidden;
    margin-right: var(--ag-main-space);
    scrollbar-width: none;

    &::-webkit-scrollbar {
        height: 5px;
    }
    &::-webkit-scrollbar-thumb {
        background: #eaeaea;
        border-radius: var(--el-border-radius-base);
        box-shadow: none;
        -webkit-box-shadow: none;
    }
    &::-webkit-scrollbar-track {
        background: v-bind('config.layout.mode == "Default" ? "none":config.getColorValue("headerBarBackground")');
    }
    &:hover {
        &::-webkit-scrollbar-thumb:hover {
            background: #c8c9cc;
        }
    }
}
.ai-go-nav-tab {
    white-space: nowrap;
    height: 40px;
}
</style>
