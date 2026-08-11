package crud

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	"github.com/ai-go-hub/ai-go-admin/internal/repository"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	repoCrud "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/crud"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
	"github.com/ai-go-hub/ai-go-admin/pkg/airx"
	"github.com/ai-go-hub/ai-go-admin/pkg/filesystem"
	"github.com/ai-go-hub/ai-go-admin/pkg/util"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CrudLogService CRUD 记录服务
type CrudLogService struct {
	service.IService[model.CrudLog]
	repo     *repoCrud.CrudLogRepository
	repoAuth *repoAuth.AdminRuleRepository
}

// NewCrudLogService 创建 CRUD 记录服务实例
func NewCrudLogService(repo *repoCrud.CrudLogRepository, repoAuth *repoAuth.AdminRuleRepository) *CrudLogService {
	return &CrudLogService{
		IService: service.NewService(repo),
		repo:     repo,
		repoAuth: repoAuth,
	}
}

// Delete 读取 CRUD 记录，并删除对应记录中已生成的代码文件、语言包与视图目录（不删除记录本身）
func (s *CrudLogService) Delete(ctx context.Context, opts service.Options) error {
	if len(opts.PrimaryKeyValues) == 0 {
		return errors.New("缺少主键")
	}

	logs, err := s.repo.ListByIDs(ctx, opts.PrimaryKeyValues)
	if err != nil {
		return err
	}

	for _, log := range logs {
		var table dto.CRUDTable
		if err := json.Unmarshal(log.Table, &table); err != nil {
			// 记录数据损坏时跳过，不影响其余记录
			continue
		}
		removeGeneratedModule(table)
		if err := s.removeMenuRules(ctx, table); err != nil {
			return err
		}

		// 文件删除后，将该条记录的 Status 标记为已删除
		if err := s.repo.Update(ctx, map[string]any{"status": "deleted"}, repository.Options{
			PrimaryKeyValue: strconv.FormatUint(uint64(log.ID), 10),
		}); err != nil {
			return err
		}
	}
	return nil
}

// removeGeneratedModule 删除生成模块的代码文件、语言包、视图目录，并移除 router.go 中新增的空白导入（尽力而为）
func removeGeneratedModule(table dto.CRUDTable) {
	// 删除期间暂停 air 热更新，避免文件未删完进程被重启中断
	airx.Pause()
	defer airx.Resume()

	for _, file := range []string{table.ModelFile, table.HandlerFile, table.RepositoryFile, table.RouterFile, table.ServiceFile} {
		filesystem.RemoveFileWithCleanup(file)
	}

	// 路由包若因本次删除而不含 .go 文件，则移除 router.go 中对应空白导入
	if routerDir := filesystem.Dir(table.RouterFile); !dirHasGoFiles(routerDir) {
		registry.RemoveRouterImport(ModulePath + "/" + routerDir)
	}

	// 语言包文件由视图目录推导: 解析 views 目录，再据此生成 lang 文件路径
	if viewsBasic := ParseGenerateFileBasicData("views", table.WebViewsDir); viewsBasic.Path != "" {
		langBasic := GenerateFileBasicData("lang", viewsBasic.Path, viewsBasic.App)
		filesystem.RemoveFileWithCleanup(langBasic.CnFile)
		filesystem.RemoveFileWithCleanup(langBasic.EnFile)
	}

	// 前端视图文件: 单独删除 index.vue 与 dialogForm.vue
	filesystem.RemoveFileWithCleanup(table.WebViewsDir + "/index.vue")
	filesystem.RemoveFileWithCleanup(table.WebViewsDir + "/dialogForm.vue")
}

// CrudLog 记录 CRUD 日志，ID 大于 0 时更新该记录状态，否则创建新记录
func CrudLog(ctx context.Context, repo *repoCrud.CrudLogRepository, data dto.CrudLogData) (uint, error) {
	id := util.FromPtr(data.ID)
	if id > 0 {
		if err := repo.Update(ctx, map[string]any{"status": data.Status, "sql": data.Sql}, repository.Options{
			PrimaryKeyValue: strconv.FormatUint(uint64(id), 10),
		}); err != nil {
			return 0, err
		}
		return id, nil
	}

	tableData, err := json.Marshal(data.Table)
	if err != nil {
		return 0, err
	}
	fieldsData, err := json.Marshal(data.Fields)
	if err != nil {
		return 0, err
	}

	log := model.CrudLog{
		Name:    data.Table.Name,
		Comment: data.Table.Comment,
		Table:   datatypes.JSON(tableData),
		Fields:  datatypes.JSON(fieldsData),
		Sql:     data.Sql,
		Status:  data.Status,
	}
	if err := repo.Create(ctx, &log, repository.Options{}); err != nil {
		return 0, err
	}
	return log.ID, nil
}

// UpdateLogStatus 更新 CRUD 日志状态（如 failed/succeeded/deleted）
func UpdateLogStatus(ctx context.Context, repo *repoCrud.CrudLogRepository, id uint, status string) error {
	return repo.Update(ctx, map[string]any{"status": status}, repository.Options{
		PrimaryKeyValue: strconv.FormatUint(uint64(id), 10),
	})
}

// dirHasGoFiles 目录下是否仍存在 .go 文件
func dirHasGoFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".go") {
			return true
		}
	}
	return false
}

// removeMenuRules 删除 CRUD 生成时写入的后台菜单与权限节点（级联清理菜单树）
func (s *CrudLogService) removeMenuRules(ctx context.Context, table dto.CRUDTable) error {
	// 路由路径取前端提供的 RoutePath，与 CreateMenuRule 保持一致
	menuPath := table.RoutePath
	if menuPath == "" {
		return nil
	}
	menu, err := s.repoAuth.FindByName(ctx, menuPath)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return s.repoAuth.BatchDelete(ctx, []uint{menu.ID}, true)
}
