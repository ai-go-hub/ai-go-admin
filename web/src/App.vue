<template>
    <el-config-provider :value-on-clear="() => null" :locale="elLocale">
        <router-view></router-view>
    </el-config-provider>
</template>

<script setup lang="ts">
import { useConfig } from '@/stores/config'
import { init as viteInit } from '@/utils/vite'
import elEn from 'element-plus/es/locale/lang/en'
import elZhCn from 'element-plus/es/locale/lang/zh-cn'
import { computed, onMounted } from 'vue'

const config = useConfig()

const elLocales: Record<string, typeof elEn> = {
    en: elEn,
    'zh-cn': elZhCn,
}

const elLocale = computed(() => elLocales[config.lang.active] || elLocales[config.lang.fallback])

onMounted(() => {
    viteInit()
})
</script>
