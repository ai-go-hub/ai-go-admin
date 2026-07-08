import { defineStore } from 'pinia'
import type { RouteRecordRaw } from 'vue-router'
import { adminBaseRoutePath } from '/@/router/static/adminBase'
import type { Menu } from '/@/stores/interface/index'

export const useMenu = defineStore('menu', {
    state: (): Menu => {
        return {
            rawData: [],
            children: [],
            authNode: new Map(),
        }
    },
    actions: {
        /**
         * 设置从后台加载到的菜单路由列表
         */
        setRawData(data: RouteRecordRaw[]): void {
            this.rawData = encodeRoutesURI(data)
        },

        /**
         * 设置次级菜单数据
         */
        setChildren(data: RouteRecordRaw[]): void {
            this.children = data
        },

        /**
         * 以key设置权限节点
         */
        setAuthNode(key: string, data: string[]) {
            this.authNode.set(key, data)
        },

        /**
         * 覆盖设置权限节点
         */
        authNodeFill(data: Map<string, string[]>) {
            this.authNode = data
        },
    },
})

/**
 * 对iframe的url进行编码
 */
function encodeRoutesURI(data: RouteRecordRaw[]) {
    data.forEach((item) => {
        if (item.meta?.openType == 'iframe') {
            item.path = adminBaseRoutePath + '/iframe/' + encodeURIComponent(item.path)
        }

        if (item.children && item.children.length) {
            item.children = encodeRoutesURI(item.children)
        }
    })
    return data
}
