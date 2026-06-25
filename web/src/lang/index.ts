import { merge, set } from 'lodash-es'
import type { App } from 'vue'
import { createI18n } from 'vue-i18n'
import { parse as parseYaml } from 'yaml'
import { useConfig } from '/@/stores/config'

/**
 * 支持的语言类型
 */
export type LangKey = 'zh-cn' | 'en'

/**
 * 支持的语言列表
 */
export const langs: LangKey[] = ['zh-cn', 'en']

/**
 * 语言显示名称
 */
export const langNames: Record<LangKey, string> = {
    en: 'English',
    'zh-cn': '简体中文',
}

/**
 * i18n 实例
 */
const i18n = createI18n({
    legacy: false,
    locale: 'zh-cn',
    fallbackLocale: 'zh-cn',
    messages: {},
})

// 使用 vite import.meta.glob 批量导入 lang 目录下所有 .yaml 文件（包括子目录）
const langGlobs = {
    en: import.meta.glob('./en/**/*.yaml', { query: '?raw', import: 'default' }),
    'zh-cn': import.meta.glob('./zh-cn/**/*.yaml', { query: '?raw', import: 'default' }),
}

/**
 * 设置 i18n，并为 vue 安装 i18n 插件
 */
export async function setupI18n(app: App): Promise<void> {
    const config = useConfig()
    i18n.global.fallbackLocale.value = config.lang.fallback

    // 初始化当前语言包
    await setLang(config.lang.active)

    app.use(i18n)
}

/**
 * 设置语言
 * @param lang 语言标识
 */
export async function setLang(lang: LangKey): Promise<void> {
    await loadMessages(lang)

    const config = useConfig()
    i18n.global.locale.value = lang
    config.setLang(lang)
}

/**
 * 懒加载语言包
 * @param lang 语言标识
 */
export async function loadMessages(lang: LangKey): Promise<void> {
    // 如果已加载则跳过
    if (i18n.global.availableLocales.includes(lang)) {
        return
    }

    try {
        // 批量加载 lang 目录下所有 .yaml 文件
        const glob = langGlobs[lang]
        const promises = Object.entries(glob).map(async ([path, loader]) => {
            const raw = (await loader()) as string
            const data = parseYaml(raw)
            return { path, data }
        })
        const modules = await Promise.all(promises)

        // 按文件路径构建嵌套的 messages 结构
        const mergedMessages: Record<string, any> = {}
        for (const { path, data } of modules) {
            if (typeof data !== 'object' || data === null) {
                continue
            }
            const keys = filePathToKeys(lang, path)
            if (keys.length === 0) {
                // 合并到顶层
                merge(mergedMessages, data)
            } else {
                // 子模块 — 按路径嵌套
                merge(mergedMessages, set({}, keys, data))
            }
        }

        i18n.global.setLocaleMessage(lang, mergedMessages)
    } catch (error) {
        console.error(`Failed to load lang: ${lang}`, error)
    }
}

const filePathToKeys = (lang: LangKey, path: string) => {
    const langPathPrefix = `/${lang}`
    const pathName = path.slice(path.lastIndexOf(langPathPrefix) + (langPathPrefix.length + 1), path.lastIndexOf('.'))
    const keys = pathName.split('/')

    return keys
}

export default i18n
