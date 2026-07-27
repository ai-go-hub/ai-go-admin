import type {
    ButtonProps,
    ButtonType,
    ColProps,
    ElTooltipProps,
    FormInstance,
    ImageProps,
    LinkProps,
    PopconfirmProps,
    SwitchProps,
    TableColumnCtx,
    TagProps,
} from 'element-plus'
import type { Component } from 'vue'
import type { TableManagerAPI } from '/@/api/table'
import Icon from '/@/components/icon/index.vue'
import Table from '/@/components/table/index.vue'

declare global {
    /**
     * 表格数据行
     */
    interface TableRow extends AnyObj {
        children?: TableRow[]
    }

    /**
     * 表格列
     */
    interface TableColumn extends Partial<TableColumnCtx<TableRow>> {
        // 是否于表格显示此列
        show?: boolean
        // 渲染器组件名，即 \src\components\table\cellRenderer\ 中的组件之一，也可以查看 TableCellRenderer 类型定义获取渲染器列表
        render?: TableCellRenderer
        // 字典数据（值替换数据），同时用于单元格渲染和公共搜索下拉框数据，格式如: { open: '开', close: '关', disable: '已禁用' }
        dict?: Record<string, any>

        // render=slot 时，slot 的名称
        slotName?: string
        // render=customRender 时，要渲染的组件或已注册组件名称的字符串
        customRender?: string | Component
        // render=customTemplate 时，自定义渲染 html，应谨慎使用: 请返回 html 内容，务必确保返回内容是 xss 安全的
        customTemplate?: (row: TableRow, columnConfig: TableColumn, column: TableColumnCtx<TableRow>, cellValue: any, index: number) => string
        // 渲染前对字段值的预处理函数（对 el-table 的 formatter 扩展）
        formatter?: (row: TableRow, column: TableColumnCtx<TableRow>, cellValue: any, index: number) => any

        /**
         * 自定义单元格渲染属性（比如单元格渲染器内部的 tag、button 组件的属性，设计上不仅是组件属性，也可以自定义其他渲染相关属性）
         * 直接定义对应组件的属性 object，或使用一个函数返回组件属性 object
         */
        customRenderAttr?: {
            tag?: TableContextDataFun<TagProps>
            icon?: TableContextDataFun<InstanceType<typeof Icon>['$props']>
            image?: TableContextDataFun<ImageProps>
            switch?: TableContextDataFun<SwitchProps>
            tooltip?: TableContextDataFun<ElTooltipProps>
            link?: TableContextDataFun<LinkProps>
            [key: string]: any
        }

        // render=buttons 时，按钮数据数组
        buttons?: OptButton[]

        /**
         * 单元格渲染器需要的其他任意自定义数据
         * 1. render=tag 时，可单独指定每个不同的值 tag 的 type 属性 { open: 'success', close: 'info', disable: 'danger' }
         * 2. render=datetime 时，可指定时间日期的格式化模板（dayjs().format 模板，如 { format: 'YYYY-MM-DD HH:mm:ss' }）
         */
        custom?: {
            format?: string
            [key: string]: any
        }

        // 默认值（单元格值为 undefined,null,'' 时取默认值，仅使用了 render 时有效）
        default?: any
        // 作为快速搜索字段之一
        quickSearch?: boolean
        // 是否允许动态控制字段是否显示，默认为 true
        columnDisplayControl?: boolean
        // 单元格渲染组件的 key，默认将根据列配置等属性自动生成（此 key 值改变时单元格将自动重新渲染）
        getRenderKey?: (row: TableRow, columnConfig: TableColumn, column: TableColumnCtx<TableRow>, index: number) => string

        // 操作符（一般用于公共搜索），默认值为 = ，值为 false 禁用此字段公共搜索，支持的操作符见下类型定义
        operator?: boolean | OperatorStr
        // 公共搜索框的 placeholder
        comSearchPlaceholder?: string | string[]
        // 公共搜索渲染方式，render=tag|switch 时公共搜索也会渲染为下拉，数字会渲染为范围筛选等
        comSearchRender?: 'string' | 'remoteSelect' | 'select' | 'time' | 'date' | 'datetime' | 'customRender' | 'slot'
        // 公共搜索自定义渲染为 slot 时，slot 的名称
        comSearchSlotName?: string
        // 公共搜索自定义组件/函数渲染
        comSearchCustomRender?: string | Component
        // 公共搜索自定义渲染时，外层 el-col 的属性（仅 customRender、slot 支持）
        comSearchColAttr?: Partial<ColProps>
        // 公共搜索是否显示字段的 label
        comSearchShowLabel?: boolean
        // 公共搜索输入组件的扩展属性
        comSearchInputAttr?: AnyObj
        // 公共搜索渲染为远程下拉时，远程下拉组件的必要属性
        comSearchRemote?: {
            pk?: string
            field?: string
            params?: AnyObj
            multiple?: boolean
            remoteURL: string
        }
    }

    /**
     * 表格右侧操作按钮
     */
    interface OptButton {
        /**
         * 渲染方式:basic=普通按钮,tip=带提示的按钮,confirm=带确认框的按钮,sort=拖动排序按钮
         */
        render: 'basic' | 'tip' | 'confirm' | 'sort'

        /**
         * 按钮名称，同时将作为触发表格内事件（onTableAction）时的事件名
         */
        name: string

        /**
         * 鼠标 hover 时的提示
         */
        title?: string

        /**
         * 直接在按钮内显示的文字，可为空
         */
        text?: string

        /**
         * 自定义按钮的点击事件
         * @param row 当前行数据
         * @param column 当前列数据
         */
        click?: (row: TableRow, column: TableColumn) => void

        /**
         * 按钮是否显示（请返回布尔值，比如: display: auth('create')）
         * @param row 当前行数据
         * @param column 当前列数据
         */
        display?: (row: TableRow, column: TableColumn) => boolean

        /**
         * 按钮是否禁用（请返回布尔值）
         * @param row 当前行数据
         * @param column 当前列数据
         */
        disabled?: (row: TableRow, column: TableColumn) => boolean

        /**
         * 按钮是否正在加载中（请返回布尔值）
         * @param row 当前行数据
         * @param column 当前列数据
         */
        loading?: (row: TableRow, column: TableColumn) => boolean

        /**
         * 自定义 el-button 的其他属性（格式为属性 object 或一个返回属性 object 的函数）
         */
        attr?: TableContextDataFun<ButtonProps>

        // 按钮 class
        class?: string
        // 按钮 type
        type: ButtonType
        // 按钮 icon 的名称
        icon: string
        // 确认按钮的气泡确认框的属性（el-popconfirm 的属性，格式为属性 object 或一个返回属性 object 的函数）
        popconfirm?: TableContextDataFun<PopconfirmProps>
        // 是否禁用 title 提示，此值通常由系统动态调整以确保提示的显示效果
        disabledTip?: boolean
    }

    /**
     * 公共搜索操作符支持的值
     */
    type OperatorStr =
        | 'eq' // 等于，默认值
        | 'ne' // 不等于
        | 'gt' // 大于
        | 'egt' // 大于等于
        | 'lt' // 小于
        | 'elt' // 小于等于
        | 'LIKE'
        | 'NOT LIKE'
        | 'ILIKE'
        | 'NOT ILIKE'
        | 'IN'
        | 'NOT IN'
        | 'BETWEEN' // 范围，将生成两个输入框，可以输入最小值和最大值
        | 'NOT BETWEEN'
        | 'NULL' // 是否为NULL，将生成单个复选框
        | 'NOT NULL'

    /**
     * 公共搜索事件返回的 Data
     * 同时也是发送给服务端的单条 Where 条件类型定义
     */
    interface ComSearchData {
        field: string
        value: string | string[] | number | number[]
        operator: string
    }

    /**
     * Where 查询条件分组
     * 服务端对 `字段是否存在、操作符合是否合法` 进行检查，并对 `值` 使用 `预处理语句参数占位符` 拼接
     */
    interface WhereGroup {
        wheres: ComSearchData[]
        or?: boolean // 组内条件是否使用 OR 连接，值为 false 则使用 AND 连接条件
    }

    interface TableManagerListAPIDefaultData<T = any> {
        list: T
        total: number
    }

    /**
     * useTableManager.table
     */
    interface TableInterface {
        /**
         * 表格数据，通过 useTableManager.getData 获取
         * 刷新数据可使用 useTableManager.refresh({ event: 'custom' })
         */
        data?: TableRow[]

        /**
         * 表格列定义
         */
        column: TableColumn[]

        /**
         * 获取表格数据时的过滤条件（含公共搜索、快速搜索、分页、排序等数据）
         * 公共搜索数据可使用 useTableManager.setComSearchData 和 useTableManager.getComSearchData 进行管理
         */
        filter?: {
            page?: number
            limit?: number
            sort?: string
            order?: 'desc' | 'asc'
            wheres?: WhereGroup[]
            quickSearchKeywords?: string
            [key: string]: any
        }

        /**
         * 不需要双击编辑的列；
         * 禁用全部列的双击编辑，可使用 ['all']；
         * type=selection 的列为 undefined，将自动禁用
         */
        dblClickNotEditColumn?: string[]

        /**
         * 表格扩展数据，随意定义，以便一些自定义数据可以随 useTableManager 实例传递
         */
        extend?: AnyObj

        // 表格 ref，通常在页面 onMounted 时赋值，可选的
        ref?: InstanceType<typeof Table> | null
        // 表格对应数据表的主键字段，默认 id
        pk?: string
        // 表格加载状态
        loading?: boolean
        // 当前选中行
        selections?: TableRow[]
        // 数据总量
        total?: number
        // 接受 url 的 query 参数并自动触发公共搜索
        acceptQuery?: boolean
        // 显示公共搜索
        showComSearch?: boolean
        // 是否展开所有子项，树状表格专用属性
        expandAll?: boolean
        // 当前表格所在页面的路由 path
        routePath?: string
        // 拖动排序限位字段，例如拖动行 pid=1，那么拖动目的行 pid 也需要为 1
        dragSortLimitField?: string
        // 拖动排序权重字段，进行拖拽排序时，必需先以此字段排序，系统将修改此字段的值来完成新顺序落盘，留空则取 `weigh`
        dragSortWeighField?: string
    }

    /**
     * useTableManager.form
     */
    interface FormInterface {
        /**
         * 当前表单项数据
         */
        items?: AnyObj

        /**
         * 当前表单操作标识:create=添加,update=更新
         */
        operate?: string

        /**
         * 添加表单字段默认值，打开表单时会使用 cloneDeep 赋值给 useTableManager.form.items 对象
         */
        defaultItems?: AnyObj

        /**
         * 表单扩展数据，可随意定义，以便一些自定义数据可以随 useTableManager 实例传递
         */
        extend?: AnyObj

        // 表单 ref，实例化表格时通常无需传递
        ref?: FormInstance | null
        // 表单项 label 的宽度
        labelWidth?: number
        // 被操作数据主键，支持批量更新:create=[],update=[1,2,n]
        operatePKs?: string[]
        // 提交按钮状态
        submitLoading?: boolean
        // 表单数据的加载状态
        loading?: boolean
    }

    /**
     * 表格公共搜索数据
     */
    interface ComSearchInterface {
        // 表单项数据
        form: AnyObj
        // 字段搜索配置，搜索操作符（operator）、字段渲染方式（render）等
        fieldData: Map<string, any>
    }

    /**
     * 创建表格管家的选项
     */
    interface UseTableManagerOptions {
        api: TableManagerAPI
        table?: TableInterface
        form?: FormInterface
        before?: TableManagerBefore
        after?: TableManagerAfter
        auth?: (node: string) => boolean
    }

    /**
     * 表格管家实例
     */
    interface TableManagerInstance {
        api: TableManagerAPI
        form: FormInterface
        table: TableInterface
        opts: UseTableManagerOptions
        comSearch: ComSearchInterface
        del: (pks: string[]) => void
        auth: (node: string) => boolean
        refresh: (custom: AnyObj) => void
        initCtx: () => void
        getData: () => Promise<void> | undefined
        submitForm: (formEl?: FormInstance | null) => void
        toggleForm: (operate = '', operatePKs: string[] = []) => void
        handleEvent: (event: TableEventName, data: AnyObj) => boolean | undefined
        getEditData: (pk: string) => Promise<void> | undefined
        initDragSort: () => void
        getSelectionPKs: () => string[]
        setComSearchData: (query: AnyObj) => void
        getComSearchData: () => WhereGroup | false
        getQuickSearchData: (keywords: string) => WhereGroup | false
        setFilterWheres: (groups: WhereGroup[]) => void
    }

    /**
     * useTableManager 内事件名称
     * selection-change=选中项改变,page-size-change=每页数量改变,current-page-change=翻页,sort-change=排序,cell-change=单元格值改变,com-search=公共搜索
     * quick-search=快捷搜索,toggle-expansion=切换树形表格折叠展开,show-column-change=改变列显示状态,column-dblclick=双击列单元格
     * refresh=表头刷新,add=表头添加,edit-selected=编辑选中项,delete-selected=删除选中项,edit=编辑,delete=删除
     * reinit-com-search=重新初始化公共搜索,toggle-com-search=切换公共搜索显示状态
     */
    type TableEventName =
        | 'selection-change'
        | 'page-size-change'
        | 'current-page-change'
        | 'sort-change'
        | 'cell-change'
        | 'com-search'
        | 'quick-search'
        | 'toggle-expansion'
        | 'show-column-change'
        | 'column-dblclick'
        | 'refresh'
        | 'add'
        | 'edit'
        | 'edit-selected'
        | 'delete'
        | 'delete-selected'
        | 'reinit-com-search'
        | 'toggle-com-search'

    /**
     * useTableManager 前置处理函数（前置埋点）
     */
    interface TableManagerBefore {
        /**
         * 获取表格数据前的钩子（返回 false 可取消原操作）
         */
        getData?: () => boolean | void

        /**
         * 获取被编辑行数据前的钩子（返回 false 可取消原操作）
         * @param object.pk 被编辑行主键
         */
        getEditData?: ({ pk }: { pk: string }) => boolean | void

        /**
         * 删除前的钩子（返回 false 可取消原操作）
         * @param object.pks 被删除数据的主键集合
         */
        delete?: ({ pks }: { pks: string[] }) => boolean | void

        /**
         * 双击表格具体操作执行前钩子（返回 false 可取消原操作）
         * @param object.row 被双击行数据
         * @param object.column 被双击列数据
         */
        columnDblclick?: ({ row, column }: { row: TableRow; column: TableColumn }) => boolean | void

        /**
         * 表单切换前钩子（返回 false 可取消默认行为）
         * @param object.operate 当前操作标识:create=添加,update=更新
         * @param object.operatePKs 被操作的行的主键集合
         */
        toggleForm?: ({ operate, operatePKs }: { operate: string; operatePKs: string[] }) => boolean | void

        /**
         * 表单提交前钩子（返回 false 可取消原操作）
         * @param object.formEl 表单组件ref
         * @param object.operate 当前操作标识:create=添加,update=更新
         * @param object.items 表单数据
         */
        submitForm?: ({ formEl, operate, items }: { formEl?: FormInstance | null; operate: string; items: AnyObj }) => boolean | void

        /**
         * 表格内事件响应前钩子（返回 false 可取消原操作）
         * @param object.event 事件名称
         * @param object.data 事件携带的数据
         */
        tableEvent?: ({ event, data }: { event: TableEventName; data: AnyObj }) => boolean | void

        /**
         * 表格上下文数据初始化前钩子
         */
        initCtx?: () => boolean | void

        // 可自定义其他钩子
        [key: string]: Function | undefined
    }

    /**
     * useTableManager 后置处理函数（后置埋点）
     */
    interface TableManagerAfter {
        /**
         * 请求到表格数据后钩子
         * 此时 useTableManager.table.data 已赋值
         * @param object.res 请求完整响应
         */
        getData?: ({ res }: { res: ApiResponse }) => void

        /**
         * 获取到编辑行数据后钩子
         * 此时 useTableManager.form.items 已赋值
         * @param object.res 请求完整响应
         */
        getEditData?: ({ res }: { res: ApiResponse }) => void

        /**
         * 删除请求后钩子
         * @param object.res 请求完整响应
         */
        delete?: ({ res }: { res: ApiResponse }) => void

        /**
         * 双击单元格操作执行后钩子
         * @param object.row 当前行数据
         * @param object.column 当前列数据
         */
        columnDblclick?: ({ row, column }: { row: TableRow; column: TableColumn }) => void

        /**
         * 表单切换后钩子
         * @param object.operate 当前操作标识:create=添加,update=更新
         * @param object.operatePKs 被操作的主键集合
         */
        toggleForm?: ({ operate, operatePKs }: { operate: string; operatePKs: string[] }) => void

        /**
         * 表单提交后钩子
         * @param object.res 请求完整响应
         */
        submitForm?: ({ res }: { res: ApiResponse }) => void

        /**
         * 表格内事件响应后钩子
         * @param object.event 事件名称
         * @param object.data 事件携带的数据
         */
        tableEvent?: ({ event, data }: { event: TableEventName; data: AnyObj }) => void

        // 可自定义其他钩子
        [key: string]: Function | undefined
    }

    /**
     * 表头支持的按钮
     */
    type TableHeaderOptButton = 'refresh' | 'add' | 'edit' | 'delete' | 'rowExpansion' | 'comSearch' | 'quickSearch' | 'columnDisplay'

    /**
     * 表格上下文数据
     */
    interface TableContextData {
        row?: TableRow
        columnConfig?: TableColumn
        column?: TableColumnCtx<TableRow>
        cellValue?: any
        index?: number
    }

    /**
     * 接受表格上下文数据的任意属性计算函数
     */
    type TableContextDataFun<T> = Partial<T> | ((context: TableContextData) => Partial<T>)
}
