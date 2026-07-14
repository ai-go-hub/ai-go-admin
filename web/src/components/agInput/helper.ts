/**
 * 将逗号分隔的字符串或字符串数组格式化为字符串数组
 */
export const stringToArray = (val: string | string[]): string[] => {
    if (typeof val === 'string') {
        return val === '' ? [] : val.split(',')
    }
    return val
}
