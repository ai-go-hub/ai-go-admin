import type { AxiosRequestConfig } from 'axios'
import type { ClickRequest } from '/@/components/clickCaptcha/index'
import { uuid } from '/@/utils/random'
import request from '/@/utils/request'

// ==================== 点选验证码 ====================

export function getClickCaptcha(apiBaseURL?: string) {
    return request(
        {
            url: '/common/captcha/create',
            method: 'GET',
            ...(apiBaseURL ? { baseURL: apiBaseURL } : {}),
        },
        {
            convertCase: true,
        }
    )
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
            convertCase: true,
            showErrorMessage: false,
        }
    )
}

// ==================== 文件上传 ====================

/**
 * 文件上传请求参数
 */
export interface UploadParams {
    // 上传文件
    file: File
    // 业务主题分类，如 avatar、article
    topic?: string
    // 上传驱动，默认 local
    driver?: string
}

/**
 * 文件上传响应
 */
export interface UploadResult {
    // 资源访问地址
    url: string
    // 文件大小（字节）
    size: number
    // 文件后缀
    suffix: string
    // 是否为图片
    isImage: boolean
}

/**
 * 文件上传请求，multipart/form-data 表单
 */
export function upload(params: UploadParams, config?: AxiosRequestConfig) {
    const form = new FormData()
    form.append('file', params.file)
    if (params.topic) form.append('topic', params.topic)
    if (params.driver) form.append('driver', params.driver)

    return request<UploadResult>({
        url: '/common/upload',
        method: 'POST',
        data: form,
        ...(config ?? {}),
    })
}

// ==================== 省份地区 ====================

/**
 * 获取省份地区数据
 */
export function getArea(values: number[] = []) {
    const params: { province?: number; city?: number; uuid: string } = { uuid: uuid() }
    if (values[0]) {
        params.province = values[0]
    }
    if (values[1]) {
        params.city = values[1]
    }
    return request({
        url: '/common/area',
        method: 'GET',
        params,
    })
}

// ==================== 远程下拉 ====================

/**
 * 请求远程下拉框数据
 */
export function selectList(listURL: string, data: AnyObj = {}) {
    return request({
        url: listURL,
        method: 'POST',
        data,
    })
}
