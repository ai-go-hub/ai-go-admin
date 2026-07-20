import { TableColumnCtx } from 'element-plus'

export type DefaultOptButType = 'sort' | 'edit' | 'delete'

export interface CellRendererProps {
    row: TableRow
    columnConfig: TableColumn
    column: TableColumnCtx<TableRow>
    cellValue: any
    index: number
    manager: TableManagerInstance
}
