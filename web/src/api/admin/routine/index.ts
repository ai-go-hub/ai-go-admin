import request from '/@/utils/request'

/**
 * 获取当前管理员最新个人信息
 * @param id 管理员 ID
 */
export function getAdminInfo(id: number) {
    return request({
        url: `/admin/auth/admin/get/${id}`,
        method: 'GET',
    })
}

/**
 * 更新当前管理员个人信息
 * @param id 管理员 ID
 * @param data 更新数据
 */
export function updateAdminInfo(id: number, data: AnyObj) {
    return request(
        {
            url: `/admin/auth/admin/update/${id}`,
            method: 'POST',
            data,
        },
        {
            showSuccessMessage: true,
        }
    )
}

/**
 * 获取管理员操作日志
 * @param data 筛选条件
 */
export function getAdminLog(data: AnyObj) {
    return request({
        url: '/admin/auth/admin-log/list',
        method: 'POST',
        data,
    })
}
