<template>
    <component :is="config.layout.mode"></component>
</template>

<script setup lang="ts">
import { useEventListener } from '@vueuse/core'
import { isEmpty } from 'lodash-es'
import { onBeforeMount, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { getInit } from '/@/api/admin/index'
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
    getInit().then((res) => {
        adminInfo.dataFill({ ...res.data.data.admin, super: res.data.data.super }, ['token'])

        config.siteDataFill(res.data.data.siteConfig)
        config.setSiteInitStatus(true)

        if (res.data.data.rules.length) {
            handleAdminRoute(res.data.data.rules)

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
    })
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
