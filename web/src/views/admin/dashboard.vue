<template>
    <div class="default-main dashboard">
        <!-- ==================== 欢迎条 ==================== -->
        <div class="welcome-row">
            <div class="welcome-bar welcome-bar-left">
                <img :src="DashboardHeaderIllustration" class="welcome-illustration" alt="" />
                <div class="welcome-content">
                    <div class="hello">{{ greeting }}，{{ adminInfo.nickname }} 👋</div>
                    <div class="tip">
                        开源等于互助；需要大家一起来支持，比如使用、推荐、写教程、保护生态、贡献代码、回答问题、分享经验、打赏赞助等；欢迎您加入我们！
                    </div>
                    <div class="welcome-meta">
                        <div class="meta-item">
                            <Icon :size="14" name="el-calendar" />
                            <span>{{ today }}</span>
                        </div>
                        <div class="meta-item">
                            <Icon :size="14" name="el-timer" />
                            <span>{{ nowTime }}</span>
                        </div>
                    </div>
                </div>
            </div>

            <div class="welcome-timer">
                <div class="timer-label">
                    <Icon :size="14" name="el-clock" />
                    <span>您今天已工作了</span>
                </div>
                <div class="timer-time" :class="{ paused: !working }">
                    <div v-for="(unit, idx) in workUnits" :key="unit.label" class="timer-group">
                        <div class="timer-unit">
                            <div class="num">{{ pad(unit.value) }}</div>
                            <div class="txt">{{ unit.label }}</div>
                        </div>
                        <div v-if="idx < workUnits.length - 1" class="colon">:</div>
                    </div>
                </div>
                <button type="button" class="timer-btn" :class="{ resting: !working }" @click="toggleWork">
                    <Icon :size="14" v-if="working" name="el-coffee" />
                    <Icon :size="14" v-else name="el-video-play" />
                    {{ working ? '休息片刻' : '继续工作' }}
                </button>
            </div>
        </div>

        <!-- ==================== KPI 卡片行 ==================== -->
        <div class="kpi-row">
            <div v-for="item in kpiList" :key="item.key" class="kpi-card">
                <div class="kpi-top">
                    <div class="kpi-title">{{ item.title }}</div>
                    <div class="kpi-icon">
                        <Icon :size="16" :name="item.icon" />
                    </div>
                </div>
                <div class="kpi-value">{{ item.value }}</div>
                <div class="kpi-bottom">
                    <span class="kpi-trend" :class="item.trend > 0 ? 'up' : 'down'">
                        <Icon :size="12" :name="item.trend > 0 ? 'el-caret-top' : 'el-caret-bottom'" />
                        {{ Math.abs(item.trend) }}%
                    </span>
                    <span class="kpi-sub">较上周</span>
                    <div :ref="(el) => setSparkRef(el as HTMLElement, item.key)" class="kpi-spark"></div>
                </div>
            </div>
        </div>

        <!-- ==================== 图表行 ==================== -->
        <div class="chart-row">
            <div class="chart-card chart-card-span-2">
                <div class="card-header">
                    <div class="card-title">访问趋势</div>
                    <div class="trend-ranges">
                        <button
                            v-for="r in rangeList"
                            :key="r.value"
                            type="button"
                            class="range-chip"
                            :class="{ 'is-active': trendRange === r.value }"
                            @click="setRange(r.value)"
                        >
                            {{ r.label }}
                        </button>
                    </div>
                </div>
                <div ref="trendRef" class="chart-body chart-body-lg"></div>
            </div>

            <div class="chart-card">
                <div class="card-header">
                    <div class="card-title">流量来源</div>
                    <el-tag size="small" round type="primary">今日</el-tag>
                </div>
                <div ref="sourceRef" class="chart-body chart-body-md"></div>
            </div>
        </div>

        <!-- ==================== 图表行 ==================== -->
        <div class="chart-row">
            <div class="chart-card">
                <div class="card-header">
                    <div class="card-title">分类销量对比</div>
                    <el-tag size="small" round>本月</el-tag>
                </div>
                <div ref="salesRef" class="chart-body chart-body-md"></div>
            </div>

            <div class="chart-card">
                <div class="card-header">
                    <div class="card-title">服务器负载</div>
                    <el-tag size="small" round type="success">稳定</el-tag>
                </div>
                <div ref="gaugeRef" class="chart-body chart-body-md"></div>
            </div>

            <div class="chart-card">
                <div class="card-header">
                    <div class="card-title">能力雷达</div>
                    <el-tag size="small" round type="warning">最新</el-tag>
                </div>
                <div ref="radarRef" class="chart-body chart-body-md"></div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useEventListener } from '@vueuse/core'
import type { ECharts, EChartsOption } from 'echarts'
import * as echarts from 'echarts'
import { debounce, padStart } from 'lodash-es'
import { computed, nextTick, onBeforeUnmount, onMounted, reactive, ref, watch } from 'vue'
import DashboardHeaderIllustration from '/@/assets/svg/dashboard-illustration.svg'
import { useAdminInfo } from '/@/stores/adminInfo'
import { useConfig } from '/@/stores/config'
import { Local } from '/@/utils/storage'

const config = useConfig()
const adminInfo = useAdminInfo()
const workTimerKey = 'ai-go-dashboard-work-timer'

// ==================== 工具函数 ====================

/**
 * 在一位数字的前面补零
 */
const pad = (n: number) => padStart(String(n), 2, '0')

/**
 * 构造一条底部渐变面积色（同色系深浅过渡）
 */
const makeAreaGradient = (color: string, from = 0.35, to = 0) =>
    new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: withAlpha(color, from) },
        { offset: 1, color: withAlpha(color, to) },
    ])

/**
 * 构造一条从上到下渐变的柱状色（顶饱和 → 底通透）
 */
const makeBarGradient = (color: string) =>
    new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color },
        { offset: 1, color: color + '55' },
    ])

/**
 * 把 hex 或 rgb 色转成 rgba 字符串
 */
const withAlpha = (color: string, alpha: number) => {
    if (color.startsWith('#')) {
        const hex = color.slice(1)
        const bigint = parseInt(hex.length === 3 ? hex.replace(/./g, (c) => c + c) : hex, 16)
        const r = (bigint >> 16) & 255
        const g = (bigint >> 8) & 255
        const b = bigint & 255
        return `rgba(${r}, ${g}, ${b}, ${alpha})`
    }
    return color
}

// ==================== 顶部欢迎条 ====================

const nowTime = ref('')
const today = new Date().toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', weekday: 'long' })
const greeting = computed(() => {
    const h = new Date().getHours()
    if (h < 6) return '夜深了'
    if (h < 11) return '上午好'
    if (h < 13) return '中午好'
    if (h < 18) return '下午好'
    return '晚上好'
})

// ==================== 今日工作计时 ====================

// 持久化到 localStorage，跨页面刷新保持；跨 0 点自动重置
type WorkState = {
    date: string // yyyy-mm-dd
    accumulated: number // 已累计工作秒数（不含正在计时的当前段）
    startAt: number | null // 当前段起点时间戳（ms），null 表示暂停
}

const todayKey = () => {
    const d = new Date()
    return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`
}

const loadWorkState = (): WorkState => {
    const cached = Local.get(workTimerKey) as WorkState | null
    if (cached && cached.date === todayKey()) return cached
    // 新的一天，或首次进入：默认启动计时
    return { date: todayKey(), accumulated: 0, startAt: Date.now() }
}

const saveWorkState = () => Local.set(workTimerKey, workState)

const workState: WorkState = reactive(loadWorkState())
const workSeconds = ref(0)
const working = computed(() => workState.startAt !== null)
const workUnits = computed(() => [
    { label: '小时', value: Math.floor(workSeconds.value / 3600) },
    { label: '分', value: Math.floor((workSeconds.value % 3600) / 60) },
    { label: '秒', value: workSeconds.value % 60 },
])

const updateWorkTime = (now: Date) => {
    // 跨 0 点重置
    if (workState.date !== todayKey()) {
        workState.date = todayKey()
        workState.accumulated = 0
        workState.startAt = now.getTime()
        saveWorkState()
    }
    const running = workState.startAt !== null ? Math.floor((now.getTime() - workState.startAt) / 1000) : 0
    workSeconds.value = workState.accumulated + running
}

const toggleWork = () => {
    const now = Date.now()
    if (workState.startAt !== null) {
        // 暂停：把当前段并入累计
        workState.accumulated += Math.floor((now - workState.startAt) / 1000)
        workState.startAt = null
    } else {
        // 继续
        workState.startAt = now
    }
    saveWorkState()
    updateWorkTime(new Date())
}

// ==================== 每秒 tick（同时更新时钟与工作计时） ====================

let clockTimer: ReturnType<typeof setInterval> | null = null

const tick = () => {
    const d = new Date()
    nowTime.value = `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
    updateWorkTime(d)
}

// ==================== KPI 卡片数据 ====================

type KpiItem = {
    key: 'user' | 'order' | 'sales' | 'visit'
    title: string
    value: string
    trend: number
    icon: string
    spark: number[]
}
const kpiList: KpiItem[] = [
    {
        key: 'user',
        title: '总用户数',
        value: '128,436',
        trend: 12.5,
        icon: 'el-user-filled',
        spark: [12, 18, 22, 20, 25, 28, 32, 30, 35, 40, 38, 45],
    },
    {
        key: 'order',
        title: '今日订单',
        value: '1,286',
        trend: 8.2,
        icon: 'el-goods',
        spark: [8, 12, 10, 15, 14, 20, 18, 24, 22, 26, 30, 28],
    },
    {
        key: 'sales',
        title: '销售额（元）',
        value: '86,540',
        trend: -3.1,
        icon: 'el-money',
        spark: [40, 45, 38, 42, 50, 48, 55, 52, 60, 55, 62, 58],
    },
    {
        key: 'visit',
        title: '页面访问',
        value: '52,180',
        trend: 21.8,
        icon: 'el-trend-charts',
        spark: [20, 25, 30, 28, 35, 40, 38, 45, 50, 55, 60, 65],
    },
]

// ==================== 图表实例统一管理 ====================

type ChartKey = 'trend' | 'source' | 'sales' | 'gauge' | 'radar'
const trendRef = ref<HTMLElement | null>(null)
const sourceRef = ref<HTMLElement | null>(null)
const salesRef = ref<HTMLElement | null>(null)
const gaugeRef = ref<HTMLElement | null>(null)
const radarRef = ref<HTMLElement | null>(null)
const chartRefs: Record<ChartKey, ReturnType<typeof ref<HTMLElement | null>>> = {
    trend: trendRef,
    source: sourceRef,
    sales: salesRef,
    gauge: gaugeRef,
    radar: radarRef,
}
const charts: Partial<Record<ChartKey, ECharts>> = {}

/**
 * 惰性获取 ECharts 实例；同一 DOM 上重复调用返回同一实例
 */
const ensureChart = (key: ChartKey): ECharts | null => {
    if (charts[key]) return charts[key]!
    const el = chartRefs[key].value
    if (!el) return null
    charts[key] = echarts.init(el)
    return charts[key]!
}

// ==================== Spark（KPI 内嵌迷你图） ====================

const sparkEls: Partial<Record<KpiItem['key'], HTMLElement | null>> = {}
const sparkCharts: Partial<Record<KpiItem['key'], ECharts>> = {}
const setSparkRef = (el: HTMLElement | null, key: KpiItem['key']) => {
    sparkEls[key] = el
}

const renderSpark = (item: KpiItem) => {
    const el = sparkEls[item.key]
    if (!el) return
    sparkCharts[item.key]?.dispose()
    const chart = echarts.init(el)
    chart.setOption({
        grid: { top: 4, right: 4, bottom: 4, left: 4 },
        xAxis: { type: 'category', show: false, data: item.spark.map((_, i) => i) },
        yAxis: { type: 'value', show: false },
        tooltip: { show: false },
        series: [
            {
                type: 'line',
                data: item.spark,
                smooth: true,
                showSymbol: false,
                // 线条细一点、面积渐变淡一点，避免整排 4 个 spark 视觉过重
                lineStyle: { color: '#409eff', width: 1.5 },
                areaStyle: { color: makeAreaGradient('#409eff', 0.25, 0) },
            },
        ],
    })
    sparkCharts[item.key] = chart
}

// ==================== 主题（暗/亮） ====================
type ThemeTokens = {
    textPrimary: string
    textSecondary: string
    borderColor: string
    axisLine: string
    splitLine: string
    bg: string
    splitArea: [string, string]
}
const isDark = computed(() => config.layout.dark)
const theme = computed<ThemeTokens>(() =>
    isDark.value
        ? {
              textPrimary: '#e5eaf3',
              textSecondary: '#a3a6ad',
              borderColor: '#414243',
              axisLine: '#4c4d4f',
              splitLine: '#363637',
              bg: '#1d1e1f',
              splitArea: ['#2a2a2b', '#232324'],
          }
        : {
              textPrimary: '#303133',
              textSecondary: '#606266',
              borderColor: '#e4e7ed',
              axisLine: '#dcdfe6',
              splitLine: '#ebeef5',
              bg: '#ffffff',
              splitArea: ['#fafbfc', '#f5f7fa'],
          }
)

/**
 * 主题化的通用 tooltip 配置
 */
const buildTooltip = (extra: Partial<EChartsOption['tooltip']> = {}): EChartsOption['tooltip'] => ({
    backgroundColor: theme.value.bg,
    borderColor: theme.value.borderColor,
    textStyle: { color: theme.value.textPrimary },
    ...extra,
})

/**
 * 主题化的通用 x/y 轴配置
 */
const buildAxis = (opt: { hideLine?: boolean } = {}) => ({
    axisLine: opt.hideLine ? { show: false } : { lineStyle: { color: theme.value.axisLine } },
    axisTick: opt.hideLine ? { show: false } : undefined,
    splitLine: { lineStyle: { color: theme.value.splitLine } },
    axisLabel: { color: theme.value.textSecondary },
})

// ==================== 访问趋势 ====================

const trendRange = ref<'7' | '30'>('7')
const rangeList = [
    { value: '7' as const, label: '近 7 天' },
    { value: '30' as const, label: '近 30 天' },
]
const setRange = (v: '7' | '30') => {
    trendRange.value = v
    renderTrend()
}

/**
 * 根据当天日期生成稳定伪随机数据（避免每次刷新数值跳变）
 */
const buildTrendData = (days: number) => {
    const now = new Date()
    const dates: string[] = []
    const pv: number[] = []
    const uv: number[] = []
    for (let i = days - 1; i >= 0; i--) {
        const d = new Date(now)
        d.setDate(d.getDate() - i)
        dates.push(`${d.getMonth() + 1}/${d.getDate()}`)
        const seed = d.getMonth() * 31 + d.getDate()
        pv.push(2000 + Math.floor(Math.abs(Math.sin(seed)) * 3500))
        uv.push(600 + Math.floor(Math.abs(Math.cos(seed)) * 1200))
    }
    return { dates, pv, uv }
}

const renderTrend = () => {
    const chart = ensureChart('trend')
    if (!chart) return
    const { dates, pv, uv } = buildTrendData(Number(trendRange.value))
    const t = theme.value

    chart.setOption(
        {
            color: ['#409eff', '#67c23a'],
            tooltip: buildTooltip({ trigger: 'axis' }),
            legend: {
                data: ['PV', 'UV'],
                right: 12,
                top: 4,
                textStyle: { color: t.textSecondary },
                icon: 'roundRect',
                itemWidth: 12,
                itemHeight: 8,
                itemGap: 14,
            },
            grid: { top: 40, right: 20, bottom: 30, left: 50 },
            xAxis: { type: 'category', data: dates, boundaryGap: false, ...buildAxis() },
            yAxis: { type: 'value', ...buildAxis({ hideLine: true }) },
            series: [
                {
                    name: 'PV',
                    type: 'line',
                    data: pv,
                    smooth: true,
                    showSymbol: false,
                    // 线条更细、面积渐变更淡，减少视觉重量
                    lineStyle: { width: 2 },
                    areaStyle: { color: makeAreaGradient('#409eff', 0.2, 0) },
                    emphasis: { focus: 'series' },
                },
                {
                    name: 'UV',
                    type: 'line',
                    data: uv,
                    smooth: true,
                    showSymbol: false,
                    lineStyle: { width: 2 },
                    areaStyle: { color: makeAreaGradient('#67c23a', 0.16, 0) },
                    emphasis: { focus: 'series' },
                },
            ],
        },
        true
    )
}

// ==================== 流量来源 ====================

const renderSource = () => {
    const chart = ensureChart('source')
    if (!chart) return
    const t = theme.value
    chart.setOption({
        color: ['#409eff', '#6cb0ff', '#9dcbff', '#c8dfff', '#909399'],
        tooltip: buildTooltip({ trigger: 'item' }),
        legend: {
            bottom: 0,
            left: 'center',
            textStyle: { color: t.textSecondary },
            itemWidth: 10,
            itemHeight: 10,
        },
        series: [
            {
                name: '流量来源',
                type: 'pie',
                radius: ['48%', '72%'],
                center: ['50%', '45%'],
                avoidLabelOverlap: true,
                itemStyle: { borderRadius: 8, borderColor: t.bg, borderWidth: 2 },
                label: { show: false, position: 'center' },
                emphasis: {
                    label: { show: true, fontSize: 18, fontWeight: 'bold', color: t.textPrimary },
                },
                labelLine: { show: false },
                data: [
                    { value: 4820, name: '搜索引擎' },
                    { value: 3210, name: '直接访问' },
                    { value: 2140, name: '社交媒体' },
                    { value: 1580, name: '推荐链接' },
                    { value: 980, name: '其他' },
                ],
            },
        ],
    })
}

// ==================== 分类销量 ====================

const SALES_ITEMS = [
    { name: '数码', value: 520 },
    { name: '服饰', value: 820 },
    { name: '家居', value: 380 },
    { name: '食品', value: 680 },
    { name: '美妆', value: 240 },
    { name: '图书', value: 560 },
]
const renderSales = () => {
    const chart = ensureChart('sales')
    if (!chart) return
    const t = theme.value
    chart.setOption({
        tooltip: buildTooltip({ trigger: 'axis', axisPointer: { type: 'shadow' } }),
        grid: { top: 20, right: 16, bottom: 30, left: 50 },
        xAxis: { type: 'category', data: SALES_ITEMS.map((i) => i.name), ...buildAxis() },
        yAxis: { type: 'value', ...buildAxis({ hideLine: true }) },
        series: [
            {
                type: 'bar',
                data: SALES_ITEMS.map((i) => i.value),
                itemStyle: {
                    color: makeBarGradient('#409eff'),
                    borderRadius: [6, 6, 0, 0],
                },
                barWidth: 22,
                label: { show: true, position: 'top', color: t.textSecondary, fontSize: 11 },
            },
        ],
    })
}

// ==================== 服务器负载仪表盘 ====================

const renderGauge = () => {
    const chart = ensureChart('gauge')
    if (!chart) return
    const t = theme.value
    chart.setOption({
        series: [
            {
                type: 'gauge',
                center: ['50%', '58%'],
                startAngle: 200,
                endAngle: -20,
                min: 0,
                max: 100,
                splitNumber: 10,
                itemStyle: { color: '#409eff' },
                progress: { show: true, width: 20, roundCap: true },
                pointer: { show: false },
                axisLine: {
                    lineStyle: {
                        width: 20,
                        color: [
                            [0.5, '#67c23a'],
                            [0.8, '#e6a23c'],
                            [1, '#f56c6c'],
                        ],
                    },
                },
                axisTick: { show: false },
                splitLine: { show: false },
                axisLabel: { show: false },
                anchor: { show: false },
                title: { show: false },
                detail: {
                    valueAnimation: true,
                    fontSize: 32,
                    fontWeight: 700,
                    offsetCenter: [0, '-10%'],
                    formatter: '{value}%',
                    color: t.textPrimary,
                },
                data: [{ value: 62, name: 'CPU 使用率' }],
            },
        ],
    })
}

// ==================== 能力雷达 ====================

const RADAR_INDICATORS = ['稳定性', '性能', '安全性', '易用性', '扩展性']
const renderRadar = () => {
    const chart = ensureChart('radar')
    if (!chart) return
    const t = theme.value
    chart.setOption({
        tooltip: buildTooltip(),
        radar: {
            indicator: RADAR_INDICATORS.map((name) => ({ name, max: 100 })),
            radius: '65%',
            center: ['50%', '52%'],
            splitLine: { lineStyle: { color: t.splitLine } },
            splitArea: { areaStyle: { color: t.splitArea as string[] } },
            axisLine: { lineStyle: { color: t.axisLine } },
            axisName: { color: t.textSecondary, fontSize: 12 },
        },
        series: [
            {
                type: 'radar',
                data: [
                    {
                        value: [88, 82, 92, 76, 84],
                        name: '当前版本',
                        lineStyle: { color: '#409eff', width: 2 },
                        itemStyle: { color: '#409eff' },
                        areaStyle: { color: withAlpha('#409eff', 0.25) },
                    },
                    {
                        value: [72, 66, 80, 82, 74],
                        name: '上一版本',
                        lineStyle: { color: '#e6a23c', width: 2 },
                        itemStyle: { color: '#e6a23c' },
                        areaStyle: { color: withAlpha('#e6a23c', 0.18) },
                    },
                ],
            },
        ],
    })
}

// ==================== 生命周期 ====================

const renderAll = () => {
    renderTrend()
    renderSource()
    renderSales()
    renderGauge()
    renderRadar()
    kpiList.forEach(renderSpark)
}

const resizeAll = () => {
    Object.values(charts).forEach((c) => c?.resize())
    Object.values(sparkCharts).forEach((c) => c?.resize())
}
// 窗口拖拽期间会连续触发 resize，防抖 100ms 减少全量重排
const debouncedResize = debounce(resizeAll, 100)

const disposeAll = () => {
    Object.values(charts).forEach((c) => c?.dispose())
    Object.values(sparkCharts).forEach((c) => c?.dispose())
}

// 窗口 resize 与页面卸载时自动清理
useEventListener(window, 'resize', debouncedResize)
useEventListener(window, 'beforeunload', saveWorkState)

onMounted(() => {
    tick()
    clockTimer = setInterval(tick, 1000)

    nextTick(() => {
        renderAll()
    })
})

onBeforeUnmount(() => {
    if (clockTimer) clearInterval(clockTimer)
    debouncedResize.cancel()
    saveWorkState()
    disposeAll()
})

// 暗/亮模式切换
watch(
    () => config.layout.dark,
    () => nextTick(renderAll)
)
</script>

<style scoped lang="scss">
$card-radius: 12px;

.dashboard {
    display: flex;
    flex-direction: column;
    gap: var(--ag-main-space);
}

// ============================================================
// 欢迎条
// ============================================================
.welcome-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--ag-main-space);
    align-items: stretch;
}

// 左侧文案卡
.welcome-bar-left {
    grid-column: span 3;
}

// 通用欢迎卡容器
.welcome-bar {
    position: relative;
    display: flex;
    align-items: center;
    min-width: 0;
    min-height: 160px;
    padding: 14px 28px;
    border-radius: $card-radius;
    background: #e1eaf9;
    color: var(--el-text-color-primary);
    box-shadow: 0 2px 12px rgba(64, 158, 255, 0.08);
    overflow: hidden;
    transition:
        box-shadow 0.35s ease,
        transform 0.35s ease;

    // 背景两处径向光斑装饰
    &::before,
    &::after {
        content: '';
        position: absolute;
        border-radius: 50%;
        pointer-events: none;
        opacity: 0.7;
        transition:
            transform 0.6s ease,
            opacity 0.6s ease;
    }
    &::before {
        top: -120px;
        right: -60px;
        width: 220px;
        height: 220px;
        background: radial-gradient(circle, rgba(64, 158, 255, 0.35) 0%, transparent 70%);
    }
    &::after {
        bottom: -100px;
        left: -40px;
        width: 180px;
        height: 180px;
        background: radial-gradient(circle, rgba(64, 158, 255, 0.22) 0%, transparent 70%);
    }

    &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 18px rgba(64, 158, 255, 0.16);
        &::before,
        &::after {
            opacity: 0.4;
        }
        &::before {
            transform: translate(-8px, 6px) scale(1.08);
        }
        &::after {
            transform: translate(8px, -6px) scale(1.08);
        }
        .hello {
            transform: translateX(2px);
        }
        .welcome-illustration {
            transform: translateY(-53%) scale(1.03);
        }
    }
}

.welcome-illustration {
    position: absolute;
    top: 50%;
    right: 24px;
    z-index: 0;
    width: 180px;
    max-width: 30%;
    height: auto;
    pointer-events: none;
    opacity: 0.9;
    transition: transform 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
    animation: welcome-float 4s ease-in-out infinite;
}

.welcome-content {
    position: relative;
    z-index: 1;
    flex: 1;
    min-width: 0;
    padding-right: 220px;
}

.hello {
    margin-bottom: 8px;
    font-size: 20px;
    font-weight: 700;
    color: var(--el-text-color-primary);
    transition: transform 0.3s ease;
}
.tip {
    font-size: 13px;
    line-height: 1.7;
    color: var(--el-text-color-primary);
}

.welcome-meta {
    display: flex;
    gap: 12px;
    margin-top: 12px;
}
.meta-item {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 3px 10px;
    font-size: 12px;
    color: var(--el-text-color-regular);
    background: rgba(255, 255, 255, 0.6);
    border-radius: 20px;
    transition:
        background 0.25s ease,
        transform 0.25s ease;

    .el-icon {
        font-size: 14px;
        color: var(--el-color-primary);
    }
    &:hover {
        background: rgba(255, 255, 255, 0.9);
        transform: translateY(-1px);
    }
}

@keyframes welcome-float {
    0%,
    100% {
        transform: translateY(-50%);
    }
    50% {
        transform: translateY(calc(-50% - 6px));
    }
}

// 右侧工作计时
.welcome-timer {
    grid-column: span 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 10px;
    padding: 0 12px;
}
.timer-label {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    font-weight: 500;
    color: var(--el-text-color-regular);

    .el-icon {
        font-size: 15px;
        color: var(--el-color-primary);
    }
}
.timer-time {
    display: flex;
    align-items: baseline;
    gap: 6px;

    .timer-group {
        display: contents;
    }
    .timer-unit {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 3px;
    }
    .num {
        font-size: 30px;
        font-weight: 700;
        font-variant-numeric: tabular-nums;
        line-height: 1;
        letter-spacing: 0.5px;
        color: var(--el-color-primary);
        transition: color 0.3s ease;
    }
    .txt {
        font-size: 11px;
        font-weight: 500;
        color: var(--el-text-color-secondary);
    }
    .colon {
        padding-bottom: 14px;
        font-size: 24px;
        font-weight: 700;
        line-height: 1;
        color: var(--el-color-primary-light-3);
        opacity: 0.7;
        animation: timer-blink 1s ease-in-out infinite;
    }

    &.paused {
        .num {
            color: var(--el-text-color-secondary);
        }
        .colon {
            animation: none;
            opacity: 0.25;
        }
    }
}

@keyframes timer-blink {
    0%,
    100% {
        opacity: 0.7;
    }
    50% {
        opacity: 0.2;
    }
}

.timer-btn {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 6px 16px;
    font-size: 13px;
    font-weight: 500;
    color: #ffffff;
    background: var(--el-color-primary);
    border: none;
    border-radius: 20px;
    outline: none;
    cursor: pointer;
    box-shadow: 0 3px 10px color-mix(in srgb, var(--el-color-primary) 25%, transparent);
    transition:
        transform 0.25s ease,
        box-shadow 0.25s ease,
        filter 0.25s ease;

    &:hover {
        transform: translateY(-1px);
        filter: brightness(1.08);
        box-shadow: 0 5px 14px color-mix(in srgb, var(--el-color-primary) 35%, transparent);
    }
    &:active {
        transform: translateY(0);
    }

    &.resting {
        background: var(--el-color-success);
        box-shadow: 0 3px 10px color-mix(in srgb, var(--el-color-success) 25%, transparent);
        &:hover {
            box-shadow: 0 5px 14px color-mix(in srgb, var(--el-color-success) 35%, transparent);
        }
    }
}

// ============================================================
// KPI 卡片
// ============================================================
.kpi-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: var(--ag-main-space);
}
.kpi-card {
    position: relative;
    padding: 18px 20px;
    background: var(--el-bg-color-overlay);
    border-radius: $card-radius;
    overflow: hidden;
    transition:
        transform 0.25s ease,
        box-shadow 0.25s ease;

    &:hover {
        transform: translateY(-1px);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    }
}
.kpi-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 13px;
    color: var(--el-text-color-secondary);
}
.kpi-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    font-size: 16px;
    color: var(--el-color-primary);
    background: color-mix(in srgb, var(--el-color-primary) 12%, transparent);
    border-radius: 8px;
}
.kpi-value {
    margin: 10px 0 6px;
    font-size: 28px;
    font-weight: 700;
    letter-spacing: 0.5px;
    color: var(--el-text-color-primary);
}
.kpi-bottom {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    white-space: nowrap;

    .kpi-trend {
        display: inline-flex;
        align-items: center;
        gap: 2px;
        flex-shrink: 0;
        padding: 2px 6px;
        font-weight: 600;
        border-radius: 4px;

        &.up {
            color: var(--el-color-success);
            background: color-mix(in srgb, var(--el-color-success) 12%, transparent);
        }
        &.down {
            color: var(--el-color-danger);
            background: color-mix(in srgb, var(--el-color-danger) 12%, transparent);
        }
    }
    .kpi-sub {
        flex-shrink: 0;
    }
    .kpi-spark {
        flex: 1;
        min-width: 0;
        height: 32px;
        margin-left: 8px;
    }
}

// ============================================================
// 图表卡片
// ============================================================
.chart-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: var(--ag-main-space);
}
.chart-card-span-2 {
    grid-column: span 2;
}
.chart-card {
    padding: 16px 18px;
    background: var(--el-bg-color-overlay);
    border-radius: $card-radius;
    transition: box-shadow 0.25s ease;

    &:hover {
        box-shadow: 0 6px 20px rgba(0, 0, 0, 0.06);
    }
}
.card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 8px;

    .card-title {
        position: relative;
        padding-left: 10px;
        font-size: 15px;
        font-weight: 600;
        color: var(--el-text-color-primary);

        &::before {
            content: '';
            position: absolute;
            left: 0;
            top: 50%;
            width: 3px;
            height: 14px;
            background: var(--el-color-primary);
            border-radius: 2px;
            transform: translateY(-50%);
        }
    }
}
.chart-body-lg {
    height: 320px;
}
.chart-body-md {
    height: 280px;
}

// 访问趋势的时间范围切换
.trend-ranges {
    display: inline-flex;
    padding: 2px;
    background: var(--el-fill-color-light);
    border-radius: 8px;
}
.range-chip {
    padding: 4px 12px;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    background: transparent;
    border: none;
    border-radius: 6px;
    outline: none;
    cursor: pointer;
    transition:
        color 0.25s ease,
        background 0.25s ease,
        box-shadow 0.25s ease;

    &:hover {
        color: var(--el-text-color-primary);
    }
    &.is-active {
        font-weight: 600;
        color: var(--el-color-primary);
        background: var(--el-bg-color);
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
    }
}

// ============================================================
// 暗黑模式
// ============================================================
html.dark {
    .welcome-bar {
        background: var(--ag-bg-color-overlay);
    }
}

// ============================================================
// 响应式
// ============================================================
@media (max-width: 1200px) {
    .kpi-row {
        grid-template-columns: repeat(2, 1fr);
    }
    .chart-row {
        grid-template-columns: 1fr;
    }
    .chart-card-span-2 {
        grid-column: auto;
    }
    .welcome-row {
        grid-template-columns: 1fr;
    }
    .welcome-bar-left,
    .welcome-timer {
        grid-column: auto;
    }
    .welcome-timer {
        min-height: 120px;
    }
}
@media (max-width: 640px) {
    .kpi-row {
        grid-template-columns: 1fr;
    }
    .welcome-bar {
        padding: 16px 20px;
    }
    .welcome-content {
        padding-right: 0;
    }
    .welcome-illustration {
        display: none;
    }
    .timer-time .num {
        font-size: 26px;
    }
}
</style>
