<script lang="ts">
import { inputTypes } from '@/components/agInput'
import agUpload from '@/components/agInput/components/agUpload.vue'
import AreaSelect from '@/components/agInput/components/areaSelect.vue'
import Array from '@/components/agInput/components/array.vue'
import Editor from '@/components/agInput/components/editor.vue'
import IconSelect from '@/components/agInput/components/iconSelect.vue'
import RemoteSelect from '@/components/agInput/components/remoteSelect.vue'
import type { InputAttr, ModelValueTypes } from '@/components/agInput/index'
import { isArray, isString } from 'lodash-es'
import type { PropType, VNode } from 'vue'
import { computed, createVNode, defineComponent, resolveComponent } from 'vue'

export default defineComponent({
    name: 'agInput',
    props: {
        // 输入框类型，支持的输入框见 inputTypes 定义
        type: {
            type: String,
            required: true,
            validator: (value: string) => {
                return inputTypes.includes(value)
            },
        },
        // 双向绑定值
        modelValue: {
            type: null,
            required: true,
        },
        // 输入框的附加属性
        attr: {
            type: Object as PropType<InputAttr>,
            default: () => {},
        },
    },
    emits: ['update:modelValue'],
    setup(props, { emit, slots }) {
        // 可以在此拼装默认 attr、对传入 attr 做额外处理等预留
        const attrs = computed(() => {
            return props.attr
        })

        // 通用值更新函数
        const onValueUpdate = (value: ModelValueTypes) => {
            emit('update:modelValue', value)
        }

        // 基础用法 string textarea password
        const bases = () => {
            return () =>
                createVNode(
                    resolveComponent('el-input'),
                    {
                        type: props.type == 'string' ? 'text' : props.type,
                        ...attrs.value,
                        modelValue: props.modelValue,
                        'onUpdate:modelValue': onValueUpdate,
                    },
                    slots
                )
        }
        // radio checkbox
        const rc = () => {
            if (!attrs.value.dict) {
                console.warn('请传递 ' + props.type + ' 的 dict')
            }

            const vNodes = computed(() => {
                const vNode: VNode[] = []
                const dictIsArray = isArray(attrs.value.dict)
                const type = attrs.value.button ? props.type + '-button' : props.type
                for (const key in attrs.value.dict) {
                    let nodeProps = {}
                    if (dictIsArray) {
                        if (typeof attrs.value.dict[key].value === 'number') {
                            console.warn(props.type + ' 的 dict.value 不能是数字')
                        }

                        nodeProps = {
                            ...attrs.value.dict[key],
                            border: attrs.value.border ? attrs.value.border : false,
                            ...(attrs.value.subAttr || {}),
                        }
                    } else {
                        nodeProps = {
                            value: key,
                            label: attrs.value.dict[key],
                            border: attrs.value.border ? attrs.value.border : false,
                            ...(attrs.value.subAttr || {}),
                        }
                    }
                    vNode.push(createVNode(resolveComponent('el-' + type), nodeProps, slots))
                }
                return vNode
            })

            return () => {
                const valueComputed = computed(() => {
                    if (props.type == 'radio') {
                        if (props.modelValue === undefined || props.modelValue === null) return ''
                        return '' + props.modelValue
                    } else {
                        let modelValueArr: AnyObj = []
                        for (const key in props.modelValue) {
                            modelValueArr[key] = '' + props.modelValue[key]
                        }
                        return modelValueArr
                    }
                })
                return createVNode(
                    resolveComponent('el-' + props.type + '-group'),
                    {
                        ...attrs.value,
                        modelValue: valueComputed.value,
                        'onUpdate:modelValue': onValueUpdate,
                    },
                    () => vNodes.value
                )
            }
        }
        // select selects
        const select = () => {
            if (!attrs.value.dict) {
                console.warn('请传递 ' + props.type + '的 dict')
            }

            const vNodes = computed(() => {
                const vNode: VNode[] = []
                for (const key in attrs.value.dict) {
                    vNode.push(
                        createVNode(
                            resolveComponent('el-option'),
                            {
                                key: key,
                                label: attrs.value.dict[key],
                                value: key,
                                ...(attrs.value.subAttr || {}),
                            },
                            slots
                        )
                    )
                }
                return vNode
            })

            return () => {
                const valueComputed = computed(() => {
                    if (props.type == 'select') {
                        if (props.modelValue === undefined || props.modelValue === null) return ''
                        return '' + props.modelValue
                    } else {
                        let modelValueArr: AnyObj = []
                        for (const key in props.modelValue) {
                            modelValueArr[key] = '' + props.modelValue[key]
                        }
                        return modelValueArr
                    }
                })
                return createVNode(
                    resolveComponent('el-select'),
                    {
                        class: 'w100',
                        multiple: props.type == 'select' ? false : true,
                        clearable: true,
                        ...attrs.value,
                        modelValue: valueComputed.value,
                        'onUpdate:modelValue': onValueUpdate,
                    },
                    () => vNodes.value
                )
            }
        }
        // datetime
        const datetime = () => {
            const valueComputed = computed(() => (!props.modelValue ? null : '' + props.modelValue))
            let valueFormat = 'YYYY-MM-DD HH:mm:ss'
            switch (props.type) {
                case 'date':
                    valueFormat = 'YYYY-MM-DD'
                    break
                case 'year':
                    valueFormat = 'YYYY'
                    break
            }
            return () =>
                createVNode(
                    resolveComponent('el-date-picker'),
                    {
                        class: 'w100',
                        type: props.type,
                        'value-format': valueFormat,
                        ...attrs.value,
                        modelValue: valueComputed.value,
                        'onUpdate:modelValue': onValueUpdate,
                    },
                    slots
                )
        }
        // upload
        const upload = () => {
            return () =>
                createVNode(
                    agUpload,
                    {
                        type: props.type,
                        modelValue: props.modelValue,
                        'onUpdate:modelValue': onValueUpdate,
                        ...attrs.value,
                    },
                    slots
                )
        }

        // remoteSelect remoteSelects
        const remoteSelect = () => {
            return () =>
                createVNode(
                    RemoteSelect,
                    {
                        modelValue: props.modelValue,
                        'onUpdate:modelValue': onValueUpdate,
                        multiple: props.type == 'remoteSelect' ? false : true,
                        ...attrs.value,
                    },
                    slots
                )
        }

        const buildFun = new Map([
            ['text', bases],
            ['string', bases],
            [
                'number',
                () => {
                    return () =>
                        createVNode(
                            resolveComponent('el-input-number'),
                            {
                                class: 'w100',
                                'controls-position': 'right',
                                ...attrs.value,
                                modelValue: isString(props.modelValue) ? Number(props.modelValue) : props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                            },
                            slots
                        )
                },
            ],
            ['textarea', bases],
            ['password', bases],
            ['radio', rc],
            ['checkbox', rc],
            [
                'switch',
                () => {
                    // 值类型:string,number,boolean,custom
                    const valueType = computed(() => {
                        if (attrs.value && typeof attrs.value.activeValue !== 'undefined' && typeof attrs.value.inactiveValue !== 'undefined') {
                            return 'custom'
                        }
                        return typeof props.modelValue
                    })

                    // 要传递给 el-switch 组件的绑定值，该组件对传入值有限制，先做处理
                    const valueComputed = computed(() => {
                        if (valueType.value === 'boolean' || valueType.value === 'custom') {
                            return props.modelValue
                        } else {
                            let valueTmp = parseInt(props.modelValue as string)
                            return isNaN(valueTmp) || valueTmp <= 0 ? false : true
                        }
                    })
                    return () =>
                        createVNode(
                            resolveComponent('el-switch'),
                            {
                                ...attrs.value,
                                modelValue: valueComputed.value,
                                'onUpdate:modelValue': (value: boolean) => {
                                    let newValue: boolean | string | number = value
                                    switch (valueType.value) {
                                        case 'string':
                                            newValue = value ? '1' : '0'
                                            break
                                        case 'number':
                                            newValue = value ? 1 : 0
                                    }
                                    emit('update:modelValue', newValue)
                                },
                            },
                            slots
                        )
                },
            ],
            ['datetime', datetime],
            ['year', datetime],
            ['date', datetime],
            [
                'time',
                () => {
                    return () =>
                        createVNode(
                            resolveComponent('el-time-picker'),
                            {
                                class: 'w100',
                                clearable: true,
                                format: 'HH:mm:ss',
                                valueFormat: 'HH:mm:ss',
                                ...attrs.value,
                                modelValue: props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                            },
                            slots
                        )
                },
            ],
            ['select', select],
            ['selects', select],
            [
                'array',
                () => {
                    return () =>
                        createVNode(
                            Array,
                            {
                                modelValue: props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                                ...attrs.value,
                            },
                            slots
                        )
                },
            ],
            ['remoteSelect', remoteSelect],
            ['remoteSelects', remoteSelect],
            [
                'areaSelect',
                () => {
                    return () =>
                        createVNode(
                            AreaSelect,
                            {
                                modelValue: props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                                ...attrs.value,
                            },
                            slots
                        )
                },
            ],
            ['image', upload],
            ['images', upload],
            ['file', upload],
            ['files', upload],
            [
                'iconSelect',
                () => {
                    return () =>
                        createVNode(
                            IconSelect,
                            {
                                modelValue: props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                                ...attrs.value,
                            },
                            slots
                        )
                },
            ],
            [
                'color',
                () => {
                    return () =>
                        createVNode(
                            resolveComponent('el-color-picker'),
                            {
                                modelValue: props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                                ...attrs.value,
                            },
                            slots
                        )
                },
            ],
            [
                'editor',
                () => {
                    return () =>
                        createVNode(
                            Editor,
                            {
                                class: 'w100',
                                modelValue: props.modelValue,
                                'onUpdate:modelValue': onValueUpdate,
                                ...attrs.value,
                            },
                            slots
                        )
                },
            ],
            [
                'default',
                () => {
                    console.warn('暂不支持' + props.type + '的输入框类型，你可以自行在 agInput 组件内添加逻辑')
                },
            ],
        ])

        let action = buildFun.get(props.type) || buildFun.get('default')
        return action!.call(this)
    },
})
</script>

<style scoped lang="scss">
.ag-upload-image :deep(.el-upload--picture-card) {
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.ag-upload-file :deep(.el-upload-list) {
    margin-left: -10px;
}
</style>
