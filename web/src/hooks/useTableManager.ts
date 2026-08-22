import { findIndexRow } from '@/components/table/index'
import i18n from '@/lang/index'
import { auth as authentication, getArrayKey } from '@/utils/common'
import { dayjs, ElNotification, FormInstance } from 'element-plus'
import { assign, cloneDeep, defaults, isArray, isEmpty } from 'lodash-es'
import Sortable from 'sortablejs'
import { reactive } from 'vue'
import { useRoute } from 'vue-router'

/**
 * 表格管家工厂函数
 */
export function useTableManager(opts: UseTableManagerOptions): TableManagerInstance {
    const table: TableInterface = reactive({
        data: [],
        column: [],
        filter: {},

        dblClickNotEditColumn: [],
        extend: {},

        ref: null,
        pk: 'id',
        loading: false,
        selections: [],
        total: 0,
        acceptQuery: true,
        showComSearch: false,
        expandAll: false,
        routePath: '',
        dragSortLimitField: 'pid',
    })

    const form: FormInterface = reactive({
        items: {},
        operate: '',
        defaultItems: {},

        extend: {},

        ref: null,
        labelWidth: 160,
        operatePKs: [],
        submitLoading: false,
        loading: false,
    })

    const comSearch: ComSearchInterface = reactive({
        form: {},
        fieldData: new Map(),
    })

    // 初始化 opts，好处例如: 若 opts.before 未定义，则无法使用 opts.before!.getData 这种语法快速新增钩子
    opts = defaults(opts, {
        before: {},
        after: {},
        table: {},
        form: {},
    })

    // 合并 opts
    assign(table, opts.table)
    assign(form, opts.form)

    /**
     * 表格内部鉴权方法
     * 此方法在表头或表行组件内部自动调用，传递权限节点名，如: create、delete
     */
    const auth = (node: string): boolean => {
        return typeof opts.auth === 'function' ? opts.auth(node) : authentication(node)
    }

    /**
     * 运行前置函数
     * @param funName 函数名
     * @param args 参数
     */
    const runBefore = (funName: string, args: any = {}) => {
        if (opts.before && opts.before[funName] && typeof opts.before[funName] == 'function') {
            return opts.before[funName]({ ...args }) === false ? false : true
        }
        return true
    }

    /**
     * 运行后置函数
     * @param funName 函数名
     * @param args 参数
     */
    const runAfter = (funName: string, args: any = {}) => {
        if (opts.after && opts.after[funName] && typeof opts.after[funName] == 'function') {
            return opts.after[funName]({ ...args }) === false ? false : true
        }
        return true
    }

    /**
     * 表格数据获取（请求表格对应控制器的 list 方法）
     */
    const getData = () => {
        if (runBefore('getData') === false) return
        table.loading = true
        return opts.api
            .list(table.filter)
            .then((res) => {
                table.data = res.data.data.list
                table.total = res.data.data.total
                runAfter('getData', { res })
            })
            .catch((err) => {
                runAfter('getData', { err })
            })
            .finally(() => {
                table.loading = false
            })
    }

    /**
     * 批量删除数据
     */
    const del = (pks: string[]) => {
        if (runBefore('delete', { pks }) === false) return
        opts.api.delete(pks).then((res) => {
            refresh({ event: 'delete', pks })
            runAfter('delete', { res })
        })
    }

    /**
     * 获取被编辑行数据
     */
    const getEditData = (pk: string) => {
        if (runBefore('getEditData', { pk }) === false) return
        form.loading = true
        form.items = {}
        return opts.api
            .get(pk)
            .then((res) => {
                form.items = res.data.data.row
                runAfter('getEditData', { res })
            })
            .catch((err) => {
                toggleForm()
                runAfter('getEditData', { err })
            })
            .finally(() => {
                form.loading = false
            })
    }

    /**
     * 切换表单
     * @param operate 操作:create=添加,update=更新
     * @param operatePKs 被操作项的数组:create=[],update=[1,2,...]
     */
    const toggleForm = (operate = '', operatePKs: string[] = []) => {
        if (runBefore('toggleForm', { operate, operatePKs }) === false) return
        if (operate == 'update') {
            if (!operatePKs.length) {
                return false
            }
            // 批量编辑支持，请求下一条被编辑数据
            getEditData(operatePKs[0])
        } else if (operate == 'create') {
            form.items = cloneDeep(form.defaultItems)
        }
        form.operate = operate
        form.operatePKs = operatePKs
        runAfter('toggleForm', { operate, operatePKs })
    }

    /**
     * 提交表单
     * @param formEl 表单组件 ref
     */
    const submitForm = (formEl?: FormInstance | null) => {
        if (runBefore('submitForm', { formEl: formEl, operate: form.operate, items: form.items! }) === false) return

        // 表单验证通过后执行的 api 请求操作
        const submitCallback = () => {
            form.submitLoading = true
            opts.api
                .post(form.operate!, form.items![table.pk!], form.items!)
                .then((res) => {
                    refresh({ event: 'submit-form', operate: form.operate, items: form.items })
                    form.operatePKs?.shift()
                    if (form.operatePKs!.length > 0) {
                        // 批量编辑支持，继续编辑下一条数据
                        toggleForm('update', form.operatePKs)
                    } else {
                        toggleForm()
                    }
                    runAfter('submitForm', { res })
                })
                .finally(() => {
                    form.submitLoading = false
                })
        }

        if (formEl) {
            form.ref = formEl
            formEl.validate((valid: boolean) => {
                if (valid) {
                    submitCallback()
                }
            })
        } else {
            submitCallback()
        }
    }

    /**
     * 获取表格选择项的主键数组
     */
    const getSelectionPKs = () => {
        const pks: string[] = []
        table.selections?.forEach((item) => {
            pks.push(item[table.pk!])
        })
        return pks
    }

    /**
     * 处理表格内部事件
     * @param event 事件名称，含义请参考其类型定义
     * @param data 携带数据
     */
    const handleEvent = (event: TableEventName, data: AnyObj) => {
        if (runBefore('tableEvent', { event, data }) === false) return

        const actionFun = new Map([
            [
                'selection-change',
                () => {
                    table.selections = data as TableRow[]
                },
            ],
            [
                'page-size-change',
                () => {
                    table.filter!.limit = data.size
                    refresh({ event: 'page-size-change', ...data })
                },
            ],
            [
                'current-page-change',
                () => {
                    table.filter!.page = data.page
                    refresh({ event: 'current-page-change', ...data })
                },
            ],
            [
                'sort-change',
                () => {
                    if (data.order === '') {
                        delete table.filter!.sort
                        delete table.filter!.order
                    } else {
                        table.filter!.sort = data.prop
                        table.filter!.order = data.order
                    }
                    refresh({ event: 'sort-change', ...data })
                },
            ],
            [
                'edit',
                () => {
                    toggleForm('update', [data.row[table.pk!]])
                },
            ],
            [
                'delete',
                () => {
                    del([data.row[table.pk!]])
                },
            ],
            [
                'cell-change',
                () => {
                    if (data.field && data.field.prop && table.data![data.index]) {
                        table.data![data.index][data.field.prop!] = data.value
                    }
                },
            ],
            [
                'com-search',
                () => {
                    // 公共搜索
                    const groups: WhereGroup[] = []
                    if (table.filter!.quickSearchKeywords) {
                        const quick = getQuickSearchData(table.filter!.quickSearchKeywords)
                        if (quick !== false) groups.push(quick)
                    }
                    const com = getComSearchData()
                    if (com !== false) groups.push(com)

                    setFilterWheres(groups)

                    // 刷新表格
                    refresh({ event: 'com-search', data: table.filter?.wheres })
                },
            ],
            [
                'column-dblclick',
                () => {
                    if (data.column.property === undefined) return
                    if (!table.dblClickNotEditColumn!.includes('all') && !table.dblClickNotEditColumn!.includes(data.column.property)) {
                        if (runBefore('columnDblclick', { row: data.row, column: data.column }) === false) return
                        toggleForm('update', [data.row[table.pk!]])
                        runAfter('columnDblclick', { row: data.row, column: data.column })
                    }
                },
            ],
            [
                'refresh',
                () => {
                    // 刷新表格在大多数情况下无需置空 data，但任需防范表格列组件的 :key 不会被更新的问题，比如关联表的数据列
                    table.data = []
                    getData()
                },
            ],
            [
                'add',
                () => {
                    toggleForm('create')
                },
            ],
            [
                'edit-selected',
                () => {
                    toggleForm('update', getSelectionPKs())
                },
            ],
            [
                'delete-selected',
                () => {
                    del(getSelectionPKs())
                },
            ],
            [
                'toggle-expansion',
                () => {
                    if (!table.ref) {
                        console.warn('useTableManager.table.ref is undefined')
                        return
                    }
                    table.expandAll = data.expanded
                    table.ref.toggleExpansionAll(data.expanded)
                },
            ],
            [
                'quick-search',
                () => {
                    table.filter!.quickSearchKeywords = data.keywords

                    const groups: WhereGroup[] = []
                    const com = getComSearchData()
                    if (com !== false) groups.push(com)
                    if (data.keywords) {
                        const quick = getQuickSearchData(data.keywords)
                        if (quick !== false) groups.push(quick)
                    }

                    setFilterWheres(groups)

                    refresh({ event: 'quick-search', ...data })
                },
            ],
            [
                'show-column-change',
                () => {
                    const columnKey = getArrayKey(table.column, 'prop', data.field)
                    table.column[columnKey].show = data.value
                },
            ],
            [
                'reinit-com-search',
                () => {
                    initComSearch()
                },
            ],
            [
                'toggle-com-search',
                () => {
                    table.showComSearch = data.value
                },
            ],
            [
                'default',
                () => {
                    console.warn('No action defined')
                },
            ],
        ])

        const action = actionFun.get(event) || actionFun.get('default')
        typeof action === 'function' && action()
        return runAfter('tableEvent', { event, data })
    }

    /**
     * 刷新表格列表数据
     * @param custom 任意自定义数据，如 { event: 'delete', pks }
     */
    const refresh = (custom: AnyObj) => {
        handleEvent('refresh', custom)
    }

    /**
     * 初始化表格拖动排序，按使用场景和需要手动调用
     * 1. 需要在 getData（获取表格数据）的回调函数中调用，即表格数据加载完成之前调用无效
     * 2. 需要在 onMounted 回调函数内调用，即在表格渲染出来之前调用无效
     */
    const initDragSort = () => {
        const buttonsKey = getArrayKey(table.column, 'render', 'buttons')
        if (buttonsKey === false) return
        const moveButton = getArrayKey(table.column[buttonsKey].buttons!, 'render', 'sort')
        if (moveButton === false) return
        if (!table.ref) {
            console.warn('useTableManager.table.ref is undefined')
            return
        }

        const el = (table.ref.getElTableRef() as any)?.$el.querySelector('.el-table__body-wrapper .el-table__body tbody')
        const disabledTip = table.column[buttonsKey].buttons![moveButton].disabledTip ? true : false
        Sortable.create(el, {
            animation: 200,
            handle: '.table-row-sort',
            ghostClass: 'ag-table-row',
            onStart: () => {
                table.column[buttonsKey].buttons![moveButton].disabledTip = true
            },
            onEnd: (evt: Sortable.SortableEvent) => {
                table.column[buttonsKey].buttons![moveButton].disabledTip = disabledTip

                const dragSortWeighField = table.dragSortWeighField ? table.dragSortWeighField : 'weigh'
                if (!table.filter?.sort || table.filter.sort != dragSortWeighField) {
                    ElNotification({ type: 'error', message: i18n.global.t('common.dragSortWeighFieldError', { field: dragSortWeighField }) })
                    return
                }

                // 目标位置不变
                if (evt.oldIndex == evt.newIndex || typeof evt.newIndex == 'undefined' || typeof evt.oldIndex == 'undefined') return

                // 找到对应行
                const moveRow = findIndexRow(table.data!, evt.oldIndex) as TableRow
                const targetRow = findIndexRow(table.data!, evt.newIndex) as TableRow
                const targetWeigh = targetRow[dragSortWeighField]

                if (targetWeigh === undefined || targetWeigh == null) {
                    ElNotification({ type: 'error', message: i18n.global.t('common.dragSortTargetWeighError') })
                    return
                }

                const eventData = {
                    move: String(moveRow[table.pk!]),
                    target: String(targetRow[table.pk!]),
                    sort: table.filter?.sort,
                    order: table.filter?.order,
                    wheres: table.filter?.wheres,
                    direction: evt.newIndex > evt.oldIndex ? 'down' : 'up',
                    weigh: targetWeigh,
                }

                if (table.dragSortLimitField && moveRow[table.dragSortLimitField] != targetRow[table.dragSortLimitField]) {
                    refresh({ event: 'sort', ...eventData })
                    ElNotification({ type: 'error', message: i18n.global.t('common.beyondMovableRange') })
                    return
                }

                opts.api.sort(eventData).finally(() => {
                    refresh({ event: 'sort', ...eventData })
                })
            },
        })
    }

    /**
     * 设置 getData 请求时的过滤条件数据
     * 整体覆盖，调用方如需保留旧条件，请先读 table.filter.wheres 后再拼
     * @param groups Where 分组数组
     */
    const setFilterWheres = (groups: WhereGroup[]) => {
        table.filter!.wheres = groups
    }

    /**
     * 公共搜索初始化
     */
    const initComSearch = () => {
        const form: AnyObj = {}
        const field = table.column

        if (field.length <= 0) return

        for (const key in field) {
            // 关闭搜索的字段
            if (field[key].operator === false) continue

            // 取默认操作符号
            if (typeof field[key].operator == 'undefined') {
                field[key].operator = 'eq'
            }

            // 公共搜索表单字段初始化
            const prop = field[key].prop
            if (prop) {
                if (field[key].operator == 'BETWEEN' || field[key].operator == 'NOT BETWEEN') {
                    // 范围查询
                    form[prop] = ''
                    form[prop + '-start'] = ''
                    form[prop + '-end'] = ''
                } else if (field[key].operator == 'NULL' || field[key].operator == 'NOT NULL') {
                    // 复选框
                    form[prop] = false
                } else {
                    // 普通文本框
                    form[prop] = ''
                }

                // 初始化字段的公共搜索数据
                comSearch.fieldData.set(prop, {
                    operator: field[key].operator,
                    render: field[key].render,
                    comSearchRender: field[key].comSearchRender,
                })
            }
        }

        comSearch.form = Object.assign(comSearch.form, form)
    }

    /**
     * 设置公共搜索表单数据
     */
    const setComSearchData = (query: AnyObj) => {
        // 必需已经完成公共搜索数据的初始化
        if (comSearch.fieldData.size === 0) {
            initComSearch()
        }

        for (const key in table.column) {
            const prop = table.column[key].prop
            if (prop && typeof query[prop] !== 'undefined') {
                const queryProp = query[prop] ?? ''
                if (table.column[key].operator == 'BETWEEN' || table.column[key].operator == 'NOT BETWEEN') {
                    const range = queryProp.split(',')
                    if (table.column[key].comSearchRender == 'datetime' || table.column[key].comSearchRender == 'date') {
                        if (range && range.length >= 2) {
                            const rangeDayJs = [dayjs(range[0]), dayjs(range[1])]
                            if (rangeDayJs[0].isValid() && rangeDayJs[1].isValid()) {
                                if (table.column[key].comSearchRender == 'date') {
                                    comSearch.form[prop] = [rangeDayJs[0].format('YYYY-MM-DD'), rangeDayJs[1].format('YYYY-MM-DD')]
                                } else {
                                    comSearch.form[prop] = [rangeDayJs[0].format('YYYY-MM-DD HH:mm:ss'), rangeDayJs[1].format('YYYY-MM-DD HH:mm:ss')]
                                }
                            }
                        }
                    } else if (table.column[key].comSearchRender == 'time') {
                        if (range && range.length >= 2) {
                            comSearch.form[prop] = [range[0], range[1]]
                        }
                    } else {
                        comSearch.form[prop + '-start'] = range[0] ?? ''
                        comSearch.form[prop + '-end'] = range[1] ?? ''
                    }
                } else if (table.column[key].operator == 'NULL' || table.column[key].operator == 'NOT NULL') {
                    comSearch.form[prop] = queryProp ? true : false
                } else if (table.column[key].comSearchRender == 'datetime' || table.column[key].comSearchRender == 'date') {
                    const propDayJs = dayjs(queryProp)
                    if (propDayJs.isValid()) {
                        comSearch.form[prop] = propDayJs.format(table.column[key].comSearchRender == 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss')
                    }
                } else {
                    comSearch.form[prop] = queryProp
                }
            }
        }
    }

    /**
     * 获取快速搜索数据
     */
    const getQuickSearchData = (keywords: string): WhereGroup | false => {
        const wheres: ComSearchData[] = []
        for (const key in table.column) {
            if (keywords && table.column[key]['quickSearch'] === true && table.column[key]['prop']) {
                wheres.push({
                    field: table.column[key]['prop'],
                    value: keywords,
                    operator: 'ILIKE',
                })
            }
        }
        if (!wheres.length) return false
        return { wheres, or: true }
    }

    /**
     * 获取公共搜索表单数据
     */
    const getComSearchData = (): WhereGroup | false => {
        // 必需已经完成公共搜索数据的初始化
        if (comSearch.fieldData.size === 0) {
            initComSearch()
        }

        const wheres: ComSearchData[] = []

        for (const key in comSearch.form) {
            if (!comSearch.fieldData.has(key)) continue

            let val = null
            const fieldDataTemp = comSearch.fieldData.get(key)
            if (
                ['datetime', 'date', 'time'].includes(fieldDataTemp.comSearchRender) &&
                (fieldDataTemp.operator == 'BETWEEN' || fieldDataTemp.operator == 'NOT BETWEEN')
            ) {
                if (comSearch.form[key] && comSearch.form[key].length >= 2) {
                    // 日期范围
                    if (fieldDataTemp.comSearchRender == 'date') {
                        val = comSearch.form[key][0] + ' 00:00:00' + ',' + comSearch.form[key][1] + ' 23:59:59'
                    } else {
                        // 时间范围、时间日期范围
                        val = comSearch.form[key][0] + ',' + comSearch.form[key][1]
                    }
                }
            } else if (fieldDataTemp.operator == 'BETWEEN' || fieldDataTemp.operator == 'NOT BETWEEN') {
                // 普通的范围筛选，公共搜索初始化时已准备好 start 和 end 字段
                if (!comSearch.form[key + '-start'] && !comSearch.form[key + '-end']) {
                    continue
                }
                val = comSearch.form[key + '-start'] + ',' + comSearch.form[key + '-end']
            } else if (comSearch.form[key]) {
                val = comSearch.form[key]
            }

            if (val === null) continue
            if (isArray(val) && !val.length) continue

            wheres.push({
                field: key,
                value: val,
                operator: fieldDataTemp.operator,
            })
        }

        if (!wheres.length) return false
        return { wheres, or: false }
    }

    /**
     * 表格管家的基础上下文初始化，按使用场景和需要手动调用
     * 1. 记录表格的路由路径
     * 2. 初始化公共搜索表单和字段数据
     * 3. 从 URL 读取公共搜索默认值
     */
    const initCtx = () => {
        if (runBefore('initCtx') === false) return

        // 记录表格的路由路径
        const route = useRoute()
        table.routePath = route.fullPath

        // 按需初始化公共搜索表单数据和字段 Map
        if (comSearch.fieldData.size === 0) {
            initComSearch()
        }

        // 公共搜索默认值
        if (table.acceptQuery && !isEmpty(route.query)) {
            // 根据当前 URL 的 query 初始化公共搜索默认值
            setComSearchData(route.query)

            // 获取公共搜索和快速搜索数据合并至表格筛选条件
            const groups: WhereGroup[] = []
            const com = getComSearchData()
            if (com !== false) groups.push(com)
            if (table.filter!.quickSearchKeywords) {
                const quick = getQuickSearchData(table.filter!.quickSearchKeywords)
                if (quick !== false) groups.push(quick)
            }

            setFilterWheres(groups)
        }
    }

    return {
        opts,
        form,
        table,
        comSearch,
        api: opts.api,
        del,
        auth,
        initCtx,
        getData,
        refresh,
        runAfter,
        runBefore,
        toggleForm,
        submitForm,
        getEditData,
        handleEvent,
        initDragSort,
        getSelectionPKs,
        setFilterWheres,
        setComSearchData,
        getComSearchData,
        getQuickSearchData,
    }
}
