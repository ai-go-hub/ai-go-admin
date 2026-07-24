<template>
    <div>
        <el-popover :placement="placement" trigger="focus" :hide-after="0" :width="state.selectWidth" :visible="state.popoverVisible">
            <div @mouseover.stop="state.iconSelectMouseover = true" @mouseout.stop="state.iconSelectMouseover = false" class="icon-select">
                <div class="icon-select-box">
                    <div class="select-header">
                        <div class="select-title">{{ title ? title : $t('common.pleaseSelectIcon') }}</div>
                        <div class="select-tab">
                            <span
                                :title="'Element Plus ' + $t('common.icon')"
                                @click="onChangeTab('ele')"
                                :class="state.iconType == 'ele' ? 'active' : ''"
                            >
                                ele
                            </span>
                            <span
                                :title="'Lucide ' + $t('common.icon')"
                                @click="onChangeTab('lucide')"
                                :class="state.iconType == 'lucide' ? 'active' : ''"
                            >
                                lucide
                            </span>
                        </div>
                    </div>
                    <el-scrollbar class="select-body">
                        <div v-if="state.loading" class="icon-select-loading">
                            <Icon name="el-loading" />
                        </div>
                        <div v-else-if="renderFontIconNames.length > 0">
                            <div class="icon-select-item" :title="item" @click="onIcon(item)" v-for="(item, key) in renderFontIconNames" :key="key">
                                <Icon :name="item" />
                            </div>
                        </div>
                    </el-scrollbar>
                </div>
            </div>
            <template #reference>
                <el-input
                    v-model="state.inputValue"
                    :size="size"
                    :disabled="disabled"
                    :placeholder="$t('common.search') + $t('common.icon')"
                    ref="selectInput"
                    @focus="onInputFocus"
                    @blur="onInputBlur"
                    :class="'size-' + size"
                >
                    <template #prepend>
                        <div class="icon-prepend ag-flex-center">
                            <Icon :size="16" :key="'icon' + state.iconKey" :name="state.prependIcon ? state.prependIcon : state.defaultModelValue" />
                            <div v-if="showIconName" class="name" @click="onCopy(state.prependIcon ? state.prependIcon : state.defaultModelValue)">
                                {{ state.prependIcon ? state.prependIcon : state.defaultModelValue }}
                            </div>
                        </div>
                    </template>
                    <template #append>
                        <Icon :size="16" @click="onInputRefresh" name="el-refresh-right" />
                    </template>
                </el-input>
            </template>
        </el-popover>
    </div>
</template>

<script setup lang="ts">
import { useEventListener } from '@vueuse/core'
import { ElMessage, type Placement } from 'element-plus'
import { computed, nextTick, onMounted, reactive, useTemplateRef, watch } from 'vue'
import { getElementPlusIconNames, getLucideIconNames } from '/@/components/agInput/components/iconSelect'
import { copy } from '/@/utils/common'

type IconType = 'ele' | 'lucide'

interface Props {
    size?: 'default' | 'small' | 'large'
    disabled?: boolean
    title?: string
    type?: IconType
    placement?: Placement
    modelValue?: string
    showIconName?: boolean
}
const props = withDefaults(defineProps<Props>(), {
    size: 'default',
    disabled: false,
    title: '',
    type: 'ele',
    placement: 'bottom',
    modelValue: '',
    showIconName: false,
})

const emits = defineEmits<{
    (e: 'update:modelValue', value: string): void
    (e: 'change', value: string): void
}>()

const selectInput = useTemplateRef('selectInput')
const state: {
    iconType: IconType
    selectWidth: number
    popoverVisible: boolean
    inputFocus: boolean
    iconSelectMouseover: boolean
    fontIconNames: string[]
    loading: boolean
    inputValue: string
    prependIcon: string
    defaultModelValue: string
    iconKey: number
} = reactive({
    iconType: props.type,
    selectWidth: 0,
    popoverVisible: false,
    inputFocus: false,
    iconSelectMouseover: false,
    fontIconNames: [],
    loading: false,
    inputValue: '',
    prependIcon: props.modelValue,
    defaultModelValue: props.modelValue || 'lucide-circle-small',
    iconKey: 0, // 给 icon 标签准备个key，以随时使用 h 函数重新生成元素
})

const onInputFocus = () => {
    state.inputFocus = state.popoverVisible = true
}
const onInputBlur = () => {
    state.inputFocus = false
    state.popoverVisible = state.iconSelectMouseover
}
const onInputRefresh = () => {
    state.iconKey++
    state.prependIcon = state.defaultModelValue
    state.inputValue = ''
    emits('update:modelValue', state.defaultModelValue)
    emits('change', state.defaultModelValue)
}

// 图标列表缓存，避免重复加载
const iconListCache: Record<string, string[]> = {}

const loadIcons = (name: IconType) => {
    state.fontIconNames = []
    nextTick(() => {
        if (iconListCache[name]) {
            state.fontIconNames = iconListCache[name]
            return
        }
        state.loading = true
        const loader = name == 'ele' ? getElementPlusIconNames : getLucideIconNames
        loader()
            .then((res) => {
                iconListCache[name] = res
                state.fontIconNames = res
            })
            .catch(() => {
                state.fontIconNames = []
            })
            .finally(() => {
                state.loading = false
            })
    })
}

const onChangeTab = (name: IconType) => {
    state.iconType = name
    loadIcons(name)
}

const onCopy = (name: string) => {
    copy(name).then((success) => {
        if (success) ElMessage.success(name)
    })
}

const onIcon = (icon: string) => {
    state.iconSelectMouseover = state.popoverVisible = false
    state.iconKey++
    state.prependIcon = icon
    state.inputValue = ''
    emits('update:modelValue', icon)
    emits('change', icon)
    nextTick(() => {
        selectInput.value?.blur()
    })
}

const renderFontIconNames = computed(() => {
    if (!state.inputValue) return state.fontIconNames

    const inputValue = state.inputValue.trim().toLowerCase()
    return state.fontIconNames.filter((icon: string) => icon.toLowerCase().includes(inputValue))
})

// 获取 input 的宽度设定 popover 的宽度
const getInputWidth = () => {
    nextTick(() => {
        state.selectWidth = selectInput.value?.$el.offsetWidth < 260 ? 260 : selectInput.value?.$el.offsetWidth
    })
}

const popoverVisible = () => {
    state.popoverVisible = state.inputFocus || state.iconSelectMouseover
}

watch(
    () => props.modelValue,
    () => {
        state.iconKey++
        if (props.modelValue != state.prependIcon) state.defaultModelValue = props.modelValue
        if (props.modelValue == '') state.defaultModelValue = 'el-refresh-right'
        state.prependIcon = props.modelValue
    }
)

/**
 * 1. 图标选择面板一旦显示就监听 document 的点击事件
 * 2. 点击后输入框和面板会失去焦点，面板将自动隐藏
 * 3. 面板隐藏后删除点击事件监听
 */
let removeClickHidePopoverListenerFn = () => {}
watch(
    () => state.popoverVisible,
    () => {
        if (state.popoverVisible) {
            removeClickHidePopoverListenerFn = useEventListener(document, 'click', popoverVisible)
        } else {
            removeClickHidePopoverListenerFn()
        }
    }
)

onMounted(() => {
    getInputWidth()
    loadIcons(props.type)
})
</script>

<style scoped lang="scss">
.size-small {
    height: 24px;
}
.size-large {
    height: 40px;
}
.size-default {
    height: 32px;
}
.icon-prepend {
    .name {
        padding-left: 5px;
    }
}
.select-header {
    display: flex;
    align-items: center;
    margin-bottom: 12px;
}
.select-tab {
    margin-left: auto;
    span {
        padding: 0 5px;
        cursor: pointer;
        user-select: none;
        &.active,
        &:hover {
            color: var(--el-color-primary);
            text-decoration: underline;
        }
    }
}
.select-body {
    height: 250px;
}
.icon-select-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 100%;
    font-size: 24px;
    color: var(--el-color-primary);
    animation: rotate 1s linear infinite;
}
.icon-select-item {
    display: inline-block;
    padding: 10px 10px 6px 10px;
    margin: 3px;
    border: 1px solid var(--ag-border-color);
    border-radius: var(--el-border-radius-base);
    cursor: pointer;
    font-size: 18px;
    .icon {
        height: 18px;
        width: 18px;
    }
    &:hover {
        border: 1px solid var(--el-color-primary);
    }
}
:deep(.el-input-group__prepend) {
    padding: 0 10px;
}
:deep(.el-input-group__append) {
    padding: 0 10px;
}
</style>
