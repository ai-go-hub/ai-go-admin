<template>
    <div class="iframe-main" v-loading="state.loading">
        <iframe :src="state.iframeSrc" :style="iframeStyle(35)" height="100%" width="100%" id="iframe" @load="hideLoading"></iframe>
    </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { mainHeight as iframeStyle } from '/@/utils/layout'

const router = useRouter()

const state = reactive({
    loading: true,
    iframeSrc: router.currentRoute.value.meta.url as string,
})

const hideLoading = () => {
    state.loading = false
}
</script>

<style scoped lang="scss">
.iframe-main {
    margin: var(--ag-main-space);
    iframe {
        border: none;
    }
}
</style>
