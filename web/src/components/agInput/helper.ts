import type { FieldData } from './index'

/**
 * 将逗号分隔的字符串或字符串数组格式化为字符串数组
 */
export const stringToArray = (val: string | string[]): string[] => {
    if (typeof val === 'string') {
        return val === '' ? [] : val.split(',')
    }
    return val
}

export const defaultNUPPD = {
    null: true,
    unique: false,
    primaryKey: false,
    precision: 0,
    defaultType: 'NULL' as FieldData['defaultType'],
}

/**
 * 所有 Input 组件对应的数据表字段设计数据
 * 根据组件绑定值类型、服务端接受方便、数据查询比较方便等层面设计，属于默认/示例设计，主要供可视化 CRUD 设计器使用
 */
export const fieldData: Record<string, FieldData> = {
    // ==================== varchar ====================
    text: {
        type: 'varchar',
        length: 255,
        ...defaultNUPPD,
    },
    string: {
        type: 'varchar',
        length: 255,
        ...defaultNUPPD,
    },
    password: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },
    radio: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },
    select: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },
    color: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },
    iconSelect: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },
    areaSelect: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },
    image: {
        type: 'varchar',
        length: 255,
        ...defaultNUPPD,
    },
    file: {
        type: 'varchar',
        length: 255,
        ...defaultNUPPD,
    },
    time: {
        type: 'varchar',
        length: 64,
        ...defaultNUPPD,
    },

    // ==================== bool ====================

    switch: {
        type: 'boolean',
        length: 0,
        ...defaultNUPPD,
    },

    // ==================== 数字 ====================
    int: {
        type: 'bigint',
        length: 64,
        ...defaultNUPPD,
    },
    number: {
        type: 'numeric',
        length: 10,
        ...defaultNUPPD,
        precision: 2,
    },
    remoteSelect: {
        type: 'bigint',
        length: 64,
        ...defaultNUPPD,
    },
    year: {
        // 年份选择器 el-date-picker 的绑定值为 string，但数据库层面最好使用 int 来方便查询和比较
        // 接口可直接接受 string，数据表字段类型设为 smallint 即可（pgsql 能自动转换），或者使用 tag 指定输入类型完成接受: Year *uint16 `year,string`
        type: 'smallint',
        length: 16,
        ...defaultNUPPD,
    },

    // ==================== jsonb ====================
    checkbox: {
        type: 'jsonb',
        length: 0,
        ...defaultNUPPD,
    },
    array: {
        type: 'jsonb',
        length: 0,
        ...defaultNUPPD,
    },
    selects: {
        type: 'jsonb',
        length: 0,
        ...defaultNUPPD,
    },
    remoteSelects: {
        type: 'jsonb',
        length: 0,
        ...defaultNUPPD,
    },
    images: {
        type: 'jsonb',
        length: 0,
        ...defaultNUPPD,
    },
    files: {
        type: 'jsonb',
        length: 0,
        ...defaultNUPPD,
    },

    // ==================== text ====================
    editor: {
        type: 'text',
        length: 0,
        ...defaultNUPPD,
    },
    textarea: {
        type: 'text',
        length: 0,
        ...defaultNUPPD,
    },

    // ==================== date / datetime ====================
    date: {
        type: 'date',
        length: 0,
        ...defaultNUPPD,
    },
    datetime: {
        type: 'timestamptz',
        length: 6,
        ...defaultNUPPD,
    },
}
