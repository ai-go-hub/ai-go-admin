<template>
    <component :is="config.layout.mode"></component>
</template>

<script setup lang="ts">
import { useEventListener } from '@vueuse/core'
import { isEmpty } from 'lodash-es'
import { onBeforeMount, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'
import Classic from '/@/layouts/admin/container/classic.vue'
import Default from '/@/layouts/admin/container/default.vue'
import Double from '/@/layouts/admin/container/double.vue'
import LeftSplit from '/@/layouts/admin/container/leftSplit.vue'
import Streamline from '/@/layouts/admin/container/streamline.vue'
import router from '/@/router/index'
import { adminBaseRoutePath } from '/@/router/static/adminBase'
import { useAdminInfo } from '/@/stores/adminInfo'
import { useConfig } from '/@/stores/config'
import { BEFORE_RESIZE_LAYOUT } from '/@/stores/constant/cacheKey'
import { useMenu } from '/@/stores/menu'
import { setNavTabsWidth } from '/@/utils/layout'
import { getFirstMenu, handleAdminRoute } from '/@/utils/router'
import { Session } from '/@/utils/storage'
import { keysToCamelCase } from '/@/utils/common'

defineOptions({
    components: { Default, Classic, Streamline, Double, LeftSplit },
})

const menu = useMenu()
const route = useRoute()
const config = useConfig()
const adminInfo = useAdminInfo()

const state = reactive({
    autoMenuCollapseLock: false,
})

onMounted(() => {
    if (!adminInfo.token) {
        return router.push({ name: 'adminLogin' })
    }

    init()
    setNavTabsWidth()
    useEventListener(window, 'resize', setNavTabsWidth)
})
onBeforeMount(() => {
    onAdaptiveLayout()
    useEventListener(window, 'resize', onAdaptiveLayout)
})

/**
 * 后台初始化请求，获取站点配置，动态路由等信息
 */
const init = () => {
    // 临时 MOCK 数据
    const res = keysToCamelCase({
        data: {
            site: {
                name: 'AI GO MALL',
                version: '1.0.0',
                timezone: 'Asia/Shanghai',
                record_number: '渝ICP备8888888号',
                initialized: true,
            },
            adminInfo: {
                id: 1,
                username: 'admin',
                nickname: '超级管理员',
                avatar: '/static/avatar.png',
                last_login_at: '2021-01-01 00:00:00',
                last_login_ip: '127.0.0.1',
                token: 'token',
            },
            menus: [
                {
                    id: 1,
                    type: 'menu',
                    title: '控制台',
                    name: 'dashboard',
                    path: 'dashboard',
                    icon: 'lucide-layout-dashboard',
                    menu_type: 'tab',
                    component: '/src/views/admin/dashboard.vue',
                    keepalive: 'dashboard',
                    extend: 'none',
                    children: [
                        {
                            id: 89,
                            pid: 1,
                            type: 'button',
                            title: '查看',
                            name: 'dashboard/index',
                        },
                    ],
                },
                {
                    id: 2,
                    type: 'menu_dir',
                    title: '权限管理',
                    name: 'auth',
                    path: 'auth',
                    icon: 'lucide-users',
                    extend: 'none',
                    children: [
                        {
                            id: 3,
                            pid: 2,
                            type: 'menu',
                            title: '角色组管理',
                            name: 'auth/group',
                            path: 'auth/group',
                            icon: 'lucide-users-round',
                            menu_type: 'tab',
                            component: '/src/views/admin/dashboard.vue',
                            keepalive: 'auth/group',
                            extend: 'none',
                            children: [
                                {
                                    id: 4,
                                    pid: 3,
                                    type: 'button',
                                    title: '查看',
                                    name: 'auth/group/index',
                                },
                                {
                                    id: 5,
                                    pid: 3,
                                    type: 'button',
                                    title: '添加',
                                    name: 'auth/group/add',
                                },
                            ],
                        },
                        {
                            id: 8,
                            pid: 2,
                            type: 'menu',
                            title: '管理员管理',
                            name: 'auth/admin',
                            path: 'auth/admin',
                            icon: 'lucide-user-lock',
                            menu_type: 'tab',
                            component: '/src/views/backend/auth/admin/index.vue',
                            keepalive: 'auth/admin',
                            extend: 'none',
                            children: [
                                {
                                    id: 9,
                                    pid: 8,
                                    type: 'button',
                                    title: '查看',
                                    name: 'auth/admin/index',
                                },
                                {
                                    id: 11,
                                    pid: 8,
                                    type: 'button',
                                    title: '编辑',
                                    name: 'auth/admin/edit',
                                },
                                {
                                    id: 12,
                                    pid: 8,
                                    type: 'button',
                                    title: '删除',
                                    name: 'auth/admin/del',
                                },
                            ],
                        },
                    ],
                },
            ],
        },
    })

    config.siteDataFill(res.data.site)
    config.setSiteInitStatus(true)

    if (!isEmpty(res.data.adminInfo)) {
        adminInfo.dataFill({ ...res.data.adminInfo, super: true }, ['token'])
    }

    if (res.data.menus) {
        handleAdminRoute(res.data.menus)

        // 显示布局引导
        if (config.layout.tourUnfinished) {
            setTimeout(() => {
                config.setLayoutValue('showTour', true)
            }, 1000)
        }

        // 预跳转到上次路径
        if (route.params.to) {
            const lastRoute = JSON.parse(route.params.to as string)
            if (lastRoute.path != adminBaseRoutePath) {
                let query = !isEmpty(lastRoute.query) ? lastRoute.query : {}
                router.push({ path: lastRoute.path, query: query })
                return
            }
        }

        // 跳转到第一个菜单
        let firstRoute = getFirstMenu(menu.rawData)
        if (firstRoute) {
            router.push(firstRoute.path)
        }
    }
}

const onAdaptiveLayout = () => {
    let defaultBeforeResizeLayout = {
        menuCollapse: config.layout.menuCollapse,
    }
    let beforeResizeLayout = Session.get(BEFORE_RESIZE_LAYOUT)
    if (!beforeResizeLayout) {
        Session.set(BEFORE_RESIZE_LAYOUT, defaultBeforeResizeLayout)
    }

    const clientWidth = document.body.clientWidth
    if (clientWidth < 1024) {
        /**
         * 锁定窗口改变自动调整 menuCollapse
         * 避免已是小窗且打开了菜单栏时，意外的自动关闭菜单栏
         */
        if (!state.autoMenuCollapseLock) {
            state.autoMenuCollapseLock = true
            config.setLayoutValue('menuCollapse', true)
        }
        config.setLayoutValue('shrink', true)
    } else {
        state.autoMenuCollapseLock = false
        let beforeResizeLayoutTemp = beforeResizeLayout || defaultBeforeResizeLayout

        config.setLayoutValue('menuCollapse', beforeResizeLayoutTemp.menuCollapse)
        config.setLayoutValue('shrink', false)
    }
}
</script>
