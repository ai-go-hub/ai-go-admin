<template>
    <el-container class="layout-container">
        <Aside />
        <el-container class="content-wrapper">
            <Header />
            <Main />
        </el-container>
        <CloseFullScreen v-if="navTab.state.activeFullScreen" />

        <el-tour v-model="config.layout.showTour" :gap="{ offset: 0, radius: 2 }" @close="onTourClose('tourUnfinished')">
            <el-tour-step
                placement="bottom"
                target=".nav-tabs .ai-go-nav-tab.active"
                :title="t('layouts.contextmenu')"
                :description="t('layouts.contextmenuTips')"
            />
            <el-tour-step
                placement="left-end"
                target=".ai-go-layout-config-btn"
                :title="t('layouts.layoutConfiguration')"
                :description="t('layouts.layoutConfigurationTips')"
            />
            <el-tour-step
                placement="right-start"
                target=".aside-footer-toolbar-wrap .toolbar-search"
                :title="t('layouts.menuSearch')"
                :description="t('layouts.menuSearchTips')"
            />
        </el-tour>
    </el-container>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import Aside from '/@/layouts/admin/components/aside/aside.vue'
import CloseFullScreen from '/@/layouts/admin/components/closeFullScreen.vue'
import Header from '/@/layouts/admin/components/header.vue'
import Main from '/@/layouts/admin/router-view/main.vue'
import { useConfig } from '/@/stores/config'
import { Layout } from '/@/stores/interface/config'
import { useNavTab } from '/@/stores/navTab'

const { t } = useI18n()
const config = useConfig()
const navTab = useNavTab()

const onTourClose = (key: keyof Layout) => {
    config.setLayoutValue(key, false)
}
</script>

<style scoped>
.layout-container {
    height: 100%;
    width: 100%;
}
.content-wrapper {
    flex-direction: column;
    width: 100%;
    height: 100%;
}
</style>
