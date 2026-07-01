import type { RouteRecordRaw } from 'vue-router'

/**
 * 后台基础路由路径
 */
export const adminBaseRoutePath = '/admin'

/*
 * 后台基础静态路由
 */
const adminBaseRoute: RouteRecordRaw = {
    path: adminBaseRoutePath,
    name: 'admin',
    component: () => import('/@/layouts/admin/index.vue'),
    // 直接重定向到 loading 路由
    redirect: adminBaseRoutePath + '/loading',
    meta: {
        title: `common.loading`,
    },
    children: [
        {
            path: 'loading/:to?',
            name: 'adminMainLoading',
            component: () => import('/@/layouts/admin/loading.vue'),
            meta: {
                title: `common.loading`,
            },
        },
    ],
}

export default adminBaseRoute
