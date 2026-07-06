import type { AxiosInstance, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'
import axios from 'axios'
import { ElLoading, ElMessage } from 'element-plus'
import i18n from '/@/lang'
import { useAdminInfo } from '/@/stores/adminInfo'
import { keysToCamelCase, keysToSnakeCase } from '/@/utils/common'

// ==================== 类型定义 ====================

export interface RequestOptions {
    // 是否显示全屏 loading，默认 false
    loading?: boolean
    // 是否自动取消重复请求，默认 true
    cancelDuplicate?: boolean
    // 是否显示操作成功提示，默认 false
    showSuccessMessage?: boolean
    // 是否显示业务错误提示（code !== 0），默认 true
    showErrorMessage?: boolean
    // 是否显示网络错误提示（HTTP 4xx/5xx 等），默认 true
    showNetworkErrorMessage?: boolean
    // 是否做数据 key 的命名方式转换（请求 camelCase → snake_case，响应 snake_case → camelCase），默认 true
    convertCase?: boolean
}

interface InternalRequestConfig extends InternalAxiosRequestConfig {
    requestOptions?: RequestOptions
}

// ==================== Base URL 辅助函数 ====================

export function getBaseUrl(): string {
    return import.meta.env.VITE_AXIOS_BASE_URL as string
}

export function getBaseUrlPort(): string {
    return new URL(getBaseUrl()).port
}

// ==================== Loading ====================

const loadingState = {
    count: 0,
    instance: null as ReturnType<typeof ElLoading.service> | null,
}

function showLoading() {
    if (loadingState.count === 0) {
        loadingState.instance = ElLoading.service({ fullscreen: true })
    }
    loadingState.count++
}

function hideLoading() {
    loadingState.count = Math.max(0, loadingState.count - 1)
    if (loadingState.count === 0) {
        loadingState.instance?.close()
        loadingState.instance = null
    }
}

// ==================== 重复请求取消 ====================

interface PendingEntry {
    controller: AbortController
    hasLoading: boolean
}

const pendingMap = new Map<string, PendingEntry>()

function sortedStringify(obj: Record<string, any>): string {
    return JSON.stringify(Object.fromEntries(Object.entries(obj).sort(([a], [b]) => a.localeCompare(b))))
}

/**
 * 根据请求参数，为请求生成唯一标识
 */
function buildRequestKey(config: InternalAxiosRequestConfig): string {
    const url = config.url ?? ''
    const method = (config.method ?? 'get').toLowerCase()

    let params = ''
    if (config.params != null) {
        if (config.params instanceof URLSearchParams) {
            const copy = new URLSearchParams(config.params)
            copy.sort()
            params = copy.toString()
        } else {
            params = sortedStringify(config.params as Record<string, any>)
        }
    }

    let data = ''
    if (config.data != null) {
        if (typeof config.data === 'string') {
            data = config.data
        } else if (config.data instanceof FormData || config.data instanceof Blob || config.data instanceof ArrayBuffer) {
            // 无法稳定序列化，用类型名标记；如需区分多个并发上传请将 cancelDuplicate 设为 false
            data = Object.prototype.toString.call(config.data)
        } else {
            data = sortedStringify(config.data as Record<string, any>)
        }
    }

    return JSON.stringify({ method, url, params, data })
}

function addPending(config: InternalRequestConfig): void {
    const key = buildRequestKey(config)
    if (pendingMap.has(key)) {
        const { controller, hasLoading } = pendingMap.get(key)!
        controller.abort()
        if (hasLoading) hideLoading()
        pendingMap.delete(key)
        console.warn('[Request] The repeated request has been canceled:', key)
    }
    const controller = new AbortController()
    config.signal = controller.signal
    pendingMap.set(key, {
        controller,
        hasLoading: config.requestOptions?.loading ?? false,
    })
}

function removePending(config: InternalAxiosRequestConfig): void {
    pendingMap.delete(buildRequestKey(config))
}

// ==================== Axios 实例 ====================

const instance: AxiosInstance = axios.create({
    baseURL: getBaseUrl(),
    timeout: 10000,
})

instance.interceptors.request.use(
    (config: InternalRequestConfig) => {
        const opts = config.requestOptions ?? {}

        if (opts.cancelDuplicate !== false) addPending(config)
        if (opts.loading) showLoading()

        const adminInfo = useAdminInfo()
        if (adminInfo.token) {
            config.headers.set('Authorization', `Bearer ${adminInfo.token}`)
        }

        // 将请求数据的 key 从 camelCase 转为 snake_case
        // 不转换 FormData 和 URLSearchParams，因为它们本质上是类数组或二进制容器，转换可能丢失数据
        if (opts.convertCase !== false) {
            if (
                config.data &&
                typeof config.data === 'object' &&
                !(
                    config.data instanceof FormData ||
                    config.data instanceof Blob ||
                    config.data instanceof ArrayBuffer ||
                    config.data instanceof URLSearchParams
                )
            ) {
                config.data = keysToSnakeCase(config.data)
            }
            if (config.params && typeof config.params === 'object' && !(config.params instanceof URLSearchParams)) {
                config.params = keysToSnakeCase(config.params)
            }
        }

        return config
    },
    (error) => Promise.reject(error)
)

instance.interceptors.response.use(
    (response) => {
        const config = response.config as InternalRequestConfig
        const opts = config.requestOptions ?? {}

        removePending(response.config)
        if (opts.loading) hideLoading()

        if (response.data.code !== 0) {
            if (opts.showErrorMessage !== false) {
                ElMessage.error(response.data.message || i18n.global.t('common.operationFailed'))
            }
            return Promise.reject(response.data.message || i18n.global.t('common.operationFailed'))
        }

        if (opts.showSuccessMessage) {
            ElMessage.success(response.data.message || i18n.global.t('common.operationSuccess'))
        }

        // 将响应数据的 key 从 snake_case 转为 camelCase
        if (opts.convertCase !== false && response.data?.data) {
            response.data.data = keysToCamelCase(response.data.data)
        }

        return response
    },
    (error) => {
        if (axios.isCancel(error)) return Promise.reject(error)

        const config = error.config as InternalRequestConfig | undefined
        const opts = config?.requestOptions ?? {}

        if (config) {
            removePending(error.config)
            if (opts.loading) hideLoading()
        }

        if (opts.showNetworkErrorMessage !== false) {
            const msg = (error.response?.data as ApiResponse | undefined)?.message ?? error.message ?? i18n.global.t('common.networkError')
            ElMessage.error(msg)
        }

        return Promise.reject(error)
    }
)

// ==================== 对外 API ====================

function request<T = any>(config: AxiosRequestConfig, options?: RequestOptions) {
    const merged = { ...config } as InternalRequestConfig
    if (options) merged.requestOptions = options
    return instance<ApiResponse<T>>(merged)
}

export default request
