package gen

import (
	_ "embed"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
)

//go:embed tmpl/service.tmpl
var serviceTmplStr string

// serviceTmpl 服务文件内容模板
var serviceTmpl = template.Must(template.New("service").Parse(serviceTmplStr))

// serviceFileData 服务文件模板数据
type serviceFileData struct {
	Package     string // 文件包名
	ModelPkg    string // 模型所在包名
	ModelImport string // 模型包导入路径
	ModelName   string // 模型名
	Comment     string // 表注释
	Imports     string // import 块内容，空串表示无 import
	RepoAlias   string // 仓储导入别名
	Methods     string // 追加的关联查询方法块，空串表示无
}

// buildServiceFileData 组装服务文件模板数据
func buildServiceFileData(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) serviceFileData {
	data := serviceFileData{
		Package:     basic["service"].Package,
		ModelPkg:    basic["model"].Package,
		ModelImport: crud.ModulePath + "/" + basic["model"].Dir,
		ModelName:   basic["model"].Name,
		Comment:     table.Comment,
	}
	repo := crud.ModuleImportOf(basic["repository"], "repo")
	imports := []string{
		crud.ImportSpec(data.ModelImport),
		crud.ImportSpec(repo.Path, repo.Alias),
		crud.ImportSpec(crud.ModulePath + "/internal/service"),
	}
	data.Methods = buildRemoteSelectsServiceMethods(fields, data.ModelPkg, data.ModelName)
	if data.Methods != "" {
		imports = append(imports, crud.ImportSpec("context"))
	}
	data.Imports = crud.BuildImportBlock(imports)
	data.RepoAlias = repo.Alias
	return data
}

// CreateServiceFile 创建服务文件
func CreateServiceFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	content, err := crud.RenderTmpl(serviceTmpl, buildServiceFileData(basic, table, fields))
	if err != nil {
		return err
	}
	return crud.WriteGeneratedFile(basic["service"].File, content)
}
