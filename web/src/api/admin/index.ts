import request from '/@/utils/request'

/**
 * 管理员登录请求
 */
export function adminLogin(data: { username: string; password: string; remember: boolean }) {
    return request({
        url: '/admin/login',
        method: 'POST',
        data: data,
    })
}
