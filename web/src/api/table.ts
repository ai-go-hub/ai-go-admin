import request from '/@/utils/request'

/**
 * 快速生成一个控制器的 增、删、改、查、排序 接口的请求方法；
 * 本 class 实例通常直接传递给 useTableManager 使用，开发者可重写本类的方法，亦可直接向 useTableManager 传递自定义的 API 请求类；
 * 表格相关网络请求无需局限于本类，开发者可于 api/ 创建自定义的接口请求函数，并于需要的地方导入使用即可。
 */
export class TableManagerAPI {
    public controllerUrl
    public actionUrl

    constructor(controllerUrl: string) {
        this.controllerUrl = controllerUrl
        this.actionUrl = new Map([
            ['get', controllerUrl + 'get'],
            ['list', controllerUrl + 'list'],
            ['create', controllerUrl + 'create'],
            ['update', controllerUrl + 'update'],
            ['delete', controllerUrl + 'delete'],
            ['sort', controllerUrl + 'sort'],
        ])
    }

    /**
     * 获取行数据
     * @param params 行主键等
     */
    get(params: AnyObj) {
        return request({
            url: this.actionUrl.get('get'),
            method: 'GET',
            params,
        })
    }

    /**
     * 表格数据列表接口的请求方法
     * @param filter 数据过滤条件
     */
    list(filter: TableInterface['filter'] = {}) {
        return request<TableManagerAPIDefaultData>({
            url: this.actionUrl.get('list'),
            method: 'GET',
            params: filter,
        })
    }

    /**
     * 表格删除接口的请求方法
     * @param ids 被删除数据的主键数组
     */
    delete(ids: string[]) {
        return request(
            {
                url: this.actionUrl.get('delete'),
                method: 'POST',
                data: { ids },
            },
            {
                showSuccessMessage: true,
            }
        )
    }

    /**
     * 向指定接口 POST 数据，本方法虽然较为通用，但请不要局限于此，开发者可于 api/ 创建自定义的接口请求函数，并于需要的地方导入使用即可
     * @param action 请求的接口，比如 create、update
     * @param data 要 POST 的数据
     */
    post(action: string, data: AnyObj) {
        return request(
            {
                url: this.actionUrl.has(action) ? this.actionUrl.get(action) : this.controllerUrl + action,
                method: 'POST',
                data,
            },
            {
                showSuccessMessage: true,
            }
        )
    }

    /**
     * 表格行排序接口的请求方法
     */
    sort(data: AnyObj) {
        return request({
            url: this.actionUrl.get('sort'),
            method: 'POST',
            data,
        })
    }
}
