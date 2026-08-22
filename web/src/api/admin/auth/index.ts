import request from '@/utils/request'

export function getAdminRules() {
    return request({
        url: '/admin/auth/rule/list',
        method: 'POST',
        data: {
            order: 'desc',
            sort: 'weigh',
        },
    })
}
