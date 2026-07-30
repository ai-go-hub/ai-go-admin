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
                        :value="getOptionValue(item)"
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
import { computed, getCurrentInstance, nextTick, onMounted, onUnmounted, reactive, ref, toRaw, useAttrs, useTemplateRef, watch } from 'vue'
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
    modelValue: null,
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
    let clearVal = attrs.valueOnClear as string | number | boolean | Function
    if (clearVal === undefined) {
        clearVal = attrs.multiple ? () => [] : () => null
    }
    return typeof clearVal == 'function' ? clearVal() : clearVal
})

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
 * 记录 modelValue 当前应保持的类型（跟随外部传入的类型变化）
 */
const valueKind = ref<'string' | 'number' | 'string[]' | 'number[]' | null>(null)

/**
 * 检测值类型，绑定空值时:0=number,''=string
 */
const detectValueKind = (val: ValueTypes | null) => {
    if (val === null || val === undefined) return
    if (Array.isArray(val)) {
        if (val.length === 0) return
        valueKind.value = typeof val[0] === 'number' ? 'number[]' : 'string[]'
    } else {
        valueKind.value = typeof val === 'number' ? 'number' : 'string'
    }
}
const isArrayKind = computed(() => valueKind.value === 'string[]' || valueKind.value === 'number[]')
const isNumberKind = computed(() => valueKind.value === 'number' || valueKind.value === 'number[]')

/**
 * 获取选项的绑定值，类型跟随 valueKind
 */
const getOptionValue = (item: AnyObj): string | number => {
    return isNumberKind.value ? Number(item[state.primaryKey]) : String(item[state.primaryKey])
}

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
        if (isArrayKind.value && Array.isArray(val)) {
            const valueArr = []
            for (const key in val) {
                const dataKey = getArrayKey(state.options, state.primaryKey, val[key])
                if (dataKey !== false) {
                    valueArr.push(toRaw(state.options[dataKey]))
                }
            }
            emits('row', valueArr)
        } else {
            const dataKey = getArrayKey(state.options, state.primaryKey, val)
            emits('row', dataKey !== false ? toRaw(state.options[dataKey]) : {})
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

    const wheres: WhereGroup[] = []
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
        const kwGroup: WhereGroup = { or: true, wheres: [] }
        for (const key in remoteSearchFields) {
            kwGroup.wheres.push({
                field: remoteSearchFields[key],
                value: state.keywords,
                operator: 'ILIKE',
            })
        }
        wheres.push(kwGroup)
    }

    // 绑定值初始化，使用精确匹配
    if (!isEmpty(initValue)) {
        state.currentPage = 1
        wheres.push({
            wheres: [
                {
                    field: props.pk,
                    value: initValue,
                    operator: 'eq',
                },
            ],
            or: false,
        })
    }

    state.remoteParams.selector = true
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

            // 首次加载时，若 valueKind 未确定（modelValue 为空），以第一个选项的 pk 类型作为默认值类型
            if (valueKind.value === null && opts.length > 0) {
                const firstPK = opts[0][state.primaryKey]
                const baseKind = typeof firstPK === 'number' ? 'number' : 'string'
                valueKind.value = attrs.multiple ? (`${baseKind}[]` as 'number[]' | 'string[]') : baseKind
            }
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
    // 按 valueKind 保持类型一致
    if (isArrayKind.value && Array.isArray(newVal)) {
        state.value = isNumberKind.value ? newVal.map(Number) : newVal.map(String)
    } else if (!isArrayKind.value && (typeof newVal === 'string' || typeof newVal === 'number')) {
        state.value = isNumberKind.value ? Number(newVal) : String(newVal)
    } else {
        state.value = newVal
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

    // 初始化类型 & 值
    detectValueKind(props.modelValue)
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
        // 跟踪外部传入值的类型变化
        detectValueKind(newVal)

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
    position: relative;
}
</style>
