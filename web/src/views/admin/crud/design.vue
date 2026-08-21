<template>
    <div class="default-main">
        <div class="header-config-box">
            <el-row class="header-box">
                <div class="header">
                    <div class="header-item-box">
                        <el-form-item class="mr-20 table-name-item" :label="t('crud.index.tableName')" :error="state.error.tableName">
                            <el-input v-model="state.table.name" :placeholder="t('crud.index.tableNameTip')" @change="onTableNameChange"></el-input>
                        </el-form-item>

                        <el-form-item class="table-comment-item" :label="t('crud.index.tableComment')">
                            <el-input v-model="state.table.comment" :placeholder="t('crud.index.tableCommentPlaceholder')"></el-input>
                        </el-form-item>
                    </div>
                    <div class="header-right">
                        <el-link v-if="crudState.type != 'create'" @click="state.showDesignChangeLog = true" class="design-change-log" type="primary">
                            {{ t('crud.index.tableDesignChange') }}
                        </el-link>
                        <el-button type="primary" :loading="state.loading.generate" @click="onGenerate">
                            {{ t('crud.index.generateCrudCode') }}
                        </el-button>
                        <el-button @click="onAbandonDesign" type="danger">{{ t('crud.index.giveUp') }}</el-button>
                    </div>
                </div>
            </el-row>
            <transition :name="state.showHeaderSeniorConfig ? 'el-zoom-in-top' : 'el-zoom-in-bottom'">
                <div v-if="state.showHeaderSeniorConfig" class="header-senior-config-box">
                    <div class="header-senior-config-form">
                        <el-form-item :label-width="180" :label="t('crud.index.tableQuickSearchFields')">
                            <el-select :clearable="true" :multiple="true" class="w100" v-model="state.table.quickSearchField" placement="bottom">
                                <el-option
                                    v-for="(item, idx) in state.fields"
                                    :key="idx + item.uuid!"
                                    :label="item.name + (item.comment ? '-' + item.comment : item.title)"
                                    :value="item.uuid!"
                                />
                            </el-select>
                        </el-form-item>
                        <div class="default-sort-field-box">
                            <el-form-item :label-width="180" class="default-sort-field mr-20" :label="t('crud.index.tableDefaultSortFields')">
                                <el-select :clearable="true" v-model="state.table.defaultSortField" placement="bottom">
                                    <el-option
                                        v-for="(item, idx) in state.fields"
                                        :key="idx + item.uuid!"
                                        :label="item.name + (item.comment ? '-' + item.comment : item.title)"
                                        :value="item.uuid!"
                                    />
                                </el-select>
                            </el-form-item>
                            <el-form-item class="default-sort-field-type" :label="t('crud.index.sortMethod')">
                                <el-select v-model="state.table.defaultSortType">
                                    <el-option label="DESC" value="desc"></el-option>
                                    <el-option label="ASC" value="asc"></el-option>
                                </el-select>
                            </el-form-item>
                        </div>
                        <el-form-item :label-width="180" :label="t('crud.index.fieldsAsTableColumns')">
                            <el-select :clearable="true" :multiple="true" class="w100" v-model="state.table.columnFields" placement="bottom">
                                <el-option
                                    v-for="(item, idx) in state.fields"
                                    :key="idx + item.uuid!"
                                    :label="item.name + (item.comment ? '-' + item.comment : item.title)"
                                    :value="item.uuid!"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item :label-width="180" :label="t('crud.index.fieldsAsFormItems')">
                            <el-select :clearable="true" :multiple="true" class="w100" v-model="state.table.formFields" placement="bottom">
                                <el-option
                                    v-for="(item, idx) in state.fields"
                                    :key="idx + item.uuid!"
                                    :label="item.name + (item.comment ? '-' + item.comment : item.title)"
                                    :value="item.uuid!"
                                />
                            </el-select>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.theRelativePathToTheGeneratedCode')" :label-width="180">
                            <el-input v-model="state.table.generateRelativePath" @change="onGenerateRelativePath"></el-input>
                            <div class="block-help">{{ t('crud.index.codeLocationTip') }}</div>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.generatedHandlerLocation')" :label-width="180">
                            <el-input v-model="state.table.handlerFile"></el-input>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.generatedDataModelLocation')" :label-width="180">
                            <el-input v-model="state.table.modelFile"></el-input>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.generatedRepositoryLocation')" :label-width="180">
                            <el-input v-model="state.table.repositoryFile">
                                <template #append>
                                    <el-checkbox
                                        @change="onChangeCommonModel"
                                        v-model="state.table.isRepositoryModel"
                                        :label="t('crud.index.commonRepository')"
                                        size="small"
                                        :true-value="1"
                                        :false-value="0"
                                    />
                                </template>
                            </el-input>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.generatedServiceLocation')" :label-width="180">
                            <el-input v-model="state.table.serviceFile"></el-input>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.generatedRouterLocation')" :label-width="180">
                            <el-input v-model="state.table.routerFile"></el-input>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.webEndViewDirectory')" :label-width="180">
                            <el-input v-model="state.table.webViewsDir"></el-input>
                        </el-form-item>
                        <el-form-item :label="t('crud.index.routePath')" :label-width="180">
                            <el-input v-model="state.table.routePath"></el-input>
                        </el-form-item>
                    </div>
                </div>
            </transition>
            <div @click="state.showHeaderSeniorConfig = !state.showHeaderSeniorConfig" class="header-senior-config">
                <span>{{ t('crud.index.advancedConfiguration') }}</span>
                <Icon
                    class="senior-config-arrow-icon"
                    size="14"
                    color="var(--el-text-color-primary)"
                    :name="state.showHeaderSeniorConfig ? 'el-arrow-up' : 'el-arrow-down'"
                />
            </div>
        </div>
        <el-row v-loading="state.loading.init" class="fields-box" :gutter="20">
            <el-col :xs="24" :span="6">
                <el-collapse class="field-collapse" v-model="state.fieldCollapseName">
                    <el-collapse-item :title="t('crud.index.commonFields')" name="common">
                        <div class="field-box" :ref="tabsRefs.set">
                            <div v-for="(field, index) in fieldItem.common" :key="index" class="field-item">
                                <span>{{ field.title }}</span>
                            </div>
                        </div>
                    </el-collapse-item>
                    <el-collapse-item :title="t('crud.index.baseFields')" name="base">
                        <div class="field-box" :ref="tabsRefs.set">
                            <div v-for="(field, index) in fieldItem.base" :key="index" class="field-item">
                                <span>{{ field.title }}</span>
                            </div>
                        </div>
                    </el-collapse-item>
                    <el-collapse-item :title="t('crud.index.advancedFields')" name="senior">
                        <div class="field-box" :ref="tabsRefs.set">
                            <div v-for="(field, index) in fieldItem.senior" :key="index" class="field-item">
                                <span>{{ field.title }}</span>
                            </div>
                        </div>
                    </el-collapse-item>
                </el-collapse>
            </el-col>
            <el-col :xs="24" :span="12">
                <div ref="designWindowRef" class="design-window ag-scroll-style">
                    <div
                        v-for="(field, index) in state.fields"
                        :key="index"
                        :class="index === state.activateField ? 'activate' : ''"
                        @click="onActivateField(index)"
                        class="design-field-box"
                        :data-id="index"
                    >
                        <div class="design-field">
                            <span>{{ t('crud.index.fieldName') }}: </span>
                            <el-input
                                @pointerdown.stop
                                class="design-field-name-input"
                                :model-value="field.name"
                                type="string"
                                size="small"
                                @input="($event: string) => onFieldNameChange($event, index)"
                            ></el-input>
                        </div>
                        <div class="design-field">
                            <span>{{ t('crud.index.fieldComment') }}: </span>
                            <el-input
                                @pointerdown.stop
                                class="design-field-name-comment"
                                v-model="field.comment"
                                type="string"
                                size="small"
                                @change="onFieldCommentChange"
                            ></el-input>
                        </div>
                        <div class="design-field-right">
                            <el-button
                                v-if="['remoteSelect', 'remoteSelects'].includes(field.designType)"
                                @click.stop="onEditField(index, field)"
                                type="primary"
                                size="small"
                                circle
                            >
                                <Icon color="var(--el-color-white)" size="15" name="lucide-pencil" />
                            </el-button>
                            <el-button @click.stop="onDelField(index)" type="danger" size="small" circle>
                                <Icon color="var(--el-color-white)" size="15" name="lucide-trash'" />
                            </el-button>
                        </div>
                    </div>
                    <div class="design-field-empty" v-if="!state.fields.length && !state.draggingField">
                        {{ t('crud.index.dragToDesign') }}
                    </div>
                </div>
            </el-col>
            <el-col :xs="24" :span="6">
                <div class="field-config ag-scroll-style">
                    <div v-if="state.activateField === -1" class="design-field-empty">
                        {{ t('crud.index.selectFieldFirstTip') }}
                    </div>
                    <div v-else :key="'activate-field-' + state.activateField">
                        <el-form label-position="top">
                            <el-divider content-position="left">{{ t('crud.index.common') }}</el-divider>
                            <el-form-item :label="t('crud.index.generateType')">
                                <el-select
                                    @change="onFieldDesignTypeChange($event)"
                                    class="w100"
                                    :model-value="state.fields[state.activateField].designType"
                                    placement="bottom"
                                >
                                    <el-option v-for="(item, idx) in designTypes" :key="idx" :label="item.name" :value="idx" />
                                </el-select>
                            </el-form-item>
                            <el-form-item :label="t('crud.index.fieldCommentsCrudDictionary')">
                                <el-input
                                    type="textarea"
                                    :rows="2"
                                    :placeholder="t('crud.index.fieldCommentTip')"
                                    v-model="state.fields[state.activateField].comment"
                                    @change="onFieldCommentChange"
                                ></el-input>
                            </el-form-item>
                            <el-divider content-position="left">{{ t('crud.index.fieldProperties') }}</el-divider>
                            <el-form-item :label="t('crud.index.fieldName')">
                                <el-input
                                    :model-value="state.fields[state.activateField].name"
                                    @input="($event: string) => onFieldNameChange($event, state.activateField)"
                                ></el-input>
                            </el-form-item>
                            <template v-if="state.fields[state.activateField].dataType">
                                <el-form-item :label="t('crud.index.fieldType')">
                                    <el-input
                                        type="textarea"
                                        v-model="state.fields[state.activateField].dataType"
                                        @change="onFieldAttrChange"
                                    ></el-input>
                                </el-form-item>
                            </template>
                            <template v-else>
                                <el-form-item :label="t('crud.index.fieldType')">
                                    <el-input v-model="state.fields[state.activateField].type" @change="onFieldAttrChange"></el-input>
                                </el-form-item>
                                <div class="field-inline">
                                    <el-form-item :label="t('crud.index.length')">
                                        <el-input-number
                                            class="w100"
                                            controls-position="right"
                                            v-model="state.fields[state.activateField].length"
                                            @change="onFieldAttrChange"
                                        ></el-input-number>
                                    </el-form-item>
                                    <el-form-item :label="t('crud.index.decimalPoint')">
                                        <el-input-number
                                            class="w100"
                                            controls-position="right"
                                            v-model="state.fields[state.activateField].precision"
                                            @change="onFieldAttrChange"
                                        ></el-input-number>
                                    </el-form-item>
                                </div>
                            </template>
                            <el-form-item :label="t('crud.index.fieldDefaults')">
                                <el-select v-model="state.fields[state.activateField].defaultType">
                                    <el-option label="手动输入" value="INPUT" />
                                    <el-option :label="`EMPTY STRING - ${t('crud.index.emptyString')}`" value="EMPTY STRING" />
                                    <el-option label="NULL" value="NULL" />
                                    <el-option :label="t('crud.index.noDefaultValue')" value="NONE" />
                                </el-select>
                                <el-input
                                    v-if="state.fields[state.activateField].defaultType == 'INPUT'"
                                    :placeholder="t('crud.index.pleaseInputTheDefaultValue')"
                                    type="text"
                                    v-model="state.fields[state.activateField].default"
                                    @change="onFieldAttrChange"
                                    class="default-input"
                                />
                            </el-form-item>
                            <el-form-item :label="t('crud.index.autoIncrement')">
                                <el-select v-model="state.fields[state.activateField].generated" @change="onFieldAttrChange">
                                    <el-option :label="`GENERATED ALWAYS - ${t('crud.index.always')}`" value="GENERATED ALWAYS" />
                                    <el-option :label="`GENERATED BY DEFAULT - ${t('crud.index.generatedByDefault')}`" value="GENERATED BY DEFAULT" />
                                </el-select>
                            </el-form-item>
                            <div class="field-inline">
                                <el-form-item class="form-item-position-right" :label="t('crud.index.primaryKey')">
                                    <el-switch v-model="state.fields[state.activateField].primaryKey" @change="onFieldAttrChange" />
                                </el-form-item>
                                <el-form-item class="form-item-position-right" :label="t('crud.index.allowNull')">
                                    <el-switch v-model="state.fields[state.activateField].null" @change="onFieldAttrChange" />
                                </el-form-item>
                            </div>
                            <div class="field-inline">
                                <el-form-item class="form-item-position-right" :label="t('crud.index.unique')">
                                    <el-switch v-model="state.fields[state.activateField].unique" @change="onFieldAttrChange" />
                                </el-form-item>
                            </div>
                            <template v-if="!isEmpty(state.fields[state.activateField].table)">
                                <el-divider content-position="left">{{ t('crud.index.fieldTableProperties') }}</el-divider>
                                <template v-for="(item, idx) in state.fields[state.activateField].table" :key="idx">
                                    <el-form-item :label="$t('crud.index.' + idx)">
                                        <AgInput
                                            :type="item.type"
                                            v-model="state.fields[state.activateField].table[idx].value"
                                            :placeholder="state.fields[state.activateField].table[idx].placeholder ?? ''"
                                            :attr="{
                                                dict: state.fields[state.activateField].table[idx].options ?? {},
                                                ...(state.fields[state.activateField].table[idx].attr ?? {}),
                                            }"
                                        />
                                    </el-form-item>
                                </template>
                            </template>
                            <template v-if="!isEmpty(state.fields[state.activateField].form)">
                                <el-divider content-position="left">{{ t('crud.index.fieldFormProperties') }}</el-divider>
                                <template v-for="(item, idx) in state.fields[state.activateField].form" :key="idx">
                                    <el-form-item v-if="item.type != 'hidden'" :label="$t('crud.index.' + idx)">
                                        <AgInput
                                            :type="item.type"
                                            v-model="state.fields[state.activateField].form[idx].value"
                                            :placeholder="state.fields[state.activateField].form[idx].placeholder ?? ''"
                                            :attr="{
                                                dict: state.fields[state.activateField].form[idx].options ?? {},
                                                ...(state.fields[state.activateField].form[idx].attr ?? {}),
                                            }"
                                        />
                                    </el-form-item>
                                </template>
                            </template>
                        </el-form>
                    </div>
                </div>
            </el-col>
        </el-row>
        <el-dialog
            @close="onCancelRemoteSelect"
            class="ag-operate-dialog"
            :model-value="state.remoteSelectPre.show"
            :title="t('crud.index.remoteSelectInfo')"
            :close-on-click-modal="false"
            :destroy-on-close="true"
            @keyup.enter="onSaveRemoteSelect"
        >
            <el-scrollbar max-height="60vh">
                <div class="ag-operate-form" :style="'width: calc(100% - 80px)'">
                    <el-form
                        ref="formRef"
                        :model="state.remoteSelectPre.form"
                        :rules="remoteSelectPreFormRules"
                        v-loading="state.remoteSelectPre.loading"
                        label-position="right"
                        label-width="160px"
                        v-if="state.remoteSelectPre.index != -1 && state.fields[state.remoteSelectPre.index]"
                    >
                        <el-form-item :label="t('crud.index.associatedDataTable')" prop="table">
                            <RemoteSelect
                                v-model="state.remoteSelectPre.form.table"
                                pk="table"
                                field="comment"
                                :remote-params="{
                                    exclusions: [
                                        'schema_migrations',
                                        'areas',
                                        'tokens',
                                        'captchas',
                                        'admin_group_access',
                                        'configs',
                                        'admin_logs',
                                        'crud_logs',
                                    ],
                                }"
                                :remote-url="tableListUrl"
                                @change="onJoinTableChange"
                            />
                        </el-form-item>
                        <div v-loading="state.loading.remoteSelect">
                            <el-form-item prop="pk" :label="t('crud.index.dropDownValueField')">
                                <el-select
                                    :key="'select-value' + JSON.stringify(state.remoteSelectPre.fieldList)"
                                    class="w100"
                                    clearable
                                    :placeholder="t('crud.index.valueFieldPlaceholder')"
                                    v-model="state.remoteSelectPre.form.pk"
                                >
                                    <el-option v-for="(label, key) in state.remoteSelectPre.fieldList" :key="key" :label="label" :value="key" />
                                </el-select>
                            </el-form-item>
                            <el-form-item prop="label" :label="t('crud.index.dropDownLabelField')">
                                <el-select
                                    :key="'select-label' + JSON.stringify(state.remoteSelectPre.fieldList)"
                                    class="w100"
                                    clearable
                                    :placeholder="t('crud.index.labelFieldPlaceholder')"
                                    v-model="state.remoteSelectPre.form.label"
                                >
                                    <el-option v-for="(label, key) in state.remoteSelectPre.fieldList" :key="key" :label="label" :value="key" />
                                </el-select>
                            </el-form-item>
                            <el-form-item
                                v-if="state.fields[state.remoteSelectPre.index].designType == 'remoteSelect'"
                                prop="joinField"
                                :label="t('crud.index.fieldsDisplayedInTheTable')"
                            >
                                <el-select
                                    :key="'join-field' + JSON.stringify(state.remoteSelectPre.fieldList)"
                                    class="w100"
                                    multiple
                                    clearable
                                    :placeholder="t('crud.index.tableFieldsPlaceholder')"
                                    v-model="state.remoteSelectPre.form.joinField"
                                >
                                    <el-option v-for="(label, key) in state.remoteSelectPre.fieldList" :key="key" :label="label" :value="key" />
                                </el-select>
                            </el-form-item>

                            <el-form-item prop="modelName" :label="t('crud.index.remoteModelName')">
                                <RemoteSelect
                                    v-model="state.remoteSelectPre.form.modelName"
                                    pk="name"
                                    field="comment"
                                    :remote-url="modelListUrl"
                                    @row="onJoinModelChange"
                                    :remote-params="{
                                        exclusions: ['AdminLog', 'AdminGroupAccess', 'Area', 'Token', 'Captcha', 'Config', 'CrudLog'],
                                    }"
                                />
                                <div class="block-help">{{ t('crud.index.joinTableTip') }}</div>
                            </el-form-item>

                            <template v-if="state.remoteSelectPre.form.modelName">
                                <el-form-item prop="modelPackage" :label="t('crud.index.remoteModelPackage')">
                                    <el-input
                                        disabled
                                        v-model="state.remoteSelectPre.form.modelPackage"
                                        :placeholder="t('common.pleaseEnter', { field: t('crud.index.remoteModelPackage') })"
                                    ></el-input>
                                </el-form-item>

                                <el-form-item prop="modelFile" :label="t('crud.index.remoteModelFile')">
                                    <el-input
                                        disabled
                                        v-model="state.remoteSelectPre.form.modelFile"
                                        :placeholder="t('common.pleaseEnter', { field: t('crud.index.remoteModelFile') })"
                                    ></el-input>
                                </el-form-item>
                            </template>

                            <el-form-item :label="t('crud.index.dataSourceConfigurationType')">
                                <el-radio-group v-model="state.remoteSelectPre.form.sourceConfigType">
                                    <el-radio border value="crud">{{ t('crud.index.quickConfig') }}</el-radio>
                                    <el-radio border value="custom">{{ t('crud.index.customConfiguration') }}</el-radio>
                                </el-radio-group>
                            </el-form-item>

                            <el-form-item :label="t('crud.index.remoteCrudLog')" v-if="state.remoteSelectPre.form.sourceConfigType == 'crud'">
                                <RemoteSelect
                                    v-model="state.remoteSelectPre.form.crudLog"
                                    pk="id"
                                    field="label"
                                    :remote-url="crudLogListUrl"
                                    @row="onJoinCrudLogChange"
                                    :empty-values="[0, null]"
                                    :value-on-clear="0"
                                />
                            </el-form-item>

                            <el-form-item prop="remoteUrl" :label="t('crud.index.apiUrl')">
                                <el-input v-model="state.remoteSelectPre.form.remoteUrl" :placeholder="t('crud.index.apiUrlExample')"></el-input>
                            </el-form-item>
                        </div>
                    </el-form>
                </div>
            </el-scrollbar>
            <template #footer>
                <div :style="'width: calc(100% - 88px)'">
                    <el-button @click="onCancelRemoteSelect">{{ $t('common.cancel') }}</el-button>
                    <el-button @click="onSaveRemoteSelect" type="primary">
                        {{ $t('common.save') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
        <el-dialog
            @close="closeConfirmGenerate"
            class="ag-operate-dialog confirm-generate-dialog"
            :model-value="state.confirmGenerate.show"
            :title="t('crud.index.confirmCrudCodeGeneration')"
        >
            <div class="confirm-generate-dialog-body">
                <el-alert v-if="state.confirmGenerate.handler" :title="t('crud.index.handlerExistsTip')" center type="error" />
                <el-alert v-if="showTableConflictConfirmGenerate()" :title="t('crud.index.tableExistsTip')" class="mt-10" center type="error" />
                <el-alert v-if="state.confirmGenerate.menu" :title="t('crud.index.menuRuleExistsTip')" class="mt-10" center type="error" />
            </div>
            <template #footer>
                <div class="confirm-generate-dialog-footer">
                    <el-button @click="closeConfirmGenerate">{{ $t('common.cancel') }}</el-button>
                    <el-button :loading="state.loading.generate" @click="startGenerate" type="primary">
                        {{ t('crud.index.continueBuilding') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
        <el-dialog class="ag-operate-dialog design-change-log-dialog" width="20%" v-model="state.showDesignChangeLog">
            <template #header>
                <div>
                    {{ t('crud.index.dataTableDesignChangesPreview') }}
                </div>
            </template>
            <el-scrollbar max-height="400px">
                <template v-if="state.table.designChange.length">
                    <el-timeline class="design-change-log-timeline">
                        <el-timeline-item
                            v-for="(item, idx) in state.table.designChange"
                            :key="idx"
                            :type="getTableDesignTimelineType(item.type)"
                            :hollow="true"
                            :hide-timestamp="true"
                        >
                            <div class="design-timeline-box">
                                <el-checkbox v-model="item.sync" :label="getTableDesignChangeContent(item)" size="small" />
                            </div>
                        </el-timeline-item>
                    </el-timeline>
                    <span class="design-change-tips">{{ t('crud.index.designChangeTips') }}</span>
                </template>
                <div class="design-change-tips" v-else>暂无表设计变更</div>
                <el-form-item :label="t('crud.index.tableReBuild')" class="rebuild-form-item">
                    <el-radio-group v-model="state.table.rebuild">
                        <el-radio border value="No">{{ t('crud.index.no') }}</el-radio>
                        <el-radio border value="Yes">{{ t('crud.index.yes') }}</el-radio>
                    </el-radio-group>
                    <div class="block-help">{{ t('crud.index.tableReBuildBlockHelp') }}</div>
                </el-form-item>
            </el-scrollbar>
            <template #footer>
                <div class="confirm-generate-dialog-footer">
                    <el-button @click="state.showDesignChangeLog = false">
                        {{ t('common.confirm') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { useTemplateRefsList } from '@vueuse/core'
import type { FormItemRule, MessageHandler, TimelineItemProps } from 'element-plus'
import { ElLoading, ElMessage, ElMessageBox, ElNotification } from 'element-plus'
import { cloneDeep, isEmpty, range } from 'lodash-es'
import pluralize from 'pluralize-esm'
import type { SortableEvent } from 'sortablejs'
import Sortable from 'sortablejs'
import { nextTick, onMounted, reactive, useTemplateRef } from 'vue'
import { useI18n } from 'vue-i18n'
import {
    checkGenerate,
    crudLogListUrl,
    generate,
    getGenerateBasicData,
    getTableFieldList,
    logStart,
    modelListUrl,
    parseTableData,
    tableListUrl,
} from '/@/api/admin/crud'
import { ping } from '/@/api/common'
import RemoteSelect from '/@/components/agInput/components/remoteSelect.vue'
import AgInput from '/@/components/agInput/index.vue'
import { getArrayKey, parseStrAttr } from '/@/utils/common'
import { uuid } from '/@/utils/random'
import { buildValidatorRule, regularVarName } from '/@/utils/validate'
import { reloadServer } from '/@/utils/vite'
import type { FieldItem, TableDesignChange, TableDesignChangeType } from '/@/views/admin/crud/index'
import { changeStep, state as crudState, designTypes, fieldItem, getTableAttr, tableFieldsKey } from '/@/views/admin/crud/index'

let pingTimer: number
let nameRepeatCount = 1
const { t } = useI18n()
const formRef = useTemplateRef('formRef')
const tabsRefs = useTemplateRefsList<HTMLElement>()
const designWindowRef = useTemplateRef('designWindowRef')

const state: {
    loading: {
        init: boolean
        generate: boolean
        remoteSelect: boolean
    }
    table: {
        name: string
        comment: string
        quickSearchField: string[]
        defaultSortField: string
        formFields: string[]
        columnFields: string[]
        defaultSortType: string
        generateRelativePath: string
        isRepositoryModel: number
        modelFile: string
        handlerFile: string
        repositoryFile: string
        routerFile: string
        serviceFile: string
        webViewsDir: string
        routePath: string
        designChange: TableDesignChange[]
        rebuild: string
    }
    fields: FieldItem[]
    activateField: number
    fieldCollapseName: string[]
    remoteSelectPre: {
        show: boolean
        index: number
        fieldList: AnyObj
        loading: boolean
        hideDelField: boolean
        form: {
            table: string
            pk: string
            label: string
            joinField: string[]
            sourceConfigType: 'crud' | 'custom'
            remoteUrl: string
            modelFile: string
            modelName: string
            modelPackage: string
            crudLog: number
        }
    }
    showHeaderSeniorConfig: boolean
    confirmGenerate: {
        show: boolean
        menu: boolean
        table: boolean
        handler: boolean
    }
    draggingField: boolean
    showDesignChangeLog: boolean
    error: {
        tableName: string
        fieldName: MessageHandler | null
        fieldNameDuplication: MessageHandler | null
    }
} = reactive({
    loading: {
        init: false,
        generate: false,
        remoteSelect: false,
    },
    table: {
        name: '',
        comment: '',
        quickSearchField: [],
        defaultSortField: '',
        formFields: [],
        columnFields: [],
        defaultSortType: 'desc',
        generateRelativePath: '',
        isRepositoryModel: 0,
        modelFile: '',
        handlerFile: '',
        repositoryFile: '',
        routerFile: '',
        serviceFile: '',
        webViewsDir: '',
        routePath: '',
        designChange: [],
        rebuild: 'No',
    },
    fields: [],
    activateField: -1,
    fieldCollapseName: ['common', 'base', 'senior'],
    remoteSelectPre: {
        show: false,
        index: -1,
        fieldList: [],
        loading: false,
        hideDelField: false,
        form: {
            table: '',
            pk: '',
            label: '',
            joinField: [],
            sourceConfigType: 'crud',
            remoteUrl: '',
            modelFile: '',
            modelName: '',
            modelPackage: '',
            crudLog: 0,
        },
    },
    showHeaderSeniorConfig: false,
    confirmGenerate: {
        show: false,
        menu: false,
        table: false,
        handler: false,
    },
    draggingField: false,
    showDesignChangeLog: false,
    error: {
        tableName: '',
        fieldName: null,
        fieldNameDuplication: null,
    },
})

type TableKey = keyof typeof state.table

const onActivateField = (idx: number) => {
    state.activateField = idx
}

const onFieldDesignTypeChange = (designType: string) => {
    // 获取新的类型的数据
    let fieldDesignData: FieldItem | null = null
    for (const key in fieldItem) {
        const fieldItemIndex = getArrayKey(fieldItem[key as keyof typeof fieldItem], 'designType', designType)
        if (fieldItemIndex !== false) {
            fieldDesignData = cloneDeep(fieldItem[key as keyof typeof fieldItem][fieldItemIndex])
            break
        }
    }

    if (!fieldDesignData) return false

    // 主键重复检查
    if (!primaryKeyRepeatCheck(fieldDesignData, state.activateField)) {
        return false
    }

    // 选中字段数据
    const field = cloneDeep(state.fields[state.activateField])

    // 赋值新类型
    field.designType = designType

    // 保留字段的 table 和 form 数据，此处额外处理以便交付给 handleFieldAttr 函数
    for (const tKey in field.table) {
        field.table[tKey] = field.table[tKey].value
    }
    for (const tKey in field.form) {
        field.form[tKey] = field.form[tKey].value
    }
    state.fields[state.activateField] = handleFieldAttr(field)

    // 保留字段的 uuid
    state.fields[state.activateField].uuid = field.uuid

    // 询问是否切换至预设方案（除了字段名的属性全部重置）
    ElMessageBox.confirm(t('crud.index.resetGenerateTypeAttr'), t('common.reminder'), {
        confirmButtonText: t('common.confirm') + t('common.reset'),
        cancelButtonText: t('crud.index.designEfficiency'),
        type: 'warning',
        closeOnClickModal: false,
    })
        .then(() => {
            // 记录字段属性更新
            onFieldAttrChange()

            // 删除快速搜索和排序，根据新类型重新赋值
            clearFieldTableData(state.fields[state.activateField].uuid!)

            // 重置属性，除了 name
            const name = state.fields[state.activateField].name
            state.fields[state.activateField] = handleFieldAttr(fieldDesignData)
            state.fields[state.activateField].name = name

            if (fieldDesignData.primaryKey) {
                // 设置为默认排序字段、快速搜索字段
                state.table.quickSearchField.push(state.fields[state.activateField].uuid!)
                if (!state.table.defaultSortField) {
                    state.table.defaultSortField = state.fields[state.activateField].uuid!
                }
            }

            if (fieldDesignData.designType == 'weigh') {
                state.table.defaultSortField = state.fields[state.activateField].uuid!
            }

            // 远程下拉参数预填
            if (['remoteSelect', 'remoteSelects'].includes(fieldDesignData.designType)) {
                showRemoteSelectPre(state.activateField, true)
            }

            // 表单表格字段预定义
            if (!fieldDesignData.formBuildExclude) {
                state.table.formFields.push(state.fields[state.activateField].uuid!)
            }
            if (!fieldDesignData.tableBuildExclude) {
                state.table.columnFields.push(state.fields[state.activateField].uuid!)
            }
        })
        .catch(() => {})
}

/**
 * 字段名修改
 */
const onFieldNameChange = (val: string, index: number) => {
    const name = state.fields[index].name
    state.fields[index].name = val
    logTableDesignChange({
        type: 'change-field-name',
        index: index,
        name: name,
        newName: val,
    })
}

/**
 * 主键字段重复检测
 */
const primaryKeyRepeatCheck = (field: FieldItem, excludeIndex: number = -1) => {
    if (field.primaryKey === true) {
        const primaryKeyField = state.fields.find((item, index) => {
            if (excludeIndex > -1 && index == excludeIndex) {
                return false
            }
            return item.primaryKey
        })
        if (primaryKeyField) {
            ElNotification({
                type: 'error',
                message: t('crud.index.thereCanOnlyBeOnePrimaryKeyField'),
            })
            return false
        }
    }
    return true
}

/**
 * 全部字段的名称命名规则检测
 */
const fieldNameCheck = (showErrorType: 'ElNotification' | 'ElMessage') => {
    if (state.error.fieldName) {
        state.error.fieldName.close()
        state.error.fieldName = null
    }
    for (const key in state.fields) {
        if (!regularVarName(state.fields[key].name)) {
            let msg = t('crud.index.invalidFieldName', { field: state.fields[key].name })
            if (showErrorType == 'ElMessage') {
                state.error.fieldName = ElMessage({
                    message: msg,
                    type: 'error',
                    duration: 0,
                })
            } else {
                ElNotification({
                    type: 'error',
                    message: msg,
                })
            }
            return false
        }
    }
    return true
}

/**
 * 全部字段的名称重复检测
 */
const fieldNameDuplicationCheck = (showErrorType: 'ElNotification' | 'ElMessage') => {
    if (state.error.fieldNameDuplication) {
        state.error.fieldNameDuplication.close()
        state.error.fieldNameDuplication = null
    }
    for (const key in state.fields) {
        let count = 0
        for (const checkKey in state.fields) {
            if (state.fields[key].name == state.fields[checkKey].name) {
                count++
            }
            if (count > 1) {
                let msg = t('crud.index.fieldNameDuplication', { field: state.fields[key].name })
                if (showErrorType == 'ElMessage') {
                    state.error.fieldNameDuplication = ElMessage({
                        message: msg,
                        type: 'error',
                        duration: 0,
                    })
                } else {
                    ElNotification({
                        type: 'error',
                        message: msg,
                    })
                }
                return false
            }
        }
    }
    return true
}

const onFieldAttrChange = () => {
    logTableDesignChange({
        type: 'change-field-attr',
        index: state.activateField,
        name: state.fields[state.activateField].name,
        newName: '',
    })
}

/**
 * 从 state.table.* 清理某个字段的数据
 */
const clearFieldTableData = (uuid: string) => {
    if (uuid == state.table.defaultSortField) {
        state.table.defaultSortField = ''
    }

    for (const key in tableFieldsKey) {
        const delIdx = (state.table[tableFieldsKey[key] as TableKey] as string[]).findIndex((item) => {
            return item == uuid
        })
        if (delIdx != -1) {
            ;(state.table[tableFieldsKey[key] as TableKey] as string[]).splice(delIdx, 1)
        }
    }
}

const onDelField = (index: number) => {
    if (!state.fields[index]) return
    state.activateField = -1

    clearFieldTableData(state.fields[index].uuid!)

    logTableDesignChange({
        type: 'del-field',
        name: state.fields[index].name,
        newName: '',
    })

    // 删除权重字段时，重设默认排序字段
    if (state.fields[index].designType == 'weigh') {
        const pkField = state.fields.find((item) => {
            return ['pk', 'spk'].includes(item.designType)
        })
        if (pkField) {
            state.table.defaultSortField = pkField.uuid!
        }
    }

    state.fields.splice(index, 1)
}

const showRemoteSelectPre = (index: number, hideDelField = false) => {
    state.remoteSelectPre.show = true
    state.remoteSelectPre.loading = true
    state.remoteSelectPre.index = index
    state.remoteSelectPre.hideDelField = hideDelField

    // 编辑
    if (state.fields[index] && state.fields[index].form.remoteTable.value) {
        state.remoteSelectPre.form.pk = state.fields[index].form.remotePk.value
        state.remoteSelectPre.form.label = state.fields[index].form.remoteField.value
        state.remoteSelectPre.form.table = state.fields[index].form.remoteTable.value
        state.remoteSelectPre.form.modelFile = state.fields[index].form.remoteModelFile.value
        state.remoteSelectPre.form.modelName = state.fields[index].form.remoteModelName.value
        state.remoteSelectPre.form.modelPackage = state.fields[index].form.remoteModelPackage.value
        state.remoteSelectPre.form.crudLog = state.fields[index].form.remoteCrudLog.value
        state.remoteSelectPre.form.joinField = state.fields[index].form.relationFields.value.split(',')
        state.remoteSelectPre.form.remoteUrl = state.fields[index].form.remoteUrl.value
        state.remoteSelectPre.form.sourceConfigType = state.fields[index].form.remoteSourceConfigType.value
        getTableFieldList(state.fields[index].form.remoteTable.value).then((res) => {
            const fieldSelect: AnyObj = {}
            for (const key in res.data.data.list) {
                fieldSelect[res.data.data.list[key].name] = `${res.data.data.list[key].name} - ${res.data.data.list[key].comment}`
            }
            state.remoteSelectPre.fieldList = fieldSelect
        })
    }

    state.remoteSelectPre.loading = false
}

const onEditField = (index: number, field: FieldItem) => {
    if (['remoteSelect', 'remoteSelects'].includes(field.designType)) return showRemoteSelectPre(index)
}

const closeConfirmGenerate = () => {
    state.confirmGenerate.show = false
}

const startGenerate = () => {
    state.loading.generate = true

    // 简化设计字段数据，只保留 key => value
    const fields = cloneDeep(state.fields)
    for (const key in fields) {
        for (const tKey in fields[key].table) {
            fields[key].table[tKey] = fields[key].table[tKey].value
            if (tKey == 'comSearchInputAttr') {
                fields[key].table[`${tKey}Parsed`] = parseStrAttr(fields[key].table[tKey])
            }
        }
        for (const tKey in fields[key].form) {
            fields[key].form[tKey] = fields[key].form[tKey].value
        }
    }

    // 通过 uuid 获取字段 name
    const table = cloneDeep(state.table)
    if (table.defaultSortField) {
        const defaultSortFieldIndex = getArrayKey(state.fields, 'uuid', table.defaultSortField)
        if (defaultSortFieldIndex !== false) {
            table.defaultSortField = state.fields[defaultSortFieldIndex].name
        }
    }
    for (const key in tableFieldsKey) {
        const names: string[] = []
        const uuids = table[tableFieldsKey[key] as TableKey] as string[]
        for (const uKey in uuids) {
            const uuidFieldIndex = getArrayKey(state.fields, 'uuid', uuids[uKey])
            if (uuidFieldIndex !== false) {
                names.push(state.fields[uuidFieldIndex].name)
            }
        }

        ;(table[tableFieldsKey[key] as TableKey] as string[]) = names
    }

    // 直接去掉末尾的 `表` 字
    table.comment = table.comment.replace(/表$/, '')

    generate({
        type: crudState.type,
        table,
        fields,
    })
        .then((res) => {
            nextTick(() => {
                // 要求 Vite 服务端重启
                if (import.meta.hot) {
                    if (res.data.data.air) {
                        const loadingInstance = ElLoading.service({
                            text: t('crud.index.serverRestarting'),
                        })
                        pingTimer = window.setInterval(() => {
                            ping()
                                .then(() => {
                                    clearInterval(pingTimer)
                                    loadingInstance.close()
                                    reloadServer('crud')
                                })
                                .catch(() => {})
                        }, 3000)
                    } else {
                        reloadServer('crud')
                    }
                } else {
                    ElNotification({
                        type: 'error',
                        message: t('crud.index.viteHotWarning'),
                    })
                }
            })
        })
        .finally(() => {
            state.loading.generate = false
            closeConfirmGenerate()
        })
}

const onGenerate = () => {
    // 字段名称检查
    if (!fieldNameCheck('ElNotification')) return
    if (!fieldNameDuplicationCheck('ElNotification')) return

    let msg = ''

    // 主键检查
    const pkIndex = state.fields.findIndex((item) => {
        return item.primaryKey
    })
    if (pkIndex === -1) {
        msg = t('crud.index.pleaseDesignThePrimaryKeyField')
    }

    // 表名检查
    if (!state.table.name) msg = t('crud.index.pleaseEnterTheDataTableName')
    if (state.error.tableName) msg = t('crud.index.pleaseEnterTheCorrectTableName')

    if (msg) {
        ElNotification({
            type: 'error',
            message: msg,
        })
        return
    }

    state.loading.generate = true
    checkGenerate({
        table: state.table.name,
        handler: state.table.handlerFile,
        route: state.table.routePath,
    })
        .then(() => {
            startGenerate()
        })
        .catch((res) => {
            state.loading.generate = false
            if (res.data.code == -1) {
                state.confirmGenerate.menu = res.data.data.menu
                state.confirmGenerate.table = res.data.data.table
                state.confirmGenerate.handler = res.data.data.handler
                if (showTableConflictConfirmGenerate() || state.confirmGenerate.handler || state.confirmGenerate.menu) {
                    state.confirmGenerate.show = true
                } else {
                    startGenerate()
                }
            } else {
                ElNotification({
                    type: 'error',
                    message: res.msg,
                })
            }
        })
}

const showTableConflictConfirmGenerate = () => state.confirmGenerate.table && (crudState.type == 'create' || state.table.rebuild == 'Yes')

const onAbandonDesign = () => {
    if (!state.table.name && !state.table.comment && !state.fields.length) {
        return changeStep('start')
    }
    ElMessageBox.confirm(t('crud.index.giveUpConfirm'), t('common.reminder'), {
        confirmButtonText: t('crud.index.giveUp'),
        cancelButtonText: t('common.cancel'),
        type: 'warning',
    })
        .then(() => {
            changeStep('start')
        })
        .catch(() => {})
}

interface SortableEvt extends SortableEvent {
    originalEvent?: DragEvent
}

/**
 * 处理字段的属性
 */
const handleFieldAttr = (field: FieldItem) => {
    field = cloneDeep(field)
    const designTypeAttr = cloneDeep(designTypes[field.designType])
    for (const tKey in field.form) {
        if (designTypeAttr.form[tKey]) designTypeAttr.form[tKey].value = field.form[tKey]
        if (tKey == 'imageMultiple' && field.form[tKey]) {
            designTypeAttr.table['render'] = getTableAttr('render', 'images')
        }
    }
    for (const tKey in field.table) {
        if (designTypeAttr.table[tKey]) designTypeAttr.table[tKey].value = field.table[tKey]
    }
    field.form = designTypeAttr.form
    field.table = designTypeAttr.table
    field.uuid = uuid()
    return field
}

/**
 * 根据字段字典重新生成字段的数据类型
 */
const onFieldCommentChange = (comment: string) => {
    onFieldAttrChange()
    if (['enum', 'set'].includes(state.fields[state.activateField].type)) {
        if (!comment) {
            state.fields[state.activateField].dataType = `${state.fields[state.activateField].type}()`
            return
        }
        comment = comment.replaceAll('：', ':')
        comment = comment.replaceAll('，', ',')
        let comments = comment.split(':')
        if (comments[1]) {
            comments = comments[1].split(',')
            comments = comments
                .map((value) => {
                    if (!value) return ''
                    let temp = value.split('=')
                    if (temp[0] && temp[1]) {
                        return `'${temp[0]}'`
                    }
                    return ''
                })
                .filter((str: string) => str != '')

            // 字段数据类型
            state.fields[state.activateField].dataType = `${state.fields[state.activateField].type}(${comments.join(',')})`
        }
    }
}

const loadData = () => {
    tableDesignChangeInit()
    if (!['table', 'log'].includes(crudState.type)) return

    state.loading.init = true

    // 从历史记录开始
    if (crudState.type == 'log') {
        logStart(crudState.log.id!, crudState.log.type)
            .then((res) => {
                // 字段数据
                const fields = res.data.data.fields
                for (const key in fields) {
                    const field = handleFieldAttr(fields[key])

                    // 默认值和默认值类型分析
                    if (typeof field.defaultType == 'undefined') {
                        if (field.default && ['none', 'null', 'empty string'].includes(field.default)) {
                            field.defaultType = field.default.toUpperCase() as 'EMPTY STRING' | 'NULL' | 'NONE'
                            field.default = ''
                        } else {
                            field.defaultType = 'INPUT'
                        }
                    }

                    state.fields.push(field)
                }

                // 表数据
                if (res.data.data.table.defaultSortField) {
                    const defaultSortFieldNameIndex = getArrayKey(state.fields, 'name', res.data.data.table.defaultSortField)
                    if (defaultSortFieldNameIndex !== false) {
                        res.data.data.table.defaultSortField = state.fields[defaultSortFieldNameIndex].uuid!
                    }
                }
                for (const key in tableFieldsKey) {
                    const uuids: string[] = []
                    const names = res.data.data.table[tableFieldsKey[key] as TableKey] as string[]
                    for (const nKey in names) {
                        const nameFieldIndex = getArrayKey(state.fields, 'name', names[nKey])
                        if (nameFieldIndex !== false) {
                            uuids.push(state.fields[nameFieldIndex].uuid!)
                        }
                    }

                    ;(res.data.data.table[tableFieldsKey[key] as TableKey] as string[]) = uuids
                }

                state.table = res.data.data.table
                tableDesignChangeInit()
                if (res.data.data.empty) {
                    state.table.rebuild = 'Yes'
                }
                state.table.isRepositoryModel = parseInt(res.data.data.table.isRepositoryModel)
            })
            .finally(() => {
                state.loading.init = false
            })
        return
    }

    // 从数据表或sql开始
    parseTableData({
        table: crudState.table,
    })
        .then((res) => {
            let fields = []
            for (const key in res.data.data.columns) {
                const field = handleFieldAttr(res.data.data.columns[key])
                if (!['id', 'updated_at', 'created_at', 'deleted_at'].includes(field.name)) {
                    state.table.formFields.push(field.uuid!)
                }
                if (!['textarea', 'file', 'files', 'editor', 'password', 'array'].includes(field.designType)) {
                    state.table.columnFields.push(field.uuid!)
                }
                if (field.designType == 'pk') {
                    state.table.defaultSortField = field.uuid!
                    state.table.quickSearchField.push(field.uuid!)
                }
                if (field.designType == 'weigh') {
                    state.table.defaultSortField = field.uuid!
                }
                fields.push(field)
            }
            state.fields = fields
            state.table.comment = res.data.data.comment
            if (res.data.data.empty) {
                state.table.rebuild = 'Yes'
            }
            if (crudState.type == 'table' && crudState.table) {
                state.table.name = crudState.table
                onTableChange(crudState.table)
            }
        })
        .finally(() => {
            state.loading.init = false
        })
}

/**
 * 字段名称重复时自动重命名
 */
const autoRenameRepeatField = (fieldName: string) => {
    const nameRepeatKey = getArrayKey(state.fields, 'name', fieldName)
    if (nameRepeatKey !== false) {
        fieldName += nameRepeatCount
        nameRepeatCount++
        return autoRenameRepeatField(fieldName)
    } else {
        return fieldName
    }
}

onMounted(() => {
    loadData()
    const sortable = Sortable.create(designWindowRef.value!, {
        group: 'design-field',
        animation: 200,
        filter: '.design-field-empty',
        onAdd: (evt: SortableEvt) => {
            const name = evt.originalEvent?.dataTransfer?.getData('name')
            const field = fieldItem[name as keyof typeof fieldItem]
            if (field && field[evt.oldIndex!]) {
                const data = handleFieldAttr(field[evt.oldIndex!])

                // 主键重复检测
                if (data.primaryKey) {
                    if (primaryKeyRepeatCheck(data)) {
                        // 设置为默认排序字段、快速搜索字段
                        state.table.quickSearchField.push(data.uuid!)
                        if (!state.table.defaultSortField) {
                            state.table.defaultSortField = data.uuid!
                        }
                    } else {
                        return evt.item.remove()
                    }
                }

                // 出现权重字段则以其排序
                if (data.designType == 'weigh') {
                    state.table.defaultSortField = data.uuid!
                }

                // name 重复时，自动重命名
                data.name = autoRenameRepeatField(data.name)

                // 插入字段
                state.fields.splice(evt.newIndex!, 0, data)

                logTableDesignChange({
                    type: 'add-field',
                    index: evt.newIndex!,
                    newName: data.name,
                    name: '',
                })

                // 远程下拉参数预填
                if (['remoteSelect', 'remoteSelects'].includes(data.designType)) {
                    showRemoteSelectPre(evt.newIndex!, true)
                }

                // 表单表格字段预定义
                if (!data.formBuildExclude) {
                    state.table.formFields.push(data.uuid!)
                }
                if (!data.tableBuildExclude) {
                    state.table.columnFields.push(data.uuid!)
                }
            }
            evt.item.remove()
            nextTick(() => {
                sortable.sort(range(state.fields.length).map((value) => value.toString()))
            })
        },
        onEnd: (evt) => {
            const temp = state.fields[evt.oldIndex!]
            state.fields.splice(evt.oldIndex!, 1)
            state.fields.splice(evt.newIndex!, 0, temp)

            nextTick(() => {
                sortable.sort(range(state.fields.length).map((value) => value.toString()))
            })
        },
    })

    tabsRefs.value.forEach((item, index) => {
        Sortable.create(item, {
            sort: false,
            group: {
                name: 'design-field',
                pull: 'clone',
                put: false,
            },
            animation: 200,
            setData: (dataTransfer) => {
                dataTransfer.setData('name', Object.keys(fieldItem)[index])
            },
            onStart: () => {
                state.draggingField = true
            },
            onEnd: () => {
                state.draggingField = false
            },
        })
    })
})

/**
 * 修改表名
 * @param val 新表名
 */
const onTableNameChange = (val: string) => {
    if (!val) return (state.error.tableName = '')
    if (/^[a-z_][a-z0-9_]*$/.test(val)) {
        state.error.tableName = ''
        onTableChange(val)
    } else {
        state.error.tableName = t('crud.index.tableNameFormatTip')
    }
    tableDesignChangeInit()
}

const tableDesignChangeInit = () => {
    state.table.rebuild = 'No'
    state.table.designChange = []
}

const onGenerateRelativePath = (val: string) => {
    // 关闭单数化
    onTableChange(val, false)
}

/**
 * 预获取一个表的生成数据
 * @param val 新表名
 */
const onTableChange = (val: string, singular = true) => {
    if (!val) return

    // 根据表名处理相对生成路径
    state.table.generateRelativePath = val.replaceAll('\\', '/')
    if (singular) {
        // 单数化
        state.table.generateRelativePath = pluralize.singular(state.table.generateRelativePath)
    }

    getGenerateBasicData(state.table.generateRelativePath, 'admin').then((res) => {
        const files = res.data.data.files
        state.table.modelFile = files.model.file
        state.table.handlerFile = files.handler.file
        state.table.routerFile = files.router.file
        state.table.serviceFile = files.service.file
        state.table.repositoryFile = files.repository.file
        state.table.routePath = res.data.data.route
        state.table.webViewsDir = files.views.dir

        // 如果标记了公共仓储，额外请求公共仓储的默认文件路径
        if (state.table.isRepositoryModel) {
            getGenerateBasicData(val, 'common').then((res) => {
                state.table.repositoryFile = res.data.data.files.repository.file
            })
        }
    })
}

const onChangeCommonModel = () => {
    onTableChange(state.table.generateRelativePath)
}

const onJoinModelChange = (row: AnyObj) => {
    state.remoteSelectPre.form.modelFile = row.file
    state.remoteSelectPre.form.modelPackage = row.package
}

const onJoinCrudLogChange = (row: AnyObj) => {
    if (!isEmpty(row)) {
        const app = !isEmpty(row.router_basic_data) ? `/${row.router_basic_data.app}/` : ''
        state.remoteSelectPre.form.remoteUrl = app + row.table.routePath + '/list'
    }
}

const onJoinTableChange = () => {
    if (!state.remoteSelectPre.form.table) return

    // 重置远程下拉信息表单
    resetRemoteSelectForm(['table'])

    state.loading.remoteSelect = true
    getTableFieldList(state.remoteSelectPre.form.table)
        .then((res) => {
            state.remoteSelectPre.form.pk = res.data.data.pk

            const preLabel = ['name', 'title', 'username', 'nickname']
            for (const key in res.data.data.list) {
                if (preLabel.includes(res.data.data.list[key].name)) {
                    state.remoteSelectPre.form.label = res.data.data.list[key].name
                    state.remoteSelectPre.form.joinField.push(res.data.data.list[key].name)
                    break
                }
            }

            const fieldSelect: AnyObj = {}
            for (const key in res.data.data.list) {
                fieldSelect[res.data.data.list[key].name] = `${res.data.data.list[key].name} - ${res.data.data.list[key].comment}`
            }
            state.remoteSelectPre.fieldList = fieldSelect
        })
        .finally(() => {
            state.loading.remoteSelect = false
        })
}

const getRemoteSelectFieldName = (table: string, designType: string) => {
    const fieldName = pluralize.singular(table)
    return fieldName + (designType == 'remoteSelect' ? '_id' : '_ids')
}

const onSaveRemoteSelect = () => {
    const submitCallback = () => {
        // 修改字段名
        if (state.fields[state.remoteSelectPre.index].name == 'remote_select') {
            const newName = getRemoteSelectFieldName(state.remoteSelectPre.form.table, state.fields[state.remoteSelectPre.index].designType)
            onFieldNameChange(newName, state.remoteSelectPre.index)
        }

        state.fields[state.remoteSelectPre.index].form.remotePk.value = state.remoteSelectPre.form.pk
        state.fields[state.remoteSelectPre.index].form.remoteField.value = state.remoteSelectPre.form.label
        state.fields[state.remoteSelectPre.index].form.remoteTable.value = state.remoteSelectPre.form.table
        state.fields[state.remoteSelectPre.index].form.remoteModelFile.value = state.remoteSelectPre.form.modelFile
        state.fields[state.remoteSelectPre.index].form.remoteModelName.value = state.remoteSelectPre.form.modelName
        state.fields[state.remoteSelectPre.index].form.remoteModelPackage.value = state.remoteSelectPre.form.modelPackage
        state.fields[state.remoteSelectPre.index].form.remoteCrudLog.value = state.remoteSelectPre.form.crudLog
        state.fields[state.remoteSelectPre.index].form.remoteUrl.value = state.remoteSelectPre.form.remoteUrl
        state.fields[state.remoteSelectPre.index].form.remoteSourceConfigType.value = state.remoteSelectPre.form.sourceConfigType

        state.fields[state.remoteSelectPre.index].form.relationFields.value =
            state.fields[state.remoteSelectPre.index].designType == 'remoteSelect'
                ? state.remoteSelectPre.form.joinField.join(',')
                : state.remoteSelectPre.form.label

        state.remoteSelectPre.index = -1
        state.remoteSelectPre.show = false
        resetRemoteSelectForm()
    }

    if (formRef.value) {
        formRef.value.validate((valid) => {
            if (valid) {
                submitCallback()
            }
        })
    }
}

const onCancelRemoteSelect = () => {
    state.remoteSelectPre.show = false
    resetRemoteSelectForm()
    if (state.remoteSelectPre.index !== -1 && state.remoteSelectPre.hideDelField) {
        onDelField(state.remoteSelectPre.index)
    }
}

/**
 * 重置远程下拉预填表单
 */
const resetRemoteSelectForm = (excludes: string[] = []) => {
    for (const key in state.remoteSelectPre.form) {
        if (excludes.includes(key)) continue
        if (key == 'joinField') {
            state.remoteSelectPre.form[key] = []
        } else if (key == 'sourceConfigType') {
            state.remoteSelectPre.form[key] = 'crud'
        } else if (key == 'crudLog') {
            state.remoteSelectPre.form[key] = 0
        } else {
            ;(state.remoteSelectPre.form[key as keyof typeof state.remoteSelectPre.form] as string) = ''
        }
    }
}

const remoteSelectPreFormRules: Partial<Record<string, FormItemRule[]>> = {
    table: [buildValidatorRule({ name: 'required', title: t('crud.index.remoteTable') })],
    pk: [buildValidatorRule({ name: 'required', title: t('crud.index.dropDownValueField') })],
    label: [buildValidatorRule({ name: 'required', title: t('crud.index.dropDownLabelField') })],
    joinField: [buildValidatorRule({ name: 'required', title: t('crud.index.fieldsDisplayedInTheTable') })],
    modelFile: [buildValidatorRule({ name: 'required', title: t('crud.index.remoteModelFile') })],
    modelName: [buildValidatorRule({ name: 'required', title: t('crud.index.remoteModelName') })],
    modelPackage: [buildValidatorRule({ name: 'required', title: t('crud.index.remoteModelPackage') })],
    remoteUrl: [buildValidatorRule({ name: 'required', title: t('crud.index.remoteUrl') })],
}

const logTableDesignChange = (data: TableDesignChange) => {
    if (crudState.type == 'create') return
    let push = true
    if (data.type == 'change-field-name') {
        for (const key in state.table.designChange) {
            // 有属性修改记录的字段被改名-单独循环防止字段再次改名后造成找不到属性修改记录
            if (state.table.designChange[key].type == 'change-field-attr' && data.name == state.table.designChange[key].name) {
                state.table.designChange[key].name = data.newName
            }
        }
        for (const key in state.table.designChange) {
            // 新增字段改名
            if (state.table.designChange[key].type == 'add-field' && state.table.designChange[key].newName == data.name) {
                state.table.designChange[key].newName = data.newName
                push = false
                // 同一字段不会有两条新增记录
                break
            }
            // 字段再次改名
            if (state.table.designChange[key].type == 'change-field-name' && state.table.designChange[key].newName == data.name) {
                data.name = state.table.designChange[key].name
                state.table.designChange[key] = data

                // 取消字段改名
                if (state.table.designChange[key].newName == state.table.designChange[key].name) {
                    state.table.designChange.splice(key as any, 1)
                }

                push = false
                break
            }
        }
    } else if (data.type == 'del-field') {
        let add = false
        state.table.designChange = state.table.designChange.filter((item) => {
            // 新增的字段被删除
            add = item.type == 'add-field' && item.newName == data.name
            // 有属性修改记录的字段被删除
            const attr = item.type == 'change-field-attr' && item.name == data.name

            return !add && !attr
        })

        // 有改名记录的字段被删除（延后单独处理避免先改名再改属性的情况）
        state.table.designChange = state.table.designChange.filter((item) => {
            const name = item.type == 'change-field-name' && item.newName == data.name
            if (name) data.name = item.name
            return !name
        })

        // 添加的字段需要过滤掉记录同时不记录删除操作
        if (add) push = false

        for (const key in state.table.designChange) {
            // 同一字段名称多次删除（删除后添加再删除）
            if (state.table.designChange[key].type == 'del-field' && state.table.designChange[key].name == data.name) {
                push = false
                break
            }
        }
    } else if (data.type == 'change-field-attr') {
        // 先改名再改属性无需处理
        for (const key in state.table.designChange) {
            // 重复修改属性只记录一次
            if (state.table.designChange[key].type == 'change-field-attr' && state.table.designChange[key].name == data.name) {
                push = false
                break
            }
            // 新增的字段无需记录属性修改
            if (state.table.designChange[key].type == 'add-field' && state.table.designChange[key].newName == data.name) {
                push = false
                break
            }
        }
    }
    data.sync = true
    if (push) state.table.designChange.push(data)
}

const getTableDesignChangeContent = (data: TableDesignChange): string => {
    switch (data.type) {
        case 'add-field':
            return t('crud.index.addField') + ' ' + data.newName
        case 'change-field-attr':
            return t('crud.index.modifyFieldProperties') + ' ' + data.name
        case 'change-field-name':
            return t('crud.index.modifyFieldName') + ' ' + data.name + ' => ' + data.newName
        case 'del-field':
            return t('crud.index.deleteField') + ' ' + data.name
        default:
            return t('common.unknown')
    }
}

const getTableDesignTimelineType = (type: TableDesignChangeType): TimelineItemProps['type'] => {
    let timeline = ''
    switch (type) {
        case 'change-field-name':
            timeline = 'warning'
            break
        case 'del-field':
            timeline = 'danger'
            break
        case 'add-field':
            timeline = 'primary'
            break
        case 'change-field-attr':
            timeline = 'success'
            break
        default:
            timeline = 'success'
            break
    }
    return timeline as TimelineItemProps['type']
}
</script>

<style scoped lang="scss">
.form-item-position-right {
    display: flex !important;
    align-items: center;
    :deep(.el-form-item__label) {
        margin-right: 4px;
        margin-bottom: 0 !important;
    }
}
.default-main {
    margin-bottom: 0;
}
.mt-10 {
    margin-top: 10px;
}
.mr-20 {
    margin-right: 20px;
}
.field-collapse :deep(.el-collapse-item__header) {
    padding-left: 10px;
    user-select: none;
}
.field-box {
    padding: 10px;
}
.field-item {
    display: inline-block;
    padding: 3px 16px;
    border: 1px dashed var(--el-border-color);
    border-radius: var(--el-border-radius-base);
    margin: 6px;
    cursor: pointer;
    user-select: none;
    &:hover {
        border-color: var(--el-color-primary);
    }
}
.header-config-box {
    position: relative;
    .header-senior-config {
        display: flex;
        align-items: center;
        justify-content: center;
        position: absolute;
        height: 24px;
        bottom: -24px;
        padding: 4px 20px;
        padding-top: 0;
        left: calc(50% - 10px);
        font-size: var(--el-font-size-small);
        border-bottom-left-radius: 50px;
        border-bottom-right-radius: 50px;
        background-color: var(--ag-bg-color-overlay);
        color: var(--el-text-color-primary);
        cursor: pointer;
        user-select: none;
        .senior-config-arrow-icon {
            margin-left: 4px;
        }
    }
}
.header-senior-config-box {
    width: 100%;
    padding: 10px;
    background-color: var(--ag-bg-color-overlay);
}
.header-senior-config-form {
    width: 50%;
    :deep(.el-form-item__label) {
        justify-content: flex-start;
    }
}
.header-box {
    display: flex;
    align-items: center;
    height: v-bind("state.error.tableName ? '70px':'60px'");
    padding: 10px;
    background-color: var(--ag-bg-color-overlay);
    border-radius: var(--el-border-radius-base);
    transition: 0.1s;
    .header,
    .header-item-box {
        display: flex;
        width: 100%;
        align-items: center;
        justify-content: center;
        white-space: nowrap;
        :deep(.el-form-item) {
            margin-bottom: 0;
        }
    }
    .header-item-box {
        width: 50%;
    }
    .table-name-item {
        flex: 4;
    }
    .table-comment-item {
        flex: 4;
    }
    .header-right {
        margin-left: auto;
        .design-change-log {
            margin-right: 10px;
        }
    }
}
.default-sort-field-box {
    display: flex;
    .default-sort-field {
        flex: 6;
    }
    .default-sort-field-type {
        flex: 3;
    }
}
.fields-box {
    margin-top: 36px;
}
.design-field-empty {
    display: flex;
    height: 100%;
    color: var(--el-color-info);
    font-size: var(--el-font-size-medium);
    align-items: center;
    justify-content: center;
}
.design-window {
    overflow-x: auto;
    height: calc(100vh - 200px);
    border-radius: var(--el-border-radius-base);
    background-color: var(--ag-bg-color-overlay);
    border: v-bind('state.draggingField ? "1px dashed var(--el-color-primary)":(state.fields.length ? "none":"1px dashed var(--el-border-color)")');
    .design-field-box {
        display: flex;
        padding: 10px;
        align-items: center;
        border: 1px dashed var(--el-border-color);
        border-radius: var(--el-border-radius-base);
        margin-bottom: 2px;
        cursor: pointer;
        user-select: none;
        .design-field {
            padding-right: 10px;
        }
        .design-field-name-input {
            width: 200px;
        }
        .design-field-name-comment {
            width: 100px;
        }
        .design-field-right {
            margin-left: auto;
        }
        &:hover {
            border-color: var(--el-color-primary);
        }
    }
    .design-field-box.activate {
        border-color: var(--el-color-primary);
    }
}
.field-inline {
    display: flex;
    :deep(.el-form-item) {
        width: 46%;
        margin-right: 2%;
    }
}
.default-input {
    margin-top: 10px;
}
.field-config {
    overflow-x: auto;
    height: calc(100vh - 200px);
    padding: 20px;
    background-color: var(--ag-bg-color-overlay);
}
:deep(.confirm-generate-dialog) .el-dialog__body {
    height: unset;
}
.confirm-generate-dialog-body {
    padding: 30px;
}
.confirm-generate-dialog-footer {
    display: flex;
    align-items: center;
    justify-content: center;
}
:deep(.design-change-log-dialog) .el-dialog__body {
    height: unset;
    padding-top: 20px;
    .design-change-log-timeline {
        padding-left: 10px;
        .el-timeline-item .el-timeline-item__node {
            top: 3px;
        }
    }
    .design-change-tips {
        display: block;
        margin-bottom: 20px;
        color: var(--el-color-info);
        font-size: var(--el-font-size-small);
    }
    .rebuild-form-item {
        padding-top: 20px;
        border-top: 1px solid var(--el-border-color-lighter);
    }
}
</style>
