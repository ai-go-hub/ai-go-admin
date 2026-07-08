import type { ClickRequest } from '/@/components/clickCaptcha/index'
import request from '/@/utils/request'

/**
 * 管理员登录请求参数
 */
export interface AdminLoginParams {
    username: string
    password: string
    remember: boolean
    captcha?: ClickRequest
}

/**
 * 管理员登录请求
 */
export function login(data: AdminLoginParams) {
    return request({
        url: '/admin/login',
        method: 'POST',
        data,
    })
}

/**
 * 管理员注销请求
 */
export function logout() {
    return request({
        url: '/admin/logout',
        method: 'POST',
    })
}

/**
 * 获取管理员登录配置
 */
export function getLoginConfig() {
    return request({
        url: '/admin/login-config',
        method: 'GET',
    })
}

/**
 * 后台初始化请求
 */
export function getInit() {
    return request({
        url: '/admin/init',
        method: 'GET',
    })
}
