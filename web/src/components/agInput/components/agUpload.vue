<template>
    <div
        ref="wrapperRef"
        class="w100 ag-upload-wrapper"
        :class="{ 'is-drag-over': state.isDragOver, 'is-disabled': state.attrs.disabled }"
        tabindex="0"
        @paste="onPaste"
    >
        <el-upload
            ref="uploadRef"
            class="ag-upload"
            :class="[
                type,
                state.attrs.disabled ? 'is-disabled' : '',
                hideImagePlusOnOverLimit && state.attrs.limit && state.fileList.length >= state.attrs.limit ? 'is-limit-reached' : '',
            ]"
            v-model:file-list="state.fileList"
            :auto-upload="false"
            @change="onElChange"
            @remove="onElRemove"
            @preview="onElPreview"
            @exceed="onElExceed"
            v-bind="state.attrs"
            :key="state.key"
        >
            <template v-if="!$slots.default" #default>
                <template v-if="isImageType">
                    <el-tooltip :content="t('common.copyPasteUploadTip')" placement="top" :disabled="!!state.attrs.disabled || state.isDragOver">
                        <div
                            class="ag-upload-trigger"
                            @mouseenter="onTriggerFocus"
                            @dragenter.prevent="onDragEnter"
                            @dragover.prevent="onDragOver"
                            @dragleave.prevent="onDragLeave"
                            @drop.prevent="onDrop"
                        >
                            <Icon class="ag-upload-icon" name="el-plus" size="30" color="#c0c4cc" />
                            <div @click.stop="onScreenshotUpload" class="ag-upload-screenshot">
                                {{ t('common.screenshotUpload') }}
                            </div>
                            <div v-if="state.isDragOver && !state.attrs.disabled" class="ag-upload-drag-mask">
                                <Icon name="el-upload-filled" size="40" color="var(--el-color-primary)" />
                                <span>{{ t('common.releaseToUpload') }}</span>
                            </div>
                        </div>
                    </el-tooltip>
                </template>
                <template v-else>
                    <el-tooltip :content="t('common.copyPasteUploadTip')" placement="top" :disabled="!!state.attrs.disabled || state.isDragOver">
                        <div
                            class="ag-upload-trigger ag-upload-trigger-file"
                            @mouseenter="onTriggerFocus"
                            @dragenter.prevent="onDragEnter"
                            @dragover.prevent="onDragOver"
                            @dragleave.prevent="onDragLeave"
                            @drop.prevent="onDrop"
                        >
                            <el-button type="primary">
                                <Icon name="el-plus" color="#ffffff" />
                                <span>{{ t('common.upload') }}</span>
                            </el-button>
                            <div v-if="state.isDragOver && !state.attrs.disabled" class="ag-upload-drag-mask">
                                <span>{{ t('common.releaseToUpload') }}</span>
                            </div>
                        </div>
                    </el-tooltip>
                </template>
            </template>

            <template v-for="(_slot, name) in $slots" #[name]="scopedData">
                <slot :name="name" v-bind="scopedData"></slot>
            </template>
        </el-upload>
        <el-dialog v-model="state.preview.show" :append-to-body="true" :destroy-on-close="true" class="ag-upload-preview">
            <div class="ag-upload-preview-scroll">
                <img :src="state.preview.url" class="ag-upload-preview-img" alt="" />
            </div>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import type { AxiosProgressEvent } from 'axios'
import type { UploadFiles, UploadProps, UploadRawFile, UploadUserFile } from 'element-plus'
import { ElMessage, genFileId } from 'element-plus'
import { cloneDeep } from 'lodash-es'
import Sortable from 'sortablejs'
import { computed, nextTick, onMounted, reactive, useAttrs, useTemplateRef, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { upload as uploadFile } from '/@/api/common'
import { stringToArray } from '/@/components/agInput/helper'
import { arrayFullURL, fullURL, getArrayKey, getFileNameFromPath } from '/@/utils/common'
import { uuid } from '/@/utils/random'

// 禁用 Attributes 自动继承
defineOptions({
    inheritAttrs: false,
})

interface Props extends /* @vue-ignore */ Partial<UploadProps> {
    type: 'image' | 'images' | 'file' | 'files'
    modelValue: string | string[]
    // 业务分类，透传到上传接口 topic
    topic?: string
    // 存储驱动，透传到上传接口 driver
    driver?: string
    // 返回带域名的完整路径
    returnFullUrl?: boolean
    // 强制使用本地存储驱动
    forceLocal?: boolean
    // 在上传数量达到限制时隐藏图片上传按钮
    hideImagePlusOnOverLimit?: boolean
}
interface UploadFileExt extends UploadUserFile {
    serverUrl?: string
}
interface UploadProgressEvent extends AxiosProgressEvent {
    percent: number
}

type ReturnType = 'string' | 'array'

const props = withDefaults(defineProps<Props>(), {
    type: 'image',
    modelValue: () => [],
    topic: '',
    driver: '',
    returnFullUrl: false,
    forceLocal: false,
    hideImagePlusOnOverLimit: false,
})

const emits = defineEmits<{
    (e: 'update:modelValue', value: string | string[]): void
}>()

const { t } = useI18n()
const attrs = useAttrs()
const uploadRef = useTemplateRef('uploadRef')
const wrapperRef = useTemplateRef('wrapperRef')

const isImageType = computed(() => props.type == 'image' || props.type == 'images')

const state: {
    key: string
    // 绑定到 el-upload 的属性对象
    attrs: Partial<UploadProps>
    events: anyObj
    isDragOver: boolean
    dragCounter: number
    uploadCounter: number
    fileList: UploadFileExt[]
    defaultReturnType: ReturnType
    preview: { show: boolean; url: string }
} = reactive({
    key: uuid(),
    attrs: {},
    events: {},
    fileList: [],
    dragCounter: 0,
    uploadCounter: 0,
    isDragOver: false,
    defaultReturnType: 'string',
    preview: { show: false, url: '' },
})

/**
 * 需要管理的事件列表（使用 triggerEvent 触发）
 */
const eventNameMap = {
    // el-upload 的钩子函数（它们是 props 而非 emit，el-upload 组件在 template 上已使用，故需手动触发）
    change: ['onChange', 'on-change'],
    remove: ['onRemove', 'on-remove'],
    preview: ['onPreview', 'on-preview'],
    exceed: ['onExceed', 'on-exceed'],

    // 自定义上传方法需要手动触发的钩子
    beforeUpload: ['onBeforeUpload', 'before-upload'],
    progress: ['onProgress', 'on-progress'],
    success: ['onSuccess', 'on-success'],
    error: ['onError', 'on-error'],
}

const triggerEvent = (name: keyof typeof eventNameMap, args: any[]) => {
    const evtNames = eventNameMap[name]
    for (const evtName of evtNames) {
        // 有一个返回 false 即中断
        if (typeof state.events[evtName] === 'function' && state.events[evtName](...args) === false) return false
    }
}

/**
 * 将文件加入 el-upload 并开始上传
 */
const startUploadFile = (rawFile: UploadRawFile) => {
    rawFile.uid = genFileId()
    if (state.attrs.limit && state.fileList.length >= state.attrs.limit) {
        onElExceed([rawFile])
    } else {
        uploadRef.value!.handleStart(rawFile)
    }
}

/**
 * 从拖拽/粘贴数据中收集可上传的文件（图片类型仅收集图片）
 */
const collectValidFiles = (fileList: FileList | File[]) => {
    const files = Array.from(fileList)
    if (!files.length) return []
    if (isImageType.value) return files.filter((file) => file.type.startsWith('image/'))
    return files
}

/**
 * 批量将文件加入上传队列
 */
const uploadInBatch = (validFiles: File[]) => {
    if (!validFiles.length) return false

    const limit = state.attrs.limit
    let remaining = limit ? Math.max(limit - state.fileList.length, 0) : validFiles.length

    for (const file of validFiles) {
        if (limit && remaining <= 0) {
            // limit=1 时依然触发一次以走 exceed 逻辑替换旧文件
            if (limit === 1) startUploadFile(file as UploadRawFile)
            break
        }
        startUploadFile(file as UploadRawFile)
        if (limit) remaining--
    }
    return true
}

/**
 * 从粘贴板事件中收集可上传的文件
 */
const collectPasteFiles = (event: ClipboardEvent) => {
    const files: File[] = []
    const clipboardData = event.clipboardData
    if (!clipboardData) return []

    if (clipboardData.files.length) {
        files.push(...Array.from(clipboardData.files))
    } else {
        Array.from(clipboardData.items).forEach((item) => {
            if (item.kind === 'file') {
                const file = item.getAsFile()
                if (file) files.push(file)
            }
        })
    }
    return collectValidFiles(files)
}

const onTriggerFocus = () => {
    if (state.attrs.disabled) return
    wrapperRef.value?.focus({ preventScroll: true })
}

const onPaste = (event: ClipboardEvent) => {
    if (state.attrs.disabled) return

    const validFiles = collectPasteFiles(event)
    if (!validFiles.length) {
        ElMessage.warning(isImageType.value ? t('common.noImageInClipboard') : t('common.noValidFileInPaste'))
        return
    }
    event.preventDefault()
    uploadInBatch(validFiles)
}

const resetDragState = () => {
    state.isDragOver = false
    state.dragCounter = 0
}

const onDragEnter = () => {
    if (state.attrs.disabled) return
    state.dragCounter++
    state.isDragOver = true
}
const onDragOver = () => {
    if (state.attrs.disabled) return
    state.isDragOver = true
}
const onDragLeave = () => {
    if (state.attrs.disabled) return
    state.dragCounter--
    if (state.dragCounter <= 0) resetDragState()
}
const onDrop = (event: DragEvent) => {
    resetDragState()
    if (state.attrs.disabled) return

    const validFiles = collectValidFiles(event.dataTransfer?.files || [])
    if (!uploadInBatch(validFiles)) {
        ElMessage.warning(isImageType.value ? t('common.noValidImageInDrop') : t('common.noValidFileInDrop'))
    }
}

const onElChange = (file: UploadFileExt, files: UploadFiles) => {
    // 用 files 中的对象替换 file，以便修改属性等操作
    const fileIndex = getArrayKey(files, 'uid', file.uid!)
    if (fileIndex === false) return

    file = files[fileIndex] as UploadFileExt
    if (!file || !file.raw) return
    if (triggerEvent('beforeUpload', [file]) === false) return

    const topic = props.topic || undefined
    const driver = props.forceLocal ? 'local' : props.driver || undefined

    file.status = 'uploading'
    state.uploadCounter++
    uploadFile(
        { file: file.raw as File, topic, driver },
        {
            onUploadProgress: (evt: AxiosProgressEvent) => {
                const progressEvt = evt as UploadProgressEvent
                if (evt.total && evt.total > 0 && ['ready', 'uploading'].includes(file.status!)) {
                    progressEvt.percent = (evt.loaded / evt.total) * 100
                    file.status = 'uploading'
                    file.percentage = Math.round(progressEvt.percent)
                    triggerEvent('progress', [progressEvt, file, files])
                }
            },
        }
    )
        .then((res) => {
            const body = res.data
            if (body.code === 0) {
                file.serverUrl = body.data.url
                file.status = 'success'
                emits('update:modelValue', getAllUrls())
                triggerEvent('success', [body, file, files])
            } else {
                file.status = 'fail'
                files.splice(fileIndex, 1)
                triggerEvent('error', [body, file, files])
            }
        })
        .catch((err) => {
            file.status = 'fail'
            files.splice(fileIndex, 1)
            triggerEvent('error', [err, file, files])
        })
        .finally(() => {
            state.uploadCounter--
            afterFileChange(file, files)
        })
}

const onElRemove = (file: UploadUserFile, files: UploadFiles) => {
    triggerEvent('remove', [file, files])
    afterFileChange(file, files)
    nextTick(() => {
        emits('update:modelValue', getAllUrls())
    })
}

const onElPreview = (file: UploadFileExt) => {
    triggerEvent('preview', [file])
    if (!file || !file.serverUrl) return
    if (props.type == 'file' || props.type == 'files') {
        window.open(fullURL(file.serverUrl))
        return
    }
    state.preview.show = true
    state.preview.url = fullURL(file.serverUrl)
}

const onElExceed = (files: UploadUserFile[]) => {
    const file = files[0] as UploadRawFile
    file.uid = genFileId()
    uploadRef.value!.handleStart(file)
    triggerEvent('exceed', [file, state.fileList])
}

/**
 * 初始化文件/图片的拖拽排序
 */
const initSort = () => {
    if (state.attrs.showFileList === false) return
    nextTick(() => {
        const uploadListEl = uploadRef.value?.$el.querySelector('.el-upload-list') as HTMLElement | null
        if (!uploadListEl) return
        const uploadItemEls = uploadListEl.getElementsByClassName('el-upload-list__item')
        if (uploadItemEls.length < 2) return
        Sortable.create(uploadListEl, {
            animation: 200,
            draggable: '.el-upload-list__item',
            onEnd: (evt: Sortable.SortableEvent) => {
                if (evt.oldIndex == evt.newIndex) return
                const { oldIndex, newIndex } = evt
                ;[state.fileList[oldIndex!], state.fileList[newIndex!]] = [state.fileList[newIndex!], state.fileList[oldIndex!]]
                emits('update:modelValue', getAllUrls())
            },
        })
    })
}

onMounted(() => {
    // 收集透传属性中的事件钩子和普通属性
    const passthroughEvents = new Set<string>()
    for (const key in eventNameMap) {
        for (const name of eventNameMap[key as keyof typeof eventNameMap]) passthroughEvents.add(name)
    }

    let mergedAttrs: anyObj = {}
    for (const attrKey in attrs) {
        if (passthroughEvents.has(attrKey)) {
            state.events[attrKey] = attrs[attrKey]
        } else {
            mergedAttrs[attrKey] = attrs[attrKey]
        }
    }

    // 单文件/单图默认 limit=1，多文件/多图默认 multiple
    if (props.type == 'image' || props.type == 'file') {
        mergedAttrs = { ...mergedAttrs, limit: 1 }
    } else {
        mergedAttrs = { ...mergedAttrs, multiple: true }
    }

    // 图片类型默认使用 picture-card 卡片列表
    if (isImageType.value) {
        mergedAttrs = { ...mergedAttrs, accept: 'image/*', listType: 'picture-card' }
    }

    state.attrs = mergedAttrs

    init(props.modelValue)
    initSort()
})

/**
 * 截断超过 limit 的文件，返回是否发生了截断
 */
const enforceLimit = () => {
    if (state.attrs.limit && state.fileList.length > state.attrs.limit) {
        state.fileList = state.fileList.slice(state.fileList.length - state.attrs.limit)
        return true
    }
    return false
}

const init = (modelValue: string | string[]) => {
    const urls = stringToArray(modelValue as string)
    const isSingleType = props.type == 'file' || props.type == 'image'

    state.fileList = urls.map((url) => ({
        name: getFileNameFromPath(url),
        url: fullURL(url),
        serverUrl: url,
    }))
    state.defaultReturnType = typeof modelValue === 'string' || isSingleType ? 'string' : 'array'

    // 超出过滤 || 需要返回完整 URL
    if (enforceLimit() || props.returnFullUrl) {
        emits('update:modelValue', getAllUrls())
    }
    state.key = uuid()
}

/**
 * 获取当前所有资源路径的列表
 */
const getAllUrls = (returnType: ReturnType = state.defaultReturnType) => {
    enforceLimit()
    let urlList = state.fileList.map((f) => f.serverUrl).filter((u): u is string => !!u)
    if (props.returnFullUrl) urlList = arrayFullURL(urlList)
    return returnType === 'string' ? urlList.join(',') : urlList
}

/**
 * 文件状态改变（成功/失败/移除）后的通用回调
 */
const afterFileChange = (file: string | string[] | UploadFileExt, files: UploadFileExt[]) => {
    initSort()
    triggerEvent('change', [file, files])
}

/**
 * 从剪贴板读取图片并上传（图片类型专用的截图上传按钮）
 */
const onScreenshotUpload = async () => {
    if (state.attrs.disabled) return

    if (!navigator.clipboard?.read) {
        ElMessage.warning(t('common.clipboardReadNotSupported'))
        return
    }

    try {
        const clipboardItems = await navigator.clipboard.read()
        let imageBlob: Blob | null = null
        let mimeType = ''

        for (const item of clipboardItems) {
            const imageType = item.types.find((type) => type.startsWith('image/'))
            if (imageType) {
                imageBlob = await item.getType(imageType)
                mimeType = imageType
                break
            }
        }

        if (!imageBlob) {
            ElMessage.warning(t('common.noImageInClipboard'))
            return
        }

        const ext = mimeType.split('/')[1]?.replace('jpeg', 'jpg') || 'png'
        const file = new File([imageBlob], `paste-${Date.now()}.${ext}`, { type: mimeType }) as UploadRawFile
        startUploadFile(file)
    } catch {
        ElMessage.warning(t('common.readClipboardFailed'))
    }
}

const getRef = () => uploadRef.value

defineExpose({
    getRef,
})

watch(
    () => props.modelValue,
    (newVal) => {
        if (state.uploadCounter > 0) return
        if (newVal === undefined || newVal === null) return init('')

        const newValArr = arrayFullURL(stringToArray(cloneDeep(newVal)))
        const oldValArr = arrayFullURL(getAllUrls('array') as string[])
        if (newValArr.sort().toString() != oldValArr.sort().toString()) {
            init(newVal)
        }
    }
)
</script>

<style scoped lang="scss">
.ag-upload-wrapper {
    position: relative;
    outline: none;
    &.is-drag-over:not(.is-disabled) :deep(.el-upload--picture-card),
    &.is-drag-over:not(.is-disabled) :deep(.el-upload-dragger) {
        border-color: var(--el-color-primary);
    }
}
.ag-upload-trigger {
    position: absolute;
    inset: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
    &.ag-upload-trigger-file {
        position: relative;
        inset: unset;
        gap: 6px;
        padding: 0 10px;
        min-height: 32px;
        .ag-upload-drag-mask {
            margin: 0 10px;
        }
    }
}
.ag-upload-drag-mask {
    position: absolute;
    inset: 0;
    z-index: 10;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    background: rgba(255, 255, 255, 0.92);
    border: 2px dashed var(--el-color-primary);
    border-radius: 6px;
    color: var(--el-color-primary);
    font-size: var(--el-font-size-small);
    pointer-events: none;
}
.ag-upload-screenshot {
    position: absolute;
    bottom: 0;
    width: var(--el-upload-picture-card-size);
    height: 30px;
    line-height: 30px;
    text-align: center;
    font-size: var(--el-font-size-extra-small);
    color: var(--el-text-color-regular);
    border: 1px dashed var(--el-border-color);
    border-bottom: 1px dashed transparent;
    border-radius: 6px;
    border-top-right-radius: 20px;
    border-top-left-radius: 20px;
    user-select: none;
    &:hover {
        color: var(--el-color-primary);
        border: 1px dashed var(--el-color-primary);
    }
}
.ag-upload :deep(.el-upload:hover .ag-upload-icon) {
    color: var(--el-color-primary) !important;
}
:deep(.ag-upload-preview) .el-dialog__body {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10px;
    height: auto;
}
.ag-upload-preview-scroll {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10px;
    overflow: auto;
    max-height: 70vh;
}
.ag-upload-preview-img {
    max-width: 100%;
    max-height: 100%;
}
:deep(.el-dialog__headerbtn) {
    top: 2px;
    width: 37px;
    height: 37px;
}
.ag-upload.image :deep(.el-upload--picture-card),
.ag-upload.images :deep(.el-upload--picture-card) {
    position: relative;
    display: inline-flex;
    align-items: center;
    justify-content: center;
}
.ag-upload.image :deep(.el-upload--picture-card > .el-tooltip__trigger),
.ag-upload.images :deep(.el-upload--picture-card > .el-tooltip__trigger) {
    position: absolute;
    inset: 0;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    height: 100%;
}
.ag-upload.file :deep(.el-upload),
.ag-upload.files :deep(.el-upload) {
    position: relative;
}
.ag-upload.file :deep(.el-upload-list),
.ag-upload.files :deep(.el-upload-list) {
    margin-left: -10px;
}
.ag-upload.files,
.ag-upload.images {
    :deep(.el-upload-list__item) {
        user-select: none;
        .el-upload-list__item-actions,
        .el-upload-list__item-name {
            cursor: move;
        }
    }
}
.ag-upload.is-limit-reached :deep(.el-upload--picture-card) {
    display: none;
}
.ag-upload.is-disabled :deep(.el-upload),
.ag-upload.is-disabled :deep(.el-upload) .el-button,
.ag-upload.is-disabled :deep(.el-upload--picture-card) {
    cursor: not-allowed;
}
</style>
