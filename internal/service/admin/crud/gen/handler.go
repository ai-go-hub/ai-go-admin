package gen

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
)

//go:embed tmpl/handler.tmpl
var handlerTmplStr string

// handlerTmpl 控制器文件内容模板
var handlerTmpl = template.Must(template.New("handler").Parse(handlerTmplStr))

// handlerFileData 控制器文件模板数据
type handlerFileData struct {
	Package        string // 文件包名
	ModelPkg       string // 模型所在包名
	ModelImport    string // 模型包导入路径
	ModelName      string // 模型名
	Comment        string // 表注释
	Imports        string // import 块内容，空串表示无 import
	SvcAlias       string // 服务导入别名
	HandlerOptions string // NewHandler 的预加载/列表适配器选项块，空串表示无
}

// buildHandlerFileData 组装控制器文件模板数据
func buildHandlerFileData(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) handlerFileData {
	data := handlerFileData{
		Package:     basic["handler"].Package,
		ModelPkg:    basic["model"].Package,
		ModelImport: crud.ModulePath + "/" + basic["model"].Dir,
		ModelName:   basic["model"].Name,
		Comment:     table.Comment,
	}
	svc := crud.ModuleImportOf(basic["service"], "svc")

	imports := []string{
		crud.ImportSpec(crud.ModulePath + "/internal/handler"),
		crud.ImportSpec(data.ModelImport),
		crud.ImportSpec(svc.Path, svc.Alias),
	}
	var options []string
	preloads := buildHandlerPreloads(fields)
	if len(preloads) > 0 {
		imports = append(imports, crud.ImportSpec(crud.ModulePath+"/internal/repository"))
		options = append(options, buildHandlerPreloadsOptions(preloads))
	}
	if crud.HasRemoteSelectsDesign(fields) {
		imports = append(imports,
			crud.ImportSpec("context"),
			crud.ImportSpec("errors"),
			crud.ImportSpec("github.com/ai-go-hub/ai-go-admin/internal/service"),
		)
		options = append(options, buildHandlerRemoteSelectsAdapterOptions(fields, data.ModelPkg, data.ModelName))
	}
	if omit := buildHandlerOmitFieldsOptions(fields); omit != "" {
		options = append(options, omit)
	}
	imports = append(imports, crud.ImportSpec("github.com/gin-gonic/gin"))
	data.Imports = crud.BuildImportBlock(imports)

	if len(options) > 0 {
		data.HandlerOptions = ",\n\t\t" + strings.Join(options, ",\n\t\t")
	}
	data.SvcAlias = svc.Alias
	return data
}

// CreateHandlerFile 创建控制器文件
func CreateHandlerFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	content, err := crud.RenderTmpl(handlerTmpl, buildHandlerFileData(basic, table, fields))
	if err != nil {
		return err
	}
	return crud.WriteGeneratedFile(basic["handler"].File, content)
}
