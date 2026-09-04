import { adminBaseRoutePath } from '@/router/static/adminBase'
import type { RouteRecordRaw } from 'vue-router'

/*
 * 静态路由（支持自动扩展）
 * 系统会自动加载 ./static 目录及其子目录中的所有 .ts 文件
 * 每个模块的 default 导出可以是 RouteRecordRaw 或 RouteRecordRaw[]，自动 push 到以下 staticRoutes 数组
 * name 全小写使用 / 和 - 分割为了和动态路由统一；meta.title 只填写翻译 key，因为语言包加载是异步的，路由初始化时无法获取到翻译后的文本。
 */
const staticRoutes: Array<RouteRecordRaw> = [
    {
        // 首页
        path: '/',
        name: 'home',
        component: () => import('@/views/index.vue'),
        meta: {
            title: `pageTitles.Home`,
        },
    },
    {
        // 管理员登录页 - 路由不放在 adminBaseRoute.children 因为登录页不需要使用后台的布局
        path: adminBaseRoutePath + '/login',
        name: 'admin/login',
        component: () => import('@/views/admin/login.vue'),
        meta: {
            title: `pageTitles.Login`,
        },
    },
    {
        path: '/:path(.*)*',
        redirect: '/404',
    },
    {
        // 404
        path: '/404',
        name: 'not-found',
        component: () => import('@/views/common/error/404.vue'),
        meta: {
            title: `pageTitles.NotFound`, // 页面不存在
        },
    },
    {
        // 后台找不到页面了-可能是路由未加载上
        path: adminBaseRoutePath + ':path(.*)*',
        redirect: (to) => {
            return {
                name: 'admin/main-loading',
                params: {
                    to: JSON.stringify({
                        path: to.path,
                        query: to.query,
                    }),
                },
            }
        },
    },
]

const staticFiles = import.meta.glob('./static/**/*.ts', { eager: true })
for (const path in staticFiles) {
    const module = staticFiles[path] as any
    if (!module.default) {
        console.warn(`[Router] Static route module ${path} does not export default, skipped`)
        continue
    }
    const route = module.default as RouteRecordRaw | RouteRecordRaw[]
    if (Array.isArray(route)) {
        staticRoutes.push(...route)
    } else {
        staticRoutes.push(route)
    }
}

export default staticRoutes
