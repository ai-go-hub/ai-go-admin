<template>
    <template v-for="menu in props.menus">
        <template v-if="menu.children && menu.children.length > 0">
            <el-sub-menu @click="onClickSubMenu(menu)" :index="getMenuKey(menu)" :key="getMenuKey(menu)">
                <template #title>
                    <Icon :size="18" :name="menu.meta?.icon ? menu.meta?.icon : config.layout.menuDefaultIcon" />
                    <span>{{ menu.meta?.title ? menu.meta?.title : $t('layouts.untitled') }}</span>
                </template>
                <MenuTree :extends="{ ...props.extends, level: props.extends.level + 1 }" :menus="menu.children" />
            </el-sub-menu>
        </template>
        <template v-else>
            <el-menu-item @click="openMenu(menu)" :index="getMenuKey(menu)" :key="getMenuKey(menu)">
                <Icon :size="18" :name="menu.meta?.icon ? menu.meta?.icon : config.layout.menuDefaultIcon" />
                <span>{{ menu.meta?.title ? menu.meta?.title : $t('layouts.untitled') }}</span>
            </el-menu-item>
        </template>
    </template>
</template>

<script setup lang="ts">
import { ElNotification } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { RouteRecordRaw } from 'vue-router'
import { useConfig } from '/@/stores/config'
import { getFirstMenu, getMenuKey, openMenu } from '/@/utils/router'

const { t } = useI18n()
const config = useConfig()

interface Props {
    menus: RouteRecordRaw[]
    extends?: {
        level: number
        [key: string]: any
    }
}
const props = withDefaults(defineProps<Props>(), {
    menus: () => [],
    extends: () => {
        return {
            level: 1,
        }
    },
})

/**
 * sub-menu-item 被点击 - 用于单栏布局和双栏布局
 * 顶栏菜单：点击时打开第一个菜单
 * 侧边菜单（若有）：点击只展开收缩
 *
 * sub-menu-item 被点击时，也会触发到 menu-item 的点击事件，由 el-menu 内部触发，无法很好的排除，在此检查 level 值
 */
const onClickSubMenu = (menu: RouteRecordRaw) => {
    if (props.extends?.position == 'horizontal' && props.extends.level <= 1 && menu.children?.length) {
        const firstRoute = getFirstMenu(menu.children)
        if (firstRoute) {
            openMenu(firstRoute)
        } else {
            ElNotification({ type: 'error', message: t('layouts.noChildMenu') })
        }
    }
}
</script>

<style scoped lang="scss">
.el-sub-menu,
.el-menu-item {
    color: var(--el-menu-text-color);
    .icon {
        vertical-align: middle;
        margin-right: 5px;
        text-align: center;
        flex-shrink: 0;
    }
    &.is-active {
        background-color: var(--el-menu-active-bg-color);
    }
    &.is-active {
        color: var(--el-menu-active-color);
    }
}
</style>
