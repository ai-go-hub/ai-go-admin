import { ElNotification } from 'element-plus'
import { compact, reverse } from 'lodash-es'
import type { RouteRecordRaw } from 'vue-router'
import i18n from '/@/lang/index'
import router from '/@/router/index'
import adminBaseRoute from '/@/router/static/adminBase'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { isAdminApp } from '/@/utils/common'
import { closeMask } from '/@/utils/mask'

/**
 * 获取第一个菜单
 */
export const getFirstMenu = (routes: RouteRecordRaw[]): false | RouteRecordRaw => {
    const routerPaths: string[] = []
    const routers = router.getRoutes()
    routers.forEach((item) => {
        if (item.path) {
            routerPaths.push(item.path)
        }
    })
    let find: boolean | RouteRecordRaw = false
    for (const key in routes) {
        if (routes[key].meta?.type == 'menu' && routerPaths.indexOf(routes[key].path) !== -1) {
            return routes[key]
        } else if (routes[key].children && routes[key].children?.length) {
            find = getFirstMenu(routes[key].children!)
            if (find) return find
        }
    }
    return find
}

/**
 * 打开菜单
 * @param menu 菜单数据
 */
export const openMenu = (menu: RouteRecordRaw) => {
    switch (menu.meta?.openType) {
        case 'iframe':
        case 'tab':
            router.push(menu.path)
            break
        case 'link':
            window.open(menu.path, '_blank')
            break

        default:
            ElNotification({ message: i18n.global.t('layouts.menuTypeUnrecognized'), type: 'error' })
            break
    }

    const config = useConfig()
    if (config.layout.shrink) {
        closeMask(() => {
            config.setLayoutValue('menuCollapse', true)
        })
    }
}

/**
 * 处理后台的路由
 */
export const handleAdminRoute = (routes: any) => {
    const viewsComponent = import.meta.glob('/src/views/admin/**/*.vue')
    addRouteAll(viewsComponent, routes, adminBaseRoute.name as string)
    const menuAdminBaseRoute = (adminBaseRoute.path as string) + '/'

    // 更新stores中的路由菜单数据
    const menu = useMenu()
    menu.setRawData(handleMenuRule(routes, menuAdminBaseRoute))
    menu.authNodeFill(handleAuthNode(routes, menuAdminBaseRoute))
}

/**
 * 获取菜单的paths
 */
export const getMenuPaths = (menus: RouteRecordRaw[]): string[] => {
    let menuPaths: string[] = []
    menus.forEach((item) => {
        menuPaths.push(item.path)
        if (item.children && item.children.length > 0) {
            menuPaths = menuPaths.concat(getMenuPaths(item.children))
        }
    })
    return menuPaths
}

/**
 * 获取菜单唯一标识
 * @param menu 菜单数据
 * @param prefix 前缀
 */
export const getMenuKey = (menu: RouteRecordRaw, prefix = '') => {
    if (prefix === '') {
        prefix = menu.path
    }
    return `${prefix}-${menu.name as string}-${menu.meta && menu.meta.id ? menu.meta.id : ''}`
}

/**
 * 菜单处理
 */
const handleMenuRule = (routes: any, pathPrefix = '/', type = ['menu', 'dir']) => {
    const menuRule: RouteRecordRaw[] = []
    for (const key in routes) {
        if (routes[key].extend == 'add_route_only') {
            continue
        }
        if (!type.includes(routes[key].type)) {
            continue
        }
        if (routes[key].type == 'dir' && routes[key].children && !routes[key].children.length) {
            continue
        }
        if (
            ['route', 'menu'].includes(routes[key].type) &&
            ((routes[key].openType == 'tab' && !routes[key].component) || (['link', 'iframe'].includes(routes[key].openType) && !routes[key].url))
        ) {
            continue
        }
        const currentPath = ['link', 'iframe'].includes(routes[key].openType) ? routes[key].url : pathPrefix + routes[key].path
        let children: RouteRecordRaw[] = []
        if (routes[key].children && routes[key].children.length > 0) {
            children = handleMenuRule(routes[key].children, pathPrefix, type)
        }
        menuRule.push({
            path: currentPath,
            name: routes[key].name,
            component: routes[key].component,
            meta: {
                id: routes[key].id,
                title: routes[key].title,
                icon: routes[key].icon,
                keepalive: routes[key].keepalive,
                openType: routes[key].openType,
                type: routes[key].type,
            },
            children: children,
        })
    }
    return menuRule
}

/**
 * 处理权限节点
 * @param routes 路由数据
 * @param prefix 节点前缀
 * @returns 组装好的权限节点
 */
const handleAuthNode = (routes: any, prefix = '/') => {
    const authNode: Map<string, string[]> = new Map([])
    assembleAuthNode(routes, authNode, prefix, prefix)
    return authNode
}
const assembleAuthNode = (routes: any, authNode: Map<string, string[]>, prefix = '/', parent = '/') => {
    const authNodeTemp = []
    for (const key in routes) {
        if (routes[key].type == 'node') {
            authNodeTemp.push(prefix + routes[key].name)
        }
        if (routes[key].children && routes[key].children.length > 0) {
            assembleAuthNode(routes[key].children, authNode, prefix, prefix + routes[key].name)
        }
    }
    if (authNodeTemp && authNodeTemp.length > 0) {
        authNode.set(parent, authNodeTemp)
    }
}

/**
 * 动态添加路由-带子路由
 * @param viewsComponent
 * @param routes
 * @param parentName
 * @param analyticRelation 根据 name 从已注册路由分析父级路由
 */
export const addRouteAll = (viewsComponent: Record<string, any>, routes: any, parentName: string, analyticRelation = false) => {
    for (const idx in routes) {
        if (routes[idx].extend == 'add_menu_only') {
            continue
        }
        if ((routes[idx].openType == 'tab' && viewsComponent[routes[idx].component]) || routes[idx].openType == 'iframe') {
            addRouteItem(viewsComponent, routes[idx], parentName, analyticRelation)
        }

        if (routes[idx].children && routes[idx].children.length > 0) {
            addRouteAll(viewsComponent, routes[idx].children, parentName, analyticRelation)
        }
    }
}

/**
 * 动态添加路由
 * @param viewsComponent
 * @param route
 * @param parentName
 * @param analyticRelation 根据 name 从已注册路由分析父级路由
 */
export const addRouteItem = (viewsComponent: Record<string, any>, route: any, parentName: string, analyticRelation: boolean) => {
    let path = '',
        component
    if (route.openType == 'iframe') {
        path = (isAdminApp() ? adminBaseRoute.path : '/') + '/iframe/' + encodeURIComponent(route.url)
        component = () => import('/@/layouts/common/router-view/iframe.vue')
    } else {
        path = parentName ? route.path : '/' + route.path
        component = viewsComponent[route.component]
    }

    if (route.openType == 'tab' && analyticRelation) {
        const parentNames = getParentNames(route.name)
        if (parentNames.length) {
            for (const key in parentNames) {
                if (router.hasRoute(parentNames[key])) {
                    parentName = parentNames[key]
                    break
                }
            }
        }
    }

    const routeBaseInfo: RouteRecordRaw = {
        path: path,
        name: route.name,
        component: component,
        meta: {
            title: route.title,
            extend: route.extend,
            icon: route.icon,
            keepalive: route.keepalive,
            openType: route.openType,
            type: route.type,
            url: route.url,
            addtab: true,
        },
    }
    if (parentName) {
        router.addRoute(parentName, routeBaseInfo)
    } else {
        router.addRoute(routeBaseInfo)
    }
}

/**
 * 根据name字符串，获取父级name组合的数组
 * @param name
 */
const getParentNames = (name: string) => {
    const names = compact(name.split('/'))
    const tempNames = []
    const parentNames = []
    for (const key in names) {
        tempNames.push(names[key])
        if (parseInt(key) != names.length - 1) {
            parentNames.push(tempNames.join('/'))
        }
    }
    return reverse(parentNames)
}
