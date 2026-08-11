package gen

import (
	"strings"
	"testing"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service/admin/crud"
)

// relationTable 含关联字段的示例表
func relationTable() dto.CRUDTable {
	return dto.CRUDTable{
		Name:           "user_log",
		Comment:        "用户日志",
		ModelFile:      "internal/model/user_log.go",
		RepositoryFile: "internal/repository/admin/user/log.go",
		ServiceFile:    "internal/service/admin/user/log.go",
		HandlerFile:    "internal/handler/admin/user/log.go",
		RouterFile:     "internal/router/admin/user/log.go",
		ColumnFields:   []string{"id", "admin_id", "user_ids", "created_at"},
	}
}

// relationFields 含 remoteSelect 与 remoteSelects 字段
func relationFields() []dto.CRUDFields {
	return []dto.CRUDFields{
		{Name: "id", Type: "bigint", Comment: "ID", PrimaryKey: true, Generated: "GENERATED ALWAYS"},
		{Name: "admin_id", Type: "bigint", Comment: "管理员", DesignType: "remoteSelect", Form: map[string]any{
			"remotePk": "id", "remoteField": "username", "remoteTable": "admins",
			"remoteModelFile": "internal/model/admin.go", "remoteModelName": "Admin", "remoteModelPackage": "model", "relationFields": "username,nickname",
		}},
		{Name: "user_ids", Type: "jsonb", Comment: "用户", DesignType: "remoteSelects", Form: map[string]any{
			"remotePk": "id", "remoteField": "username", "remoteTable": "admins",
			"remoteModelFile": "internal/model/admin.go", "remoteModelName": "Admin", "remoteModelPackage": "model", "relationFields": "username,nickname",
		}},
		{Name: "created_at", Type: "timestamptz", Comment: "创建时间", Null: true},
	}
}

// TestBuildModelAssociation 验证 remoteSelect 字段生成 belongs to 关联字段
func TestBuildModelAssociation(t *testing.T) {
	content, err := crud.RenderTmpl(modelTmpl, buildModelFileData(tableBasic(relationTable()), relationTable(), relationFields()))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Admin *Admin",
		`gorm:"foreignKey:AdminID;references:ID" json:"admin,omitempty"`,
		"AdminUsernameList []string",
		`gorm:"-" json:"admin_username_list"`,
		`"time"`, // created_at 仍需要 time
	} {
		if !strings.Contains(content, want) {
			t.Errorf("模型缺少 %q:\n%s", want, content)
		}
	}
	t.Logf("模型渲染:\n%s", content)
}

// TestBuildHandlerPreloads 验证 remoteSelect 字段生成 WithPreloads
func TestBuildHandlerPreloads(t *testing.T) {
	content, err := crud.RenderTmpl(handlerTmpl, buildHandlerFileData(tableBasic(relationTable()), relationTable(), relationFields()))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"handler.WithPreloads([]repository.Preload{",
		`{Association: "Admin"},`,
		"\"github.com/ai-go-hub/ai-go-admin/internal/repository\"",
		"handler.WithAdapter(handler.Adapter{",
		"h.svc.AdminUsernameByUserIds(ctx, items)",
		"items[i].AdminUsernameList = adminUsernameList[items[i].ID]",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("控制器缺少 %q:\n%s", want, content)
		}
	}
	t.Logf("控制器渲染:\n%s", content)
}

// TestBuildRepositoryRemoteSelects 验证 remoteSelects 字段生成 IN 联查方法
func TestBuildRepositoryRemoteSelects(t *testing.T) {
	content, err := crud.RenderTmpl(repositoryTmpl, buildRepositoryFileData(tableBasic(relationTable()), relationTable(), relationFields()))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"func (r *UserLogRepository) AdminUsernameByUserIds(ctx context.Context, records []model.UserLog) (map[uint][]string, error) {",
		"json.Unmarshal(*rec.UserIds, &ids)",
		`Model(&model.Admin{})`,
		`Where("id IN ?", relationIDs)`,
		`Find(&relationFields)`,
		`switch v := any(f.Username).(type) {`,
		`case *string:`,
		`fieldByID[f.ID] = *v`,
		`"context"`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("仓储缺少 %q:\n%s", want, content)
		}
	}
	t.Logf("仓储渲染:\n%s", content)
}

// TestBuildModelAssociationMultiWord 验证多词模型名（如 AdminRule）关联字段 JSONTag 转 snake_case
func TestBuildModelAssociationMultiWord(t *testing.T) {
	table := relationTable()
	fields := []dto.CRUDFields{
		{Name: "admin_rule_id", Type: "bigint", DesignType: "remoteSelect", Form: map[string]any{
			"remotePk": "id", "remoteModelName": "AdminRule", "remoteModelPackage": "model",
		}},
	}
	content, err := crud.RenderTmpl(modelTmpl, buildModelFileData(tableBasic(table), table, fields))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"AdminRule *AdminRule",
		`gorm:"foreignKey:AdminRuleID;references:ID" json:"admin_rule,omitempty"`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("模型缺少 %q:\n%s", want, content)
		}
	}
	t.Logf("多词模型渲染:\n%s", content)
}

// TestBuildRepositoryRemoteSelectsCrossPackage 验证 remoteSelects 远程模型跨包时额外导入关联模型包
func TestBuildRepositoryRemoteSelectsCrossPackage(t *testing.T) {
	table := relationTable()
	fields := []dto.CRUDFields{
		{Name: "user_ids", Type: "jsonb", DesignType: "remoteSelects", Form: map[string]any{
			"remotePk": "id", "remoteField": "username", "remoteTable": "exp_logs",
			"remoteModelFile": "internal/model/exp/log.go", "remoteModelName": "Log", "remoteModelPackage": "exp",
		}},
	}
	content, err := crud.RenderTmpl(repositoryTmpl, buildRepositoryFileData(tableBasic(table), table, fields))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Model(&exp.Log{})",
		`"github.com/ai-go-hub/ai-go-admin/internal/model/exp"`,
		"func (r *UserLogRepository) LogUsernameByUserIds(ctx context.Context, records []model.UserLog) (map[uint][]string, error) {",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("仓储缺少 %q:\n%s", want, content)
		}
	}
	t.Logf("跨包仓储渲染:\n%s", content)
}

// TestBuildIndexColumnsRelations 验证远程下拉字段生成关联展示列
func TestBuildIndexColumnsRelations(t *testing.T) {
	table := relationTable()
	fields := relationFields()
	cols := buildIndexColumns("test", table, fields, false)
	for _, want := range []string{
		`prop: 'admin.username'`,
		`prop: 'admin.nickname'`,
		`label: t('test.adminUsername')`,
		`label: t('test.adminUsernameList')`,
		`operator: 'ILIKE'`,
		`render: 'tags'`,
		`prop: 'admin_username_list'`,
		`operator: false`,
	} {
		if !strings.Contains(cols, want) {
			t.Errorf("列缺少 %q:\n%s", want, cols)
		}
	}
	t.Logf("列渲染:\n%s", cols)
}

// TestBuildRemoteSelectsServiceMethods 验证 remoteSelects 字段生成服务层委托方法
func TestBuildRemoteSelectsServiceMethods(t *testing.T) {
	table := relationTable()
	fields := relationFields()
	content, err := crud.RenderTmpl(serviceTmpl, buildServiceFileData(tableBasic(table), table, fields))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"context"`,
		"func (s *UserLogService) AdminUsernameByUserIds(ctx context.Context, records []model.UserLog) (map[uint][]string, error) {",
		"return s.repo.AdminUsernameByUserIds(ctx, records)",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("服务缺少 %q:\n%s", want, content)
		}
	}
	t.Logf("服务渲染:\n%s", content)
}

// tableBasic 按文件路径解析各生成文件的基础数据
func tableBasic(table dto.CRUDTable) map[string]dto.GenerateFileBasicDataInfo {
	return map[string]dto.GenerateFileBasicDataInfo{
		"model":      crud.ParseGenerateFileBasicData("model", table.ModelFile),
		"repository": crud.ParseGenerateFileBasicData("repository", table.RepositoryFile),
		"service":    crud.ParseGenerateFileBasicData("service", table.ServiceFile),
		"handler":    crud.ParseGenerateFileBasicData("handler", table.HandlerFile),
		"router":     crud.ParseGenerateFileBasicData("router", table.RouterFile),
	}
}

// TestBuildHandlerOmitFieldsOptions 验证 Update 默认忽略主键与关联字段
func TestBuildHandlerOmitFieldsOptions(t *testing.T) {
	out := buildHandlerOmitFieldsOptions(relationFields())
	// Update: 主键 + remoteSelect 关联(admin) + remoteSelects 展示字段(admin_username_list)
	if !strings.Contains(out, "Update: []string{\"id\", \"admin\", \"admin_username_list\"}") {
		t.Fatalf("Update 缺少关联忽略: %s", out)
	}

	// 端到端: 渲染的控制器应包含 WithOmitFields（仅 Update）
	content, err := crud.RenderTmpl(handlerTmpl, buildHandlerFileData(tableBasic(relationTable()), relationTable(), relationFields()))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"handler.WithOmitFields(handler.ActionFields{",
		`Update: []string{"id", "admin", "admin_username_list"},`,
	} {
		if !strings.Contains(content, want) {
			t.Errorf("渲染的控制器缺少 %q:\\n%s", want, content)
		}
	}
}
