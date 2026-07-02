<template>
    <div class="login-page">
        <!-- ==================== 左侧：角色动画面板 ==================== -->
        <div class="left-panel" :style="{ background: `linear-gradient(135deg, ${primaryColor}e6, ${primaryColor}, ${primaryColor}cc)` }">
            <!-- 站点名称 -->
            <div class="site-name">
                <span>AI GO MALL</span>
            </div>

            <!-- 角色区域 -->
            <div class="characters-area">
                <div class="characters-container">
                    <!-- Purple 角色 -->
                    <div ref="purpleRef" class="char purple" :style="purpleBodyStyle">
                        <div class="eyes" :style="purpleEyesStyle">
                            <div
                                v-for="i in 2"
                                :key="'pe' + i"
                                :ref="(el) => setEyeRef(el as HTMLElement, 'purple', i - 1)"
                                class="eyeball"
                                :class="{ blinking: isPurpleBlinking }"
                                data-max-dist="5"
                            >
                                <div v-if="!isPurpleBlinking" class="pupil-inner" :data-eye-key="'purple-' + (i - 1)" />
                            </div>
                        </div>
                    </div>

                    <!-- Black 角色 -->
                    <div ref="blackRef" class="char black" :style="blackBodyStyle">
                        <div class="eyes" :style="blackEyesStyle">
                            <div
                                v-for="i in 2"
                                :key="'be' + i"
                                :ref="(el) => setEyeRef(el as HTMLElement, 'black', i - 1)"
                                class="eyeball"
                                :class="{ blinking: isBlackBlinking }"
                                data-max-dist="4"
                            >
                                <div v-if="!isBlackBlinking" class="pupil-inner" :data-eye-key="'black-' + (i - 1)" />
                            </div>
                        </div>
                    </div>

                    <!-- Orange 角色 -->
                    <div ref="orangeRef" class="char orange" :style="orangeBodyStyle">
                        <div class="eyes" :style="orangeEyesStyle">
                            <div
                                v-for="i in 2"
                                :key="'oe' + i"
                                :ref="(el) => setEyeRef(el as HTMLElement, 'orange', i - 1)"
                                class="pupil-dot"
                                :data-eye-key="'orange-' + (i - 1)"
                                data-max-dist="5"
                            />
                        </div>
                    </div>

                    <!-- Yellow 角色 -->
                    <div ref="yellowRef" class="char yellow" :style="yellowBodyStyle">
                        <div class="eyes" :style="yellowEyesStyle">
                            <div
                                v-for="i in 2"
                                :key="'ye' + i"
                                :ref="(el) => setEyeRef(el as HTMLElement, 'yellow', i - 1)"
                                class="pupil-dot"
                                :data-eye-key="'yellow-' + (i - 1)"
                                data-max-dist="5"
                            />
                        </div>
                        <div class="mouth" :style="yellowMouthStyle" />
                    </div>
                </div>
            </div>

            <!-- 底部链接 -->
            <div class="footer-links">
                <a href="#">©{{ new Date().getFullYear() }} AI GO MALL</a>
                <a href="http://beian.miit.gov.cn/" target="_blank">渝ICP备8888888号-1</a>
            </div>

            <!-- 装饰 -->
            <div class="deco-grid" />
            <div class="deco-circle deco-circle-1" />
            <div class="deco-circle deco-circle-2" />
        </div>

        <!-- ==================== 右侧：登录表单 ==================== -->
        <div class="right-panel">
            <div class="form-wrapper">
                <div class="form-header">
                    <h1>{{ $t('login.login') }}</h1>
                    <p>{{ $t('login.welcomePrompt', { siteName: 'AI GO MALL' }) }}</p>
                </div>

                <el-form ref="formRef" :model="loginForm" :rules="rules" size="large" @keyup.enter="onSubmit">
                    <el-form-item prop="username">
                        <el-input
                            ref="usernameInputRef"
                            v-model="loginForm.username"
                            :placeholder="t('login.username')"
                            :prefix-icon="User"
                            clearable
                            @focus="isTyping = true"
                            @blur="isTyping = false"
                        />
                    </el-form-item>

                    <el-form-item prop="password">
                        <el-input
                            ref="passwordInputRef"
                            v-model="loginForm.password"
                            type="password"
                            show-password
                            :placeholder="t('login.password')"
                            :prefix-icon="Lock"
                            clearable
                            @focus="isTyping = true"
                            @blur="isTyping = false"
                        />
                    </el-form-item>

                    <div class="options">
                        <label class="remember">
                            <el-checkbox v-model="loginForm.remember" :label="t('login.remember')" />
                        </label>
                    </div>

                    <el-form-item>
                        <el-button size="large" type="primary" :loading="submitLoading" class="submit-btn" @click="onSubmit">
                            {{ submitLoading ? t('login.loggingIn') : t('login.login') }}
                        </el-button>
                    </el-form-item>
                </el-form>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Lock, User } from '@element-plus/icons-vue'
import type { FormInstance, FormRules, InputInstance } from 'element-plus'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { adminLogin } from '/@/api/admin'
import { getCaptchaConfig } from '/@/api/admin/index'
import clickCaptcha from '/@/components/clickCaptcha'
import type { ClickRequest } from '/@/components/clickCaptcha/index'
import { useAdminInfo } from '/@/stores/adminInfo'

const { t } = useI18n()
const router = useRouter()
const adminInfo = useAdminInfo()

// ==================== 表单相关 ====================

interface LoginForm {
    username: string
    password: string
    remember: boolean
}

const loginForm = reactive<LoginForm>({
    username: '',
    password: '',
    remember: false,
})

const formRef = ref<FormInstance>()
const submitLoading = ref(false)
const captchaEnabled = ref(true)
const usernameInputRef = ref<InputInstance>()
const passwordInputRef = ref<InputInstance>()

const rules: FormRules<LoginForm> = {
    username: [{ required: true, message: t('login.usernameRequired'), trigger: 'blur' }],
    password: [
        { required: true, message: t('login.passwordRequired'), trigger: 'blur' },
        { min: 6, message: t('login.passwordMinLength'), trigger: 'blur' },
    ],
}

async function onSubmit() {
    if (!formRef.value) return

    try {
        await formRef.value.validate()
    } catch {
        return
    }

    const doLogin = (captcha?: ClickRequest) => {
        submitLoading.value = true
        adminLogin({ ...loginForm, captcha })
            .then((res) => {
                adminInfo.dataFill(res.data.data, false)
                router.push('/admin')
            })
            .finally(() => {
                submitLoading.value = false
            })
    }

    if (captchaEnabled.value) {
        clickCaptcha((data) => doLogin(data))
    } else {
        doLogin()
    }
}

// ==================== 角色动画状态 ====================

const primaryColor = '#409eff'

// 核心状态
const isTyping = ref(false)
const isPurpleBlinking = ref(false)
const isBlackBlinking = ref(false)
const isLookingAtEachOther = ref(false)
const isPurplePeeking = ref(false)

// 派生状态
const showPassword = computed(() => passwordInputRef.value?.passwordVisible ?? false)
const hasSecret = computed(() => !!loginForm.password)
const hiding = computed(() => hasSecret.value && showPassword.value)
const leaning = computed(() => isTyping.value || (hasSecret.value && !showPassword.value))

// 鼠标位置和角色 DOM refs
const mouseX = ref(0)
const mouseY = ref(0)
const purpleRef = ref<HTMLElement | null>(null)
const blackRef = ref<HTMLElement | null>(null)
const orangeRef = ref<HTMLElement | null>(null)
const yellowRef = ref<HTMLElement | null>(null)

// 角色身体位置状态
const purplePos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })
const blackPos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })
const orangePos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })
const yellowPos = reactive({ faceX: 0, faceY: 0, bodySkew: 0 })

// 眼睛元素引用收集
const eyeElements = new Map<string, HTMLElement>()

function setEyeRef(el: HTMLElement | null, char: string, index: number) {
    const key = `${char}-${index}`
    if (el) {
        eyeElements.set(key, el)
    } else {
        eyeElements.delete(key)
    }
}

// ==================== 鼠标追踪与 rAF 循环 ====================

function calcPos(el: HTMLElement | null, target: { faceX: number; faceY: number; bodySkew: number }) {
    if (!el) return
    const r = el.getBoundingClientRect()
    const dx = mouseX.value - (r.left + r.width / 2)
    const dy = mouseY.value - (r.top + r.height / 3)
    target.faceX = Math.max(-15, Math.min(15, dx / 20))
    target.faceY = Math.max(-10, Math.min(10, dy / 30))
    target.bodySkew = Math.max(-6, Math.min(6, -dx / 120))
}

let rafId = 0
function tick() {
    // 更新角色身体位置
    calcPos(purpleRef.value, purplePos)
    calcPos(blackRef.value, blackPos)
    calcPos(orangeRef.value, orangePos)
    calcPos(yellowRef.value, yellowPos)

    // 更新所有眼睛瞳孔位置
    for (const [key, el] of eyeElements) {
        // 检查是否需要强制看向某个方向
        let forceX: number | undefined
        let forceY: number | undefined

        if (key.startsWith('purple-')) {
            if (hiding.value) {
                forceX = isPurplePeeking.value ? 4 : -4
                forceY = isPurplePeeking.value ? 5 : -4
            } else if (isLookingAtEachOther.value) {
                forceX = 3
                forceY = 4
            }
        } else if (key.startsWith('black-')) {
            if (hiding.value) {
                forceX = -4
                forceY = -4
            } else if (isLookingAtEachOther.value) {
                forceX = 0
                forceY = -4
            }
        } else if (hiding.value) {
            forceX = -5
            forceY = -4
        }

        const maxDist = parseFloat(el.dataset.maxDist || '5')
        let pupilEl: HTMLElement | null

        if (el.classList.contains('eyeball')) {
            pupilEl = el.querySelector('.pupil-inner')
        } else {
            pupilEl = el
        }

        if (!pupilEl) continue

        if (forceX !== undefined && forceY !== undefined) {
            pupilEl.style.transform = `translate(${forceX}px, ${forceY}px)`
        } else {
            const r = el.getBoundingClientRect()
            const dx = mouseX.value - (r.left + r.width / 2)
            const dy = mouseY.value - (r.top + r.height / 2)
            const dist = Math.min(Math.sqrt(dx * dx + dy * dy), maxDist)
            const angle = Math.atan2(dy, dx)
            pupilEl.style.transform = `translate(${Math.cos(angle) * dist}px, ${Math.sin(angle) * dist}px)`
        }
    }

    rafId = requestAnimationFrame(tick)
}

function onMouseMove(e: MouseEvent) {
    mouseX.value = e.clientX
    mouseY.value = e.clientY
}

// ==================== 眨眼逻辑 ====================

function setupBlink(target: { value: boolean }) {
    let timer: number
    const go = () => {
        timer = window.setTimeout(
            () => {
                target.value = true
                window.setTimeout(() => {
                    target.value = false
                    go()
                }, 150)
            },
            Math.random() * 4000 + 3000
        )
    }
    go()
    return () => clearTimeout(timer)
}

let stopPurpleBlink: (() => void) | undefined
let stopBlackBlink: (() => void) | undefined

// ==================== 相互对视 ====================

watch(isTyping, (v) => {
    if (v) {
        isLookingAtEachOther.value = true
        setTimeout(() => {
            isLookingAtEachOther.value = false
        }, 800)
    } else {
        isLookingAtEachOther.value = false
    }
})

// ==================== 紫色角色偷看 ====================

let peekTimer: number | undefined
watch([hasSecret, showPassword], () => {
    clearTimeout(peekTimer)
    if (hasSecret.value && showPassword.value) {
        peekTimer = window.setTimeout(
            () => {
                isPurplePeeking.value = true
                setTimeout(() => {
                    isPurplePeeking.value = false
                }, 800)
            },
            Math.random() * 3000 + 2000
        )
    } else {
        isPurplePeeking.value = false
    }
})

// ==================== 角色样式计算 ====================

const purpleBodyStyle = computed(() => ({
    height: leaning.value ? '440px' : '400px',
    transform: hiding.value
        ? 'skewX(0deg)'
        : leaning.value
          ? `skewX(${purplePos.bodySkew - 12}deg) translateX(40px)`
          : `skewX(${purplePos.bodySkew}deg)`,
}))

const purpleEyesStyle = computed(() => ({
    left: hiding.value ? '20px' : isLookingAtEachOther.value ? '55px' : `${45 + purplePos.faceX}px`,
    top: hiding.value ? '35px' : isLookingAtEachOther.value ? '65px' : `${40 + purplePos.faceY}px`,
    gap: '32px',
}))

const blackBodyStyle = computed(() => ({
    transform: hiding.value
        ? 'skewX(0deg)'
        : isLookingAtEachOther.value
          ? `skewX(${blackPos.bodySkew * 1.5 + 10}deg) translateX(20px)`
          : leaning.value
            ? `skewX(${blackPos.bodySkew * 1.5}deg)`
            : `skewX(${blackPos.bodySkew}deg)`,
}))

const blackEyesStyle = computed(() => ({
    left: hiding.value ? '10px' : isLookingAtEachOther.value ? '32px' : `${26 + blackPos.faceX}px`,
    top: hiding.value ? '28px' : isLookingAtEachOther.value ? '12px' : `${32 + blackPos.faceY}px`,
    gap: '24px',
}))

const orangeBodyStyle = computed(() => ({
    transform: hiding.value ? 'skewX(0deg)' : `skewX(${orangePos.bodySkew}deg)`,
}))

const orangeEyesStyle = computed(() => ({
    left: hiding.value ? '50px' : `${82 + orangePos.faceX}px`,
    top: hiding.value ? '85px' : `${90 + orangePos.faceY}px`,
    gap: '32px',
}))

const yellowBodyStyle = computed(() => ({
    transform: hiding.value ? 'skewX(0deg)' : `skewX(${yellowPos.bodySkew}deg)`,
}))

const yellowEyesStyle = computed(() => ({
    left: hiding.value ? '20px' : `${52 + yellowPos.faceX}px`,
    top: hiding.value ? '35px' : `${40 + yellowPos.faceY}px`,
    gap: '24px',
}))

const yellowMouthStyle = computed(() => ({
    left: hiding.value ? '10px' : `${40 + yellowPos.faceX}px`,
    top: hiding.value ? '88px' : `${88 + yellowPos.faceY}px`,
}))

// ==================== 生命周期 ====================

onMounted(() => {
    // 初始化鼠标位置（角色眼睛方向）为网页中心
    // 浏览器安全限制无法直接获取初始鼠标位置，需等待用户交互
    mouseX.value = window.innerWidth / 2
    mouseY.value = window.innerHeight / 2

    window.addEventListener('mousemove', onMouseMove)
    stopPurpleBlink = setupBlink(isPurpleBlinking)
    stopBlackBlink = setupBlink(isBlackBlinking)
    rafId = requestAnimationFrame(tick)

    // 获取人机验证码开关配置
    getCaptchaConfig()
        .then((res) => {
            captchaEnabled.value = res.data.data.adminLogin
        })
        .catch(() => {
            captchaEnabled.value = true
        })

    nextTick(() => {
        if (!loginForm.username) {
            usernameInputRef.value?.focus()
        } else if (!loginForm.password) {
            passwordInputRef.value?.focus()
        }
    })
})

onBeforeUnmount(() => {
    window.removeEventListener('mousemove', onMouseMove)
    cancelAnimationFrame(rafId)
    stopPurpleBlink?.()
    stopBlackBlink?.()
    clearTimeout(peekTimer)
})
</script>

<style scoped lang="scss">
// ==================== 整体布局 ====================

.login-page {
    display: grid;
    grid-template-columns: 1fr 1fr;
    height: 100vh;
}

// ==================== 左侧面板 ====================

.left-panel {
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 48px;
    color: white;
    overflow: hidden;
}

.site-name {
    position: relative;
    font-size: 18px;
    font-weight: 600;
}

// ==================== 角色容器 ====================

.characters-area {
    position: relative;
    z-index: 20;
    display: flex;
    align-items: flex-end;
    justify-content: center;
    height: 500px;
}

.characters-container {
    position: relative;
    width: 550px;
    height: 400px;
}

// ==================== 角色通用样式 ====================

.char {
    position: absolute;
    bottom: 0;
    transition: all 0.7s ease-in-out;
    transform-origin: bottom center;
}

.eyes {
    position: absolute;
    display: flex;
    transition: all 0.7s ease-in-out;
}

// ==================== Purple 角色 ====================

.purple {
    left: 70px;
    width: 180px;
    background: #6c3ff5;
    border-radius: 10px 10px 0 0;
    z-index: 1;
}

// ==================== Black 角色 ====================

.black {
    left: 240px;
    width: 120px;
    height: 310px;
    background: #2d2d2d;
    border-radius: 8px 8px 0 0;
    z-index: 2;
}

// ==================== Orange 角色 ====================

.orange {
    left: 0;
    width: 240px;
    height: 200px;
    background: #ff9b6b;
    border-radius: 120px 120px 0 0;
    z-index: 3;
}

// ==================== Yellow 角色 ====================

.yellow {
    left: 310px;
    width: 140px;
    height: 230px;
    background: #e8d754;
    border-radius: 70px 70px 0 0;
    z-index: 4;
}

.mouth {
    position: absolute;
    width: 80px;
    height: 4px;
    background: #2d2d2d;
    border-radius: 4px;
    transition: all 0.2s ease-out;
}

// ==================== 眼睛组件 ====================

.eyeball {
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.15s;
    overflow: hidden;
    background: white;

    &.blinking {
        height: 2px !important;
    }
}

.purple .eyeball {
    width: 18px;
    height: 18px;
}

.black .eyeball {
    width: 16px;
    height: 16px;
}

.pupil-inner {
    border-radius: 50%;
    transition: transform 0.1s ease-out;
}

.purple .pupil-inner {
    width: 7px;
    height: 7px;
    background: #2d2d2d;
}

.black .pupil-inner {
    width: 6px;
    height: 6px;
    background: #2d2d2d;
}

.pupil-dot {
    border-radius: 50%;
    transition: transform 0.1s ease-out;
    width: 12px;
    height: 12px;
    background: #2d2d2d;
}

// ==================== 底部链接 ====================

.footer-links {
    position: relative;
    z-index: 20;
    display: flex;
    gap: 16px;
    font-size: 14px;
    align-items: center;

    a {
        color: rgba(255, 255, 255, 0.6);
        transition: color 0.2s;
        text-decoration: none;

        &:hover {
            color: white;
        }
    }
}

// ==================== 装饰元素 ====================

.deco-grid {
    position: absolute;
    inset: 0;
    background-image:
        linear-gradient(to right, rgba(255, 255, 255, 0.05) 1px, transparent 1px),
        linear-gradient(to bottom, rgba(255, 255, 255, 0.05) 1px, transparent 1px);
    background-size: 20px 20px;
}

.deco-circle {
    position: absolute;
    border-radius: 50%;
    filter: blur(48px);
}

.deco-circle-1 {
    top: 25%;
    right: 25%;
    width: 256px;
    height: 256px;
    background: rgba(255, 255, 255, 0.1);
}

.deco-circle-2 {
    bottom: 25%;
    left: 25%;
    width: 384px;
    height: 384px;
    background: rgba(255, 255, 255, 0.05);
}

// ==================== 右侧面板 ====================

.right-panel {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32px;
    background: #fff;
}

.form-wrapper {
    width: 100%;
    max-width: 420px;
}

.form-header {
    text-align: center;
    margin-bottom: 40px;

    h1 {
        font-size: 30px;
        font-weight: 700;
        letter-spacing: -0.5px;
        margin: 0 0 8px;
    }

    p {
        color: #71717a;
        font-size: 14px;
        margin: 0;
    }
}

.options {
    display: flex;
    align-items: center;
    font-size: 14px;
    margin-bottom: 16px;

    .remember {
        display: flex;
        align-items: center;
        color: #d4d4d8;
        cursor: pointer;
        :deep(.el-checkbox) {
            --el-checkbox-text-color: var(--el-text-color-placeholder);
        }
    }
}

.submit-btn {
    width: 100%;
    --el-font-size-base: var(--el-font-size-medium);
    --el-button-size: 48px;
}

// ==================== 响应式 ====================

@media (max-width: 1023px) {
    .login-page {
        grid-template-columns: 1fr;
    }

    .left-panel {
        display: none;
    }
}
</style>
