<script lang="ts">
import { createVNode, defineComponent, h, resolveComponent } from 'vue'
import { getLucideComponent } from './index'

export default defineComponent({
    name: 'Icon',
    props: {
        name: {
            type: String,
            required: true,
        },
        size: {
            type: [Number, String],
            default: 24,
        },
        color: {
            type: String,
            default: undefined,
        },
        strokeWidth: {
            type: [Number, String],
            default: undefined,
        },
    },
    setup(props, { attrs }) {
        return () => {
            // Element Plus 的 icon，通过 registerIcons 全局批量注册为 `el-icon-name` 的组件
            if (props.name.indexOf('el-') === 0) {
                const name = props.name.replace('el-', 'el-icon-')
                return createVNode(
                    resolveComponent('el-icon'),
                    { class: 'icon ai-go-icon', ...props, ...attrs },
                    { default: () => h(resolveComponent(name)) }
                )
            }

            // lucide 的 icon，lucideIconSplitPlugin 将按首字母将全部 lucide icon 分为虚拟包，并按需加载
            if (props.name.indexOf('lucide-') === 0) {
                const name = props.name.replace('lucide-', '')
                const component = getLucideComponent(name)
                if (component) {
                    return h(component, { class: 'icon ai-go-icon', ...props, ...attrs })
                }
            }
        }
    },
})
</script>
