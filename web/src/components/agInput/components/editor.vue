<!-- 多编辑器共存支持 -->
<!-- 所有编辑器的代码位于 /@/components/agInput/components/editor/ 文件夹，一个文件为一种编辑器，文件名则为编辑器名称 -->
<!-- 开发者也可以直接导入 /@/components/agInput/components/editor/ 中的编辑器原组件直接使用 -->
<!-- 向本组件传递 name（文件名/编辑器名称）自动加载对应的编辑器进行渲染，编辑器自定义属性通过 attrs 传递 -->
<template>
    <div>
        <component v-bind="$attrs" :is="mixins[state.name]" />
    </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import { reactive } from 'vue'

interface Props {
    name?: string
}

const props = withDefaults(defineProps<Props>(), {
    name: 'default',
})

const state = reactive({
    name: props.name,
})

const mixins: Record<string, Component> = {}
const mixinComponents: Record<string, any> = import.meta.glob('./editor/**.vue', { eager: true })
for (const key in mixinComponents) {
    const fileName = key.replace('./editor/', '').replace('.vue', '')
    mixins[fileName] = mixinComponents[key].default

    // 未安装富文本编辑器时，值为 default，安装之后，则值为最后一个编辑器的名称
    if (props.name == 'default' && fileName != 'default') {
        state.name = fileName
    }
}
</script>

<style scoped lang="scss"></style>
