<template>
    <!-- 公共搜索 -->
    <el-collapse-transition>
        <ComSearch
            :manager="props.manager"
            v-model:com-search="comSearch!"
            v-if="props.buttons.includes('comSearch') && props.manager.table.showComSearch"
        >
            <template v-for="(slot, idx) in $slots" :key="idx" #[idx]>
                <slot :name="idx"></slot>
            </template>
        </ComSearch>
    </el-collapse-transition>

    <!-- 操作按钮组 -->
    <div v-bind="$attrs" class="table-header ag-scroll-style">
        <slot name="refreshPrepend"></slot>

        <el-tooltip v-if="props.buttons.includes('refresh')" :content="t('common.refresh')" placement="top">
            <el-button
                @click="props.manager.handleEvent('refresh', { event: 'header-btn' })"
                color="#40485b"
                class="table-header-operate btns-ml-12"
                type="info"
            >
                <Icon size="14" name="lucide-refresh-cw" color="#ffffff" />
            </el-button>
        </el-tooltip>

        <slot name="refreshAppend"></slot>

        <el-tooltip v-if="props.buttons.includes('add') && props.manager.auth('create')" :content="t('common.add')" placement="top">
            <el-button @click="props.manager.handleEvent('add', { event: 'header-btn' })" class="table-header-operate btns-ml-12" type="primary">
                <Icon size="16" :stroke-width="3" name="lucide-plus" color="#ffffff" />
                <span class="table-header-operate-text">{{ t('common.add') }}</span>
            </el-button>
        </el-tooltip>

        <el-tooltip v-if="props.buttons.includes('edit') && props.manager.auth('update')" :content="t('common.editSelection')" placement="top">
            <el-button
                @click="props.manager.handleEvent('edit-selected', { event: 'header-btn' })"
                :disabled="!enableBatchOpt"
                class="table-header-operate btns-ml-12"
                type="primary"
            >
                <Icon size="16" name="lucide-pencil" color="#ffffff" />
                <span class="table-header-operate-text">{{ t('common.edit') }}</span>
            </el-button>
        </el-tooltip>

        <el-popconfirm
            v-if="props.buttons.includes('delete') && props.manager.auth('delete')"
            @confirm="props.manager.handleEvent('delete-selected', { event: 'header-btn' })"
            :confirm-button-text="t('common.delete')"
            :cancel-button-text="t('common.cancel')"
            confirmButtonType="danger"
            :title="t('common.deleteSelectedRecords')"
            :disabled="!enableBatchOpt"
        >
            <template #reference>
                <div class="btns-ml-12">
                    <el-tooltip :content="t('common.deleteSelection')" placement="top">
                        <el-button :disabled="!enableBatchOpt" class="table-header-operate" type="danger">
                            <Icon size="16" name="lucide-trash" color="#ffffff" />
                            <span class="table-header-operate-text">{{ t('common.delete') }}</span>
                        </el-button>
                    </el-tooltip>
                </div>
            </template>
        </el-popconfirm>

        <el-tooltip
            v-if="props.buttons.includes('rowExpansion')"
            :content="(props.manager.table.expandAll ? t('common.collapse') : t('common.expand')) + t('common.all') + t('common.submenu')"
            placement="top"
        >
            <el-button
                @click="props.manager.handleEvent('toggle-expansion', { expanded: !props.manager.table.expandAll })"
                class="table-header-operate btns-ml-12"
                :type="props.manager.table.expandAll ? 'danger' : 'warning'"
            >
                <span class="table-header-operate-text">
                    {{ (props.manager.table.expandAll ? t('common.collapse') : t('common.expand')) + t('common.all') }}
                </span>
            </el-button>
        </el-tooltip>

        <!-- slot -->
        <slot></slot>

        <!-- 右侧搜索框和工具按钮 -->
        <div class="table-search">
            <slot name="quickSearchPrepend"></slot>

            <el-input
                v-if="props.buttons.includes('quickSearch')"
                v-model="quickSearchKeywords"
                class="xs-hidden quick-search"
                @input="onSearchInput"
                :placeholder="quickSearchPlaceholder ? quickSearchPlaceholder : t('common.search')"
                clearable
            />

            <div class="table-search-button-group" v-if="props.buttons.includes('columnDisplay') || props.buttons.includes('comSearch')">
                <el-dropdown v-if="props.buttons.includes('columnDisplay')" :max-height="380" :hide-on-click="false">
                    <div class="table-search-button-item" :class="props.buttons.includes('comSearch') ? 'right-border' : ''">
                        <Icon size="14" name="el-grid" color="var(--el-text-color-primary)" />
                    </div>
                    <template #dropdown>
                        <el-dropdown-menu>
                            <el-dropdown-item v-for="(item, idx) in columnDisplay" :key="idx">
                                <el-checkbox
                                    v-if="item.prop"
                                    @change="onChangeShowColumn($event, item.prop!)"
                                    :model-value="item.show !== false"
                                    size="small"
                                    :label="item.label"
                                />
                            </el-dropdown-item>
                        </el-dropdown-menu>
                    </template>
                </el-dropdown>

                <el-tooltip
                    v-if="props.buttons.includes('comSearch')"
                    :disabled="props.manager.table.showComSearch"
                    :content="t('common.openCommonSearch')"
                    placement="top"
                >
                    <div
                        class="table-search-button-item"
                        @click="props.manager.handleEvent('toggle-com-search', { value: !props.manager.table.showComSearch })"
                    >
                        <Icon size="14" name="el-search" color="var(--el-text-color-primary)" />
                    </div>
                </el-tooltip>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { debounce } from 'lodash-es'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import ComSearch from '/@/components/table/header/comSearch.vue'

const { t } = useI18n()
const quickSearchKeywords = ref()

interface Props {
    manager: TableManagerInstance
    buttons: TableHeaderOptButton[]
    quickSearchPlaceholder?: string
}

const props = withDefaults(defineProps<Props>(), {
    buttons: () => {
        return ['refresh', 'add', 'edit', 'delete']
    },
    quickSearchPlaceholder: '',
})
const comSearch = defineModel<ComSearchInterface>('comSearch')

if (props.buttons.includes('comSearch') && !comSearch.value) {
    console.error(`[TableHeader] 在启用公共搜索的情况下，必须传递 v-model:com-search="tableManager.comSearch"`)
}

const columnDisplay = computed(() => {
    let columnDisplayArr = []
    for (let item of props.manager.table.column) {
        item.type === 'selection' || item.render === 'buttons' || item.columnDisplayControl === false ? '' : columnDisplayArr.push(item)
    }
    return columnDisplayArr
})

const enableBatchOpt = computed(() => props.manager.table.selections!.length > 0)

const onSearchInput = debounce(() => {
    props.manager.handleEvent('quick-search', { keywords: quickSearchKeywords.value })
}, 500)

const onChangeShowColumn = (value: string | number | boolean, field: string) => {
    props.manager.handleEvent('show-column-change', { field, value })
}
</script>

<style scoped lang="scss">
.table-header {
    position: relative;
    overflow-x: auto;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    width: 100%;
    max-width: 100%;
    background-color: var(--ag-bg-color-overlay);
    border: 1px solid var(--ag-border-color);
    border-bottom: none;
    padding: 13px 15px;
    font-size: 14px;
    .table-header-operate-text {
        margin-left: 6px;
    }
}
.btns-ml-12 + .btns-ml-12 {
    margin-left: 12px;
}
.table-search {
    display: flex;
    margin-left: auto;
    .quick-search {
        width: auto;
    }
}
.table-search-button-group {
    display: flex;
    margin-left: 12px;
    border: 1px solid var(--el-border-color);
    border-radius: var(--el-border-radius-base);
    overflow: hidden;
    .table-search-button-item {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 44px;
        height: 30px;
        border: none;
        outline: none;
        cursor: pointer;
        &:hover {
            background-color: var(--el-color-info-light-7);
        }
    }
    .right-border {
        border-right: 1px solid var(--el-border-color);
    }
}
</style>
