<template>
    <div>
        <div class="ai-go-click-captcha" :class="props.class">
            <div v-if="state.loading" class="loading">{{ i18n.global.t('common.loading') }}...</div>
            <div v-else class="captcha-img-box">
                <img
                    ref="captchaImgRef"
                    class="captcha-img"
                    :src="state.captcha.imageBase64"
                    :alt="i18n.global.t('common.captchaLoadFailed')"
                    @click.prevent="onRecord($event)"
                />
                <span
                    v-for="(item, index) in state.clicks"
                    :key="index"
                    class="step"
                    @click="onCancelRecord(index)"
                    :style="`left:${item.x - 13}px;top:${item.y - 13}px`"
                >
                    {{ index + 1 }}
                </span>
            </div>
            <div class="captcha-prompt" v-if="state.tip">
                {{ state.tip }}
            </div>
            <div v-else class="captcha-prompt">
                {{ i18n.global.t('common.pleaseClick') }}
                <span v-for="(text, index) in state.captcha.elements" :key="index" :class="state.clicks.length > index ? 'clicked' : ''">
                    {{ text }}
                </span>
            </div>
            <div class="captcha-refresh-box">
                <div class="captcha-refresh-line captcha-refresh-line-l"></div>
                <span class="captcha-refresh-btn" :title="i18n.global.t('common.refresh')" @click="load">⟳</span>
                <div class="captcha-refresh-line captcha-refresh-line-r"></div>
            </div>
        </div>
        <div class="ai-go-mask" @click="onClose"></div>
    </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { checkClickCaptcha, getClickCaptcha } from '/@/api/common'
import type { ClickPoint, ClickRequest, Props } from '/@/components/clickCaptcha/index'
import i18n from '/@/lang'
import { SYSTEM_ZINDEX } from '/@/stores/constant/common'

const props = withDefaults(defineProps<Props>(), {
    class: '',
    callback: () => {},
    error: i18n.global.t('common.captchaAreaNotClicked'),
    success: i18n.global.t('common.verificationSuccess'),
})

const captchaImgRef = ref<HTMLImageElement | null>(null)

const state = reactive({
    tip: '',
    loading: true,
    clicks: [] as ClickPoint[],
    captcha: {
        key: '',
        elements: [] as string[],
        imageWidth: 350,
        imageHeight: 200,
        imageBase64: '',
    },
})

const emits = defineEmits<{
    (e: 'destroy'): void
}>()

const load = () => {
    state.loading = true
    getClickCaptcha(props.apiBaseURL)
        .then((res) => {
            state.tip = ''
            state.clicks = []
            state.captcha = res.data.data
        })
        .catch(() => {
            onClose()
        })
        .finally(() => {
            state.loading = false
        })
}

const onRecord = (event: MouseEvent) => {
    if (state.clicks.length < state.captcha.elements.length) {
        state.clicks.push({
            x: event.offsetX,
            y: event.offsetY,
            element: state.captcha.elements[state.clicks.length],
        })
        if (state.clicks.length === state.captcha.elements.length) {
            const data: ClickRequest = {
                key: state.captcha.key,
                clicks: [...state.clicks],
                renderedWidth: captchaImgRef.value!.width,
                renderedHeight: captchaImgRef.value!.height,
            }
            checkClickCaptcha(data, props.apiBaseURL)
                .then(() => {
                    state.tip = props.success
                    setTimeout(() => {
                        props.callback?.(data)
                        onClose()
                    }, 1500)
                })
                .catch(() => {
                    state.tip = props.error
                    setTimeout(() => {
                        load()
                    }, 1500)
                })
        }
    }
}

const onCancelRecord = (index: number) => {
    state.clicks.splice(index, 1)
}

const onClose = () => {
    emits('destroy')
}

load()
</script>

<style scoped lang="scss">
.ai-go-click-captcha {
    padding: 12px;
    border: 1px solid var(--el-border-color-extra-light);
    background-color: var(--el-color-white);
    position: fixed;
    z-index: v-bind('SYSTEM_ZINDEX');
    left: calc(50% - v-bind('state.captcha.imageWidth + 24') / 2 * 1px);
    top: calc(50% - v-bind('state.captcha.imageHeight + 200') / 2 * 1px);
    border-radius: 10px;
    box-shadow:
        0 0 0 1px hsla(0, 0%, 100%, 0.3) inset,
        0 0.5em 1em rgba(0, 0, 0, 0.6);
    .loading {
        color: var(--el-color-info);
        width: 350px;
        text-align: center;
        line-height: 200px;
    }
    .captcha-img-box {
        position: relative;
        .captcha-img {
            width: v-bind('state.captcha.imageWidth') px;
            height: v-bind('state.captcha.imageHeight') px;
            border: none;
            cursor: pointer;
        }
        .step {
            box-sizing: border-box;
            position: absolute;
            width: 20px;
            height: 20px;
            line-height: 20px;
            font-size: var(--el-font-size-small);
            font-weight: bold;
            text-align: center;
            color: var(--el-color-white);
            border: 1px solid var(--el-border-color-extra-light);
            background-color: var(--el-color-primary);
            border-radius: 30px;
            box-shadow: 0 0 10px var(--el-color-white);
            user-select: none;
            cursor: pointer;
        }
    }
    .captcha-prompt {
        height: 40px;
        line-height: 40px;
        font-size: var(--el-font-size-base);
        text-align: center;
        color: var(--el-color-info);
        span {
            margin-left: 10px;
            font-size: var(--el-font-size-medium);
            font-weight: bold;
            color: var(--el-color-error);
            &.clicked {
                color: var(--el-color-primary);
            }
        }
    }
    .captcha-refresh-box {
        position: relative;
        margin-top: 10px;
        .captcha-refresh-line {
            position: absolute;
            top: 16px;
            width: 140px;
            height: 1px;
            background-color: #ccc;
        }
        .captcha-refresh-line-l {
            left: 5px;
        }
        .captcha-refresh-line-r {
            right: 5px;
        }
        .captcha-refresh-btn {
            cursor: pointer;
            display: block;
            margin: 0 auto;
            width: 32px;
            height: 32px;
            font-size: 24px;
            line-height: 32px;
            text-align: center;
            color: var(--el-color-info);
        }
    }
}
</style>
