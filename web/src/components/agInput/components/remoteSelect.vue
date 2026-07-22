<template>
    <div class="w100">
        <!-- el-select 的远程下拉只在有搜索词时，才会加载数据（显示出 option 列表），-->
        <!-- 所以使用 el-popover 在无数据/无搜索词时，显示一个无数据的提醒 -->
        <el-popover
            width="100%"
            placement="bottom"
            popper-class="remote-select-popper"
            :visible="state.isFocused && !state.loading && !state.keywords && !state.options.length"
            :teleported="false"
            :content="$t('common.noData')"
            :hide-after="0"
        >
            <template #reference>
                <el-select
                    ref="selectRef"
                    class="w100"
                    remote
                    clearable
                    filterable
                    automatic-dropdown
                    remote-show-suffix
                    v-model="state.value"
                    :loading="state.loading"
                    :disabled="props.disabled || !state.isInitialized"
                    @blur="onBlur"
                    @focus="onFocus"
                    @clear="onClear"
                    @change="onChangeSelect"
                    @keydown.esc.capture="onKeyDownEsc"
                    :remote-method="onRemoteMethod"
                    v-bind="$attrs"
                >
                    <el-option
                        class="remote-select-option"
                        v-for="item in state.options"
                        :label="item[field]"
                        :value="item[state.primaryKey].toString()"
                        :key="item[state.primaryKey]"
                    >
                        <el-tooltip placement="right" effect="light" v-if="!isEmpty(tooltipParams)">
                            <template #content>
                                <p v-for="(pItem, pKey) in tooltipParams" :key="pKey">{{ pItem }}: {{ item[pKey] }}</p>
                            </template>
                            <div>{{ item[field] }}</div>
                        </el-tooltip>
                    </el-option>
                    <template v-if="state.total && props.pagination" #footer>
                        <el-pagination class="select-pagination" @current-change="onSelectCurrentPageChange" v-bind="getPaginationAttr()" />
                    </template>
                </el-select>
            </template>
        </el-popover>
    </div>
</template>

<script lang="ts" setup>
import type { ElSelect, PaginationProps } from 'element-plus'
import { debounce, isEmpty } from 'lodash-es'
import { computed, getCurrentInstance, nextTick, onMounted, onUnmounted, reactive, toRaw, useAttrs, useTemplateRef, watch } from 'vue'
import { selectList } from '/@/api/common'
import { useConfig } from '/@/stores/config'
import { getArrayKey } from '/@/utils/common'
import { uuid } from '/@/utils/random'

const attrs = useAttrs()
const config = useConfig()
const selectRef = useTemplateRef('selectRef')
type ElSelectProps = Omit<Partial<InstanceType<typeof ElSelect>['$props']>, 'modelValue'>
type ValueTypes = string | number | string[] | number[]

interface Props extends /* @vue-ignore */ ElSelectProps {
    pk?: string
    field?: string
    remoteUrl: string
    remoteParams?: AnyObj
    remoteSearchFields?: string[]
    modelValue: ValueTypes | null
    pagination?: boolean | PaginationProps
    tooltipParams?: Record<string, string>
    labelFormatter?: (optionData: AnyObj, optionKey: string) => string
    // 按下 ESC 键时直接使下拉框脱焦（默认是清理搜索词或关闭下拉面板，并且不会脱焦，造成 dialog 的按下 ESC 关闭失效）
    escBlur?: boolean
}
const props = withDefaults(defineProps<Props>(), {
    pk: 'id',
    field: 'name',
    remoteUrl: '',
    remoteParams: () => {
        return {}
    },
    remoteSearchFields: () => {
        return []
    },
    modelValue: '',
    tooltipParams: () => {
        return {}
    },
    pagination: true,
    disabled: false,
    escBlur: true,
})

/**
 * 点击清空按钮后的值，同时也是缺省值‌
 */
const valueOnClear = computed(() => {
    let valueOnClear = attrs.valueOnClear as string | number | boolean | Function
    if (valueOnClear === undefined) {
        valueOnClear = attrs.multiple ? () => [] : () => null
    }
    return typeof valueOnClear == 'function' ? valueOnClear() : valueOnClear
})

/**
 * 被认为是空值的值列表
 */
const emptyValues = computed<any>(() => attrs.emptyValues || [null, undefined, ''])

const state: {
    // 主表字段名（不带表别名）
    primaryKey: string
    options: AnyObj[]
    loading: boolean
    total: number
    currentPage: number
    pageSize: number
    remoteParams: AnyObj
    keywords: string
    value: ValueTypes
    isInitialized: boolean
    optionsValid: boolean
    isFocused: boolean
} = reactive({
    primaryKey: props.pk,
    options: [],
    loading: false,
    total: 0,
    currentPage: props.remoteParams.page || 1,
    pageSize: props.remoteParams.limit || 10,
    remoteParams: props.remoteParams,
    keywords: '',
    value: valueOnClear.value,
    isInitialized: false,
    optionsValid: false,
    isFocused: false,
})

let io: IntersectionObserver | null = null
const instance = getCurrentInstance()

const emits = defineEmits<{
    (e: 'update:modelValue', value: ValueTypes): void
    (e: 'row', value: any): void
}>()

/**
 * 获取分页组件属性
 */
const getPaginationAttr = (): Partial<PaginationProps> => {
    const defaultPaginationAttr: Partial<PaginationProps> = {
        pagerCount: 5,
        total: state.total,
        pageSize: state.pageSize,
        currentPage: state.currentPage,
        layout: 'total, ->, prev, pager, next',
        size: config.layout.shrink ? 'small' : 'default',
    }

    if (typeof props.pagination === 'boolean') {
        return defaultPaginationAttr
    }

    return { ...defaultPaginationAttr, ...props.pagination }
}

const onChangeSelect = (val: ValueTypes) => {
    val = updateValue(val)
    if (typeof instance?.vnode.props?.onRow == 'function') {
        if (typeof val == 'number' || typeof val == 'string') {
            const dataKey = getArrayKey(state.options, state.primaryKey, '' + val)
            emits('row', dataKey !== false ? toRaw(state.options[dataKey]) : {})
        } else {
            const valueArr = []
            for (const key in val) {
                const dataKey = getArrayKey(state.options, state.primaryKey, '' + val[key])
                if (dataKey !== false) {
                    valueArr.push(toRaw(state.options[dataKey]))
                }
            }
            emits('row', valueArr)
        }
    }
}

const onKeyDownEsc = (e: KeyboardEvent) => {
    if (props.escBlur) {
        e.stopPropagation()
        selectRef.value?.blur()
    }
}

const onFocus = () => {
    state.isFocused = true
    if (!state.optionsValid) {
        getData()
    }
}

const onClear = () => {
    // 点击清理按钮后，内部 input 呈聚焦状态，但选项面板不会展开，特此处理（element-plus@2.9.1）
    nextTick(() => {
        selectRef.value?.blur()
        selectRef.value?.focus()
    })
}

const onBlur = () => {
    state.keywords = ''
    state.isFocused = false
}

const onRemoteMethod = (q: string) => {
    if (state.keywords != q) {
        state.keywords = q
        state.currentPage = 1
        getData()
    }
}

const getData = debounce((initValue: ValueTypes = '') => {
    state.loading = true

    const wheres: ComSearchData[] = []
    let remoteSearchFields: string[] = []

    if (props.remoteSearchFields.length) {
        remoteSearchFields = props.remoteSearchFields
    } else {
        if (props.id) {
            remoteSearchFields.push(props.id)
        }
        if (props.field) {
            remoteSearchFields.push(props.field)
        }
    }

    // 关键词搜索
    if (state.keywords && remoteSearchFields.length) {
        for (const key in remoteSearchFields) {
            wheres.push({
                field: remoteSearchFields[key],
                value: state.keywords,
                operator: 'ILIKE',
            })
        }
    }

    // 绑定值初始化
    if (initValue) {
        state.currentPage = 1
        wheres.push({
            field: props.pk,
            value: initValue,
            operator: 'eq',
        })
    }

    state.remoteParams.wheres = wheres
    state.remoteParams.page = state.currentPage

    selectList(props.remoteUrl, state.remoteParams)
        .then((res) => {
            let opts = res.data.data.list
            if (typeof props.labelFormatter === 'function') {
                for (const key in opts) {
                    opts[key][props.field] = props.labelFormatter(opts[key], key)
                }
            }
            state.options = opts
            state.total = res.data.data.total ?? 0
            state.optionsValid = state.keywords || (typeof initValue === 'object' ? !isEmpty(initValue) : initValue) ? false : true
        })
        .finally(() => {
            state.loading = false
            state.isInitialized = true
        })
}, 100)

const onSelectCurrentPageChange = (val: number) => {
    state.currentPage = val
    getData()
}

const updateValue = (newVal: any) => {
    if (emptyValues.value.includes(newVal)) {
        state.value = valueOnClear.value
    } else {
        state.value = newVal

        // number[] 转 string[] 确保默认值能够选中
        if (typeof state.value === 'object') {
            for (const key in state.value) {
                state.value[key] = '' + state.value[key]
            }
        } else if (typeof state.value === 'number') {
            state.value = '' + state.value
        }
    }
    emits('update:modelValue', state.value)
    return state.value
}

onMounted(() => {
    // 避免两个远程下拉组件共存时，可能带来的重复请求自动取消
    state.remoteParams.uuid = uuid()

    // 去除主键中的表名
    let pkArr = props.pk.split('.')
    state.primaryKey = pkArr[pkArr.length - 1]

    // 初始化值
    updateValue(props.modelValue)
    getData(state.value)

    setTimeout(() => {
        if (window?.IntersectionObserver) {
            io = new IntersectionObserver((entries) => {
                for (const key in entries) {
                    if (!entries[key].isIntersecting) selectRef.value?.blur()
                }
            })
            if (selectRef.value?.$el instanceof Element) {
                io.observe(selectRef.value.$el)
            }
        }
    }, 500)
})

onUnmounted(() => {
    io?.disconnect()
})

watch(
    () => props.modelValue,
    (newVal) => {
        /**
         * 防止 number 到 string 的类型转换触发默认值多次初始化
         * 相当于忽略数据类型进行比较 [1, 2] == ['1', '2']
         */
        if (getString(state.value) != getString(newVal)) {
            updateValue(newVal)
            getData(state.value)
        }
    }
)

const getString = (val: ValueTypes | null) => {
    // 确保 [] 和 '' 的返回值不一样
    return `${typeof val}:${String(val)}`
}

const blur = () => {
    selectRef.value?.blur()
}

const focus = () => {
    selectRef.value?.focus()
}

const getElSelectRef = () => {
    return selectRef.value
}

defineExpose({
    blur,
    focus,
    getElSelectRef,
})
</script>

<style scoped lang="scss">
:deep(.remote-select-popper) {
    color: var(--el-text-color-secondary);
    font-size: 12px;
    text-align: center;
}
.remote-select-option {
    white-space: pre;
}
.w100 {
    width: 100%;
    position: relative;
}
</style>
