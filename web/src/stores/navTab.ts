import i18n from '@/lang/index'
import type { NavTab } from '@/stores/interface/index'
import { layoutNavTabsRef } from '@/stores/ref'
import { getMenuDataByRoute } from '@/utils/common'
import { isEmpty } from 'lodash-es'
import { defineStore } from 'pinia'
import { reactive } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export const useNavTab = defineStore('navTab', () => {
    const state: NavTab = reactive({
        list: [],
        activeIndex: 0,
        activeRoute: null,
        activeFullScreen: false,
    })

    /**
     * 通过路由路径关闭tab
     * @param fullPath 需要关闭的 tab 的路径
     */
    const closeTabByPath = (fullPath: string) => {
        layoutNavTabsRef.value?.closeTabByPath(fullPath)
    }

    /**
     * 关闭所有tab
     * @param menu 需要保留的标签，否则关闭全部标签并打开第一个路由
     */
    const closeAllTab = (menu?: RouteLocationNormalized) => {
        layoutNavTabsRef.value?.closeAllTab(menu)
    }

    /**
     * 修改 tab 标题
     * @param fullPath 需要修改标题的 tab 的路径
     * @param title 新的标题
     */
    const updateTabTitle = (fullPath: string, title: string) => {
        layoutNavTabsRef.value?.updateTabTitle(fullPath, title)
    }

    /**
     * 添加 tab（内部）
     * ps: router.push 时可自动完成 tab 添加，无需调用此方法
     */
    function _addTab(route: RouteLocationNormalized) {
        const tabView = { ...route, matched: [], meta: { ...route.meta } }
        if (!tabView.meta.addtab) return

        // 通过路由寻找菜单的原始数据
        const menuData = getMenuDataByRoute(tabView)
        if (menuData && menuData.meta) {
            tabView.name = menuData.name
            tabView.meta.id = menuData.meta.id
            tabView.meta.title = menuData.meta.title
        }

        for (const key in state.list) {
            // 菜单已在 tabs 存在，更新 params 和 query
            if (state.list[key].meta.id === tabView.meta.id || state.list[key].fullPath == tabView.fullPath) {
                state.list[key].fullPath = tabView.fullPath
                state.list[key].params = !isEmpty(tabView.params) ? tabView.params : state.list[key].params
                state.list[key].query = !isEmpty(tabView.query) ? tabView.query : state.list[key].query
                return
            }
        }
        if (typeof tabView.meta.title == 'string') {
            tabView.meta.title = i18n.global.te(tabView.meta.title) ? i18n.global.t(tabView.meta.title) : tabView.meta.title
        }
        state.list.push(tabView)
    }

    /**
     * 设置激活 tab（内部）
     * ps: router.push 时可自动完成 tab 激活，无需调用此方法
     */
    const _setActiveRoute = (route: RouteLocationNormalized): void => {
        const currentRouteIndex: number = state.list.findIndex((item: RouteLocationNormalized) => {
            return item.fullPath === route.fullPath
        })
        if (currentRouteIndex === -1) return
        state.activeRoute = route
        state.activeIndex = currentRouteIndex
    }

    /**
     * 关闭 tab（内部）
     * ps: 使用 closeTabByPath 代替
     */
    function _closeTab(route: RouteLocationNormalized) {
        state.list.map((v, k) => {
            if (v.fullPath == route.fullPath) {
                state.list.splice(k, 1)
                return
            }
        })
    }

    /**
     * 关闭多个标签（内部）
     * ps：使用 closeAllTab 代替
     */
    const _closeTabs = (retainMenu: RouteLocationNormalized | false = false) => {
        if (retainMenu) {
            state.list = [retainMenu]
        } else {
            state.list = []
        }
    }

    /**
     * 更新标签标题（内部）
     * ps: 使用 updateTabTitle 代替
     */
    const _updateTabTitle = (fullPath: string, title: string) => {
        for (const key in state.list) {
            if (state.list[key].fullPath == fullPath) {
                state.list[key].meta.title = title
                break
            }
        }
    }

    /**
     * 设置当前 tab 是否全屏
     * @param status 全屏状态
     */
    const setActiveFullScreen = (status: boolean): void => {
        state.activeFullScreen = status
    }

    return {
        state,
        closeAllTab,
        closeTabByPath,
        updateTabTitle,
        setActiveFullScreen,
        _addTab,
        _closeTab,
        _closeTabs,
        _setActiveRoute,
        _updateTabTitle,
    }
})
