<template>
    <div class="nav-bar" :class="config.layout.shrink ? 'shrink' : ''">
        <!-- 小屏设备下的展开菜单按钮 -->
        <div v-if="config.layout.shrink && config.layout.menuCollapse" class="unfold">
            <Icon @click="onMenuCollapse" name="lucide-list-indent-increase" :color="config.getColorValue('menuActiveColor')" :size="18" />
        </div>
        <NavTabs v-if="!config.layout.shrink" ref="layoutNavTabsRef" />
        <NavMenu />
    </div>
</template>

<script setup lang="ts">
import NavTabs from '/@/layouts/admin/components/navBar/tabs.vue'
import NavMenu from '/@/layouts/admin/components/navMenu.vue'
import { layoutNavTabsRef } from '/@/stores/ref'
import { useConfig } from '/@/stores/config'
import { showMask } from '/@/utils/mask'

const config = useConfig()

const onMenuCollapse = () => {
    showMask('ag-aside-menu-shade', () => {
        config.setLayoutValue('menuCollapse', true)
    })
    config.setLayoutValue('menuCollapse', false)
}
</script>

<style lang="scss" scoped>
.nav-bar {
    display: flex;
    height: 50px;
    margin: 20px var(--ag-main-space) 0 var(--ag-main-space);
    :deep(.nav-tabs) {
        display: flex;
        height: 100%;
        position: relative;
        .ai-go-nav-tab {
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 0 20px;
            cursor: pointer;
            z-index: 1;
            user-select: none;
            opacity: 0.7;
            color: v-bind('config.getColorValue("headerBarTabColor")');
            .close-icon {
                padding: 2px;
                margin: 2px 0 0 4px;
            }
            .close-icon:hover {
                background: var(--ag-color-primary-light);
                color: var(--el-border-color) !important;
                border-radius: 50%;
            }
            &.active {
                color: v-bind('config.getColorValue("headerBarTabActiveColor")');
            }
            &:hover {
                opacity: 1;
            }
        }
        .nav-tabs-active-box {
            position: absolute;
            height: 40px;
            border-radius: var(--el-border-radius-base);
            background-color: v-bind('config.getColorValue("headerBarTabActiveBackgroundFloating")');
            box-shadow: var(--el-box-shadow-light);
            transition: all 0.2s;
            -webkit-transition: all 0.2s;
        }
    }
}
.nav-bar.shrink {
    width: 100%;
    background-color: v-bind('config.getColorValue("headerBarBackground")');
    margin: 0;
    .unfold {
        align-self: center;
        padding-left: var(--ag-main-space);
    }
}
</style>
