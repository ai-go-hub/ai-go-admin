import { useTitle } from '@vueuse/core'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { createRouter, createWebHashHistory } from 'vue-router'
import i18n from '/@/lang/index'
import staticRoutes from '/@/router/static'
import { useConfig } from '/@/stores/config'
import { loading } from '/@/utils/loading'

const router = createRouter({
    history: createWebHashHistory(),
    routes: staticRoutes,
})

// 路由加载前
router.beforeEach(() => {
    NProgress.configure({ showSpinner: false })
    NProgress.start()

    if (!window.loading) {
        loading.show()
        window.loading = true
    }
})

// 路由加载后
router.afterEach((to) => {
    if (window.loading) {
        loading.hide()
    }
    NProgress.done()

    // 设置浏览器标题
    const config = useConfig()
    const titleKey = to?.meta?.title as string | undefined
    const title = titleKey && i18n.global.te(titleKey) ? i18n.global.t(titleKey) : ''
    useTitle().value = title ? `${title}${config.site.name ? ' - ' + config.site.name : ''}` : `${config.site.name}`
})

export default router
