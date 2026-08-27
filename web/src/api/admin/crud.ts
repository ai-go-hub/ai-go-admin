import { adminBaseRoutePath } from '@/router/static/adminBase'
import request, { getBaseURL } from '@/utils/request'

export const url = '/admin/crud/'
export const tableListUrl = url + 'table-list'
export const modelListUrl = url + 'model-list'
export const crudLogListUrl = url + 'log/list'
export const crudAIStreamURL = getBaseURL() + adminBaseRoutePath + '/crud/ai/stream'

export function getTableFieldList(table: string) {
    return request({
        url: url + 'table-field-list',
        method: 'GET',
        params: { table },
    })
}

export function checkLog(table: string) {
    return request({
        url: url + 'check-log',
        method: 'GET',
        params: { table },
    })
}

export function getGenerateBasicData(path: string, app: string) {
    return request({
        url: url + 'generate-basic-data',
        method: 'GET',
        params: { path, app },
    })
}

export function generate(data: AnyObj) {
    return request(
        {
            url: url + 'generate',
            method: 'POST',
            data,
        },
        {
            showSuccessMessage: true,
        }
    )
}

export function checkGenerate(data: AnyObj) {
    return request(
        {
            url: url + 'check-generate',
            method: 'POST',
            data,
        },
        {
            showErrorMessage: false,
        }
    )
}

export function parseTableData(params: AnyObj) {
    return request({
        url: url + 'parse-table-data',
        method: 'GET',
        params,
    })
}

export function logStart(id: number, type: string) {
    return request({
        url: url + 'log-start',
        method: 'GET',
        params: { id, type },
    })
}

export function deleteLog(id: number) {
    return request({
        url: url + 'delete',
        method: 'POST',
        data: { id },
    })
}
