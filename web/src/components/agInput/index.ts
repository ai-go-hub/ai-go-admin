/**
 * 支持的输入框类型
 */
export const inputTypes = [
    'text',
    'string',
    'password',
    'number',
    'radio',
    'checkbox',
    'switch',
    'textarea',
    'array',
    'datetime',
    'year',
    'date',
    'time',
    'select',
    'selects',
    'areaSelect',
    'iconSelect',
    'remoteSelect',
    'remoteSelects',
    'editor',
    'image',
    'images',
    'file',
    'files',
    'color',
]

/**
 * 输入框的绑定值类型
 * 组件输出值的类型，可能支持输出多种类型，此处列出的是推荐类型
 */
export const inputModelValueTypes = {
    string: [
        'text',
        'string',
        'password',
        'radio',
        'textarea',
        'datetime',
        'date',
        'time',
        'year',
        'select',
        'areaSelect',
        'iconSelect',
        'editor',
        'image',
        'images',
        'file',
        'files',
        'color',
    ],
    number: ['number', 'switch', 'remoteSelect'],
    array: ['checkbox', 'array', 'selects', 'remoteSelects'],
}

/**
 * 绑定值类型
 */
export type ModelValueTypes = string | number | boolean | object

/**
 * 输入框属性
 * 列出常用属性提供类型提示，同时通过 [key: string]: any 支持任意属性透传
 */
export interface InputAttr {
    // ===== 通用属性 =====
    placeholder?: string
    clearable?: boolean
    size?: 'large' | 'default' | 'small'
    disabled?: boolean
    readonly?: boolean
    class?: string
    type?: string
    multiple?: boolean
    dict?: any
    border?: boolean

    // ===== 事件 =====
    onChange?: (...args: any[]) => void
    onBlur?: (...args: any[]) => void
    onFocus?: (...args: any[]) => void

    // ===== el-input / el-input-number / textarea / password =====
    rows?: number
    min?: string | number
    max?: string | number
    step?: string | number
    maxlength?: string | number
    minlength?: string | number
    showPassword?: boolean

    // ===== el-date-picker =====
    format?: string
    valueFormat?: string
    startPlaceholder?: string
    endPlaceholder?: string
    rangeSeparator?: string

    // ===== remoteSelect =====
    pk?: string
    field?: string
    remoteUrl?: string
    remoteParams?: AnyObj
    remoteSearchFields?: string[]

    // ===== areaSelect =====
    level?: number

    // ===== radio / checkbox =====
    button?: boolean
    subAttr?: AnyObj

    // ===== upload =====
    topic?: string
    driver?: string
    accept?: string
    limit?: number

    // ===== array =====
    keyTitle?: string
    valueTitle?: string

    // ===== editor =====
    height?: string

    // ===== 允许任意额外属性透传以兜底 =====
    [key: string]: any
}
