import { camelCase, snakeCase, trimStart } from 'lodash-es'
import { RouteLocationNormalized, RouteRecordRaw } from 'vue-router'
import router from '/@/router/index'
import { adminBaseRoutePath } from '/@/router/static/adminBase'
import { useConfig } from '/@/stores/config'
import { useMenu } from '/@/stores/menu'
import { getBaseURL } from '/@/utils/request'

/**
 * 获取资源完整地址
 * @param resource 资源相对地址
 * @param domain 指定域名
 */
export const fullURL = (resource: string, domain = '') => {
    const config = useConfig()
    if (!domain) {
        domain = config.site.cdnUrl ? config.site.cdnUrl : getBaseURL()
    }
    if (!resource) {
        return ''
    }

    const regURL = new RegExp(/^http(s)?:\/\//)
    const regexData = new RegExp(/^data:/i)
    if (!domain || regURL.test(resource) || regexData.test(resource)) {
        return resource
    }

    let url = domain + resource
    if (domain === config.site.cdnUrl && config.site.cdnUrlParams) {
        const separator = url.includes('?') ? '&' : '?'
        url += separator + config.site.cdnUrlParams
    }
    return url
}

export function auth(node: string): boolean
export function auth(node: { name: string; subNodeName?: string }): boolean

/**
 * 鉴权
 * 提供 string 将根据当前路由 path 自动拼接和鉴权，还可以提供路由的 name 对象进行鉴权
 */
export function auth(node: string | { name: string; subNodeName?: string }) {
    const menu = useMenu()
    if (typeof node === 'string') {
        const path = getCurrentRoutePath()
        if (menu.authNode.has(path)) {
            const subNodeName = path + (path == '/' ? '' : '/') + node
            if (menu.authNode.get(path)!.some((v: string) => v == subNodeName)) {
                return true
            }
        }
    } else {
        // 节点列表中没有找到 name
        if (!node.name || !menu.authNode.has(node.name)) return false

        // 无需继续检查子节点或未找到子节点
        if (!node.subNodeName || menu.authNode.get(node.name)?.includes(node.subNodeName)) return true
    }
    return false
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
 * 以 pk 字段从数组中获取对应的索引值
 * @param arr
 * @param pk
 * @param value
 */
export const getArrayKey = (arr: any[], pk: string, value: any): any => {
    for (const key in arr) {
        if (arr[key][pk] == value) {
            return key
        }
    }
    return false
}

/**
 * 从一个文件路径中获取文件名
 * @param path 文件路径
 */
export const getFileNameFromPath = (path: string) => {
    const paths = path.split('/')
    return paths[paths.length - 1]
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

/**
 * 获取一组资源的完整地址
 * @param resources 资源相对地址
 * @param domain 指定域名
 */
export const arrayFullURL = (resources: string | string[], domain = '') => {
    if (typeof resources === 'string') {
        resources = resources == '' ? [] : resources.split(',')
    }
    for (const key in resources) {
        resources[key] = fullURL(resources[key], domain)
    }
    return resources
}

interface ElTreeData {
    label: string
    children?: ElTreeData[]
}

/**
 * 将数据构建为 ElTree 的 data {label:'', children: []}
 * @param data
 */
export const buildJsonToElTreeData = (data: any): ElTreeData[] => {
    if (typeof data == 'object') {
        const childrens = []
        for (const key in data) {
            childrens.push({
                label: key + ': ' + data[key],
                children: buildJsonToElTreeData(data[key]),
            })
        }
        return childrens
    } else {
        return []
    }
}
