import { getCurrentInstance } from 'vue'

/**
 * 从 appContext 获取 globalProperties
 */
export const useGlobalProperties = () => {
    if (!getCurrentInstance()) {
        throw new Error('useGlobalProperties() can only be used inside setup() or functional components!')
    }
    const instance = getCurrentInstance()
    if (instance) {
        const { appContext } = instance
        return appContext.config.globalProperties
    } else {
        return null
    }
}
