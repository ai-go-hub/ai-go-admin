import { camelCase, snakeCase, trimStart } from 'lodash-es'
import { getCurrentInstance } from 'vue'
import { RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import router from '/@/router/index'
import { adminBaseRoutePath } from '/@/router/static/adminBase'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { getBaseUrl } from '/@/utils/request'

/**
 * 获取资源完整地址
 * @param resource 资源相对地址
 * @param domain 指定域名
 */
export const fullURL = (resource: string, domain = '') => {
    const config = useConfig()
    if (!domain) {
        domain = config.site.cdnUrl ? config.site.cdnUrl : getBaseUrl()
    }
    if (!resource) {
        return ''
    }

    const regUrl = new RegExp(/^http(s)?:\/\//)
    const regexData = new RegExp(/^data:/i)
    if (!domain || regUrl.test(resource) || regexData.test(resource)) {
        return resource
    }

    let url = domain + resource
    if (domain === config.site.cdnUrl && config.site.cdnUrlParams) {
        const separator = url.includes('?') ? '&' : '?'
        url += separator + config.site.cdnUrlParams
    }
    return url
}

/**
 * 获取 globalProperties 对象
 */
export const getGlobalProperties = () => {
    if (!getCurrentInstance()) {
        throw new Error('getGlobalProperties() can only be used inside setup() or functional components!')
    }
    const instance = getCurrentInstance()
    if (instance) {
        const { appContext } = instance
        return appContext.config.globalProperties
    } else {
        return null
    }
}

/**
 * 复制文本到剪贴板
 * @param text 要复制的文本
 * @returns 复制成功返回 true，失败返回 false
 */
export async function copy(text: string): Promise<boolean> {
    try {
        await navigator.clipboard.writeText(text)
        return true
    } catch {
        // 降级方案：使用 execCommand（兼容旧浏览器或非 HTTPS 环境）
        const textarea = document.createElement('textarea')
        textarea.value = text
        textarea.style.position = 'fixed'
        textarea.style.left = '-9999px'
        textarea.style.top = '-9999px'
        document.body.appendChild(textarea)
        textarea.focus()
        textarea.select()
        try {
            document.execCommand('copy')
            return true
        } catch {
            return false
        } finally {
            document.body.removeChild(textarea)
        }
    }
}

/**
 * 获取路由 path
 */
export const getCurrentRoutePath = () => {
    let path = router.currentRoute.value.path
    if (path == '/') path = trimStart(window.location.hash, '#')
    if (path.indexOf('?') !== -1) path = path.replace(/\?.*/, '')
    return path
}

/**
 * 是否在后台应用内
 * @param path 不传递则通过当前路由 path 检查
 */
export const isAdminApp = (path = '') => {
    const regex = new RegExp(`^${adminBaseRoutePath}`)
    if (path) {
        return regex.test(path)
    }
    if (regex.test(getCurrentRoutePath())) {
        return true
    }
    return false
}

/**
 * 递归的寻找路由路径在菜单中的数据
 * @param path 路由路径
 * @param menus 菜单数据（只有 path 代表完整 url，没有 fullPath）
 * @param returnType 返回值要求:normal=返回被搜索的路径对应的菜单数据,above=返回被搜索的路径对应的上一级菜单数组
 */
export const getMenuDataByPath = (path: string, menus: RouteRecordRaw[], returnType: 'normal' | 'above'): RouteRecordRaw | false => {
    for (const key in menus) {
        // 找到目标
        if (menus[key].path === path) {
            return menus[key]
        }
        // 从子级继续寻找
        if (menus[key].children && menus[key].children.length) {
            const find = getMenuDataByPath(path, menus[key].children, returnType)
            if (find) {
                return returnType == 'above' ? menus[key] : find
            }
        }
    }
    return false
}

/**
 * 寻找路由在菜单中的数据
 * @param route 路由
 * @param returnType 返回值要求:normal=返回被搜索的路径对应的菜单数据,above=返回被搜索的路径对应的上一级菜单数组
 */
export const getMenuDataByRoute = (
    route: RouteLocationNormalized | RouteRecordRaw,
    returnType: 'normal' | 'above' = 'normal'
): RouteRecordRaw | false => {
    const menu = useMenu()
    let found: RouteRecordRaw | false = false
    const fullPath = (route as RouteLocationNormalized).fullPath
    if (fullPath) {
        // 以完整路径寻找
        found = getMenuDataByPath(fullPath, menu.rawData, returnType)
        if (found) {
            found.meta!.matched = fullPath
            return found
        }
    }

    // 以路径寻找
    found = getMenuDataByPath(route.path, menu.rawData, returnType)
    if (found) {
        found.meta!.matched = route.path
        return found
    }

    return false
}

/**
 * 递归将对象 key 从 snake_case 转为 camelCase
 */
export function keysToCamelCase(obj: any): any {
    if (Array.isArray(obj)) {
        return obj.map((item) => keysToCamelCase(item))
    }
    if (obj !== null && typeof obj === 'object') {
        const result: Record<string, any> = {}
        for (const key of Object.keys(obj)) {
            result[camelCase(key)] = keysToCamelCase(obj[key])
        }
        return result
    }
    return obj
}

/**
 * 递归将对象 key 从 camelCase 转为 snake_case
 */
export function keysToSnakeCase(obj: any): any {
    if (Array.isArray(obj)) {
        return obj.map((item) => keysToSnakeCase(item))
    }
    if (obj !== null && typeof obj === 'object') {
        const result: Record<string, any> = {}
        for (const key of Object.keys(obj)) {
            result[snakeCase(key)] = keysToSnakeCase(obj[key])
        }
        return result
    }
    return obj
}
