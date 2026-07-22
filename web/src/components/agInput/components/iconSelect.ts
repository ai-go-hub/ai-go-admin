import * as elIcons from '@element-plus/icons-vue'
import { kebabCase } from 'lodash-es'
import { Component } from 'vue'

/*
 * 获取 Element Plus 图标名称列表
 */
export function getElementPlusIconNames() {
    return new Promise<string[]>((resolve, reject) => {
        const names: string[] = []
        const icons = elIcons as any
        for (const i in icons) {
            names.push(`el-${kebabCase(icons[i].name)}`)
        }
        if (names.length > 0) {
            resolve(names)
        } else {
            reject('Element Plus icons not found')
        }
    })
}

/*
 * 获取 Lucide 图标名称列表
 *
 * 开发环境：直接从 @lucide/vue 的 icons 对象获取所有图标名称
 * 生产环境：从 virtual:lucide-icons/* 虚拟模块获取，与 Icon 组件共用已拆分的 chunk，避免重复加载全量模块
 */
export function getLucideIconNames() {
    return new Promise<string[]>((resolve, reject) => {
        const names: string[] = []

        const collectIcons = (module: Record<string, Component>) => {
            for (const key in module) {
                names.push(`lucide-${kebabCase(key)}`)
            }
        }

        if (import.meta.env.DEV) {
            import('@lucide/vue')
                .then((module) => {
                    collectIcons(module.icons as Record<string, Component>)
                    names.length > 0 ? resolve(names) : reject('Lucide icons not found')
                })
                .catch(reject)
        } else {
            Promise.all([
                import('virtual:lucide-icons/a-c'),
                import('virtual:lucide-icons/d-l'),
                import('virtual:lucide-icons/m-p'),
                import('virtual:lucide-icons/q-s'),
                import('virtual:lucide-icons/t-z'),
            ])
                .then((batches) => {
                    for (const batch of batches) {
                        collectIcons(batch as Record<string, Component>)
                    }
                    names.length > 0 ? resolve(names) : reject('Lucide icons not found')
                })
                .catch(reject)
        }
    })
}
