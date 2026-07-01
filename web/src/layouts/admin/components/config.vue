<template>
    <div class="layout-config-drawer">
        <el-drawer :model-value="configStore.layout.showConfigDrawer" :title="t('layouts.layoutConfiguration')" size="410px" @close="onCloseDrawer">
            <el-scrollbar class="layout-mode-style-scrollbar">
                <el-form :model="configStore.layout">
                    <div class="layout-mode-styles-box">
                        <el-divider content-position="left" border-style="dashed">{{ t('layouts.layoutMode') }}</el-divider>
                        <div class="layout-mode-box-style">
                            <el-row class="layout-mode-box-style-row" :gutter="10">
                                <el-col :span="12">
                                    <div
                                        @click="setLayoutMode('Default')"
                                        class="layout-mode-style default"
                                        :class="configStore.layout.mode == 'Default' ? 'active' : ''"
                                    >
                                        <div class="layout-mode-style-box">
                                            <div class="layout-mode-style-aside"></div>
                                            <div class="layout-mode-style-container-box">
                                                <div class="layout-mode-style-header"></div>
                                                <div class="layout-mode-style-container">
                                                    <div class="layout-mode-style-name">{{ t('layouts.default') }}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </el-col>
                                <el-col :span="12">
                                    <div
                                        @click="setLayoutMode('Classic')"
                                        class="layout-mode-style classic"
                                        :class="configStore.layout.mode == 'Classic' ? 'active' : ''"
                                    >
                                        <div class="layout-mode-style-box">
                                            <div class="layout-mode-style-aside"></div>
                                            <div class="layout-mode-style-container-box">
                                                <div class="layout-mode-style-header"></div>
                                                <div class="layout-mode-style-container">
                                                    <div class="layout-mode-style-name">{{ t('layouts.classic') }}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </el-col>
                            </el-row>
                            <el-row class="layout-mode-box-style-row" :gutter="10">
                                <el-col :span="12">
                                    <div
                                        @click="setLayoutMode('LeftSplit')"
                                        class="layout-mode-style left-split"
                                        :class="configStore.layout.mode == 'LeftSplit' ? 'active' : ''"
                                    >
                                        <div class="layout-mode-style-box">
                                            <div class="layout-mode-style-aside">
                                                <div class="left-split-aside"></div>
                                            </div>
                                            <div class="layout-mode-style-container-box">
                                                <div class="layout-mode-style-header"></div>
                                                <div class="layout-mode-style-container">
                                                    <div class="layout-mode-style-name">{{ t('layouts.leftSplit') }}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </el-col>
                                <el-col :span="12">
                                    <div
                                        @click="setLayoutMode('Double')"
                                        class="layout-mode-style double"
                                        :class="configStore.layout.mode == 'Double' ? 'active' : ''"
                                    >
                                        <div class="layout-mode-style-box">
                                            <div class="layout-mode-style-aside"></div>
                                            <div class="layout-mode-style-container-box">
                                                <div class="layout-mode-style-header"></div>
                                                <div class="layout-mode-style-container">
                                                    <div class="layout-mode-style-name">{{ t('layouts.doubleColumn') }}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </el-col>
                            </el-row>
                            <el-row :gutter="10">
                                <el-col :span="12">
                                    <div
                                        @click="setLayoutMode('Streamline')"
                                        class="layout-mode-style streamline"
                                        :class="configStore.layout.mode == 'Streamline' ? 'active' : ''"
                                    >
                                        <div class="layout-mode-style-box">
                                            <div class="layout-mode-style-container-box">
                                                <div class="layout-mode-style-header"></div>
                                                <div class="layout-mode-style-container">
                                                    <div class="layout-mode-style-name">{{ t('layouts.singleColumn') }}</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </el-col>
                            </el-row>
                        </div>

                        <el-divider content-position="left" border-style="dashed">{{ t('layouts.overallSituation') }}</el-divider>
                        <div class="layout-config-global">
                            <el-form-item size="large" :label="t('layouts.darkMode')">
                                <DarkSwitch @click="toggleDark()" />
                            </el-form-item>
                            <el-form-item :label="t('layouts.backgroundPageSwitchingAnimation')">
                                <el-select
                                    @change="onCommitState($event, 'mainAnimation')"
                                    :model-value="configStore.layout.mainAnimation"
                                    :placeholder="t('layouts.pleaseSelectAnimationName')"
                                >
                                    <el-option label="slide-right" value="slide-right"></el-option>
                                    <el-option label="slide-left" value="slide-left"></el-option>
                                    <el-option label="el-fade-in-linear" value="el-fade-in-linear"></el-option>
                                    <el-option label="el-fade-in" value="el-fade-in"></el-option>
                                    <el-option label="el-zoom-in-center" value="el-zoom-in-center"></el-option>
                                    <el-option label="el-zoom-in-top" value="el-zoom-in-top"></el-option>
                                    <el-option label="el-zoom-in-bottom" value="el-zoom-in-bottom"></el-option>
                                </el-select>
                            </el-form-item>
                        </div>

                        <el-divider v-if="configStore.layout.mode != 'Streamline'" content-position="left" border-style="dashed">
                            {{ t('layouts.sidebar') }}
                        </el-divider>
                        <div v-if="configStore.layout.mode != 'Streamline'" class="layout-config-aside">
                            <!-- 侧边菜单栏背景色（左分布局的主菜单） -->
                            <el-form-item v-if="configStore.layout.mode == 'LeftSplit'" :label="getLabel('menuBackgroundPrimary')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuBackgroundPrimary')"
                                    :model-value="configStore.getColorValue('menuBackgroundPrimary')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单激活项背景色（左分布局的主菜单） -->
                            <el-form-item v-if="configStore.layout.mode == 'LeftSplit'" :label="getLabel('menuActiveBackgroundPrimary')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuActiveBackgroundPrimary')"
                                    :model-value="configStore.getColorValue('menuActiveBackgroundPrimary')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单栏背景色（左分布局的次级菜单，其余布局的主菜单） -->
                            <el-form-item :label="getLabel('menuBackground')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuBackground')"
                                    :model-value="configStore.getColorValue('menuBackground')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单激活项背景色（左分布局的次级菜单，其余布局的主菜单） -->
                            <el-form-item :label="getLabel('menuActiveBackground')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuActiveBackground')"
                                    :model-value="configStore.getColorValue('menuActiveBackground')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单激活项文字色（主次菜单通用） -->
                            <el-form-item :label="getLabel('menuActiveColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuActiveColor')"
                                    :model-value="configStore.getColorValue('menuActiveColor')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单文字颜色（主次菜单通用） -->
                            <el-form-item :label="getLabel('menuColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuColor')"
                                    :model-value="configStore.getColorValue('menuColor')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单悬停时背景色（主次菜单通用） -->
                            <el-form-item v-if="configStore.layout.mode == 'LeftSplit'" :label="getLabel('menuHoverBackground')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuHoverBackgroundLeftSplit')"
                                    :model-value="configStore.getColorValue('menuHoverBackgroundLeftSplit')"
                                />
                            </el-form-item>
                            <el-form-item v-else :label="getLabel('menuHoverBackground')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuHoverBackground')"
                                    :model-value="configStore.getColorValue('menuHoverBackground')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单宽度（展开时宽度，左分布局的次级菜单，其余布局的主菜单） -->
                            <el-form-item v-if="configStore.layout.mode == 'LeftSplit'" :label="getLabel('menuWidth')">
                                <el-input
                                    @input="onCommitState($event, 'menuWidthLeftSplit')"
                                    type="number"
                                    :step="10"
                                    :model-value="configStore.layout.menuWidthLeftSplit"
                                >
                                    <template #append>px</template>
                                </el-input>
                            </el-form-item>
                            <el-form-item v-else :label="getLabel('menuWidth')">
                                <el-input
                                    @input="onCommitState($event, 'menuWidth')"
                                    type="number"
                                    :step="10"
                                    :model-value="configStore.layout.menuWidth"
                                >
                                    <template #append>px</template>
                                </el-input>
                            </el-form-item>

                            <el-form-item :label="t('layouts.sideMenuDefaultIcon')">
                                <IconSelector
                                    @change="onCommitMenuDefaultIcon($event, 'menuDefaultIcon')"
                                    :model-value="configStore.layout.menuDefaultIcon"
                                />
                            </el-form-item>
                            <el-form-item :label="t('layouts.sideMenuHorizontalCollapse')">
                                <el-switch @change="onCommitState($event, 'menuCollapse')" :model-value="configStore.layout.menuCollapse"></el-switch>
                            </el-form-item>
                            <el-form-item :label="t('layouts.sideMenuAccordion')">
                                <el-switch
                                    @change="onCommitState($event, 'menuUniqueOpened')"
                                    :model-value="configStore.layout.menuUniqueOpened"
                                ></el-switch>
                            </el-form-item>
                        </div>

                        <el-divider content-position="left" border-style="dashed">
                            {{ t('layouts.sidebarTopAndBottom') }}
                        </el-divider>
                        <div class="layout-config-aside">
                            <!-- 显示侧边菜单顶栏（站点标题栏），左分布局没有标题栏，只有 LOGO -->
                            <el-form-item
                                v-if="!['LeftSplit', 'Streamline'].includes(configStore.layout.mode)"
                                :label="t('layouts.showSideMenuTopBar')"
                            >
                                <el-switch
                                    @change="onCommitState($event, 'menuShowTopBar')"
                                    :model-value="configStore.layout.menuShowTopBar"
                                ></el-switch>
                            </el-form-item>

                            <!-- 侧边菜单顶栏背景色 -->
                            <el-form-item
                                v-if="!['LeftSplit', 'Streamline'].includes(configStore.layout.mode)"
                                :label="t('layouts.sideMenuTopBarBackgroundColor')"
                            >
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuTopBarBackground')"
                                    :model-value="configStore.getColorValue('menuTopBarBackground')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单顶栏文字色 -->
                            <el-form-item v-if="configStore.layout.mode != 'LeftSplit'" :label="t('layouts.sideMenuTopBarTextColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'menuTopBarColor')"
                                    :model-value="configStore.getColorValue('menuTopBarColor')"
                                />
                            </el-form-item>

                            <!-- 侧边菜单顶栏内容居中 -->
                            <el-form-item v-if="configStore.layout.mode != 'LeftSplit'" :label="t('layouts.sideMenuTopBarCenterContent')">
                                <el-switch
                                    @change="onCommitState($event, 'menuTopBarCenter')"
                                    :model-value="configStore.layout.menuTopBarCenter"
                                ></el-switch>
                            </el-form-item>

                            <!-- 侧边菜单顶栏显示LOGO -->
                            <el-form-item :label="t('layouts.sideMenuTopBarDisplayLogo')">
                                <el-switch
                                    @change="onCommitState($event, 'menuTopBarLogo')"
                                    :model-value="configStore.layout.menuTopBarLogo"
                                ></el-switch>
                            </el-form-item>

                            <!-- 侧边菜单底部工具栏 -->
                            <template v-if="configStore.layout.mode != 'Streamline'">
                                <el-form-item :label="t('layouts.sideMenuBottomToolbarAutoHide')">
                                    <el-switch
                                        @change="onCommitState($event, 'menuToolBarAutoHide')"
                                        :model-value="configStore.layout.menuToolBarAutoHide"
                                    ></el-switch>
                                </el-form-item>
                                <el-form-item :label="t('layouts.sideMenuBottomToolbarIconColor')">
                                    <el-color-picker
                                        @change="onCommitColorState($event, 'menuToolBarColor')"
                                        :model-value="configStore.getColorValue('menuToolBarColor')"
                                    />
                                </el-form-item>
                                <el-form-item :label="t('layouts.sideMenuBottomToolbarHoverIconColor')">
                                    <el-color-picker
                                        @change="onCommitColorState($event, 'menuToolBarHoverColor')"
                                        :model-value="configStore.getColorValue('menuToolBarHoverColor')"
                                    />
                                </el-form-item>
                                <el-form-item :label="t('layouts.sideMenuBottomToolbarHoverBackgroundColor')">
                                    <el-color-picker
                                        @change="onCommitColorState($event, 'menuToolBarHoverBackground')"
                                        :model-value="configStore.getColorValue('menuToolBarHoverBackground')"
                                    />
                                </el-form-item>
                            </template>
                        </div>

                        <el-divider content-position="left" border-style="dashed">{{ t('layouts.topBar') }}</el-divider>
                        <div class="layout-config-aside">
                            <el-form-item :label="t('layouts.topBarBackgroundColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'headerBarBackground')"
                                    :model-value="configStore.getColorValue('headerBarBackground')"
                                />
                            </el-form-item>
                            <el-form-item :label="t('layouts.topBarTextColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'headerBarTabColor')"
                                    :model-value="configStore.getColorValue('headerBarTabColor')"
                                />
                            </el-form-item>
                            <el-form-item :label="t('layouts.topBarHoverBackgroundColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'headerBarHoverBackground')"
                                    :model-value="configStore.getColorValue('headerBarHoverBackground')"
                                />
                            </el-form-item>

                            <!-- 顶栏激活项背景色 -->
                            <el-form-item
                                v-if="['Default', 'LeftSplit'].includes(configStore.layout.mode)"
                                :label="t('layouts.topBarMenuActiveItemBackgroundColor')"
                            >
                                <el-color-picker
                                    @change="onCommitColorState($event, 'headerBarTabActiveBackgroundFloating')"
                                    :model-value="configStore.getColorValue('headerBarTabActiveBackgroundFloating')"
                                />
                            </el-form-item>
                            <el-form-item v-else :label="t('layouts.topBarMenuActiveItemBackgroundColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'headerBarTabActiveBackground')"
                                    :model-value="configStore.getColorValue('headerBarTabActiveBackground')"
                                />
                            </el-form-item>

                            <el-form-item :label="t('layouts.topBarMenuActiveItemTextColor')">
                                <el-color-picker
                                    @change="onCommitColorState($event, 'headerBarTabActiveColor')"
                                    :model-value="configStore.getColorValue('headerBarTabActiveColor')"
                                />
                            </el-form-item>
                        </div>
                        <el-popconfirm @confirm="restoreDefault" :title="t('layouts.restoreConfigConfirm')">
                            <template #reference>
                                <div class="flex-center">
                                    <el-button class="w90" type="info">{{ t('layouts.restoreDefault') }}</el-button>
                                </div>
                            </template>
                        </el-popconfirm>
                    </div>
                </el-form>
            </el-scrollbar>
        </el-drawer>
    </div>
</template>

<script setup lang="ts">
import { nextTick } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import IconSelector from '/@/components/agInput/components/iconSelector.vue'
import DarkSwitch from '/@/layouts/common/components/darkSwitch.vue'
import { useConfig } from '/@/stores/config'
import { BEFORE_RESIZE_LAYOUT, CONFIG } from '/@/stores/constant/cacheKey'
import type { Layout } from '/@/stores/interface/config'
import { useMenu } from '/@/stores/menu'
import toggleDark from '/@/utils/dark'
import { Local, Session } from '/@/utils/storage'

const { t } = useI18n()
const menu = useMenu()
const router = useRouter()
const configStore = useConfig()

/**
 * 布局不同时，label 可能不同，特用函数规范化
 */
const getLabel = (name: string) => {
    const labelTranslationKey = new Map<string, string[]>([
        // 侧边菜单栏背景色（左分布局的主菜单）
        ['menuBackgroundPrimary', ['sideMenuBarBackgroundColor', 'mainMenu']],
        // 侧边菜单激活项背景色（左分布局的主菜单）
        ['menuActiveBackgroundPrimary', ['sideMenuActiveItemBackgroundColor', 'mainMenu']],
        // 侧边菜单栏背景色（左分布局的次级菜单，其余布局的主菜单）
        ['menuBackground', ['sideMenuBarBackgroundColor', 'submenu']],
        // 侧边菜单激活项背景色（左分布局的次级菜单，其余布局的主菜单）
        ['menuActiveBackground', ['sideMenuActiveItemBackgroundColor', 'submenu']],
        // 侧边菜单激活项文字色（主次菜单通用）
        ['menuActiveColor', ['sideMenuActiveItemTextColor', 'mainAndSecondaryMenus']],
        // 侧边菜单文字颜色（主次菜单通用）
        ['menuColor', ['sideMenuTextColor', 'mainAndSecondaryMenus']],
        // 侧边菜单悬停时背景色（主次菜单通用）
        ['menuHoverBackground', ['sideMenuHoverBackgroundColor', 'mainAndSecondaryMenus']],
        // 侧边菜单宽度（展开时宽度，左分布局的次级菜单，其余布局的主菜单）
        ['menuWidth', ['sideMenuWidth', 'submenu']],
    ])

    if (labelTranslationKey.has(name)) {
        const label = labelTranslationKey.get(name) as string[]

        // 左分布局下，label 可能带有后缀
        if (configStore.layout.mode == 'LeftSplit') {
            return t(`layouts.${label[0]}`) + t(`layouts.${label[1]}`)
        }

        return t(`layouts.${label[0]}`)
    }
    return name
}

const onCommitState = (value: any, name: any) => {
    configStore.setLayoutValue(name, value)
}

const onCommitColorState = (value: string | null, name: keyof Layout) => {
    if (value === null) return
    const colors = configStore.layout[name] as string[]
    if (configStore.layout.dark) {
        colors[1] = value
    } else {
        colors[0] = value
    }
    configStore.setLayoutValue(name, colors)
}

const setLayoutMode = (mode: string) => {
    configStore.setLayoutValue('mode', mode)
}

// 修改默认菜单图标
const onCommitMenuDefaultIcon = (value: any, name: any) => {
    configStore.setLayoutValue(name, value)

    // 全部菜单重新渲染
    const menus = menu.rawData
    menu.setRawData([])
    nextTick(() => {
        menu.setRawData(menus)
    })
}

const onCloseDrawer = () => {
    configStore.setLayoutValue('showConfigDrawer', false)
}

const restoreDefault = () => {
    Local.remove(CONFIG)
    Session.remove(BEFORE_RESIZE_LAYOUT)
    router.go(0)
}
</script>

<style scoped lang="scss">
.w90 {
    width: 90%;
}
.flex-center {
    display: flex;
    align-items: center;
    justify-content: center;
}
.layout-config-drawer :deep(.el-input__inner) {
    padding: 0 0 0 6px;
}
.layout-config-drawer :deep(.el-input-group__append) {
    padding: 0 10px;
}
.layout-config-drawer :deep(.el-drawer__header) {
    margin-bottom: 0 !important;
}
.layout-config-drawer :deep(.el-drawer__body) {
    padding: 0;
}
.layout-mode-styles-box {
    padding: 20px;
}
.layout-mode-box-style-row {
    margin-bottom: 10px;
}
.layout-mode-style {
    position: relative;
    height: 100px;
    border: 1px solid var(--el-border-color-light);
    border-radius: var(--el-border-radius-small);
    &:hover,
    &.active {
        border: 1px solid var(--el-color-primary);
    }
    .layout-mode-style-name {
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--el-color-primary-light-5);
        border-radius: var(--el-border-radius-base);
        height: 50px;
        width: 100px;
        border: 1px solid var(--el-color-primary-light-7);
    }
    .layout-mode-style-box {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 100%;
    }
    &.default {
        .layout-mode-style-aside {
            width: 18%;
            height: 90%;
            background-color: var(--el-border-color-lighter);
        }
        .layout-mode-style-container-box {
            width: 68%;
            height: 90%;
            margin-left: 4%;
            .layout-mode-style-header {
                width: 100%;
                height: 10%;
                background-color: var(--el-border-color-lighter);
            }
            .layout-mode-style-container {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 100%;
                height: 85%;
                background-color: var(--el-border-color-extra-light);
                margin-top: 5%;
            }
        }
    }
    &.classic {
        .layout-mode-style-aside {
            width: 18%;
            height: 100%;
            background-color: var(--el-border-color-lighter);
        }
        .layout-mode-style-container-box {
            width: 82%;
            height: 100%;
            .layout-mode-style-header {
                width: 100%;
                height: 10%;
                background-color: var(--el-border-color);
            }
            .layout-mode-style-container {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 100%;
                height: 90%;
                background-color: var(--el-border-color-extra-light);
            }
        }
    }
    &.left-split {
        .layout-mode-style-aside {
            width: 18%;
            height: 90%;
            background-color: var(--el-border-color-lighter);
            .left-split-aside {
                width: 2px;
                height: 100%;
                margin-left: 30%;
                background-color: var(--el-bg-color);
            }
        }
        .layout-mode-style-container-box {
            width: 68%;
            height: 90%;
            margin-left: 4%;
            .layout-mode-style-header {
                width: 100%;
                height: 10%;
                background-color: var(--el-border-color-lighter);
            }
            .layout-mode-style-container {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 100%;
                height: 85%;
                background-color: var(--el-border-color-extra-light);
                margin-top: 5%;
            }
        }
    }
    &.double {
        .layout-mode-style-aside {
            width: 18%;
            height: 100%;
            background-color: var(--el-border-color);
        }
        .layout-mode-style-container-box {
            width: 82%;
            height: 100%;
            .layout-mode-style-header {
                width: 100%;
                height: 10%;
                background-color: var(--el-border-color);
            }
            .layout-mode-style-container {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 100%;
                height: 90%;
                background-color: var(--el-border-color-extra-light);
            }
        }
    }
    &.streamline {
        .layout-mode-style-container-box {
            width: 100%;
            height: 100%;
            .layout-mode-style-header {
                width: 100%;
                height: 10%;
                background-color: var(--el-border-color);
            }
            .layout-mode-style-container {
                display: flex;
                align-items: center;
                justify-content: center;
                width: 100%;
                height: 90%;
                background-color: var(--el-border-color-extra-light);
            }
        }
    }
}
</style>
