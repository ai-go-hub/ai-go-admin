package gen

import (
	"context"
	_ "embed"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

//go:embed tmpl/indexvue.tmpl
var indexvueTmplStr string

// indexvueTmpl 列表页 index.vue 内容模板
var indexvueTmpl = template.Must(template.New("indexvue").Parse(indexvueTmplStr))

//go:embed tmpl/dialogform.tmpl
var dialogFormTmplStr string

// dialogFormTmpl 表单页 dialogForm.vue 内容模板
var dialogFormTmpl = template.Must(template.New("dialogform").Parse(dialogFormTmplStr))

// jsonValue 格式化 JS 属性值，true/false 与数字不加引号
func jsonValue(v string) string {
	if v == "true" || v == "false" {
		return v
	}
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return v
	}
	return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
}

// jsObject 将 map 转为 JS 对象字面量字符串
func jsObject(m map[string]any) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, k+": "+jsValue(m[k]))
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

// jsValue 格式化 JS 值: bool/数字不加引号，字符串加引号
func jsValue(v any) string {
	switch val := v.(type) {
	case bool:
		if val {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case string:
		return jsonValue(val)
	default:
		return jsonValue(fmt.Sprint(val))
	}
}

// fieldLangKey 字段的语言翻译 key（小驼峰）
func fieldLangKey(langPrefix string, f dto.CRUDFields) string {
	return langPrefix + "." + util.SnakeToLowerCamel(f.Name)
}

// commonFieldKeys 字段名命中且注释匹配时，直接复用 common 语言包，不再生成翻译
var commonFieldKeys = map[string]string{
	"created_at": "创建时间",
	"updated_at": "更新时间",
	"mobile":     "手机号",
	"content":    "内容",
	"status":     "状态",
	"enable":     "启用",
	"disable":    "禁用",
	"end_time":   "结束时间",
	"start_time": "开始时间",
	"start_date": "开始日期",
	"end_date":   "结束日期",
	"weigh":      "权重",
}

// commonFieldKey 字段名命中 common 且注释匹配时返回对应 key
func commonFieldKey(f dto.CRUDFields) (string, bool) {
	comment, ok := commonFieldKeys[f.Name]
	if !ok {
		return "", false
	}
	title, _ := crud.ParseFieldComment(f.Comment)
	if title != comment {
		return "", false
	}
	return util.SnakeToLowerCamel(f.Name), true
}

// fieldLabelKey 返回字段 label 的翻译 key
// id 用字面量 ID，命中 common 用 common.xxx，其余 langKey.字段
func fieldLabelKey(langPrefix string, f dto.CRUDFields) string {
	if f.Name == "id" {
		return "ID"
	}
	if key, ok := commonFieldKey(f); ok {
		return "common." + key
	}
	return langPrefix + "." + util.SnakeToLowerCamel(f.Name)
}

// fieldLabelExpr 返回字段 label 的表达式
func fieldLabelExpr(langPrefix string, f dto.CRUDFields) string {
	key := fieldLabelKey(langPrefix, f)
	if key == "ID" {
		return "'ID'"
	}
	return "t('" + key + "')"
}

// ============================== 语言包 ==============================

// CreateLangFile 创建中英文语言包文件
func CreateLangFile(ctx context.Context, basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	info := basic["lang"]
	// 远程关联字段注释
	remoteComments := remoteFieldComments(ctx, fields)
	if err := crud.WriteGeneratedFile(info.CnFile, buildLangContent(fields, remoteComments, true)); err != nil {
		return err
	}
	if err := crud.WriteGeneratedFile(info.EnFile, buildLangContent(fields, remoteComments, false)); err != nil {
		return err
	}
	// 执行 prettier 格式化语言包
	if err := filesystem.FormatWithPrettier(info.CnFile); err != nil {
		return err
	}
	return filesystem.FormatWithPrettier(info.EnFile)
}

// yamlQuote 需要时给 YAML 标量/键加引号
func yamlQuote(s string) string {
	if s == "" || strings.ContainsAny(s, ": #{}[],&*!|>'\"%@`") || s != strings.TrimSpace(s) {
		return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return s
}

// buildLangContent 组装语言包内容
func buildLangContent(fields []dto.CRUDFields, remoteComments map[string]string, cn bool) string {
	var b strings.Builder
	for _, f := range fields {
		if f.Name == "" {
			continue
		}
		key := util.SnakeToLowerCamel(f.Name)
		// id 且注释为 ID 时使用字面量，不在语言包中重复定义
		if f.Name == "id" && f.Comment == "ID" {
			continue
		}
		title, dict := crud.ParseFieldComment(f.Comment)
		// 命中 common 常用词的字段不再生成翻译，将直接使用 common.key
		if _, useCommon := commonFieldKey(f); !useCommon {
			val := f.Name
			if cn {
				val = title
				if val == "" {
					val = key
				}
			}
			fmt.Fprintf(&b, "%s: %s\n", key, yamlQuote(val))
		}
		for _, item := range dict {
			v := item.Key
			if cn {
				v = item.Value
				if v == "" {
					v = item.Key
				}
			}
			// 字典 key 带空格作为叶子 key，避免与字段标题 key 嵌套冲突
			fmt.Fprintf(&b, "%s: %s\n", yamlQuote(key+" "+item.Key), yamlQuote(v))
		}
	}

	// 远程关联字段注释（按 key 排序保证输出确定性）
	remoteKeys := make([]string, 0, len(remoteComments))
	for k := range remoteComments {
		remoteKeys = append(remoteKeys, k)
	}
	sort.Strings(remoteKeys)
	for _, k := range remoteKeys {
		v := k
		if cn {
			v = remoteComments[k]
			if v == "" {
				v = k
			}
		}
		fmt.Fprintf(&b, "%s: %s\n", k, yamlQuote(v))
	}
	return b.String()
}

// remoteFieldComments 收集远程关联字段的语言包条目（key -> 注释）
// remoteSelect 取 relationFields，remoteSelects 取 Field；查询失败或注释为空时跳过
func remoteFieldComments(ctx context.Context, fields []dto.CRUDFields) map[string]string {
	result := map[string]string{}
	for _, f := range fields {
		if f.DesignType != "remoteSelect" && f.DesignType != "remoteSelects" {
			continue
		}
		info := parseRemoteInfo(f)
		if info.ModelName == "" || info.Table == "" {
			continue
		}
		comments, err := crud.TableFieldComments(ctx, info.Table)
		if err != nil || len(comments) == 0 {
			continue
		}
		modelLower := util.SnakeToLowerCamel(util.PascalToSnake(info.ModelName))
		if f.DesignType == "remoteSelect" {
			for _, rf := range info.RelationFields {
				if rf == "" {
					continue
				}
				key := modelLower + snakeToPascal(rf)
				if c := comments[rf]; c != "" {
					result[key] = c
				}
			}
		} else if info.Field != "" {
			listField := RemoteSelectsListFieldName(f)
			key := util.SnakeToLowerCamel(util.PascalToSnake(listField))
			if c := comments[info.Field]; c != "" {
				result[key] = c
			}
		}
	}
	return result
}

// ============================== 表格（index.vue） ==============================

// columnPreset designType 的默认表格列属性
type columnPreset struct {
	operator            string
	render              string
	width               int
	sortable            string
	quickSearch         bool
	showOverflowTooltip bool
	comSearchRender     string
	format              string
}

var columnPresets = map[string]columnPreset{
	"pk":            {operator: "BETWEEN", width: 70},
	"string":        {operator: "ILIKE", quickSearch: true, showOverflowTooltip: true},
	"password":      {operator: "ILIKE"},
	"textarea":      {operator: "ILIKE"},
	"editor":        {operator: "ILIKE"},
	"radio":         {operator: "eq", render: "tag", comSearchRender: "select"},
	"select":        {operator: "eq", render: "tag", comSearchRender: "select"},
	"selects":       {operator: "ILIKE", render: "tags", comSearchRender: "select"},
	"color":         {operator: "", render: "color"},
	"iconSelect":    {operator: "ILIKE", render: "icon", width: 100},
	"areaSelect":    {operator: "ILIKE", width: 160},
	"file":          {operator: "false"},
	"files":         {operator: "false"},
	"image":         {operator: "false", render: "image", width: 80},
	"images":        {operator: "false", render: "images"},
	"int":           {operator: "BETWEEN", sortable: "custom"},
	"float":         {operator: "BETWEEN", sortable: "custom"},
	"weigh":         {operator: "", sortable: "custom", width: 80},
	"switch":        {operator: "false", render: "switch", sortable: "custom", width: 80},
	"remoteSelect":  {operator: "eq", width: 100},
	"remoteSelects": {operator: "false", render: "tags"},
	"array":         {operator: "false"},
	"checkbox":      {operator: "ILIKE", render: "tags"},
	"datetime":      {operator: "BETWEEN", render: "datetime", comSearchRender: "datetime", sortable: "custom", width: 160},
	"date":          {operator: "BETWEEN", render: "datetime", comSearchRender: "date", sortable: "custom", width: 120, format: "YYYY-MM-DD"},
	"year":          {operator: "BETWEEN", sortable: "custom"},
	"time":          {operator: "ILIKE", comSearchRender: "time", width: 100},
}

// columnPropsOrder 列属性输出顺序
var columnPropsOrder = []string{"label", "prop", "type", "align", "operator", "render", "width", "sortable", "quickSearch", "showOverflowTooltip", "comSearchRender", "custom", "dict", "buttons"}

// buildColumn 组装单个字段的表格列对象
func buildColumn(langPrefix string, f dto.CRUDFields) string {
	m := map[string]string{
		"label": fieldLabelExpr(langPrefix, f),
		"prop":  "'" + f.Name + "'",
		"align": "'center'",
	}

	preset := columnPresets[f.DesignType]
	if preset.operator != "" {
		if preset.operator == "false" {
			m["operator"] = "false"
		} else {
			m["operator"] = "'" + preset.operator + "'"
		}
	}
	if f.Name == "created_at" || f.Name == "updated_at" {
		m["sortable"] = "'custom'"
	} else if preset.sortable != "" {
		m["sortable"] = "'" + preset.sortable + "'"
	}
	if preset.render != "" {
		m["render"] = "'" + preset.render + "'"
	}
	if preset.width > 0 {
		m["width"] = strconv.Itoa(preset.width)
	}
	if preset.quickSearch {
		m["quickSearch"] = "true"
	}
	if preset.showOverflowTooltip {
		m["showOverflowTooltip"] = "true"
	}
	if preset.comSearchRender != "" {
		m["comSearchRender"] = "'" + preset.comSearchRender + "'"
	}
	if preset.format != "" {
		m["custom"] = "{ format: " + jsonValue(preset.format) + " }"
	}

	if _, dict := crud.ParseFieldComment(f.Comment); len(dict) > 0 {
		dictParts := make([]string, 0, len(dict))
		for _, item := range dict {
			dictParts = append(dictParts, item.Key+": t('"+fieldLangKey(langPrefix, f)+" "+item.Key+"')")
		}
		m["dict"] = "{ " + strings.Join(dictParts, ", ") + " }"

		if f.DesignType == "switch" {
			customParts := make([]string, 0, len(dict))
			for _, item := range dict {
				color := "success"
				if item.Key == "0" || item.Key == "disable" {
					color = "danger"
				}
				customParts = append(customParts, item.Key+": '"+color+"'")
			}
			m["custom"] = "{ " + strings.Join(customParts, ", ") + " }"
		}
	}

	// 合并字段自定义表格属性（覆盖同名预设）
	for k, v := range f.Table {
		switch k {
		case "comSearchInputAttrParsed":
			// 空对象忽略，非空转为 JS 对象字面量
			if obj, ok := v.(map[string]any); ok && len(obj) > 0 {
				m["comSearchInputAttr"] = jsObject(obj)
			}
			continue
		case "comSearchInputAttr":
			continue
		}
		m[fmt.Sprint(k)] = jsonValue(fmt.Sprint(v))
	}

	// render 为 none 时去除
	if v, ok := m["render"]; ok && v == "'none'" {
		delete(m, "render")
	}

	// 按规范顺序输出，其余按字母序追加
	var parts []string
	for _, k := range columnPropsOrder {
		if v, ok := m[k]; ok {
			parts = append(parts, k+": "+v)
		}
	}
	var rest []string
	for k := range m {
		if !slices.Contains(columnPropsOrder, k) {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	for _, k := range rest {
		parts = append(parts, k+": "+m[k])
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}

// buildIndexColumns 组装 index.vue 的表格列定义
func buildIndexColumns(langPrefix string, table dto.CRUDTable, fields []dto.CRUDFields, weigh bool) string {
	var b strings.Builder
	b.WriteString("{ type: 'selection', align: 'center', operator: false },\n")
	for _, f := range fields {
		// 主表字段（仅 ColumnFields 中列出的字段生成列）
		if f.Name != "" && slices.Contains(table.ColumnFields, f.Name) {
			fmt.Fprintf(&b, "%s,\n", buildColumn(langPrefix, f))
		}

		// 关联表字段（远程下拉的关联展示列，不受 ColumnFields 限制）
		if relationColumns := buildRelationColumns(langPrefix, f); relationColumns != "" {
			fmt.Fprintf(&b, "%s,\n", relationColumns)
		}
	}
	optWidth := "100"
	if weigh {
		optWidth = "140"
	}
	fmt.Fprintf(&b, "{ label: t('common.operate'), align: 'center', width: %s, render: 'buttons', buttons: optButtons, operator: false },\n", optWidth)
	return b.String()
}

// buildRelationColumns 构建远程下拉字段的关联展示列
func buildRelationColumns(langPrefix string, f dto.CRUDFields) string {
	info := parseRemoteInfo(f)
	if info.ModelName == "" {
		return ""
	}
	modelSnake := util.PascalToSnake(info.ModelName)
	modelLower := util.SnakeToLowerCamel(modelSnake)
	var cols []string
	switch f.DesignType {
	case "remoteSelect":
		for _, rf := range info.RelationFields {
			if rf == "" {
				continue
			}
			key := modelLower + snakeToPascal(rf)
			cols = append(cols, fmt.Sprintf("{ label: t('%s.%s'), prop: '%s.%s', align: 'center', operator: 'ILIKE', render: 'tags' }", langPrefix, key, modelSnake, rf))
		}
	case "remoteSelects":
		if info.Field != "" {
			listField := RemoteSelectsListFieldName(f)
			prop := util.PascalToSnake(listField)
			key := util.SnakeToLowerCamel(util.PascalToSnake(listField))
			cols = append(cols, fmt.Sprintf("{ label: t('%s.%s'), prop: '%s', align: 'center', operator: false, render: 'tags' }", langPrefix, key, prop))
		}
	}
	return strings.Join(cols, ",\n")
}

// buildFormDefaultItems 组装 form.defaultItems 内容
func buildFormDefaultItems(fields []dto.CRUDFields) string {
	var b strings.Builder
	skipZeroValueTypes := map[string]bool{"int": true, "float": true, "remoteSelect": true}
	for _, f := range fields {
		if f.Name == "" {
			continue
		}
		if skipZeroValueTypes[f.DesignType] && f.Default == "0" {
			continue
		}
		if f.DefaultType == "INPUT" && f.Default != "" {
			fmt.Fprintf(&b, "%s: %s,\n", f.Name, buildFormDefaultValue(f))
		}
	}
	return b.String()
}

// buildFormDefaultArray 将逗号分割的多值默认值转为 JS 数组字面量
func buildFormDefaultArray(f dto.CRUDFields) string {
	parts := strings.Split(f.Default, ",")
	items := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if f.DesignType == "remoteSelects" {
			items = append(items, jsonValue(p))
		} else {
			items = append(items, "'"+strings.ReplaceAll(p, "'", "\\'")+"'")
		}
	}
	return "[" + strings.Join(items, ", ") + "]"
}

// buildDblClickNotEditColumn 构建 dblClickNotEditColumn 属性
func buildDblClickNotEditColumn(fields []dto.CRUDFields) (string, bool) {
	items := make([]string, 0, len(fields))
	for _, f := range fields {
		if f.DesignType == "switch" {
			items = append(items, "'"+strings.ReplaceAll(f.Name, "'", "\\'")+"'")
		}
	}
	if len(items) == 0 {
		return "[]", false
	}
	return "[" + strings.Join(items, ", ") + "]", true
}

// buildFormDefaultValue 按字段类型格式化默认值: 数字/开关不加引号，checkbox/selects/remoteSelects 为数组，其余字符串加单引号
func buildFormDefaultValue(f dto.CRUDFields) string {
	switch f.DesignType {
	case "int", "number", "float", "weigh", "remoteSelect", "switch":
		return jsonValue(f.Default)
	case "checkbox", "selects", "remoteSelects":
		return buildFormDefaultArray(f)
	default:
		return "'" + strings.ReplaceAll(f.Default, "'", "\\'") + "'"
	}
}

// CreateIndexVueFile 创建路由入口组件（表格）index.vue
func CreateIndexVueFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	info := basic["views"]
	langPrefix := basic["lang"].LangKey
	weigh := crud.HasWeighDesign(fields)

	optButtons := "edit', 'delete"
	if weigh {
		optButtons = "sort', 'edit', 'delete"
	}
	columns := buildIndexColumns(langPrefix, table, fields, weigh)

	defSort := table.DefaultSortField
	if defSort == "" {
		defSort = "id"
	}
	defOrder := table.DefaultSortType
	if defOrder == "" {
		defOrder = "desc"
	}

	apiURL := "/" + basic["router"].App + "/" + strings.Trim(table.RoutePath, "/") + "/"

	defaultItems := buildFormDefaultItems(fields)
	dblClickNotEditColumn, hasDblClickNotEditColumn := buildDblClickNotEditColumn(fields)
	content, err := crud.RenderTmpl(indexvueTmpl, map[string]any{
		"ApiURL":                   apiURL,
		"NeedOnMounted":            weigh,
		"OptButtons":               optButtons,
		"Columns":                  columns,
		"DefaultSort":              defSort,
		"DefaultOrder":             defOrder,
		"HasFormDefaultItems":      defaultItems != "",
		"FormDefaultItems":         defaultItems,
		"HasDblClickNotEditColumn": hasDblClickNotEditColumn,
		"DblClickNotEditColumn":    dblClickNotEditColumn,
	})
	if err != nil {
		return err
	}
	return crud.WriteGeneratedFile(info.Dir+"/index.vue", content)
}

// ============================== 表单（dialogForm.vue） ==============================

// CreateDialogFormFile 创建表单页 dialogForm.vue
func CreateDialogFormFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	info := basic["views"]
	langPrefix := basic["lang"].LangKey

	var fieldsBlock strings.Builder
	importSet := map[string]bool{}
	rules := make([]string, 0)
	for _, f := range fields {
		if f.Name == "" || !slices.Contains(table.FormFields, f.Name) {
			continue
		}
		fieldsBlock.WriteString(buildFormField(langPrefix, f, importSet))
		// 验证规则取自 Form.validator（规则名数组）
		validators := crud.FormValidators(f)
		if len(validators) > 0 {
			msg := formStr(f.Form, "validatorMsg")
			parts := make([]string, 0, len(validators))
			for _, v := range validators {
				if msg != "" {
					// 用户自定义错误消息: 以 message 属性代替 title
					parts = append(parts, "buildValidatorRule({ name: '"+v+"', message: "+jsonValue(msg)+" })")
				} else {
					parts = append(parts, "buildValidatorRule({ name: '"+v+"', title: "+fieldLabelExpr(langPrefix, f)+" })")
				}
			}
			rules = append(rules, util.SnakeToLowerCamel(f.Name)+": ["+strings.Join(parts, ", ")+"],")
		}
	}

	rulesBlock := "const rules: Partial<Record<string, FormItemRule[]>> = {}\n"
	if len(rules) > 0 {
		rulesBlock = "const rules: Partial<Record<string, FormItemRule[]>> = {\n" + strings.Join(rules, "\n") + "\n}\n"
	}

	content, err := crud.RenderTmpl(dialogFormTmpl, map[string]any{
		"FormFields":       fieldsBlock.String(),
		"Rules":            rulesBlock,
		"HasAgUpload":      importSet["AgUpload"],
		"HasRemoteSelect":  importSet["RemoteSelect"],
		"HasAreaSelect":    importSet["AreaSelect"],
		"HasArrayInput":    importSet["ArrayInput"],
		"HasIconSelect":    importSet["IconSelect"],
		"HasEditor":        importSet["Editor"],
		"HasValidatorRule": len(rules) > 0,
	})
	if err != nil {
		return err
	}
	return crud.WriteGeneratedFile(info.Dir+"/dialogForm.vue", content)
}

// buildFormField 组装单个表单项
func buildFormField(langPrefix string, f dto.CRUDFields, importSet map[string]bool) string {
	label := fieldLabelExpr(langPrefix, f)
	placeholderEnter := "t('common.pleaseEnter', { field: " + label + " })"
	placeholderSelect := "t('common.pleaseSelect', { field: " + label + " })"
	_, dict := crud.ParseFieldComment(f.Comment)

	var b strings.Builder
	fmt.Fprintf(&b, "<el-form-item :label=\"%s\" prop=\"%s\">\n", label, f.Name)

	switch f.DesignType {
	case "textarea":
		fmt.Fprintf(&b, "<el-input @keyup.enter.stop=\"\" @keyup.ctrl.enter=\"manager.submitForm(formRef)\" v-model=\"formItems.%s\" type=\"textarea\" :rows=\"3\" :placeholder=\"%s\" />\n", f.Name, placeholderEnter)
	case "radio":
		fmt.Fprintf(&b, "<el-radio-group v-model=\"formItems.%s\">\n", f.Name)
		for _, item := range dict {
			fmt.Fprintf(&b, "<el-radio value=\"%s\">{{ t('%s %s') }}</el-radio>\n", item.Key, fieldLangKey(langPrefix, f), item.Key)
		}
		b.WriteString("</el-radio-group>\n")
	case "checkbox":
		fmt.Fprintf(&b, "<el-checkbox-group v-model=\"formItems.%s\">\n", f.Name)
		for _, item := range dict {
			fmt.Fprintf(&b, "<el-checkbox value=\"%s\" border>{{ t('%s %s') }}</el-checkbox>\n", item.Key, fieldLangKey(langPrefix, f), item.Key)
		}
		b.WriteString("</el-checkbox-group>\n")
	case "select", "selects":
		multiple := ""
		if f.DesignType == "selects" {
			multiple = " multiple"
		}
		fmt.Fprintf(&b, "<el-select v-model=\"formItems.%s\" :placeholder=\"%s\"%s clearable>\n", f.Name, placeholderSelect, multiple)
		for _, item := range dict {
			fmt.Fprintf(&b, "<el-option :label=\"t('%s %s')\" value=\"%s\" />\n", fieldLangKey(langPrefix, f), item.Key, item.Key)
		}
		b.WriteString("</el-select>\n")
	case "color":
		fmt.Fprintf(&b, "<el-color-picker v-model=\"formItems.%s\" />\n", f.Name)
	case "iconSelect":
		importSet["IconSelect"] = true
		fmt.Fprintf(&b, "<IconSelect v-model=\"formItems.%s\" />\n", f.Name)
	case "areaSelect":
		importSet["AreaSelect"] = true
		fmt.Fprintf(&b, "<AreaSelect class=\"w100\" v-model=\"formItems.%s\" />\n", f.Name)
	case "image", "images", "file", "files":
		importSet["AgUpload"] = true
		multiple := ""
		if f.DesignType == "images" || f.DesignType == "files" {
			multiple = " :multiple=\"true\""
		}
		fmt.Fprintf(&b, "<AgUpload type=\"%s\" v-model=\"formItems.%s\"%s />\n", f.DesignType, f.Name, multiple)
	case "number", "int", "float", "weigh":
		fmt.Fprintf(&b, "<el-input-number class=\"w100\" controls-position=\"right\" v-model=\"formItems.%s\" :placeholder=\"%s\" />\n", f.Name, placeholderEnter)
	case "switch":
		fmt.Fprintf(&b, "<el-switch v-model=\"formItems.%s\" />\n", f.Name)
	case "remoteSelect", "remoteSelects":
		importSet["RemoteSelect"] = true
		multiple := ""
		clearAttrs := ""
		if f.DesignType == "remoteSelects" {
			multiple = " :multiple=\"true\""
		} else {
			clearAttrs = " :empty-values=\"[null, 0]\" :value-on-clear=\"0\""
		}
		remote := parseRemoteInfo(f)
		field := remote.Field
		if field == "" {
			field = "name"
		}
		pk := remote.PK
		if pk == "" {
			pk = "id"
		}
		url := formStr(f.Form, "remoteUrl")
		if url == "" && remote.Table != "" {
			url = "/admin/" + remote.Table + "/list"
		}
		fmt.Fprintf(&b, "<RemoteSelect v-model=\"formItems.%s\" pk=\"%s\" field=\"%s\" remote-url=\"%s\"%s%s />\n", f.Name, pk, field, url, multiple, clearAttrs)
	case "array":
		importSet["ArrayInput"] = true
		fmt.Fprintf(&b, "<ArrayInput v-model=\"formItems.%s\" />\n", f.Name)
	case "datetime":
		fmt.Fprintf(&b, "<el-date-picker class=\"w100\" v-model=\"formItems.%s\" value-format=\"YYYY-MM-DDTHH:mm:ssZ\" type=\"datetime\" :placeholder=\"%s\" />\n", f.Name, placeholderSelect)
	case "date":
		fmt.Fprintf(&b, "<el-date-picker class=\"w100\" v-model=\"formItems.%s\" value-format=\"YYYY-MM-DDTHH:mm:ssZ\" type=\"date\" :placeholder=\"%s\" />\n", f.Name, placeholderSelect)
	case "year":
		fmt.Fprintf(&b, "<el-date-picker class=\"w100\" v-model=\"formItems.%s\" value-format=\"YYYY\" type=\"year\" :placeholder=\"%s\" />\n", f.Name, placeholderSelect)
	case "time":
		fmt.Fprintf(&b, "<el-time-picker class=\"w100\" v-model=\"formItems.%s\" value-format=\"HH:mm:ss\" :placeholder=\"%s\" />\n", f.Name, placeholderSelect)
	case "password":
		fmt.Fprintf(&b, "<el-input v-model=\"formItems.%s\" type=\"password\" autocomplete=\"new-password\" :placeholder=\"%s\" />\n", f.Name, placeholderEnter)
	case "editor":
		importSet["Editor"] = true
		fmt.Fprintf(&b, "<Editor v-model=\"formItems.%s\" @keyup.enter.stop=\"\" @keyup.ctrl.enter=\"manager.submitForm(formRef)\" />\n", f.Name)
	default:
		fmt.Fprintf(&b, "<el-input v-model=\"formItems.%s\" :placeholder=\"%s\" />\n", f.Name, placeholderEnter)
	}

	b.WriteString("</el-form-item>\n\n")
	return b.String()
}
