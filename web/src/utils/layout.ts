import { useConfig } from '@/stores/config'
import { useNavTab } from '@/stores/navTab'
import { isAdminApp } from '@/utils/common'
import type { CSSProperties } from 'vue'

/**
 * 管理员后台各个布局顶栏高度
 */
export const adminLayoutHeaderBarHeight = {
    Double: 60,
    Default: 70,
    Classic: 50,
    LeftSplit: 70,
    Streamline: 60,
}

/**
 * main高度
 * @param extra main 高度额外减去的px数,可以实现隐藏原有的滚动条
 * @returns CSSProperties
 */
export function mainHeight(extra = 0): CSSProperties {
    let height = extra
    if (isAdminApp()) {
        const config = useConfig()
        const navTab = useNavTab()
        if (!navTab.state.activeFullScreen) {
            height += adminLayoutHeaderBarHeight[config.layout.mode as keyof typeof adminLayoutHeaderBarHeight]
        }
    }
    return {
        height: 'calc(100vh - ' + height.toString() + 'px)',
    }
}

/**
 * 设置导航栏宽度
 * @returns
 */
export function setNavTabsWidth() {
    const navTabs = document.querySelector('.nav-tabs') as HTMLElement
    if (!navTabs) {
        return
    }
    const navBar = document.querySelector('.nav-bar') as HTMLElement
    const navMenus = document.querySelector('.nav-menus') as HTMLElement
    const minWidth = navBar.offsetWidth - (navMenus.offsetWidth + 20)
    navTabs.style.width = minWidth.toString() + 'px'
}
