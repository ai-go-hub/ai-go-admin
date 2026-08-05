import request from '/@/utils/request'

/**
 * 获取当前管理员最新个人信息
 * @param id 管理员 ID
 */
export function getAdminProfile(id: number) {
    return request({
        url: `/admin/routine/profile/get/${id}`,
        method: 'GET',
    })
}

/**
 * 更新当前管理员个人信息
 * @param id 管理员 ID
 * @param data 更新数据
 */
export function updateAdminProfile(id: number, data: AnyObj) {
    return request(
        {
            url: `/admin/routine/profile/update/${id}`,
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
        url: '/admin/auth/log/list',
        method: 'POST',
        data,
    })
}

/**
 * 发送测试邮件
 */
export function sendTestMail(data: AnyObj, mail: string) {
    return request(
        {
            url: '/admin/routine/config/send-test-mail',
            method: 'post',
            data: { ...data, test_mail: mail },
            timeout: 0,
        },
        {
            showSuccessMessage: true,
        }
    )
}
