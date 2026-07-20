import type { RuleType } from 'async-validator'
import type { FormItemRule } from 'element-plus'
import i18n from '../lang'

/**
 * 手机号码验证
 */
export function validatorMobile(rule: any, mobile: string | number, callback: Function) {
    // 允许空值，若需必填请额外添加必填验证规则
    if (!mobile) {
        return callback()
    }
    if (!/^(1[3-9])\d{9}$/.test(mobile.toString())) {
        return callback(new Error(i18n.global.t('common.invalidEntry', { field: i18n.global.t('common.mobile') })))
    }
    return callback()
}

/**
 * 账户名验证
 */
export function validatorAccount(rule: any, val: string, callback: Function) {
    if (!val) {
        return callback()
    }
    if (!/^[a-zA-Z][a-zA-Z0-9_]{2,15}$/.test(val)) {
        return callback(new Error(i18n.global.t('common.accountFormatError')))
    }
    return callback()
}

/**
 * 密码验证
 */
export function regularPassword(val: string) {
    return /^(?!.*[&<>"'\n\r]).{6,32}$/.test(val)
}
export function validatorPassword(rule: any, val: string, callback: Function) {
    if (!val) {
        return callback()
    }
    if (!regularPassword(val)) {
        return callback(new Error(i18n.global.t('common.passwordFormatError')))
    }
    return callback()
}

/**
 * 变量名验证
 */
export function regularVarName(val: string) {
    return /^([^\x00-\xff]|[a-zA-Z_$])([^\x00-\xff]|[a-zA-Z0-9_$])*$/.test(val)
}
export function validatorVarName(rule: any, val: string, callback: Function) {
    if (!val) {
        return callback()
    }
    if (!regularVarName(val)) {
        return callback(new Error(i18n.global.t('common.varNameFormatError')))
    }
    return callback()
}

/**
 * 富文本必填
 */
export function validatorRichTextRequired(rule: any, val: string, callback: Function) {
    if (!val || val == '<p><br></p>') {
        return callback(new Error(i18n.global.t('common.richTextRequired')))
    }
    return callback()
}

/**
 * 支持的表单验证规则
 */
export const validatorType = {
    required: i18n.global.t('common.required'),
    mobile: i18n.global.t('common.mobile'),
    account: i18n.global.t('common.account'),
    password: i18n.global.t('common.password'),
    varName: i18n.global.t('common.varName'),
    richTextRequired: i18n.global.t('common.richTextRequired'),
    url: 'URL',
    email: i18n.global.t('common.email'),
    date: i18n.global.t('common.date'),
    number: i18n.global.t('common.number'), // 数字（包括浮点和整数）
    integer: i18n.global.t('common.integer'), // 整数（不包括浮点数）
    float: i18n.global.t('common.float'), // 浮点数（不包括整数）
}

export interface BuildValidatorParams {
    // 规则名:required=必填,mobile=手机号,account=账户,password=密码,varName=变量名,richTextRequired=富文本必填,number、integer、float、date、url、email
    name: 'required' | 'mobile' | 'account' | 'password' | 'varName' | 'richTextRequired' | 'number' | 'integer' | 'float' | 'date' | 'url' | 'email'
    // 自定义验证错误消息
    message?: string
    // 验证项的标题，以下验证方式不支持此字段:mobile、account、password、varName、richTextRequired
    title?: string
    // 验证触发方式
    trigger?: 'change' | 'blur'
}

/**
 * 构建表单验证规则
 * @param {BuildValidatorParams} paramsObj 参数对象
 */
export function buildValidatorRule({ name, message, title, trigger = 'blur' }: BuildValidatorParams): FormItemRule {
    // 必填
    if (name == 'required') {
        return {
            required: true,
            message: message ? message : i18n.global.t('common.pleaseEnter', { field: title }),
            trigger: trigger,
        }
    }

    // 常见类型
    const validatorType = ['number', 'integer', 'float', 'date', 'url', 'email']
    if (validatorType.includes(name)) {
        return {
            type: name as RuleType,
            message: message ? message : i18n.global.t('common.invalidEntry', { field: title }),
            trigger: trigger,
        }
    }

    // 自定义验证方法
    const validatorCustomFun: AnyObj = {
        mobile: validatorMobile,
        account: validatorAccount,
        password: validatorPassword,
        varName: validatorVarName,
        richTextRequired: validatorRichTextRequired,
    }
    if (validatorCustomFun[name]) {
        return {
            required: name == 'richTextRequired' ? true : false,
            validator: validatorCustomFun[name],
            trigger: trigger,
            message: message,
        }
    }
    return {}
}
