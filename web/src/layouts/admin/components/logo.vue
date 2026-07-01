<template>
    <div class="layout-title-bar" :class="config.layout.menuTopBarCenter ? 'layout-title-bar-center' : ''">
        <img
            v-if="config.layout.menuCollapse || config.layout.menuTopBarLogo"
            :class="config.layout.menuCollapse ? 'collapse' : ''"
            src="/@/assets/logo.png"
            class="logo-img"
            alt="logo"
        />
        <div v-if="!config.layout.menuCollapse" class="website-name">{{ config.site.name }}</div>
    </div>
</template>

<script setup lang="ts">
import { useConfig } from '/@/stores/config'

const config = useConfig()
</script>

<style scoped lang="scss">
.layout-title-bar {
    width: 100%;
    display: flex;
    align-items: center;
    box-sizing: border-box;
    padding: 10px 20px;
    height: v-bind('config.layout.mode == "Double" ?  "60px":"50px"');
    background: v-bind('config.layout.mode == "Streamline" ?  "transparent":config.getColorValue("menuTopBarBackground")');
}
.layout-title-bar-center {
    justify-content: center;
}
.logo-img {
    width: 26px;
    &.collapse {
        margin: 0 auto;
    }
}
.website-name {
    display: block;
    padding: 0 2px;
    font-size: var(--el-font-size-extra-large);
    font-weight: 600;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    color: v-bind('config.getColorValue("menuTopBarColor")');
}
.logo-img + .website-name {
    margin-left: 6px;
}
</style>
