import type { ClickRequest } from '/@/components/clickCaptcha/index'
import request from '/@/utils/request'

export function getClickCaptcha(apiBaseURL?: string) {
    return request({
        url: '/common/captcha/create',
        method: 'GET',
        ...(apiBaseURL ? { baseURL: apiBaseURL } : {}),
    })
}

export function checkClickCaptcha(data: ClickRequest, apiBaseURL?: string) {
    return request(
        {
            url: '/common/captcha/verify',
            method: 'POST',
            data,
            ...(apiBaseURL ? { baseURL: apiBaseURL } : {}),
        },
        {
            showErrorMessage: false,
        }
    )
}

/**
 * 缓存清理接口
 */
export function clearCache(type: string) {
    return request(
        {
            url: '/admin/clear-cache',
            method: 'POST',
            data: { type },
        },
        {
            showSuccessMessage: true,
        }
    )
}
