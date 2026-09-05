import { auth } from '@/utils/common'
import { useEventListener } from '@vueuse/core'
import type { App } from 'vue'

export function registerDirectives(app: App) {
    // 鉴权指令
    authDirective(app)
    // 点击后自动失焦指令
    blurDirective(app)
}

/**
 * 鉴权指令
 * @description v-auth="'name'"，name 为权限节点，如: create、update、delete
 */
function authDirective(app: App) {
    app.directive('auth', {
        mounted(el, binding) {
            if (!binding.value) return false
            if (!auth(binding.value)) el.parentNode.removeChild(el)
        },
    })
}

/**
 * 点击后自动失焦指令
 * @description v-blur
 */
function blurDirective(app: App) {
    app.directive('blur', {
        mounted(el) {
            useEventListener(el, 'focus', () => el.blur())
        },
    })
}
