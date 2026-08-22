import { createVNode, render } from 'vue'
import ClickCaptchaConstructor from './index.vue'

/**
 * 点击坐标
 */
export interface ClickPoint {
    x: number
    y: number
    element: string
}

/**
 * 验证请求数据
 */
export interface ClickRequest {
    key: string
    clicks: ClickPoint[]
    // 渲染图片宽度
    width: number
    // 渲染图片高度
    height: number
}

/**
 * 验证成功回调
 */
export type ClickCaptchaCallback = (data: ClickRequest) => void

interface ClickCaptchaOptions {
    class?: string
    error?: string
    success?: string
    // 自定义 API 基础 URL，默认使用 VITE_AXIOS_BASE_URL
    apiBaseURL?: string
}

/**
 * 弹出点选验证码
 * @param callback 验证成功的回调
 */
const clickCaptcha = (callback?: ClickCaptchaCallback, options: ClickCaptchaOptions = {}) => {
    const container = document.createElement('div')
    const vnode = createVNode(ClickCaptchaConstructor, {
        callback,
        ...options,
        onDestroy: () => {
            render(null, container)
        },
    })
    render(vnode, container)
    document.body.appendChild(container.firstElementChild!)
}

export interface Props extends ClickCaptchaOptions {
    callback?: ClickCaptchaCallback
}

export default clickCaptcha
