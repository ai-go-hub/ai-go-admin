<template>
    <div class="home-page">
        <!-- ==================== 导航栏 ==================== -->
        <header class="navbar">
            <div class="navbar-inner">
                <span class="logo">{{ config.site.name }}</span>
                <nav class="nav-links">
                    <a href="https://github.com/ai-go-hub/ai-go-admin" target="_blank" class="nav-link">GitHub</a>
                    <a href="https://gitee.com/ai-go-hub/ai-go-admin" target="_blank" class="nav-link">Gitee</a>
                </nav>
            </div>
        </header>

        <!-- ==================== 主视觉区 ==================== -->
        <section class="hero">
            <div class="hero-bg-deco">
                <div class="deco-blob deco-blob-1" />
                <div class="deco-blob deco-blob-2" />
                <div class="deco-grid" />
            </div>
            <div class="hero-content">
                <div class="badge-row">
                    <span class="ai-badge">
                        <el-icon><Cpu /></el-icon>
                        {{ $t('index.heroBadge') }}
                    </span>
                </div>
                <h1 class="hero-title">
                    <span class="title-line">{{ config.site.name }}</span>
                </h1>
                <p class="hero-desc">{{ $t('index.heroDesc') }}</p>
                <div class="hero-actions">
                    <el-button size="large" type="primary" round @click="$router.push('/admin')"> {{ $t('index.enterAdmin') }} </el-button>
                    <el-button size="large" round @click="scrollToTech"> {{ $t('index.learnTech') }} </el-button>
                </div>
            </div>
            <div class="scroll-hint" @click="scrollToTech">
                <span>{{ $t('index.scrollDown') }}</span>
                <el-icon class="bounce-arrow"><ArrowDown /></el-icon>
            </div>
        </section>

        <!-- ==================== 技术栈 ==================== -->
        <section ref="techSection" class="tech-section">
            <div class="section-header">
                <h2>{{ $t('index.techTitle') }}</h2>
                <p>{{ $t('index.techSubtitle') }}</p>
            </div>
            <div class="tech-grid">
                <div v-for="item in techStack" :key="item.name" class="tech-card">
                    <div class="tech-icon" :style="{ background: item.gradient }">
                        <el-icon :size="28"><component :is="item.icon" /></el-icon>
                    </div>
                    <h3>{{ item.name }}</h3>
                    <p>{{ item.desc }}</p>
                </div>
            </div>
        </section>

        <!-- ==================== 开发理念 ==================== -->
        <section class="philosophy-section">
            <div class="section-header">
                <h2>{{ $t('index.philosophyTitle') }}</h2>
                <p>{{ $t('index.philosophySubtitle') }}</p>
            </div>
            <div class="philosophy-content">
                <div class="philosophy-card human">
                    <div class="card-icon">
                        <el-icon :size="40"><UserFilled /></el-icon>
                    </div>
                    <h3>{{ $t('index.humanCommand') }}</h3>
                    <p>{{ $t('index.humanDesc') }}</p>
                </div>
                <div class="philosophy-divider">
                    <div class="divider-line" />
                    <span class="divider-plus">
                        <el-icon :size="24"><Plus /></el-icon>
                    </span>
                    <div class="divider-line" />
                </div>
                <div class="philosophy-card ai">
                    <div class="card-icon">
                        <el-icon :size="40"><Cpu /></el-icon>
                    </div>
                    <h3>{{ $t('index.aiWriter') }}</h3>
                    <p>{{ $t('index.aiDesc') }}</p>
                </div>
            </div>
        </section>

        <!-- ==================== 博客链接 ==================== -->
        <section class="blog-section">
            <div class="section-header">
                <h2>{{ $t('index.blogTitle') }}</h2>
                <p>{{ $t('index.blogSubtitle') }}</p>
            </div>
            <div class="blog-cards">
                <a href="https://github.com/ai-go-hub/ai-go-blog" target="_blank" class="blog-card">
                    <div class="blog-card-icon github">
                        <svg viewBox="0 0 24 24" fill="currentColor" class="github-svg">
                            <path
                                d="M12 0C5.37 0 0 5.37 0 12c0 5.3 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61-.546-1.385-1.335-1.755-1.335-1.755-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.605-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 21.795 24 17.295 24 12 24 5.37 18.63 0 12 0z"
                            />
                        </svg>
                    </div>
                    <div class="blog-card-body">
                        <h3>GitHub</h3>
                        <p>github.com/ai-go-hub/ai-go-blog</p>
                    </div>
                    <el-icon class="external-icon"><ArrowRight /></el-icon>
                </a>

                <a href="https://gitee.com/ai-go-hub/ai-go-blog" target="_blank" class="blog-card">
                    <div class="blog-card-icon gitee">
                        <svg
                            fill="currentColor"
                            viewBox="0 0 90 90"
                            version="1.1"
                            xmlns="http://www.w3.org/2000/svg"
                            xmlns:xlink="http://www.w3.org/1999/xlink"
                        >
                            <title>Gitee Logo</title>
                            <g stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
                                <circle fill="#C71D23" cx="45" cy="45" r="45"></circle>
                                <path
                                    d="M67.558546,39.8714292 L42.0857966,39.8714292 C40.8627004,39.8720094 39.8710953,40.8633548 39.8701949,42.0864508 L39.8687448,47.623783 C39.867826,48.8471055 40.8592652,49.8390642 42.0825879,49.8393845 C42.0827874,49.8393846 42.0829869,49.8393846 42.0831864,49.8387862 L57.5909484,49.838657 C58.8142711,49.8386283 59.8059783,50.830319 59.8059885,52.0536417 C59.8059885,52.0536479 59.8059885,52.053654 59.8059701,52.0536602 L59.8059701,52.6073539 L59.8059701,52.6073539 L59.8059701,53.161115 C59.8059701,56.8310831 56.8308731,59.80618 53.160905,59.80618 L32.1165505,59.80618 C30.8934034,59.806119 29.9018373,58.8145802 29.9017425,57.5914331 L29.9011625,36.5491188 C29.9008781,32.8791508 32.8758931,29.9039718 36.5458611,29.9038706 C36.5459222,29.9038706 36.5459833,29.9038706 36.5460443,29.9040538 L67.5523638,29.9040538 C68.77515,29.9026795 69.7666266,28.9118177 69.7687593,27.6890325 L69.7721938,22.1516997 C69.774326,20.928378 68.7832423,19.9360642 67.5599198,19.9353054 C67.5594619,19.9353051 67.5590039,19.935305 67.558546,19.9366784 L36.5479677,19.9366784 C27.3730474,19.9366784 19.935305,27.3744208 19.935305,36.549341 L19.935305,67.558546 C19.935305,68.7818687 20.927004,69.7735676 22.1503267,69.7735676 L54.8224984,69.7735676 C63.079746,69.7735676 69.7735676,63.079746 69.7735676,54.8224984 L69.7735676,42.0864509 C69.7735676,40.8631282 68.7818687,39.8714292 67.558546,39.8714292 Z"
                                    fill="#FFFFFF"
                                ></path>
                            </g>
                        </svg>
                    </div>
                    <div class="blog-card-body">
                        <h3>Gitee</h3>
                        <p>gitee.com/ai-go-hub/ai-go-blog</p>
                    </div>
                    <el-icon class="external-icon"><ArrowRight /></el-icon>
                </a>
            </div>
        </section>

        <!-- ==================== 页脚 ==================== -->
        <footer class="footer">
            <div class="footer-inner">
                <div class="footer-brand">
                    <span class="footer-logo">{{ config.site.name }}</span>
                    <span class="footer-tagline">{{ $t('index.footerTagline') }}</span>
                </div>
                <div class="footer-links">
                    <a href="https://github.com/ai-go-hub/ai-go-admin" target="_blank">GitHub</a>
                    <span class="footer-divider">·</span>
                    <a href="https://gitee.com/ai-go-hub/ai-go-admin" target="_blank">Gitee</a>
                    <span class="footer-divider">·</span>
                    <a href="http://beian.miit.gov.cn/" target="_blank">{{ config.site.recordNumber }}</a>
                </div>
                <p class="footer-copy">&copy; {{ new Date().getFullYear() }} AI GO MALL. All rights reserved.</p>
            </div>
        </footer>
    </div>
</template>

<script setup lang="ts">
import { ArrowDown, ArrowRight, Coin, Connection, Cpu, DataBoard, Grid, Monitor, Plus, Setting, UserFilled } from '@element-plus/icons-vue'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfig } from '/@/stores/config'

const { t } = useI18n()
const config = useConfig()

const techSection = ref<HTMLElement | null>(null)

function scrollToTech() {
    techSection.value?.scrollIntoView({ behavior: 'smooth' })
}

interface TechItem {
    name: string
    desc: string
    icon: any
    gradient: string
}

const techStack: TechItem[] = [
    {
        name: 'Golang',
        desc: t('index.techGolang'),
        icon: Monitor,
        gradient: 'linear-gradient(135deg, #00ADD8, #00ADD8cc)',
    },
    {
        name: 'TypeScript',
        desc: t('index.techTypeScript'),
        icon: DataBoard,
        gradient: 'linear-gradient(135deg, #3178C6, #3178C6cc)',
    },
    {
        name: 'PostgreSQL',
        desc: t('index.techPostgreSQL'),
        icon: Coin,
        gradient: 'linear-gradient(135deg, #4169E1, #4169E1cc)',
    },
    {
        name: 'Vue 3 + Pinia',
        desc: t('index.techVue'),
        icon: Connection,
        gradient: 'linear-gradient(135deg, #4FC08D, #4FC08Dcc)',
    },
    {
        name: 'Element Plus',
        desc: t('index.techElementPlus'),
        icon: Grid,
        gradient: 'linear-gradient(135deg, #409EFF, #409EFFcc)',
    },
    {
        name: 'Claude Code',
        desc: t('index.techClaudeCode'),
        icon: Setting,
        gradient: 'linear-gradient(135deg, #D97706, #D97706cc)',
    },
]
</script>

<style scoped lang="scss">
// ==================== 导航栏 ====================

.navbar {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 100;
    background: rgba(255, 255, 255, 0.8);
    backdrop-filter: blur(12px);
    border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.navbar-inner {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 32px;
    height: 64px;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.logo {
    font-size: 18px;
    font-weight: 700;
    color: var(--el-text-color-primary);
    letter-spacing: -0.5px;
}

.nav-links {
    display: flex;
    align-items: center;
    gap: 16px;
}

.nav-link {
    font-size: 14px;
    color: var(--el-text-color-regular);
    text-decoration: none;
    transition: color 0.2s;

    &:hover {
        color: var(--el-color-primary);
    }
}

// ==================== 主视觉区 ====================

.hero {
    position: relative;
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    padding: 96px 32px 64px;
}

.hero-bg-deco {
    position: absolute;
    inset: 0;
    pointer-events: none;
}

.deco-blob {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.6;
}

.deco-blob-1 {
    top: -20%;
    right: -10%;
    width: 500px;
    height: 500px;
    background: color-mix(in srgb, var(--el-color-primary) 15%, transparent);
}

.deco-blob-2 {
    bottom: -10%;
    left: -10%;
    width: 400px;
    height: 400px;
    background: color-mix(in srgb, var(--el-color-primary) 10%, transparent);
}

.deco-grid {
    position: absolute;
    inset: 0;
    background-image:
        linear-gradient(to right, color-mix(in srgb, var(--el-color-primary) 4%, transparent) 1px, transparent 1px),
        linear-gradient(to bottom, color-mix(in srgb, var(--el-color-primary) 4%, transparent) 1px, transparent 1px);
    background-size: 60px 60px;
    mask-image: radial-gradient(ellipse at center, black 30%, transparent 70%);
}

.hero-content {
    position: relative;
    z-index: 10;
    text-align: center;
    max-width: 720px;
}

.badge-row {
    margin-bottom: 24px;
}

.ai-badge {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 6px 16px;
    background: color-mix(in srgb, var(--el-color-primary) 8%, transparent);
    color: var(--el-color-primary);
    font-size: 13px;
    font-weight: 500;
    border-radius: 100px;
    border: 1px solid color-mix(in srgb, var(--el-color-primary) 15%, transparent);
}

.hero-title {
    margin-bottom: 24px;

    .title-line {
        display: block;
        font-size: 56px;
        font-weight: 800;
        color: var(--el-text-color-primary);
        letter-spacing: -1.5px;
        line-height: 1.2;
    }
}

.hero-desc {
    font-size: 16px;
    color: var(--el-text-color-regular);
    line-height: 1.8;
    margin-bottom: 40px;
    max-width: 560px;
    margin-left: auto;
    margin-right: auto;
}

.hero-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
    flex-wrap: wrap;
}

.scroll-hint {
    position: absolute;
    bottom: 32px;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
    cursor: pointer;
    z-index: 10;
    transition: color 0.2s;

    &:hover {
        color: var(--el-color-primary);
    }
}

.bounce-arrow {
    animation: bounce 2s infinite;
}

@keyframes bounce {
    0%,
    100% {
        transform: translateY(0);
    }
    50% {
        transform: translateY(6px);
    }
}

// ==================== 通用 section 头部 ====================

.section-header {
    text-align: center;
    margin-bottom: 56px;

    h2 {
        font-size: 32px;
        font-weight: 700;
        color: var(--el-text-color-primary);
        letter-spacing: -0.5px;
        margin: 0 0 12px;
    }

    p {
        font-size: 16px;
        color: var(--el-text-color-regular);
        margin: 0;
    }
}

// ==================== 技术栈 ====================

.tech-section {
    padding: 100px 32px;
    background: var(--el-fill-color-light);
}

.tech-grid {
    max-width: 1000px;
    margin: 0 auto;
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 24px;
}

.tech-card {
    background: var(--el-bg-color);
    border-radius: 16px;
    padding: 32px 28px;
    border: 1px solid var(--el-border-color-light);
    transition: all 0.3s ease;

    &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.08);
    }

    .tech-icon {
        width: 56px;
        height: 56px;
        border-radius: 14px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        margin-bottom: 20px;
    }

    h3 {
        font-size: 18px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin: 0 0 8px;
    }

    p {
        font-size: 14px;
        color: var(--el-text-color-regular);
        line-height: 1.7;
        margin: 0;
    }
}

// ==================== 开发理念 ====================

.philosophy-section {
    padding: 100px 32px;
}

.philosophy-content {
    max-width: 800px;
    margin: 0 auto;
    display: flex;
    align-items: center;
    gap: 32px;
}

.philosophy-card {
    flex: 1;
    text-align: center;
    padding: 40px 28px;
    border-radius: 16px;
    border: 1px solid var(--el-border-color-light);
    transition: all 0.3s ease;

    &:hover {
        transform: translateY(-4px);
        box-shadow: 0 12px 40px rgba(0, 0, 0, 0.08);
    }

    .card-icon {
        width: 80px;
        height: 80px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 24px;
    }

    h3 {
        font-size: 20px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin: 0 0 12px;
    }

    p {
        font-size: 14px;
        color: var(--el-text-color-regular);
        line-height: 1.8;
        margin: 0;
    }
}

.philosophy-card.human {
    background: var(--el-bg-color);

    .card-icon {
        background: color-mix(in srgb, var(--el-color-primary) 8%, transparent);
        color: var(--el-color-primary);
    }
}

.philosophy-card.ai {
    background: linear-gradient(135deg, #1f2937, #111827);

    h3 {
        color: #f9fafb;
    }

    p {
        color: #9ca3af;
    }

    .card-icon {
        background: rgba(255, 255, 255, 0.1);
        color: #f9fafb;
    }
}

.philosophy-divider {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
}

.divider-line {
    width: 1px;
    height: 40px;
    background: var(--el-border-color-light);
}

.divider-plus {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: var(--el-fill-color-light);
    border: 1px solid var(--el-border-color-light);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--el-text-color-secondary);
}

// ==================== 博客链接 ====================

.blog-section {
    padding: 100px 32px;
    background: var(--el-fill-color-light);
}

.blog-cards {
    max-width: 640px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.blog-card {
    display: flex;
    align-items: center;
    gap: 20px;
    padding: 24px 28px;
    background: var(--el-bg-color);
    border-radius: 16px;
    border: 1px solid var(--el-border-color-light);
    text-decoration: none;
    transition: all 0.3s ease;

    &:hover {
        transform: translateY(-2px);
        box-shadow: 0 8px 30px rgba(0, 0, 0, 0.06);
    }
}

.blog-card-icon {
    width: 52px;
    height: 52px;
    border-radius: 14px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: white;

    &.github {
        background: #24292e;
    }

    &.gitee {
        background: #c71d23;
    }
}

.github-svg {
    width: 32px;
    height: 32px;
}

.blog-card-body {
    flex: 1;

    h3 {
        font-size: 16px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin: 0 0 4px;
    }

    p {
        font-size: 13px;
        color: var(--el-text-color-secondary);
        margin: 0;
        font-family: 'SF Mono', 'Fira Code', monospace;
    }
}

.external-icon {
    color: var(--el-text-color-secondary);
    flex-shrink: 0;
    transition: all 0.2s;

    .blog-card:hover & {
        color: var(--el-color-primary);
        transform: translateX(2px);
    }
}

// ==================== 页脚 ====================

.footer {
    padding: 48px 32px;
    border-top: 1px solid var(--el-border-color-light);
    background: var(--el-bg-color);
}

.footer-inner {
    max-width: 1200px;
    margin: 0 auto;
    text-align: center;
}

.footer-brand {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    margin-bottom: 16px;
}

.footer-logo {
    font-size: 16px;
    font-weight: 700;
    color: var(--el-text-color-primary);
}

.footer-tagline {
    font-size: 13px;
    color: var(--el-text-color-secondary);
    padding: 2px 10px;
    background: var(--el-fill-color-light);
    border-radius: 100px;
}

.footer-links {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-bottom: 12px;

    a {
        font-size: 13px;
        color: var(--el-text-color-secondary);
        text-decoration: none;
        transition: color 0.2s;

        &:hover {
            color: var(--el-color-primary);
        }
    }
}

.footer-divider {
    color: var(--el-text-color-primary);
}

.footer-copy {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin: 0;
}

// ==================== 响应式 ====================

@media (max-width: 900px) {
    .tech-grid {
        grid-template-columns: repeat(2, 1fr);
    }

    .philosophy-content {
        flex-direction: column;
    }

    .philosophy-divider {
        flex-direction: row;

        .divider-line {
            width: 40px;
            height: 1px;
        }
    }
}

@media (max-width: 640px) {
    .hero-title {
        .title-line {
            font-size: 36px;
        }
    }

    .hero-desc {
        font-size: 14px;
    }

    .tech-grid {
        grid-template-columns: 1fr;
    }

    .section-header h2 {
        font-size: 24px;
    }

    .hero-actions {
        flex-direction: column;
        align-items: center;
    }
}
</style>
