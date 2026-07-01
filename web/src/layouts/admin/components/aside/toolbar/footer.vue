<template>
    <div>
        <div
            class="aside-footer-toolbar-wrap"
            :class="[config.layout.menuCollapse ? 'collapse' : '', config.layout.menuToolBarAutoHide ? 'auto-hide' : '']"
        >
            <div class="aside-footer-toolbar">
                <div class="footer-toolbar-item" @click="onMenuCollapse">
                    <Icon
                        v-if="config.layout.menuCollapse"
                        name="lucide-list-indent-increase"
                        :color="config.getColorValue('menuToolBarColor')"
                        :size="18"
                    />
                    <Icon v-else name="lucide-list-indent-decrease" :color="config.getColorValue('menuToolBarColor')" :size="18" />
                </div>
                <div class="footer-toolbar-item toolbar-search" @click="onMenuSearch">
                    <Icon name="lucide-search" :color="config.getColorValue('menuToolBarColor')" :size="18" />
                </div>
            </div>
        </div>

        <MenuSearchDialog v-model="menuSearchDialogVisible" />
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import MenuSearchDialog from '/@/layouts/admin/components/aside/toolbar/menuSearch/dialog.vue'
import { useConfig } from '/@/stores/config'
import { BEFORE_RESIZE_LAYOUT } from '/@/stores/constant/cacheKey'
import { setNavTabsWidth } from '/@/utils/layout'
import { closeMask } from '/@/utils/mask'
import { Session } from '/@/utils/storage'

const config = useConfig()
const menuSearchDialogVisible = ref(false)

const onMenuSearch = function () {
    menuSearchDialogVisible.value = true
}

const onMenuCollapse = function () {
    if (config.layout.shrink && !config.layout.menuCollapse) {
        closeMask()
    }

    config.setLayoutValue('menuCollapse', !config.layout.menuCollapse)

    Session.set(BEFORE_RESIZE_LAYOUT, {
        menuCollapse: config.layout.menuCollapse,
    })

    // 等待侧边栏动画结束后重新计算导航栏宽度
    setTimeout(() => {
        setNavTabsWidth()
    }, 350)
}
</script>

<style scoped lang="scss">
.aside-footer-toolbar-wrap {
    position: relative;
    height: 50px;
    background-color: v-bind('config.getColorValue("menuBackground")');
    .aside-footer-toolbar {
        position: absolute;
        display: flex;
        align-items: center;
        justify-content: space-between;
        height: 50px;
        width: 100%;
        padding: 0 20px;
        transition: all 0.2s ease;
        .footer-toolbar-item {
            padding: 10px;
            border-radius: 50%;
            cursor: pointer;
            &:hover {
                color: v-bind('config.getColorValue("menuToolBarHoverColor")') !important;
                background-color: v-bind('config.getColorValue("menuToolBarHoverBackground")');
            }
        }
    }
    &.collapse {
        height: 100px;
        .aside-footer-toolbar {
            flex-direction: column-reverse;
            padding: 10px 0;
            height: 100px;
        }
    }
    &.auto-hide.collapse {
        .aside-footer-toolbar {
            top: 100px;
        }
    }
    &.auto-hide {
        cursor: pointer;
        .aside-footer-toolbar {
            top: 50px;
        }
        &:hover {
            .aside-footer-toolbar {
                top: 0;
            }
        }
    }
}
</style>
