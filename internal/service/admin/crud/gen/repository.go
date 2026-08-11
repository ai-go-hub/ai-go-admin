package gen

import (
	_ "embed"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
)

//go:embed tmpl/repository.tmpl
var repositoryTmplStr string

// repositoryTmpl 仓储文件内容模板
var repositoryTmpl = template.Must(template.New("repository").Parse(repositoryTmplStr))

// repositoryFileData 仓储文件模板数据
type repositoryFileData struct {
	Package     string // 文件包名
	ModelPkg    string // 模型所在包名
	ModelImport string // 模型包导入路径
	ModelName   string // 模型名
	Comment     string // 表注释
	Imports     string // import 块内容，空串表示无 import
	Methods     string // 追加的查询方法块，空串表示无
}

// buildRepositoryFileData 组装仓储文件模板数据
func buildRepositoryFileData(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) repositoryFileData {
	data := repositoryFileData{
		Package:     basic["repository"].Package,
		ModelPkg:    basic["model"].Package,
		ModelImport: crud.ModulePath + "/" + basic["model"].Dir,
		ModelName:   basic["model"].Name,
		Comment:     table.Comment,
	}

	imports := []string{
		crud.ImportSpec(data.ModelImport),
		crud.ImportSpec(crud.ModulePath + "/internal/repository"),
	}
	methodBlocks, remoteImports := buildRepoRemoteSelectsMethods(fields, data.ModelName, data.ModelPkg, data.ModelImport)
	methods := joinRepoMethods(methodBlocks)
	if len(methodBlocks) > 0 {
		imports = append(imports,
			crud.ImportSpec("context"),
			crud.ImportSpec("encoding/json"),
		)
		imports = append(imports, remoteImports...)
	}
	data.Imports = crud.BuildImportBlock(imports)

	data.Methods = methods
	return data
}

// CreateRepositoryFile 创建仓储文件
func CreateRepositoryFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	content, err := crud.RenderTmpl(repositoryTmpl, buildRepositoryFileData(basic, table, fields))
	if err != nil {
		return err
	}
	return crud.WriteGeneratedFile(basic["repository"].File, content)
}
