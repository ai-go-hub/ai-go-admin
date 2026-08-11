package gen

import (
	"fmt"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"
)

// remoteInfo 关联定义信息
type remoteInfo struct {
	PK             string   // 远程表主键，如 id
	Field          string   // 远程表展示字段，如 username
	Table          string   // 远程表名，如 admins
	ModelFile      string   // 远程模型文件，如 internal/model/admin.go
	ModelName      string   // 远程模型名，如 Admin
	ModelPkg       string   // 远程模型包名，如 model
	RelationFields []string // 远程表展示关联字段列表
}

// formStr 取字段表单属性中指定键的字符串值
func formStr(form map[string]any, key string) string {
	if form == nil {
		return ""
	}
	switch v := form[key].(type) {
	case string:
		return v
	case nil:
		return ""
	default:
		return fmt.Sprint(v)
	}
}

// parseRemoteInfo 解析字段的远程关联信息
func parseRemoteInfo(f dto.CRUDFields) remoteInfo {
	info := remoteInfo{
		PK:        formStr(f.Form, "remotePk"),
		Field:     formStr(f.Form, "remoteField"),
		Table:     formStr(f.Form, "remoteTable"),
		ModelFile: formStr(f.Form, "remoteModelFile"),
		ModelName: formStr(f.Form, "remoteModelName"),
		ModelPkg:  formStr(f.Form, "remoteModelPackage"),
	}

	// relationFields 为逗号分隔的远程表展示字段名列表
	for rf := range strings.SplitSeq(formStr(f.Form, "relationFields"), ",") {
		if rf = strings.TrimSpace(rf); rf != "" {
			info.RelationFields = append(info.RelationFields, rf)
		}
	}
	return info
}

// RemoteSelectsMethodName 生成 remoteSelects 字段的批量联查方法名，如 user_ids -> AdminUsernameByUserIds
func RemoteSelectsMethodName(f dto.CRUDFields) string {
	info := parseRemoteInfo(f)
	return info.ModelName + snakeToPascal(info.Field) + "By" + snakeToPascal(f.Name)
}

// RemoteSelectsListFieldName 生成 remoteSelects 字段的模型展示字段名，如 AdminUsernameList
func RemoteSelectsListFieldName(f dto.CRUDFields) string {
	info := parseRemoteInfo(f)
	return info.ModelName + snakeToPascal(info.Field) + "List"
}

// buildModelAssociationField 构建 belongs to 关联字段
// modelPkg 为当前模型文件所在包名，远程模型跨包时类型加包名前缀
func buildModelAssociationField(f dto.CRUDFields, modelPkg string) (modelFileField, bool) {
	info := parseRemoteInfo(f)
	if info.ModelName == "" || f.Name == "" {
		return modelFileField{}, false
	}
	typ := "*" + info.ModelName
	if info.ModelPkg != "" && info.ModelPkg != modelPkg {
		typ = "*" + info.ModelPkg + "." + info.ModelName
	}
	return modelFileField{
		Name:    info.ModelName,
		Type:    typ,
		JSONTag: util.PascalToSnake(info.ModelName) + ",omitempty",
		GormTag: "foreignKey:" + snakeToPascal(f.Name) + ";references:" + snakeToPascal(info.PK),
	}, true
}

// buildHandlerPreloads 收集 remoteSelect 字段的关联名列表
func buildHandlerPreloads(fields []dto.CRUDFields) []string {
	var preloads []string
	for _, f := range fields {
		if f.DesignType != "remoteSelect" {
			continue
		}
		if info := parseRemoteInfo(f); info.ModelName != "" {
			preloads = append(preloads, info.ModelName)
		}
	}
	return preloads
}

// buildHandlerPreloadsOptions 构建控制器 NewHandler 的预加载选项块，无预加载时返回空串
func buildHandlerPreloadsOptions(preloads []string) string {
	if len(preloads) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("handler.WithPreloads([]repository.Preload{\n")
	for _, p := range preloads {
		fmt.Fprintf(&b, "\t\t\t\t{Association: %q},\n", p)
	}
	b.WriteString("\t\t\t})")
	return b.String()
}

// buildHandlerOmitFieldsOptions 构建控制器 NewHandler 的 WithOmitFields 选项
func buildHandlerOmitFieldsOptions(fields []dto.CRUDFields) string {
	pk := crud.PkFieldName(fields)

	var create, update []string
	if pk != "" {
		update = append(update, pk)
	}

	for _, f := range fields {
		switch f.DesignType {
		case "remoteSelect":
			if info := parseRemoteInfo(f); info.ModelName != "" {
				update = append(update, util.PascalToSnake(info.ModelName))
			}
		case "remoteSelects":
			if info := parseRemoteInfo(f); info.ModelName != "" && info.Field != "" {
				update = append(update, util.PascalToSnake(RemoteSelectsListFieldName(f)))
			}
		}
	}

	if len(create) == 0 && len(update) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("handler.WithOmitFields(handler.ActionFields{")
	if len(create) > 0 {
		fmt.Fprintf(&b, "\n\t\t\t\tCreate: %s,", crud.QuoteStrSlice(create))
	}
	if len(update) > 0 {
		fmt.Fprintf(&b, "\n\t\t\t\tUpdate: %s,", crud.QuoteStrSlice(update))
	}
	b.WriteString("\n\t\t\t})")
	return b.String()
}

// buildRepoRemoteSelectsMethods 构建 remoteSelects 字段的仓储联查方法块列表
// repoModelPkg/repoModelImport 为仓储已导入的模型包名与路径；返回方法代码块与需要额外导入的模型包 import 行
func buildRepoRemoteSelectsMethods(fields []dto.CRUDFields, modelName, repoModelPkg, repoModelImport string) ([]string, []string) {
	var methods []string
	var imports []string
	for _, f := range fields {
		if f.DesignType != "remoteSelects" {
			continue
		}
		info := parseRemoteInfo(f)
		if info.ModelName == "" || info.Table == "" || info.Field == "" || info.PK == "" || f.Name == "" {
			continue
		}
		remotePkg := info.ModelPkg
		if remotePkg == "" {
			remotePkg = repoModelPkg
		}
		if info.ModelFile != "" {
			if remoteImport := crud.ModulePath + "/" + filesystem.Dir(info.ModelFile); remoteImport != repoModelImport {
				imports = append(imports, crud.ImportSpec(remoteImport))
			}
		}
		fieldPascal := snakeToPascal(f.Name)
		methodName := RemoteSelectsMethodName(f)
		methods = append(methods, buildRepoRemoteSelectsMethod(methodName, modelName, fieldPascal, info, remotePkg, repoModelPkg))
	}
	return methods, imports
}

// buildRepoRemoteSelectsMethod 构建单个 remoteSelects 批量联查方法
func buildRepoRemoteSelectsMethod(methodName, modelName, fieldPascal string, info remoteInfo, remotePkg, repoModelPkg string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "// %s 批量查询记录关联的 %s 表 %s 字段值，返回 map[记录主键][]string\n", methodName, info.Table, info.Field)
	b.WriteString("// 一次收齐所有关联主键并去重，仅需一次查询关联表，避免列表 N+1\n")
	fmt.Fprintf(&b, "func (r *%sRepository) %s(ctx context.Context, records []%s.%s) (map[uint][]string, error) {\n", modelName, methodName, repoModelPkg, modelName)
	b.WriteString("\t// 收集每条记录关联的远程主键\n")
	b.WriteString("\tidSet := make(map[uint]struct{})\n")
	b.WriteString("\tidsByRecord := make(map[uint][]uint, len(records))\n")
	b.WriteString("\tfor _, rec := range records {\n")
	b.WriteString("\t\tvar ids []uint\n")
	fmt.Fprintf(&b, "\t\tif rec.%s != nil {\n", fieldPascal)
	fmt.Fprintf(&b, "\t\t\tif err := json.Unmarshal(*rec.%s, &ids); err != nil {\n\t\t\t\treturn nil, err\n\t\t\t}\n\t\t}\n", fieldPascal)
	b.WriteString("\t\tidsByRecord[rec.ID] = ids\n")
	b.WriteString("\t\tfor _, id := range ids {\n\t\t\tidSet[id] = struct{}{}\n\t\t}\n\t}\n\n")
	b.WriteString("\t// 一次查询关联表: 主键 -> 关联字段值\n")
	b.WriteString("\tfieldByID := make(map[uint]string)\n")
	b.WriteString("\trelationIDs := make([]uint, 0, len(idSet))\n")
	b.WriteString("\tfor id := range idSet {\n\t\trelationIDs = append(relationIDs, id)\n\t}\n")
	b.WriteString("\tif len(relationIDs) > 0 {\n")
	fmt.Fprintf(&b, "\t\tvar relationFields []%s.%s\n", remotePkg, info.ModelName)
	fmt.Fprintf(&b, "\t\tif err := r.DB().WithContext(ctx).\n\t\t\tModel(&%s.%s{}).\n\t\t\tWhere(%q, relationIDs).\n\t\t\tFind(&relationFields).Error; err != nil {\n\t\t\treturn nil, err\n\t\t}\n", remotePkg, info.ModelName, info.PK+" IN ?")
	b.WriteString("\t\tfor _, f := range relationFields {\n")
	fmt.Fprintf(&b, "\t\t\tswitch v := any(f.%s).(type) {\n", snakeToPascal(info.Field))
	b.WriteString("\t\t\tcase string:\n")
	fmt.Fprintf(&b, "\t\t\t\tfieldByID[f.%s] = v\n", snakeToPascal(info.PK))
	b.WriteString("\t\t\tcase *string:\n")
	fmt.Fprintf(&b, "\t\t\t\tif v != nil {\n\t\t\t\t\tfieldByID[f.%s] = *v\n\t\t\t\t}\n", snakeToPascal(info.PK))
	b.WriteString("\t\t\t}\n")
	b.WriteString("\t\t}\n")
	b.WriteString("\t}\n\n")
	b.WriteString("\t// 逐条记录组装关联字段值列表\n")
	b.WriteString("\tresult := make(map[uint][]string, len(records))\n")
	b.WriteString("\tfor recID, ids := range idsByRecord {\n")
	b.WriteString("\t\tfields := make([]string, 0, len(ids))\n")
	b.WriteString("\t\tfor _, id := range ids {\n")
	b.WriteString("\t\t\tif v, ok := fieldByID[id]; ok {\n\t\t\t\tfields = append(fields, v)\n\t\t\t}\n\t\t}\n")
	b.WriteString("\t\tresult[recID] = fields\n\t}\n")
	b.WriteString("\treturn result, nil\n}\n")
	return b.String()
}

// buildRemoteSelectsServiceMethods 构建 remoteSelects 字段的服务层委托方法块，无 remoteSelects 字段时返回空串
func buildRemoteSelectsServiceMethods(fields []dto.CRUDFields, modelPkg, modelName string) string {
	var b strings.Builder
	for _, f := range fields {
		if f.DesignType != "remoteSelects" {
			continue
		}
		info := parseRemoteInfo(f)
		if info.ModelName == "" || info.Field == "" || f.Name == "" {
			continue
		}
		methodName := RemoteSelectsMethodName(f)
		fmt.Fprintf(&b, "\n// %s 批量查询记录关联的 %s 表 %s 字段值，返回 map[记录主键][]string\n", methodName, info.Table, info.Field)
		fmt.Fprintf(&b, "func (s *%sService) %s(ctx context.Context, records []%s.%s) (map[uint][]string, error) {\n", modelName, methodName, modelPkg, modelName)
		fmt.Fprintf(&b, "\treturn s.repo.%s(ctx, records)\n}\n", methodName)
	}
	return b.String()
}

// buildHandlerRemoteSelectsAdapterOptions 构建 remoteSelects 字段的列表数据适配器选项块，无 remoteSelects 字段时返回空串
func buildHandlerRemoteSelectsAdapterOptions(fields []dto.CRUDFields, modelPkg, modelName string) string {
	var blocks []string
	for _, f := range fields {
		if f.DesignType != "remoteSelects" {
			continue
		}
		info := parseRemoteInfo(f)
		if info.ModelName == "" || info.Field == "" || f.Name == "" {
			continue
		}
		methodName := RemoteSelectsMethodName(f)
		listField := RemoteSelectsListFieldName(f)
		mapVar := util.SnakeToLowerCamel(util.PascalToSnake(listField))
		blocks = append(blocks, fmt.Sprintf(`
		%s, err := h.svc.%s(ctx, items)
		if err != nil {
			return nil, err
		}
		for i := range items {
			items[i].%s = %s[items[i].ID]
		}`, mapVar, methodName, listField, mapVar)+"\n")
	}
	if len(blocks) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("handler.WithAdapter(handler.Adapter{\n")
	b.WriteString("\t\tList: func(ctx context.Context, data any, opts service.Options) (any, error) {\n")
	fmt.Fprintf(&b, "\t\t\titems, ok := data.([]%s.%s)\n", modelPkg, modelName)
	b.WriteString("\t\t\tif !ok {\n\t\t\t\treturn nil, errors.New(\"列表数据类型错误\")\n\t\t\t}\n")
	for _, blk := range blocks {
		b.WriteString(blk)
	}
	b.WriteString("\t\t\treturn items, nil\n")
	b.WriteString("\t\t},\n\t})")
	return b.String()
}

// joinRepoMethods 拼接方法块列表，空时返回空串
func joinRepoMethods(methods []string) string {
	if len(methods) == 0 {
		return ""
	}
	return "\n\n" + strings.Join(methods, "\n")
}
