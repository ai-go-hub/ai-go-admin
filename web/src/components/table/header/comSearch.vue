<template>
    <div class="table-com-search-wrapper">
        <div class="table-com-search">
            <el-form
                @submit.prevent=""
                @keyup.enter="props.manager.handleEvent('com-search', { event: 'submit-com-search-form', data: comSearch.form })"
                label-position="top"
                :model="comSearch.form"
            >
                <el-row>
                    <template v-for="(item, idx) in props.manager.table.column" :key="idx">
                        <template v-if="item.operator !== false">
                            <!-- 自定义渲染 component、slot -->
                            <el-col
                                v-if="item.comSearchRender == 'customRender' || item.comSearchRender == 'slot'"
                                v-bind="{
                                    xs: item.comSearchColAttr?.xs ? item.comSearchColAttr?.xs : 24,
                                    sm: item.comSearchColAttr?.sm ? item.comSearchColAttr?.sm : 6,
                                    ...item.comSearchColAttr,
                                }"
                            >
                                <!-- 外部可以使用 :deep() 选择器修改css样式 -->
                                <div class="com-search-col" :class="item.prop">
                                    <div class="com-search-col-label" v-if="item.comSearchShowLabel !== false">{{ item.label }}</div>
                                    <div class="com-search-col-input">
                                        <!-- 自定义组件/函数渲染 -->
                                        <component
                                            v-if="item.comSearchRender == 'customRender'"
                                            :is="item.comSearchCustomRender"
                                            :renderColumnConfig="item"
                                            :renderColumnProp="item.prop!"
                                            :renderValue="comSearch.form[item.prop!]"
                                            :manager="props.manager"
                                        />

                                        <!-- 自定义渲染-slot -->
                                        <slot v-else-if="item.comSearchRender == 'slot'" :name="item.comSearchSlotName"></slot>
                                    </div>
                                </div>
                            </el-col>

                            <!-- 时间日期范围 -->
                            <el-col
                                v-else-if="
                                    (item.comSearchRender == 'datetime' || item.comSearchRender == 'date') &&
                                    (item.operator == 'BETWEEN' || item.operator == 'NOT BETWEEN')
                                "
                                :xs="24"
                                :sm="12"
                            >
                                <div class="com-search-col" :class="item.prop">
                                    <div class="com-search-col-label w16" v-if="item.comSearchShowLabel !== false">{{ item.label }}</div>
                                    <div class="com-search-col-input-range w83">
                                        <el-date-picker
                                            class="datetime-picker w100"
                                            v-model="comSearch.form[item.prop!]"
                                            :default-time="[new Date(2000, 1, 1, 0, 0, 0), new Date(2000, 1, 1, 23, 59, 59)]"
                                            :type="item.comSearchRender == 'date' ? 'daterange' : 'datetimerange'"
                                            :range-separator="$t('common.to')"
                                            :start-placeholder="getPlaceholder(item.comSearchPlaceholder, 0, $t('common.startDate'))"
                                            :end-placeholder="getPlaceholder(item.comSearchPlaceholder, 1, $t('common.endDate'))"
                                            :value-format="item.comSearchRender == 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss'"
                                            :teleported="false"
                                            v-bind="item.comSearchInputAttr"
                                        />
                                    </div>
                                </div>
                            </el-col>

                            <!-- 时间范围 -->
                            <el-col
                                v-else-if="item.comSearchRender == 'time' && (item.operator == 'BETWEEN' || item.operator == 'NOT BETWEEN')"
                                :xs="24"
                                :sm="12"
                            >
                                <div class="com-search-col" :class="item.prop">
                                    <div class="com-search-col-label w16" v-if="item.comSearchShowLabel !== false">{{ item.label }}</div>
                                    <div class="com-search-col-input-range w83">
                                        <el-time-picker
                                            class="time-picker w100"
                                            v-model="comSearch.form[item.prop!]"
                                            is-range
                                            :default-value="[new Date(2000, 1, 1, 0, 0, 0), new Date(2000, 1, 1, 23, 59, 59)]"
                                            :range-separator="$t('common.to')"
                                            :start-placeholder="getPlaceholder(item.comSearchPlaceholder, 0, $t('common.startTime'))"
                                            :end-placeholder="getPlaceholder(item.comSearchPlaceholder, 1, $t('common.endTime'))"
                                            value-format="HH:mm:ss"
                                            v-bind="item.comSearchInputAttr"
                                        />
                                    </div>
                                </div>
                            </el-col>

                            <!-- 其他 -->
                            <el-col v-else :xs="24" :sm="6">
                                <div class="com-search-col" :class="item.prop">
                                    <div class="com-search-col-label" v-if="item.comSearchShowLabel !== false">{{ item.label }}</div>
                                    <!-- 数字范围 -->
                                    <div v-if="item.operator == 'BETWEEN' || item.operator == 'NOT BETWEEN'" class="com-search-col-input-range">
                                        <el-input
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                            type="string"
                                            v-model="comSearch.form[item.prop! + '-start']"
                                            :clearable="true"
                                            v-bind="item.comSearchInputAttr"
                                        ></el-input>
                                        <div class="range-separator">{{ $t('common.to') }}</div>
                                        <el-input
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder, 1)"
                                            type="string"
                                            v-model="comSearch.form[item.prop! + '-end']"
                                            :clearable="true"
                                            v-bind="item.comSearchInputAttr"
                                        ></el-input>
                                    </div>
                                    <!-- 是否 [NOT] NULL -->
                                    <div v-else-if="item.operator == 'NULL' || item.operator == 'NOT NULL'" class="com-search-col-input">
                                        <el-checkbox
                                            v-model="comSearch.form[item.prop!]"
                                            :label="item.operator"
                                            size="large"
                                            v-bind="item.comSearchInputAttr"
                                        ></el-checkbox>
                                    </div>
                                    <div v-else-if="item.operator" class="com-search-col-input">
                                        <!-- 时间日期筛选 -->
                                        <el-date-picker
                                            class="datetime-picker w100"
                                            v-if="item.comSearchRender == 'date' || item.comSearchRender == 'datetime'"
                                            v-model="comSearch.form[item.prop!]"
                                            :type="item.comSearchRender == 'date' ? 'date' : 'datetime'"
                                            :value-format="item.comSearchRender == 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss'"
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                            :teleported="false"
                                            v-bind="item.comSearchInputAttr"
                                        />

                                        <!-- 时间筛选 -->
                                        <el-time-picker
                                            class="time-picker w100"
                                            v-if="item.comSearchRender == 'time'"
                                            v-model="comSearch.form[item.prop!]"
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                            value-format="HH:mm:ss"
                                            v-bind="item.comSearchInputAttr"
                                        />

                                        <!-- tag、tags、select -->
                                        <el-select
                                            class="w100"
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                            v-else-if="
                                                (item.render == 'tag' || item.render == 'tags' || item.comSearchRender == 'select') && item.dict
                                            "
                                            v-model="comSearch.form[item.prop!]"
                                            :multiple="item.operator == 'IN' || item.operator == 'NOT IN'"
                                            :clearable="true"
                                            v-bind="item.comSearchInputAttr"
                                        >
                                            <el-option v-for="(opt, okey) in item.dict" :key="item.prop! + okey" :label="opt" :value="okey" />
                                        </el-select>

                                        <!-- 远程 select -->
                                        <!-- <agInput
                                            v-else-if="item.comSearchRender == 'remoteSelect'"
                                            type="remoteSelect"
                                            v-model="comSearch.form[item.prop]"
                                            :attr="{ ...item.comSearchRemote, ...item.comSearchInputAttr }"
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                        /> -->

                                        <!-- 开关 -->
                                        <el-select
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                            v-else-if="item.render == 'switch'"
                                            v-model="comSearch.form[item.prop!]"
                                            :clearable="true"
                                            class="w100"
                                            v-bind="item.comSearchInputAttr"
                                        >
                                            <template v-if="!isEmpty(item.dict)">
                                                <el-option v-for="(opt, okey) in item.dict" :key="item.prop! + okey" :label="opt" :value="okey" />
                                            </template>
                                            <template v-else>
                                                <el-option :label="$t('common.open')" value="1" />
                                                <el-option :label="$t('common.close')" value="0" />
                                            </template>
                                        </el-select>

                                        <!-- 字符串 -->
                                        <el-input
                                            :placeholder="getPlaceholder(item.comSearchPlaceholder)"
                                            v-else
                                            type="string"
                                            v-model="comSearch.form[item.prop!]"
                                            :clearable="true"
                                            v-bind="item.comSearchInputAttr"
                                        ></el-input>
                                    </div>
                                </div>
                            </el-col>
                        </template>
                    </template>

                    <el-col :xs="24" :sm="6">
                        <div class="com-search-col pl-20">
                            <el-button
                                @click="props.manager.handleEvent('com-search', { event: 'submit-com-search-form', data: comSearch.form })"
                                type="primary"
                            >
                                {{ $t('common.search') }}
                            </el-button>
                            <el-button @click="onResetForm()">{{ $t('common.reset') }}</el-button>
                        </div>
                    </el-col>
                </el-row>
            </el-form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { isArray, isEmpty, isUndefined } from 'lodash-es'

interface Props {
    manager: TableManagerInstance
}
const props = defineProps<Props>()
const comSearch = defineModel<ComSearchInterface>('comSearch', { required: true })

const onResetForm = () => {
    /**
     * 封装好的 /utils/common.ts/onResetForm 工具在此处不能使用，因为未使用 el-form-item
     * 改用公共搜索重新初始化函数
     */
    props.manager.handleEvent('reinit-com-search', { type: 'reset-com-search-form' })

    // 通知 Table 管家重新发起公共搜索
    props.manager.handleEvent('com-search', { event: 'reset-com-search-form' })
}

const getPlaceholder = (placeholder: string | string[] | undefined, key = 0, defaultValue = '') => {
    if (isUndefined(placeholder)) {
        return defaultValue
    } else if (isArray(placeholder)) {
        return placeholder[key]
    } else {
        return placeholder
    }
}
</script>

<style scoped lang="scss">
.table-com-search {
    box-sizing: border-box;
    width: 100%;
    max-width: 100%;
    background-color: var(--ag-bg-color-overlay);
    border: 1px solid var(--ag-border-color);
    border-bottom: none;
    padding: 13px 15px;
    font-size: 14px;
    .com-search-col {
        display: flex;
        align-items: center;
        padding-top: 8px;
        color: var(--el-text-color-regular);
        font-size: 13px;
    }
    .com-search-col-label {
        width: 33.33%;
        padding: 0 15px;
        text-align: right;
        overflow: hidden;
        white-space: nowrap;
    }
    .com-search-col-input {
        padding: 0 15px;
        width: 66.66%;
    }
    .com-search-col-input-range {
        display: flex;
        align-items: center;
        padding: 0 15px;
        width: 66.66%;
        .range-separator {
            padding: 0 5px;
        }
    }
}
.pl-20 {
    padding-left: 20px;
}
.w16 {
    width: 16.5% !important;
}
.w83 {
    width: 83.5% !important;
}
.w100 {
    width: 100%;
}
</style>
