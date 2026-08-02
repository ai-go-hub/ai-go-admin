<template>
    <div>
        <el-cascader
            clearable
            class="w100"
            :props="cascaderProps"
            :model-value="selectedValueNumberArray"
            :placeholder="placeholder ? placeholder : t('common.pleaseSelect', { field: '' })"
            value-on-clear=""
            @update:model-value="onChange"
        />
    </div>
</template>

<script setup lang="ts">
import { toString } from 'lodash-es'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { getArea } from '/@/api/common'

const { t } = useI18n()

interface Props {
    // 最大显示层级，默认 3（省-市-区）；设为 2 则只到市
    level?: number
    // 选中的地区值，支持数组或逗号分隔字符串
    modelValue?: string | (string | number)[]
    // 占位文本
    placeholder?: string
}

const props = withDefaults(defineProps<Props>(), {
    level: 3,
    placeholder: '',
})

const emit = defineEmits<{
    'update:modelValue': [value: string]
}>()

const leafMaxLevel = computed(() => props.level - 1)

// 绑定值支持数组或逗号分隔字符串，统一转为数组供 el-cascader 内部使用
const selectedValueNumberArray = computed(() => {
    if (!props.modelValue) return []
    if (Array.isArray(props.modelValue)) return props.modelValue.map(Number)
    return props.modelValue.split(',').map((v) => Number(v))
})

// 节点缓存: { [level]: { [key]: Node[] } }
const nodeCache: Record<number, Record<string, any[]>> = {}

// 请求快照
const lastLazy: {
    key: string
    value: string
    nodes: any[]
    pending: Promise<void> | null
} = {
    key: '',
    value: '',
    nodes: [],
    pending: null,
}

const cascaderProps = {
    lazy: true,
    lazyLoad(node: any, resolve: (nodes: any[]) => void) {
        const { level, pathValues } = node
        const key = (pathValues as number[]).join(',') || 'init'

        // 检查本地缓存
        const cached = getCache(level, key)
        if (cached) {
            return resolve(cached)
        }

        // 同一 key 且 modelValue 未变，直接复用
        const stringVal = toString(props.modelValue ?? '')
        if (lastLazy.key === key && lastLazy.value === stringVal) {
            if (lastLazy.pending) {
                return lastLazy.pending.then(() => resolve(lastLazy.nodes))
            }
            return resolve(lastLazy.nodes)
        }

        lastLazy.key = key
        lastLazy.value = stringVal
        lastLazy.pending = getArea(pathValues as number[]).then((res) => {
            const items = res.data.data || []

            const nodes = items.map((item: AnyObj) => ({
                value: item.id,
                label: item.name,
                leaf: level >= leafMaxLevel.value,
            }))

            lastLazy.nodes = nodes
            lastLazy.pending = null
            setCache(level, key, nodes)
            resolve(nodes)
        })
    },
}

function getCache(level: number, key: string): any[] | null {
    return nodeCache[level]?.[key] ?? null
}

function setCache(level: number, key: string, nodes: any[]) {
    if (!nodeCache[level]) {
        nodeCache[level] = {}
    }
    nodeCache[level][key] = nodes
}

function onChange(val: any) {
    emit('update:modelValue', toString(val))
}
</script>

<style scoped lang="scss"></style>
