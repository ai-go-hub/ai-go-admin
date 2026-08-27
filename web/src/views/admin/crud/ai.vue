<template>
    <div>
        <el-dialog class="ag-operate-dialog crud-ai-dialog" v-model="model" :title="t('crud.index.ai')" width="50%" @opened="onDialogOpened">
            <div class="session-window">
                <TransitionGroup class="chat-records ag-scroll-style" name="el-fade-in" tag="div" @scroll.passive="onChatScroll">
                    <div class="placeholder" key="placeholder"></div>
                    <div class="chat-records-grow" key="chat-records-grow"></div>
                    <div v-for="(item, idx) in state.messages" :key="item.uuid || item.content" class="chat-record-item" :class="item.type">
                        <span v-if="item.type == 'system'">{{ item.content }}</span>
                        <template v-else>
                            <div class="record-avatar" v-if="item.type == 'me'">
                                <img draggable="false" :src="fullURL(adminInfo.avatar)" />
                            </div>
                            <div class="record-content-container">
                                <div class="chat-record-nickname">{{ item.nickname }}</div>
                                <div class="record-content" :id="`content-${idx}`">
                                    <div class="record-reasoning" v-if="item.type == 'you' && item.reasoning && item.reasoning.trim()">
                                        <div
                                            class="record-reasoning-header"
                                            :class="{ 'is-collapsed': !item.reasoningShow }"
                                            @click="item.reasoningShow = !item.reasoningShow"
                                        >
                                            <span>{{ t('crud.ai.reasoning') }}</span>
                                            <Icon name="el-arrow-down" size="12" />
                                        </div>
                                        <div class="record-reasoning-content" v-show="item.reasoningShow">
                                            {{ item.reasoning }}
                                        </div>
                                    </div>
                                    <pre
                                        v-if="item.designData"
                                        class="json-view language-json"
                                        v-html="Prism.highlight(JSON.stringify(item.designData, null, 4), Prism.languages.json, 'json')"
                                    ></pre>
                                    <template v-else>{{ item.content }}</template>
                                </div>
                                <div class="content-tags">
                                    <el-tooltip v-if="!item.loading" effect="dark" :content="t('crud.ai.copy')" placement="bottom">
                                        <el-tag @click="onCopy(item.content)">
                                            <Icon name="lucide-copy" size="12" color="var(--el-color-primary)" />
                                        </el-tag>
                                    </el-tooltip>
                                    <el-tooltip
                                        v-if="!item.loading && item.designData"
                                        effect="dark"
                                        :content="t('crud.ai.startDesign')"
                                        placement="bottom"
                                    >
                                        <el-tag type="success" @click="onStartDesign(item.designData)">
                                            <Icon name="lucide-play" size="12" color="var(--el-color-success)" />
                                        </el-tag>
                                    </el-tooltip>
                                    <el-tooltip v-if="item.loading" effect="dark" :content="t('crud.ai.stop')" placement="bottom">
                                        <Icon @click="onStopGenerate" class="is-loading message-loading" name="el-loading" size="14" />
                                    </el-tooltip>
                                </div>
                            </div>
                        </template>
                    </div>
                </TransitionGroup>
                <div class="message-textarea-container">
                    <textarea
                        ref="messageTextareaRef"
                        @input="onMessageInput"
                        class="ag-scroll-style"
                        id="message-text-input"
                        :placeholder="hasDesign ? t('crud.ai.placeholderModify') : t('crud.ai.placeholderDesign')"
                        rows="3"
                        @keydown="onTextareaKeydown"
                        v-model="state.message"
                    ></textarea>
                    <div class="message-textarea-footer">
                        <el-dropdown class="model-dropdown">
                            <span>
                                {{ state.modelList[state.model] }}
                                <Icon name="el-arrow-down" size="14" />
                            </span>
                            <template #dropdown>
                                <el-dropdown-menu>
                                    <el-dropdown-item v-for="(title, name) in state.modelList" :key="name" @click="onModelChange(name)">
                                        {{ title }}
                                    </el-dropdown-item>
                                </el-dropdown-menu>
                            </template>
                        </el-dropdown>
                        <el-button
                            type="primary"
                            :disabled="state.aiOutputMessageKey != null"
                            :loading="state.aiOutputMessageKey != null"
                            @click="onSendMessage"
                            size="small"
                        >
                            {{ t('crud.ai.send') }}
                        </el-button>
                    </div>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { crudAIStreamURL } from '@/api/admin/crud'
import { getAIConfig } from '@/api/admin/index'
import { useAdminInfo } from '@/stores/adminInfo'
import { copy, fullURL } from '@/utils/common'
import { uuid } from '@/utils/random'
import { fetchEventSource, type EventSourceMessage } from '@microsoft/fetch-event-source'
import { ElMessage, ElMessageBox } from 'element-plus'
import { debounce } from 'lodash-es'
import Prism from 'prismjs'
import 'prismjs/components/prism-json'
import 'prismjs/themes/prism-tomorrow.css'
import { computed, nextTick, reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import { MAX_AI_RETRY, retryPrompt, systemPrompt, temperature, topP, welcome } from './ai'
import { changeStep, type FieldItem } from './index'

interface MessageItem {
    uuid: string
    type: 'you' | 'me' | 'system'
    nickname: string
    content: string
    loading?: boolean
    designData?: TableDesignData | null
    reasoning?: string
    reasoningShow?: boolean
}

interface ChatMessage {
    role: 'system' | 'user' | 'assistant'
    content: string
}

interface TableDesignData {
    table: string
    comment: string
    fields: FieldItem[]
}

const { t } = useI18n()
const adminInfo = useAdminInfo()
const DEFAULT_RECORDS_HEIGHT = 'calc(100% - 93px)'
const messageTextareaRef = useTemplateRef('messageTextareaRef')

const model = defineModel<boolean>()
const state = reactive({
    message: '',
    messages: [
        {
            uuid: uuid(),
            type: 'you',
            nickname: 'AI',
            content: welcome,
        },
    ] as MessageItem[],
    aiOutputMessageKey: null as number | null,
    recordsHeight: DEFAULT_RECORDS_HEIGHT,
    model: '',
    modelList: {} as AnyObj,
})

// 流式请求控制器
let abortController: AbortController | null = null
// 递增的请求序号，用于忽略过期请求（如重试前旧流）的回调
let aiRequestSeq = 0
// JSON 解析失败重试计数
let aiRetryCount = 0
// 用户上滑后，本次对话停止自动跟随底部
let userScrolledUp = false

// 是否已产出首稿设计
const hasDesign = computed(() => state.messages.some((item) => item.designData))

const onMessageInput = debounce((el: Event) => {
    const elem = el.target as HTMLInputElement
    elem.style.height = 'auto'
    const height = elem.scrollHeight > 150 ? 150 : elem.scrollHeight
    state.recordsHeight = 'calc(100% - ' + (height + 38) + 'px)'
    elem.style.height = height + 'px'
    elem.scrollTop = height

    pullMessageScrollBar()
}, 300)

const onCopy = (content: string) => {
    if (!copy(content)) {
        ElMessage.error(t('common.operationFailed'))
    } else {
        ElMessage.success(t('common.operationSuccess'))
    }
}

/**
 * 将 AI 生成的设计 JSON 传入设计器并打开
 */
const onStartDesign = (design: TableDesignData) => {
    ElMessageBox.confirm(t('crud.ai.confirmStartDesign'), t('common.reminder'), {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
        closeOnClickModal: false,
    })
        .then(() => {
            changeStep('json', design)
            model.value = false
        })
        .catch(() => {
            /* 取消跳转 */
        })
}

/**
 * 发送后重置输入框高度与聊天区高度
 */
const resetMessageTextareaHeight = () => {
    onMessageInput.cancel()
    const textarea = messageTextareaRef.value
    if (textarea) {
        textarea.style.height = ''
    }
    state.recordsHeight = DEFAULT_RECORDS_HEIGHT
}

/**
 * 输入框按键事件监听
 * 单独 Enter 发送，Ctrl + Enter 换行
 */
const onTextareaKeydown = (e: KeyboardEvent) => {
    if (e.isComposing) return
    if (e.key !== 'Enter') return

    const textarea = e.target as HTMLTextAreaElement

    // Ctrl + Enter 插入换行
    if (e.ctrlKey && !e.shiftKey && !e.altKey && !e.metaKey) {
        e.preventDefault()
        const start = textarea.selectionStart ?? textarea.value.length
        const end = textarea.selectionEnd ?? textarea.value.length
        state.message = textarea.value.slice(0, start) + '\n' + textarea.value.slice(end)
        nextTick(() => {
            textarea.selectionStart = start + 1
            textarea.selectionEnd = start + 1
        })
        onMessageInput({ target: textarea } as unknown as Event)
        return
    }

    // 单独 Enter 发送
    if (!e.ctrlKey && !e.shiftKey && !e.altKey && !e.metaKey) {
        e.preventDefault()
        onSendMessage()
    }
}

const onModelChange = (val: string) => {
    state.model = val
}

/**
 * 组装发送给上游的对话上下文
 * 跳过当前正在生成的 AI 消息（其内容尚未完整）
 */
const buildChatMessages = (): ChatMessage[] => {
    const list: ChatMessage[] = [{ role: 'system', content: systemPrompt }]
    // state.messages 最新在前，倒序遍历还原时间正序，并跳过当前正在生成的 AI 消息
    for (let i = state.messages.length - 1; i >= 0; i--) {
        if (i === state.aiOutputMessageKey) continue
        const item = state.messages[i]
        if (item.type === 'me') {
            list.push({ role: 'user', content: item.content })
        } else if (item.type === 'you') {
            list.push({ role: 'assistant', content: item.content })
        }
    }
    return list
}

/**
 * 从 AI 输出中解析设计 JSON 数据
 */
const parseDesignJSON = (content: string): TableDesignData | null => {
    if (!content) return null
    // 去掉可能存在的 markdown 代码块围栏
    let text = content
        .trim()
        .replace(/^```(?:json)?\s*[\r\n]*/i, '')
        .replace(/```\s*$/, '')
    const start = text.indexOf('{')
    const end = text.lastIndexOf('}')
    if (start < 0 || end <= start) return null
    try {
        const obj = JSON.parse(text.slice(start, end + 1))
        if (
            obj &&
            typeof obj === 'object' &&
            typeof obj.table === 'string' &&
            obj.table &&
            typeof obj.comment === 'string' &&
            Array.isArray(obj.fields) &&
            obj.fields.length > 0
        ) {
            return { table: obj.table, comment: obj.comment, fields: obj.fields }
        }
    } catch {
        /* ignore */
    }
    return null
}

/**
 * 结束当前 AI 消息渲染: 渲染完毕或 onStopGenerate 后调用
 * aiOutputMessageKey 是 state.messages 的 key，置为 null 表示当前无正在渲染的 AI 消息
 * @param seq 请求序号，过期请求的回调直接忽略
 * @param allowRetry 是否允许 JSON 解析失败时自动要求模型重试
 */
const finishGenerate = (seq: number, allowRetry: boolean) => {
    if (seq !== aiRequestSeq) return
    if (state.aiOutputMessageKey === null) return
    const aiMessage = state.messages[state.aiOutputMessageKey]
    if (aiMessage) {
        aiMessage.loading = false
    }
    state.aiOutputMessageKey = null
    abortController = null

    if (allowRetry && aiMessage) {
        const parsed = parseDesignJSON(aiMessage.content)
        if (parsed) {
            aiMessage.designData = parsed
            aiMessage.reasoningShow = false
        } else if (aiRetryCount < MAX_AI_RETRY) {
            // JSON 不符合要求，自动要求模型重试
            aiRetryCount++
            sendMessage(retryPrompt)
            return
        } else {
            ElMessage.error(t('crud.ai.designError', { count: MAX_AI_RETRY }))
        }
    }

    aiRetryCount = 0
    pullMessageScrollBar()
}

/**
 * 发送一条用户消息并开始 AI 流式回复（用户输入与自动重试共用）
 */
const sendMessage = (content: string) => {
    state.messages.unshift({
        uuid: uuid(),
        type: 'me',
        nickname: adminInfo.nickname,
        content,
    })

    // AI 占位消息插到最前，其 index 恒为 0
    state.messages.unshift({
        uuid: uuid(),
        type: 'you',
        nickname: 'AI',
        content: '',
        reasoning: '',
        reasoningShow: true,
        loading: true,
    })
    state.aiOutputMessageKey = 0

    const seq = ++aiRequestSeq
    abortController = new AbortController()
    pullMessageScrollBar()

    fetchEventSource(crudAIStreamURL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${adminInfo.token}`,
        },
        body: JSON.stringify({
            model: state.model,
            messages: buildChatMessages(),
            temperature,
            top_p: topP,
        }),
        signal: abortController!.signal,
        openWhenHidden: true,
        onmessage: (ev) => onAIEvent(ev, seq),
        onclose: () => finishGenerate(seq, true),
        onerror: (err) => {
            finishGenerate(seq, false)
            // 主动停止时无需提示
            if (err?.name !== 'AbortError') {
                ElMessage.error(t('common.networkError'))
            }
        },
    })
}

/**
 * 发送消息
 */
const onSendMessage = () => {
    const content = state.message.trim()
    if (!content || state.aiOutputMessageKey !== null) {
        return
    }

    state.message = ''
    aiRetryCount = 0
    userScrolledUp = false
    resetMessageTextareaHeight()
    sendMessage(content)
}

/**
 * 处理 AI 流式事件，增量渲染到当前 AI 消息
 */
const onAIEvent = (ev: EventSourceMessage, seq: number) => {
    if (seq !== aiRequestSeq) return
    const aiMessage = state.aiOutputMessageKey !== null ? state.messages[state.aiOutputMessageKey] : null
    if (!aiMessage) return

    // 后端 SSE 错误事件（如 AI 配置不完整、上游接口错误），不触发重试
    if (ev.event === 'error') {
        let msg = t('common.operationFailed')
        try {
            const data = JSON.parse(ev.data)
            if (data.message) msg = data.message
        } catch {
            /* ignore */
        }
        ElMessage.error(msg)
        finishGenerate(seq, false)
        return
    }

    // 流结束标记
    if (ev.data === '[DONE]') {
        finishGenerate(seq, true)
        return
    }

    let parsed: any
    try {
        parsed = JSON.parse(ev.data)
    } catch {
        return
    }

    // OpenAI Responses API: 思考过程增量
    if (parsed.type === 'response.reasoning_text.delta' || parsed.type === 'response.reasoning_summary_text.delta') {
        aiMessage.reasoning = (aiMessage.reasoning ?? '') + (parsed.delta ?? '')
        pullMessageScrollBar()
        return
    }

    // OpenAI Responses API: 文本增量事件
    if (parsed.type === 'response.output_text.delta') {
        aiMessage.content += parsed.delta ?? ''
        pullMessageScrollBar()
        return
    }

    // 输出完成 / 失败
    if (parsed.type === 'response.completed' || parsed.type === 'response.failed') {
        finishGenerate(seq, true)
    }
}

/**
 * 停止生成
 */
const onStopGenerate = () => {
    abortController?.abort()
    finishGenerate(aiRequestSeq, false)
}

/**
 * 窗口滚动条到底
 */
const pullMessageScrollBar = () => {
    // 用户上滑后，本次对话不再自动滑到底部
    if (userScrolledUp) return

    const messageEl = document.getElementsByClassName('chat-records')[0] as HTMLDivElement
    nextTick(() => {
        messageEl.scrollTop = 0
    })
}

/**
 * 聊天区上滑至离开底部则本次对话停止自动跟随
 */
const onChatScroll = (e: Event) => {
    const el = e.target as HTMLDivElement
    if (Math.abs(el.scrollTop) > 10) {
        userScrolledUp = true
    }
}

/**
 * 窗口打开后聚焦输入框
 */
const onDialogOpened = () => {
    nextTick(() => {
        messageTextareaRef.value?.focus()
    })
}

const getData = () => {
    getAIConfig().then((res) => {
        state.model = res.data.data.default_model ?? ''

        let modelList: any[] = []
        try {
            modelList = JSON.parse(res.data.data.model_list || '[]')
        } catch {
            modelList = []
        }
        const modelListKv: AnyObj = {}
        for (const key in modelList) {
            modelListKv[modelList[key]['value']] = modelList[key]['key']
        }

        state.modelList = modelListKv
    })
}
getData()
</script>

<style scoped lang="scss">
:deep(.crud-ai-dialog) {
    padding-bottom: 15px;
}
.session-window {
    position: relative;
    display: block;
    width: 100%;
    height: 100%;
    .chat-records {
        width: 100%;
        display: flex;
        flex-direction: column-reverse;
        height: v-bind('state.recordsHeight');
        margin-right: -10px;
        padding: 0;
        overflow-y: auto;
        overflow-x: hidden;
        box-sizing: border-box;
        .chat-records-grow {
            flex-grow: 1;
            flex-shrink: 1;
        }
        .chat-record-item {
            display: flex;
            padding-top: 16px;
            align-items: flex-start;
        }
        .system span {
            font-size: 12px;
            display: inline-block;
            background: var(--el-color-info-light-9);
            color: var(--el-text-color-primary);
            line-height: 12px;
            border-radius: 5px;
            padding: 3px 10px;
            text-align: center;
            word-wrap: break-word;
            word-break: break-all;
            margin: 0 auto;
            margin-bottom: 6px;
        }
        .record-avatar {
            position: relative;
            display: inline-block;
            background: transparent;
            width: 36px;
            height: 36px;
            cursor: pointer;
            user-select: none;
            img {
                height: 100%;
                width: 100%;
                border-radius: 6px;
            }
        }
        .record-content-container {
            position: relative;
            max-width: 75%;
            margin: 0 13px;
        }
        .record-reasoning {
            margin-bottom: 4px;
            .record-reasoning-header {
                display: flex;
                align-items: center;
                gap: 4px;
                font-size: 12px;
                color: var(--el-text-color-secondary);
                cursor: pointer;
                user-select: none;
                .icon {
                    transition: transform 0.2s;
                }
                &.is-collapsed .icon {
                    transform: rotate(-90deg);
                }
            }
            .record-reasoning-content {
                margin-top: 4px;
                font-size: 12px;
                line-height: 18px;
                color: var(--el-text-color-secondary);
                background: var(--el-color-info-light-9);
                border-radius: 5px;
                white-space: pre-wrap;
                word-break: break-all;
            }
        }
        .me {
            flex-direction: row-reverse;
            display: flex;
            justify-content: flex-start;
            align-content: center;
            padding-right: 6px;
        }
        .chat-record-nickname {
            color: var(--el-text-color-secondary);
            padding-bottom: 3px;
        }
        .you .chat-record-nickname {
            text-align: left;
        }
        .you .record-content {
            color: #000;
            background: var(--el-color-info-light-9);
            .el-link {
                --el-link-text-color: var(--el-color-success);
            }
        }
        .me .chat-record-nickname {
            text-align: right;
        }
        .me .record-content {
            margin-left: auto;
            color: var(--el-color-white);
            background: var(--el-color-primary);
            .el-link {
                --el-link-text-color: var(--el-color-warning);
                &:hover {
                    --el-link-hover-text-color: var(--el-color-warning);
                }
            }
        }
        .you .record-content-container:before {
            left: -4px;
            background: var(--el-color-info-light-9);
        }
        .me .record-content-container:before {
            right: -4px;
            background: var(--el-color-primary);
        }
        .record-content-container:before {
            position: absolute;
            top: 26px;
            display: block;
            width: 8px;
            height: 6px;
            content: '\00a0';
            -webkit-transform: rotate(29deg) skew(-35deg);
            transform: rotate(29deg) skew(-35deg);
        }
        .record-content {
            font-size: 14px;
            color: var(--el-text-color-primary);
            padding: 10px;
            border-radius: 5px;
            position: relative;
            width: fit-content;
            max-width: 100%;
            word-wrap: break-word;
            word-break: break-all;
        }
        .message-loading {
            margin: 4px 0;
            cursor: pointer;
        }
        .json-view {
            margin: 0;
            border-radius: 5px;
        }
        .content-tags {
            margin-top: 10px;
            margin-bottom: 5px;
            .el-tag {
                margin-right: 4px;
                margin-bottom: 4px;
                cursor: pointer;
            }
        }
        .placeholder {
            height: 20px;
            width: 100%;
        }
    }
    .message-textarea-container {
        width: 100%;
        border-top: 1px solid var(--el-color-info-light-9);
        margin-right: -10px;
        padding-top: 2px;
        textarea {
            width: 100%;
            border: none;
            outline: none;
            background-color: transparent;
            padding: 5px 0;
            line-height: 15px;
            resize: none;
        }
        .message-textarea-footer {
            display: flex;
            align-items: center;
            justify-content: flex-end;
            font-size: var(--el-font-size-small);
            color: var(--el-text-color-placeholder);
            margin-top: 8px;
            .model-dropdown {
                padding: 0 8px;
                cursor: pointer;
                span {
                    outline: none;
                    display: flex;
                    align-items: center;
                    color: var(--el-text-color-placeholder);
                    .icon {
                        margin-left: 4px;
                    }
                }
            }
        }
    }
}
</style>
