import App from '@/App.vue'
import { registerIcons } from '@/components/icon/index.js'
import { setupI18n } from '@/lang/index'
import router from '@/router/index'
import pinia from '@/stores/index'
import { registerDirectives } from '@/utils/directive'
import elementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/display.css'
import mitt from 'mitt'
import { createApp } from 'vue'
import '@/styles/index.scss'

async function start() {
    const app = createApp(App)

    app.use(pinia)
    app.use(router)
    app.use(elementPlus)

    await setupI18n(app)

    // 全局注册
    registerIcons(app) // icon
    registerDirectives(app) // 指令

    app.mount('#app')

    // 将 mitt 事件总线挂载到 vue 全局属性
    app.config.globalProperties.eventBus = mitt()
}

start()
