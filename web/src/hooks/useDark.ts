import { useConfig } from '@/stores/config'
import { useToggle, useDark as useVueUseDark } from '@vueuse/core'

/**
 * 更新 html 元素的 dark class
 */
export const updateHtmlDarkClass = (val: boolean) => {
    document.documentElement.classList.toggle('dark', val)
}

/**
 * 暗黑模式相关操作
 */
export function useDark() {
    const config = useConfig()

    const isDark = useVueUseDark({
        onChanged(dark: boolean) {
            updateHtmlDarkClass(dark)
            config.setLayoutValue('dark', dark)
        },
    })

    /**
     * 切换暗黑模式
     */
    const toggleDark = useToggle(isDark)

    return {
        isDark,
        toggleDark,
    }
}
