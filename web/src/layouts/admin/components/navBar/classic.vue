<template>
    <div class="nav-bar">
        <div v-if="config.layout.shrink && config.layout.menuCollapse" class="unfold">
            <Icon @click="onMenuCollapse" name="lucide-list-indent-increase" :color="config.getColorValue('menuActiveColor')" :size="18" />
        </div>
        <NavTabs v-if="!config.layout.shrink" ref="layoutNavTabsRef" />
        <NavMenu />
    </div>
</template>

<script setup lang="ts">
import NavMenu from '/@/layouts/admin/components/navMenu.vue'
import NavTabs from '/@/layouts/admin/components/navBar/tabs.vue'
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

<style scoped lang="scss">
.nav-bar {
    display: flex;
    height: 50px;
    width: 100%;
    background-color: v-bind('config.getColorValue("headerBarBackground")');
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
            height: 100%;
            user-select: none;
            color: v-bind('config.getColorValue("headerBarTabColor")');
            transition: all 0.2s;
            -webkit-transition: all 0.2s;
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
                background-color: v-bind('config.getColorValue("headerBarHoverBackground")');
            }
        }
        .nav-tabs-active-box {
            position: absolute;
            height: 50px;
            background-color: v-bind('config.getColorValue("headerBarTabActiveBackground")');
            transition: all 0.2s;
            -webkit-transition: all 0.2s;
        }
    }
}
.unfold {
    align-self: center;
    padding-left: var(--ag-main-space);
}
</style>
