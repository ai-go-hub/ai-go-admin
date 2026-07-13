import { pick } from 'lodash-es'
import { defineStore } from 'pinia'
import { reactive } from 'vue'
import { LangKey } from '/@/lang'
import { CONFIG } from '/@/stores/constant/cacheKey'
import type { Lang, Layout, Site } from '/@/stores/interface/config'
import { useMenu } from '/@/stores/menu'

export const useConfig = defineStore(
    'config',
    () => {
        /**
         * 多语言相关配置
         */
        const lang: Lang = reactive({
            active: 'zh-cn',
            fallback: 'zh-cn',
        })

        /**
         * 站点配置
         */
        const site: Site = reactive({
            name: '',
            cdnUrl: '',
            version: '',
            timezone: '',
            recordNumber: '',
            cdnUrlParams: '',
            initialized: false,
        })

        /**
         * 布局配置
         */
        const layout: Layout = reactive({
            // 全局
            mode: 'Default',
            dark: false,
            shrink: false,
            showConfigDrawer: false,
            mainAnimation: 'slide-right',

            // 菜单栏
            menuBackground: ['#ffffff', '#1d1e1f'],
            menuColor: ['#303133', '#CFD3DC'],
            menuActiveBackground: ['#ffffff', '#1d1e1f'],
            menuActiveColor: ['#409eff', '#3375b9'],
            menuHoverBackground: ['#ecf5ff', '#18222c'],
            menuWidth: 260,
            menuDefaultIcon: 'lucide-circle-small',
            menuCollapse: false,
            menuUniqueOpened: false,
            menuShowTopBar: true,
            menuTopBarBackground: ['#fcfcfc', '#1d1e1f'],
            menuTopBarColor: ['#409eff', '#3375b9'],
            menuTopBarCenter: false,
            menuTopBarLogo: false,
            menuToolBarAutoHide: false,
            menuToolBarColor: ['#303133', '#CFD3DC'],
            menuToolBarHoverColor: ['#409eff', '#CFD3DC'],
            menuToolBarHoverBackground: ['#ecf5ff', '#18222c'],

            // 主菜单栏额外配置项（一些布局存在主次两个菜单栏）
            menuBackgroundPrimary: ['#f5f5f5', '#18222c'],
            menuActiveBackgroundPrimary: ['#c6e2ff', '#1d1e1f'],

            // 左分布局独有菜单栏配置
            menuWidthLeftSplit: 180,
            menuHoverBackgroundLeftSplit: ['#ebebeb', '#213d5b'],

            // 顶栏
            headerBarTabColor: ['#000000', '#CFD3DC'],
            headerBarTabActiveColor: ['#000000', '#409EFF'],
            headerBarBackground: ['#ffffff', '#1d1e1f'],
            headerBarHoverBackground: ['#f5f5f5', '#18222c'],
            headerBarTabActiveBackground: ['#f5f5f5', '#141414'],
            headerBarTabActiveBackgroundFloating: ['#ffffff', '#1d1e1f'],

            // tour
            showTour: false,
            tourUnfinished: true,
        })

        /**
         * 设置激活语言
         */
        function setLang(value: LangKey) {
            lang.active = value
        }

        /**
         * 批量填充站点配置数据
         */
        function siteDataFill(data: Record<string, any>) {
            Object.assign(site, pick(data, Object.keys(site)))
        }

        function setSiteInitStatus(status: boolean) {
            site.initialized = status
        }

        /**
         * 设置布局配置项
         */
        const setLayoutValue = (name: keyof Layout, value: any) => {
            ;(layout[name] as any) = value
        }

        /**
         * 获取布局配置中的颜色值
         */
        const getColorValue = function (name: keyof Layout): string {
            const colors = layout[name] as string[]
            if (layout.dark) {
                return colors[1]
            } else {
                return colors[0]
            }
        }

        /**
         * 获取菜单宽度
         */
        function getMenuWidth() {
            // 菜单折叠时基本宽度
            const menuCollapseBaseWidth = 64

            // 左分布局特有
            if (layout.mode == 'LeftSplit') {
                const menu = useMenu()

                // 本布局带来的额外菜单宽度，主菜单宽度 80 + 次级菜单的左右内边距 16
                const modeMenuWidth = 96
                // 最终菜单宽度
                let leftSplitMenuWidth = layout.menuCollapse
                    ? menuCollapseBaseWidth + modeMenuWidth
                    : parseInt(layout.menuWidthLeftSplit.toString()) + modeMenuWidth

                // 无次级菜单，固定宽度
                if (!menu.children.length) {
                    leftSplitMenuWidth = 80
                }

                // 小屏模式
                if (layout.shrink) {
                    return layout.menuCollapse ? 0 : `${leftSplitMenuWidth}px`
                }

                return `${leftSplitMenuWidth}px`
            }

            // 小屏模式
            if (layout.shrink) {
                return layout.menuCollapse ? 0 : `${layout.menuWidth}px`
            }

            // 菜单是否折叠
            return layout.menuCollapse ? `${menuCollapseBaseWidth}px` : `${layout.menuWidth}px`
        }

        return { lang, site, layout, setLang, siteDataFill, setSiteInitStatus, setLayoutValue, getColorValue, getMenuWidth }
    },
    {
        persist: {
            key: CONFIG,
        },
    }
)
