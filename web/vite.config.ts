import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import type { ConfigEnv, UserConfig } from 'vite'
import { loadEnv } from 'vite'
import { lucideIconSplitPlugin } from './src/components/icon/vitePlugin.ts'
import { customHotUpdate } from './src/utils/vite.ts'

// https://vitejs.cn/config/
const viteConfig = ({ mode }: ConfigEnv): UserConfig => {
    const { VITE_PORT, VITE_OPEN, VITE_BASE_PATH, VITE_OUT_DIR } = loadEnv(mode, process.cwd())

    return {
        plugins: [vue(), lucideIconSplitPlugin(), customHotUpdate()],
        root: process.cwd(),
        resolve: {
            alias: {
                '@': resolve(import.meta.dirname, 'src'),
            },
        },
        base: VITE_BASE_PATH,
        server: {
            port: parseInt(VITE_PORT),
            open: VITE_OPEN != 'false',
        },
        build: {
            cssCodeSplit: false,
            sourcemap: false,
            outDir: VITE_OUT_DIR,
            emptyOutDir: true,
            chunkSizeWarningLimit: 1500,
        },
    }
}

export default viteConfig
