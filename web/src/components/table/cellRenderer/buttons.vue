<template>
    <div v-memo="[props.columnConfig]">
        <template v-for="(btn, idx) in columnConfig.buttons" :key="idx">
            <template v-if="btn.display ? btn.display(row, columnConfig) : true">
                <!-- 常规按钮 -->
                <el-button
                    v-if="btn.render == 'basic'"
                    @click="onButtonClick(btn)"
                    :class="btn.class"
                    size="small"
                    class="ag-table-render-buttons-item buttons-ml-6"
                    :type="btn.type"
                    :loading="btn.loading && btn.loading(row, columnConfig)"
                    :disabled="btn.disabled && btn.disabled(row, columnConfig)"
                    v-bind="invokeTableContextDataFun(btn.attr, { row, columnConfig, column, cellValue: btn, index })"
                >
                    <Icon v-if="btn.icon" size="14" color="var(--ag-bg-color-overlay)" :name="btn.icon" />
                    <div v-if="btn.text" class="text">{{ btn.text }}</div>
                </el-button>

                <!-- 带提示信息的按钮 -->
                <el-tooltip
                    v-if="btn.render == 'tip' && ((btn.name == 'edit' && props.manager.auth('update')) || btn.name != 'edit')"
                    :disabled="btn.title && !btn.disabledTip ? false : true"
                    :content="btn.title"
                    placement="top"
                    v-bind="invokeTableContextDataFun(columnConfig.customRenderAttr?.tooltip, { row, columnConfig, column, cellValue: btn, index })"
                >
                    <el-button
                        @click="onButtonClick(btn)"
                        :class="btn.class"
                        size="small"
                        class="ag-table-render-buttons-item buttons-ml-6"
                        :type="btn.type"
                        :loading="btn.loading && btn.loading(row, columnConfig)"
                        :disabled="btn.disabled && btn.disabled(row, columnConfig)"
                        v-bind="invokeTableContextDataFun(btn.attr, { row, columnConfig, column, cellValue: btn, index })"
                    >
                        <Icon v-if="btn.icon" size="14" color="var(--ag-bg-color-overlay)" :name="btn.icon" />
                        <div v-if="btn.text" class="text">{{ btn.text }}</div>
                    </el-button>
                </el-tooltip>

                <!-- 带确认框的按钮 -->
                <el-popconfirm
                    v-if="btn.render == 'confirm' && ((btn.name == 'delete' && props.manager.auth('delete')) || btn.name != 'delete')"
                    :disabled="btn.disabled && btn.disabled(row, columnConfig)"
                    v-bind="invokeTableContextDataFun(btn.popconfirm, { row, columnConfig, column, cellValue: btn, index })"
                    @confirm="onButtonClick(btn)"
                >
                    <template #reference>
                        <div class="buttons-popconfirm-reference-box buttons-ml-6">
                            <el-tooltip
                                :disabled="btn.title ? false : true"
                                :content="btn.title"
                                placement="top"
                                v-bind="
                                    invokeTableContextDataFun(columnConfig.customRenderAttr?.tooltip, {
                                        row,
                                        columnConfig,
                                        column,
                                        cellValue: btn,
                                        index,
                                    })
                                "
                            >
                                <el-button
                                    :class="btn.class"
                                    size="small"
                                    class="ag-table-render-buttons-item"
                                    :type="btn.type"
                                    :loading="btn.loading && btn.loading(row, columnConfig)"
                                    :disabled="btn.disabled && btn.disabled(row, columnConfig)"
                                    v-bind="invokeTableContextDataFun(btn.attr, { row, columnConfig, column, cellValue: btn, index })"
                                >
                                    <Icon v-if="btn.icon" size="14" color="var(--ag-bg-color-overlay)" :name="btn.icon" />
                                    <div v-if="btn.text" class="text">{{ btn.text }}</div>
                                </el-button>
                            </el-tooltip>
                        </div>
                    </template>
                </el-popconfirm>

                <!-- 带提示的可拖拽按钮 -->
                <el-tooltip
                    v-if="btn.render == 'sort' && ((btn.name == 'sort' && props.manager.auth('sort')) || btn.name != 'sort')"
                    :disabled="btn.title && !btn.disabledTip ? false : true"
                    :content="btn.title"
                    placement="top"
                    v-bind="invokeTableContextDataFun(columnConfig.customRenderAttr?.tooltip, { row, columnConfig, column, cellValue: btn, index })"
                >
                    <el-button
                        :class="btn.class"
                        size="small"
                        class="ag-table-render-buttons-item move-button buttons-ml-6"
                        :type="btn.type"
                        :loading="btn.loading && btn.loading(row, columnConfig)"
                        :disabled="btn.disabled && btn.disabled(row, columnConfig)"
                        v-bind="invokeTableContextDataFun(btn.attr, { row, columnConfig, column, cellValue: btn, index })"
                    >
                        <Icon v-if="btn.icon" size="14" color="var(--ag-bg-color-overlay)" :name="btn.icon" />
                        <div v-if="btn.text" class="text">{{ btn.text }}</div>
                    </el-button>
                </el-tooltip>
            </template>
        </template>
    </div>
</template>

<script setup lang="ts">
import { invokeTableContextDataFun } from '@/components/table/index'
import { CellRendererProps } from '@/components/table/types'

const props = defineProps<CellRendererProps>()

const onButtonClick = (btn: OptButton) => {
    if (typeof btn.click === 'function') {
        btn.click(props.row, props.columnConfig)
        return
    }
    props.manager.handleEvent(btn.name as TableEventName, props)
}
</script>

<style scoped lang="scss">
.ag-table-render-buttons-item {
    .text {
        font-size: 14px;
    }
    .icon + .text {
        padding-left: 5px;
    }
    &.move-button {
        cursor: move;
    }
    &.el-button--small {
        padding: 5px;
        height: auto;
    }
}
.buttons-popconfirm-reference-box {
    display: inline-flex;
    vertical-align: middle;
}
.buttons-ml-6 + .buttons-ml-6 {
    margin-left: 6px;
}
</style>
