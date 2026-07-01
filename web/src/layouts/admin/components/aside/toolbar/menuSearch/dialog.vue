<template>
    <div>
        <el-dialog
            :show-close="false"
            @close="onMenuSearchClosed"
            @opened="onMenuSearchOpened"
            :model-value="model"
            class="menu-search-dialog"
            header-class="menu-search-dialog-header"
            :destroy-on-close="true"
            :width="config.layout.shrink ? '80%' : 600"
        >
            <div>
                <el-select
                    size="large"
                    ref="menuSearchInputRef"
                    class="menu-search-input"
                    :placeholder="t('layouts.searchSimple')"
                    filterable
                    automatic-dropdown
                    @change="onMenuSearchChange"
                    @visible-change="onMenuSearchVisibleChange"
                    :popper-style="{ zIndex: SYSTEM_ZINDEX }"
                >
                    <SelectOption :menus="menu.rawData" />
                </el-select>
            </div>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ElNotification } from 'element-plus'
import { useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import SelectOption from '/@/layouts/admin/components/aside/toolbar/menuSearch/selectOption.vue'
import { useConfig } from '/@/stores/config'
import { SYSTEM_ZINDEX } from '/@/stores/constant/common'
import { useMenu } from '/@/stores/menu'
import { getMenuDataByPath } from '/@/utils/common'
import { openMenu } from '/@/utils/router'

const { t } = useI18n()
const config = useConfig()
const menu = useMenu()
const model = defineModel<boolean>({ default: false })
const menuSearchInputRef = useTemplateRef('menuSearchInputRef')

const onMenuSearchOpened = function () {
    menuSearchInputRef.value?.focus()
}

const onMenuSearchClosed = function () {
    model.value = false
}

const onMenuSearchVisibleChange = function (visible: boolean) {
    // 按下 ESC 隐藏下拉面板后，输入框不会失去焦点，导致不能再次使用 ESC 关闭外层的 dialog，element-plus@2.13.7
    if (!visible) {
        menuSearchInputRef.value?.blur()
    }
}

const onMenuSearchChange = function (path: string) {
    const menuData = getMenuDataByPath(path, menu.rawData, 'normal')
    if (menuData) {
        openMenu(menuData)
        model.value = false
    } else {
        ElNotification({ message: t('layouts.invalidRoute'), type: 'error' })
    }
}
</script>

<style scoped lang="scss">
:deep(.menu-search-dialog) {
    --el-dialog-padding-primary: 0;
    .menu-search-dialog-header {
        display: none;
    }
}
.menu-search-input {
    height: 60px;
    border: none;
    border-radius: var(--el-border-radius-base);
    box-shadow: none;
    :deep(.el-select__wrapper) {
        height: 60px;
        border: none;
        box-shadow: none;
        background-color: var(--el-fill-color-light);
        &:hover {
            border: none;
            box-shadow: none;
        }
    }
}
</style>
