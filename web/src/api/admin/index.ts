import type { ClickRequest } from '/@/components/clickCaptcha/index'
import request from '/@/utils/request'

/**
 * 管理员登录请求参数
 */
export interface AdminLoginParams {
    username: string
    password: string
    remember: boolean
    captcha: ClickRequest
}

/**
 * 管理员登录请求
 */
export function adminLogin(data: AdminLoginParams) {
    return request({
        url: '/admin/login',
        method: 'POST',
        data,
    })
}

/**
 * 管理员注销请求
 */
export function adminLogout() {
    return request({
        url: '/admin/logout',
        method: 'POST',
    })
}
