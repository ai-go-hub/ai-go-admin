import { TableColumnCtx } from 'element-plus'
import { fullURL } from '/@/utils/common'

const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp', 'svg', 'ico']

/**
 * 表格和表单中文件预览图的生成
 */
export const previewRenderFormatter = (row: TableRow, _column: TableColumnCtx<TableRow> | null, cellValue: any) => {
    const ext = cellValue?.split('.').pop()
    if (ext && imageExtensions.includes(ext)) {
        return fullURL(row.url)
    }
    return fullURL(`/common/file-svg?suffix=${ext}`)
}

/**
 * 格式化文件大小
 */
export const formatFileSize = (bytes: number): string => {
    if (!bytes || bytes === 0) return '0 B'
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let i = 0
    let size = bytes
    while (size >= 1024 && i < units.length - 1) {
        size /= 1024
        i++
    }
    return size.toFixed(1) + ' ' + units[i]
}
