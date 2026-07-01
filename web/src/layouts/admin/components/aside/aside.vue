<template>
    <el-aside v-if="!navTab.state.activeFullScreen" :class="['layout-aside-' + config.layout.mode, config.layout.shrink ? 'shrink' : '']">
        <Logo v-if="config.layout.menuShowTopBar && config.layout.mode != 'LeftSplit'" />

        <MenuVerticalChildren v-if="config.layout.mode == 'Double'" />
        <MenuLeftSplit v-else-if="config.layout.mode == 'LeftSplit'" />
        <MenuVertical v-else />

        <AsideFooterToolbar v-if="['Default', 'Classic', 'Double'].includes(config.layout.mode)" />
    </el-aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AsideFooterToolbar from '/@/layouts/admin/components/aside/toolbar/footer.vue'
import Logo from '/@/layouts/admin/components/logo.vue'
import MenuLeftSplit from '/@/layouts/admin/components/menu/menuLeftSplit.vue'
import MenuVertical from '/@/layouts/admin/components/menu/menuVertical.vue'
import MenuVerticalChildren from '/@/layouts/admin/components/menu/menuVerticalChildren.vue'
import { useConfig } from '/@/stores/config'
import { SYSTEM_ZINDEX } from '/@/stores/constant/common'
import { useNavTab } from '/@/stores/navTab'

defineOptions({
    name: 'layout/aside',
})

const config = useConfig()
const navTab = useNavTab()
const menuWidth = computed(() => config.getMenuWidth())
</script>

<style scoped lang="scss">
.layout-aside-Default,
.layout-aside-LeftSplit,
.layout-aside-Classic,
.layout-aside-Double {
    --el-aside-width: v-bind(menuWidth);
}
.layout-aside-Default:not(.shrink),
.layout-aside-LeftSplit:not(.shrink) {
    background: var(--ag-bg-color-overlay);
    margin: 16px 0 16px 16px;
    height: calc(100% - 32px);
    box-shadow: var(--el-box-shadow-light);
    border-radius: var(--el-border-radius-base);
    overflow: hidden;
    transition: width 0.3s ease;
}
.layout-aside-Default.shrink,
.layout-aside-LeftSplit.shrink,
.layout-aside-Classic,
.layout-aside-Double {
    background: var(--ag-bg-color-overlay);
    margin: 0;
    height: 100%;
    overflow: hidden;
    transition: width 0.3s ease;
}
.shrink {
    position: fixed;
    top: 0;
    left: 0;
    z-index: v-bind('SYSTEM_ZINDEX');
}
</style>
