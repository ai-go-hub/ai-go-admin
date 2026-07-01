import { LangKey } from '/@/lang/index'

export interface Lang {
    // 激活语言
    active: LangKey
    // 备用语言
    fallback: LangKey
}

/**
 * 站点配置
 */
export interface Site {
    name: string
    version: string
    timezone: string
    record_number: string
    initialized: boolean
}

export interface Layout {
    // ==================== 全局 ====================
    // 后台布局方式，可选值<Default|Classic|Streamline|Double|LeftSplit>
    mode: string
    // 是否暗黑模式
    dark: boolean
    // 是否收缩布局（小屏设备）
    shrink: boolean
    // 是否显示布局配置抽屉
    showConfigDrawer: boolean
    // 后台主页面切换动画，可选值<slide-right|slide-left|el-fade-in-linear|el-fade-in|el-zoom-in-center|el-zoom-in-top|el-zoom-in-bottom>
    mainAnimation: string

    // ==================== 侧边栏 ====================
    // 侧边菜单宽度（展开时），单位px
    menuWidth: number
    // 侧边菜单项默认图标
    menuDefaultIcon: string
    // 是否水平折叠收起菜单
    menuCollapse: boolean
    // 是否只保持一个子菜单的展开（手风琴）
    menuUniqueOpened: boolean
    // 侧边菜单背景色
    menuBackground: string[]
    // 侧边菜单文字颜色
    menuColor: string[]
    // 侧边菜单激活项背景色
    menuActiveBackground: string[]
    // 侧边菜单激活项文字色
    menuActiveColor: string[]
    // 侧边菜单悬停时背景色
    menuHoverBackground: string[]
    // 显示菜单栏顶栏（站点标题栏）
    menuShowTopBar: boolean
    // 侧边菜单顶栏背景色
    menuTopBarBackground: string[]
    // 侧边菜单顶栏文字颜色
    menuTopBarColor: string[]
    // 侧边菜单顶栏内容居中
    menuTopBarCenter: boolean
    // 侧边菜单顶栏显示LOGO
    menuTopBarLogo: boolean
    // 侧边菜单底部工具栏自动隐藏
    menuToolBarAutoHide: boolean
    // 侧边菜单底部工具栏图标颜色
    menuToolBarColor: string[]
    // 侧边菜单底部工具栏悬停时图标色
    menuToolBarHoverColor: string[]
    // 侧边菜单底部工具栏悬停时背景色
    menuToolBarHoverBackground: string[]

    // 侧边菜单背景色（一些布局存在主次两个菜单栏，此处单独配置主菜单栏的背景色）
    menuBackgroundPrimary: string[]
    // 侧边菜单激活项背景色（一些布局存在主次两个菜单栏，此处单独配置主菜单栏的激活项背景色）
    menuActiveBackgroundPrimary: string[]

    // 侧边菜单宽度，左分双栏专用（展开时），单位px
    menuWidthLeftSplit: number
    // 侧边菜单悬停时背景色，左分双栏专用
    menuHoverBackgroundLeftSplit: string[]

    // ==================== 顶栏 ====================
    // 顶栏文字色
    headerBarTabColor: string[]
    // 顶栏背景色
    headerBarBackground: string[]
    // 顶栏悬停时背景色
    headerBarHoverBackground: string[]
    // 顶栏激活项背景色
    headerBarTabActiveBackground: string[]
    // 顶栏激活项文字色
    headerBarTabActiveColor: string[]
    // 顶栏激活项背景色（悬浮顶栏专用，即默认布局和左分双栏布局）
    headerBarTabActiveBackgroundFloating: string[]

    // ==================== 漫游式引导（tour） ====================
    // 布局漫游式引导显示状态
    showTour: boolean
    // 布局漫游式引导未完成标记
    tourUnfinished: boolean
}
