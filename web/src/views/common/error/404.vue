<template>
    <div class="error-page">
        <!-- ==================== 背景装饰 ==================== -->
        <div class="error-bg-deco" aria-hidden="true">
            <span class="deco-blob deco-blob-1" />
            <span class="deco-blob deco-blob-2" />
            <span class="deco-blob deco-blob-3" />
            <span class="deco-dot deco-dot-1" />
            <span class="deco-dot deco-dot-2" />
            <span class="deco-dot deco-dot-3" />
            <span class="deco-dot deco-dot-4" />
        </div>

        <div class="error-card">
            <!-- ==================== 404 大数字 + 装饰环 ==================== -->
            <div class="error-hero">
                <span class="deco-ring deco-ring-outer" />
                <span class="deco-ring deco-ring-inner" />
                <span class="error-code">404</span>
            </div>

            <!-- ==================== 文案 ==================== -->
            <h1 class="error-title">{{ $t('pageTitles.NotFound') }}</h1>
            <p class="error-desc">
                {{ $t('common.error.404.desc1') }}
                <br />
                {{ $t('common.error.404.desc2') }}
            </p>

            <!-- ==================== 操作 ==================== -->
            <div class="error-actions">
                <el-button type="primary" size="large" round @click="goHome">
                    <Icon name="lucide-house" :size="16" class="btn-icon" />
                    <span>{{ $t('common.error.404.backHome') }}</span>
                </el-button>
                <el-button size="large" round @click="goBack">
                    <Icon name="el-back" :size="16" class="btn-icon" />
                    <span>{{ $t('common.error.404.goBack') }}</span>
                </el-button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
    router.push('/')
}

/**
 * 返回上一页；无历史记录时回首页
 */
const goBack = () => {
    if (window.history.length > 1) {
        router.back()
    } else {
        goHome()
    }
}
</script>

<style scoped lang="scss">
.error-page {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
    padding: 24px;
    overflow: hidden;
    background: radial-gradient(circle at 50% 0%, var(--el-color-primary-light-9), transparent 60%), var(--ag-bg-color);
}

.error-bg-deco {
    position: absolute;
    inset: 0;
    pointer-events: none;

    .deco-blob {
        position: absolute;
        border-radius: 50%;
        filter: blur(70px);
        opacity: 0.55;

        &-1 {
            top: -18%;
            right: -6%;
            width: 420px;
            height: 420px;
            background: color-mix(in srgb, var(--el-color-primary) 16%, transparent);
            animation: blob-drift 12s ease-in-out infinite;
        }

        &-2 {
            bottom: -14%;
            left: -6%;
            width: 360px;
            height: 360px;
            background: color-mix(in srgb, var(--el-color-primary) 11%, transparent);
            animation: blob-drift 14s ease-in-out infinite reverse;
        }

        &-3 {
            top: 20%;
            left: 10%;
            width: 220px;
            height: 220px;
            background: color-mix(in srgb, var(--el-color-primary-light-3) 20%, transparent);
            animation: blob-drift 16s ease-in-out infinite;
        }
    }

    .deco-dot {
        position: absolute;
        border-radius: 50%;
        background: var(--el-color-primary);
        opacity: 0.35;
        animation: dot-float 4.5s ease-in-out infinite;

        &-1 {
            top: 20%;
            right: 16%;
            width: 9px;
            height: 9px;
        }

        &-2 {
            bottom: 24%;
            right: 22%;
            width: 5px;
            height: 5px;
            animation-delay: -1.2s;
        }

        &-3 {
            top: 28%;
            left: 14%;
            width: 6px;
            height: 6px;
            animation-delay: -2.4s;
        }

        &-4 {
            bottom: 18%;
            left: 20%;
            width: 8px;
            height: 8px;
            animation-delay: -3.6s;
        }
    }
}

.error-card {
    position: relative;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    max-width: 560px;
    padding: 48px 40px 56px;
    text-align: center;
    background: var(--el-bg-color-overlay);
    border: 1px solid var(--el-border-color-light);
    border-radius: 16px;
    box-shadow: var(--el-box-shadow-light);
}

.error-hero {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 50px 90px;
    margin-bottom: 8px;
    animation: hero-bob 3.4s ease-in-out infinite;

    .error-code {
        position: relative;
        z-index: 1;
        font-size: 120px;
        font-weight: 800;
        line-height: 1;
        letter-spacing: 4px;
        user-select: none;
        background: linear-gradient(135deg, var(--el-color-primary), var(--el-color-primary-light-3));
        background-clip: text;
        -webkit-background-clip: text;
        color: transparent;
        filter: drop-shadow(0 10px 24px color-mix(in srgb, var(--el-color-primary) 30%, transparent));
    }

    .deco-ring {
        position: absolute;
        inset: 0;
        margin: auto;
        border-radius: 50%;
        border: 1px solid color-mix(in srgb, var(--el-color-primary) 6%, transparent);

        &::before {
            content: '';
            position: absolute;
            top: -3px;
            left: 50%;
            width: 8px;
            height: 8px;
            margin-left: -4px;
            border-radius: 50%;
            background: var(--el-color-primary);
            box-shadow: 0 0 12px color-mix(in srgb, var(--el-color-primary) 60%, transparent);
        }

        &-outer {
            width: 220px;
            height: 220px;
            animation: ring-spin 12s linear infinite;
        }

        &-inner {
            width: 158px;
            height: 158px;
            border-color: color-mix(in srgb, var(--el-color-primary) 4%, transparent);
            animation: ring-spin 8s linear infinite reverse;

            &::before {
                top: -1.5px;
                width: 5px;
                height: 5px;
                margin-left: -2.5px;
                background: var(--el-color-primary-light-3);
                box-shadow: 0 0 8px color-mix(in srgb, var(--el-color-primary-light-3) 60%, transparent);
            }
        }
    }
}

.error-title {
    margin: 0 0 12px;
    font-size: 26px;
    font-weight: 600;
    color: var(--el-text-color-primary);
}

.error-desc {
    margin: 0 0 32px;
    font-size: 14px;
    line-height: 1.8;
    color: var(--el-text-color-secondary);
}

.error-actions {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;

    .btn-icon {
        margin-right: 6px;
    }
}

@keyframes hero-bob {
    0%,
    100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-10px);
    }
}

@keyframes blob-drift {
    0%,
    100% {
        transform: translate(0, 0) scale(1);
    }
    50% {
        transform: translate(18px, -14px) scale(1.06);
    }
}

@keyframes dot-float {
    0%,
    100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(-16px);
    }
}

@keyframes ring-spin {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

@media (max-width: 480px) {
    .error-card {
        padding: 32px 24px 40px;
    }

    .error-hero {
        padding: 32px 46px;

        .error-code {
            font-size: 84px;
        }

        .deco-ring-outer {
            width: 170px;
            height: 170px;
        }

        .deco-ring-inner {
            width: 122px;
            height: 122px;
        }
    }
}
</style>
