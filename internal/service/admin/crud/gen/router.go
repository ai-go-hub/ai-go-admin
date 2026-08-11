package gen

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
)

//go:embed tmpl/router.tmpl
var routerTmplStr string

// routerTmpl 路由文件内容模板
var routerTmpl = template.Must(template.New("router").Parse(routerTmplStr))

// routerFileData 路由文件模板数据
type routerFileData struct {
	Package      string // 文件包名
	ModelPkg     string // 模型所在包名
	ModelImport  string // 模型包导入路径
	ModelName    string // 模型名
	Comment      string // 表注释
	Imports      string // import 块内容，空串表示无 import
	RepoAlias    string // 仓储导入别名
	SvcAlias     string // 服务导入别名
	HandlerAlias string // 控制器导入别名
	IsAdmin      bool   // 是否为后台（admin）路由
	GroupPath    string // 路由分组路径
}

// buildRouterFileData 组装路由文件模板数据
func buildRouterFileData(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable) routerFileData {
	data := routerFileData{
		Package:     basic["router"].Package,
		ModelPkg:    basic["model"].Package,
		ModelImport: crud.ModulePath + "/" + basic["model"].Dir,
		ModelName:   basic["model"].Name,
		Comment:     table.Comment,
	}
	handlerImp := crud.ModuleImportOf(basic["handler"], "handler")
	repo := crud.ModuleImportOf(basic["repository"], "repo")
	svc := crud.ModuleImportOf(basic["service"], "svc")

	routerInfo := basic["router"]
	isAdmin := routerInfo.App == "admin"
	// 路由分组路径取前端 RoutePath: admin 挂在 /admin 组下不拼 app，非 admin 在根路由需拼 app
	routePath := strings.Trim(table.RoutePath, "/")
	groupPath := "/" + routePath
	if !isAdmin {
		groupPath = "/" + routerInfo.App + groupPath
	}

	data.Imports = crud.BuildImportBlock([]string{
		crud.ImportSpec(handlerImp.Path, handlerImp.Alias),
		crud.ImportSpec(repo.Path, repo.Alias),
		crud.ImportSpec(crud.ModulePath + "/internal/router/registry"),
		crud.ImportSpec(svc.Path, svc.Alias),
		crud.ImportSpec("github.com/gin-gonic/gin"),
	})
	data.RepoAlias = repo.Alias
	data.SvcAlias = svc.Alias
	data.HandlerAlias = handlerImp.Alias
	data.IsAdmin = isAdmin
	data.GroupPath = groupPath
	return data
}

// CreateRouterFile 创建路由文件，并确保 router.go 空白导入该路由包（触发 init() 以完成路由注册）
func CreateRouterFile(basic map[string]dto.GenerateFileBasicDataInfo, table dto.CRUDTable, fields []dto.CRUDFields) error {
	content, err := crud.RenderTmpl(routerTmpl, buildRouterFileData(basic, table))
	if err != nil {
		return err
	}
	if err := crud.WriteGeneratedFile(basic["router"].File, content); err != nil {
		return err
	}
	return registry.AddRouterImport(crud.ModulePath + "/" + basic["router"].Dir)
}
