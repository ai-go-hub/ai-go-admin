/**
 * reference
 * 全局提供：引用（指向）一些对象（组件）的句柄，比如公共组件的 ref
 */
import NavTabs from '@/layouts/backend/components/navBar/tabs.vue'
import { mainHeight } from '@/utils/layout'
import type { ScrollbarInstance } from 'element-plus'
import type { CSSProperties, VNodeRef } from 'vue'
import { computed, ref } from 'vue'

/**
 * 后台顶栏（NavTabs）组件 ref（组件仅默认和经典布局显示）
 */
export const layoutNavTabsRef = ref<InstanceType<typeof NavTabs>>()

/**
 * 后台布局的主体的滚动条组件 ref
 */
export const layoutMainScrollbarRef = ref<VNodeRef>()

/**
 * 后台布局的主体滚动条的额外样式，包括高度
 */
export const layoutMainScrollbarStyle = computed<CSSProperties>(() => mainHeight())

/**
 * 后台布局的菜单组件的 ref
 */
export const layoutMenuRef = ref<VNodeRef>()

/**
 * 后台布局的菜单栏滚动条组件的 ref
 */
export const layoutMenuScrollbarRef = ref<ScrollbarInstance>()
