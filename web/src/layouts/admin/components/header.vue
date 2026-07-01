<template>
    <el-header v-if="!navTab.state.activeFullScreen" class="layout-header">
        <component :is="config.layout.mode + 'NavBar'"></component>
    </el-header>
</template>

<script setup lang="ts">
import StreamlineNavBar from '/@/layouts/admin/components/menu/menuHorizontal.vue'
import ClassicNavBar from '/@/layouts/admin/components/navBar/classic.vue'
import DefaultNavBar from '/@/layouts/admin/components/navBar/default.vue'
import DoubleNavBar from '/@/layouts/admin/components/navBar/double.vue'
import LeftSplitNavBar from '/@/layouts/admin/components/navBar/leftSplit.vue'
import { useConfig } from '/@/stores/config'
import { useNavTab } from '/@/stores/navTab'

defineOptions({
    name: 'layout/header',
    components: { DefaultNavBar, ClassicNavBar, LeftSplitNavBar, StreamlineNavBar, DoubleNavBar },
})

const config = useConfig()
const navTab = useNavTab()
</script>

<style scoped lang="scss">
.layout-header {
    height: auto;
    padding: 0;
}
</style>
