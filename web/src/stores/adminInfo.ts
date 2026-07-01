import { difference, pick } from 'lodash-es'
import { defineStore } from 'pinia'
import { ADMIN_INFO } from '/@/stores/constant/cacheKey'
import type { AdminInfo } from '/@/stores/interface'

export const useAdminInfo = defineStore('adminInfo', {
    state: (): AdminInfo => {
        return {
            id: 0,
            username: '',
            nickname: '',
            avatar: '',
            last_login_at: '',
            last_login_ip: '',
            token: '',
            super: false,
        }
    },
    actions: {
        /**
         * 状态批量填充
         * @param state 新状态数据
         * @param [exclude=true] 是否排除某些字段（忽略填充），默认值 true 排除 token 和 super 字段，传递 false 则不排除，还可传递 string[] 指定排除字段列表
         */
        dataFill(state: Record<string, any>, exclude: boolean | string[] = true) {
            if (exclude === true) {
                exclude = ['token', 'super']
            } else if (exclude === false) {
                exclude = []
            }

            // 白名单过滤：只保留 $state 中已定义的字段，并排除指定字段
            const allowedKeys = difference(Object.keys(this.$state), exclude)
            this.$patch(pick(state, allowedKeys))
        },
        setToken(token: string) {
            this.token = token
        },
        removeToken() {
            this.token = ''
        },
    },
    persist: {
        key: ADMIN_INFO,
    },
})
