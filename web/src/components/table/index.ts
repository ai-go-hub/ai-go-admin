import { DefaultOptButType } from '@/components/table/types'
import i18n from '@/lang/index'
import { TableColumnCtx } from 'element-plus'
import { isUndefined } from 'lodash-es'

/**
 * 获取单元格值
 */
export const getCellValue = (row: TableRow, columnConfig: TableColumn, column: TableColumnCtx<TableRow>, index: number) => {
    if (!columnConfig.prop) {
        return ''
    }

    const prop = columnConfig.prop
    let cellValue: any = row[prop]

    // 字段 prop 带 . 比如 user.nickname
    if (prop.indexOf('.') > -1) {
        const fieldNameArr = prop.split('.')
        cellValue = row[fieldNameArr[0]]
        for (let index = 1; index < fieldNameArr.length; index++) {
            cellValue = cellValue ? (cellValue[fieldNameArr[index]] ?? '') : ''
        }
    }

    // 若无值，尝试取默认值
    if ([undefined, null, ''].includes(cellValue) && columnConfig.default !== undefined) {
        cellValue = columnConfig.default
    }

    // 渲染前格式化
    if (columnConfig.formatter && typeof columnConfig.formatter == 'function') {
        cellValue = columnConfig.formatter(row, column, cellValue, index)
    }

    return cellValue
}

/*
 * 默认按钮组
 */
export const getDefaultOptButtons = (optButType: DefaultOptButType[] = ['sort', 'edit', 'delete']): OptButton[] => {
    const optButtonsPre: Map<string, OptButton> = new Map([
        [
            // 拖动排序按钮
            'sort',
            {
                render: 'sort',
                name: 'sort',
                title: i18n.global.t('common.dragSort'),
                text: '',
                type: 'info',
                icon: 'lucide-fold-vertical',
                class: 'table-row-sort',
            },
        ],
        [
            // 编辑按钮
            'edit',
            {
                render: 'tip',
                name: 'edit',
                title: i18n.global.t('common.edit'),
                text: '',
                type: 'primary',
                icon: 'lucide-pencil',
                class: 'table-row-edit',
            },
        ],
        [
            // 删除按钮
            'delete',
            {
                render: 'confirm',
                name: 'delete',
                title: i18n.global.t('common.delete'),
                text: '',
                type: 'danger',
                icon: 'lucide-trash',
                class: 'table-row-delete',
                popconfirm: {
                    confirmButtonText: i18n.global.t('common.delete'),
                    cancelButtonText: i18n.global.t('common.cancel'),
                    confirmButtonType: 'danger',
                    title: i18n.global.t('common.deleteSelectedRecords'),
                },
            },
        ],
    ])

    const optButtons: OptButton[] = []
    for (const key in optButType) {
        if (optButtonsPre.has(optButType[key])) {
            optButtons.push(optButtonsPre.get(optButType[key])!)
        }
    }
    return optButtons
}

/**
 * 将带 children 的数组降维，然后寻找 index 所在的行
 */
export const findIndexRow = (data: TableRow[], findIdx: number, keyIndex: number | TableRow = -1): number | TableRow => {
    for (const key in data) {
        if (typeof keyIndex == 'number') {
            keyIndex++
        }

        if (keyIndex == findIdx) {
            return data[key]
        }

        if (data[key].children) {
            keyIndex = findIndexRow(data[key].children, findIdx, keyIndex)
            if (typeof keyIndex != 'number') {
                return keyIndex
            }
        }
    }

    return keyIndex
}

/**
 * 调用一个接受表格上下文数据的任意属性计算函数
 */
export const invokeTableContextDataFun = <T>(
    fun: TableContextDataFun<T> | undefined,
    context: TableContextData,
    defaultValue: any = {}
): Partial<T> => {
    if (isUndefined(fun)) {
        return defaultValue
    } else if (typeof fun === 'function') {
        return fun(context)
    }
    return fun
}
