import { readdirSync, writeFile } from 'fs'
import { trimEnd } from 'lodash-es'
import { spawn } from 'node:child_process'

function formatTime(): string {
    const now = new Date()
    return now.toTimeString().slice(0, 8)
}

function gray(text: string): string {
    return `\x1b[90m${text}\x1b[0m`
}

function cyan(text: string): string {
    return `\x1b[36m${text}\x1b[0m`
}

// ========================== 生成 tcr.d.ts ===============================

function getVueFileNames(dir: string) {
    const dirents = readdirSync(dir, {
        withFileTypes: true,
    })
    const fileNames: string[] = []
    for (const dirent of dirents) {
        if (!dirent.isDirectory()) fileNames.push(dirent.name.replace('.vue', ''))
    }
    return fileNames
}

/**
 * 生成 ./types/tcr.d.ts 文件
 */
const buildTableCellRendererType = () => {
    let cellRenderer = getVueFileNames('./src/components/table/cellRenderer/')

    // 增加 slot，去除 default
    cellRenderer.push('slot')
    cellRenderer = cellRenderer.filter((item) => item !== 'default')

    let tableCellRendererContent = '/**\n * 可用的表格单元格渲染器，以 ./src/components/table/cellRenderer/ 目录中的文件名自动生成\n */'
    tableCellRendererContent += '\ntype TableCellRenderer =\n    | '
    for (const key in cellRenderer) {
        tableCellRendererContent += `'${cellRenderer[key]}'\n    | `
    }
    tableCellRendererContent = trimEnd(tableCellRendererContent, '    | ')

    writeFile('./types/tcr.d.ts', tableCellRendererContent, 'utf-8', (err) => {
        if (err) throw err
    })

    console.log(`${gray(formatTime())} ${cyan('[table]')} updated: types/tcr.d.ts`)
}

// ========================== 启动 Vite 开发服务器 ===============================

async function start() {
    // 1. 生成 tcr.d.ts
    buildTableCellRendererType()

    // 2. 启动 Vite 开发服务器
    const vite = spawn('vite', ['--force'], {
        stdio: 'inherit',
        shell: true,
    })

    vite.on('exit', (code) => {
        process.exit(code ?? 0)
    })
}

start()
