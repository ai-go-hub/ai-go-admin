import * as elIcons from '@element-plus/icons-vue'
import { camelCase, kebabCase, upperFirst } from 'lodash-es'
import { App, defineAsyncComponent, type Component } from 'vue'
import Icon from '/@/components/icon/index.vue'

type IconMap = Record<string, Component>

const lucideCache: IconMap = {}

/*
 * 归一化图标名称: 去除所有分隔符并转小写
 */
function normalizeIconName(s: string): string {
    return s.toLowerCase().replace(/[^a-z0-9]/g, '')
}

/*
 * 在图标映射中查找归一化后与目标名称一致的导出键
 */
function findIconKey(icons: IconMap, name: string): string | undefined {
    const target = normalizeIconName(name)
    return Object.keys(icons).find((k) => normalizeIconName(k) === target)
}

export function getLucideComponent(name: string): Component | null {
    const key = upperFirst(camelCase(name))
    if (lucideCache[key]) {
        return lucideCache[key]
    }

    let loader: () => Promise<Component | { render: () => null }>

    if (import.meta.env.DEV) {
        let iconsPromise: Promise<IconMap> | null = null
        loader = () =>
            (iconsPromise ??= import('@lucide/vue').then((m) => m.icons as IconMap)).then((icons) => {
                const matchedKey = icons[key] ? key : findIconKey(icons, name)
                return (matchedKey && icons[matchedKey]) || { render: () => null }
            })
    } else {
        const batchLoaders: Record<string, () => Promise<IconMap>> = {
            'a-c': () => import('virtual:lucide-icons/a-c') as Promise<IconMap>,
            'd-l': () => import('virtual:lucide-icons/d-l') as Promise<IconMap>,
            'm-p': () => import('virtual:lucide-icons/m-p') as Promise<IconMap>,
            'q-s': () => import('virtual:lucide-icons/q-s') as Promise<IconMap>,
            't-z': () => import('virtual:lucide-icons/t-z') as Promise<IconMap>,
        }

        const first = key.charAt(0).toLowerCase()
        const batch = 'abc'.includes(first)
            ? 'a-c'
            : 'defghijkl'.includes(first)
              ? 'd-l'
              : 'mnop'.includes(first)
                ? 'm-p'
                : 'qrs'.includes(first)
                  ? 'q-s'
                  : 't-z'

        loader = () =>
            batchLoaders[batch]().then((icons) => {
                const matchedKey = icons[key] ? key : findIconKey(icons, name)
                return (matchedKey && icons[matchedKey]) || { render: () => null }
            })
    }

    const asyncComp = defineAsyncComponent(loader)
    lucideCache[key] = asyncComp
    return asyncComp
}

/*
 * 全局注册 icon 组件
 */
export function registerIcons(app: App) {
    /*
     * 全局注册 Icon 组件
     * 使用方式: <Icon name="name" size="size" color="color" strokeWidth="线条宽度（仅 lucide 支持）" />
     * name 支持两种前缀 el- 和 lucide-，分别表示 element plus 图标和 lucide 图标
     */
    app.component('Icon', Icon)

    /*
     * Element Plus 的 icon
     */
    const icons = elIcons as any
    for (const i in icons) {
        app.component(`el-icon-${kebabCase(icons[i].name)}`, icons[i])
    }
}
