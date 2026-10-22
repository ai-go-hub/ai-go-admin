package auth

import (
	"context"

	"github.com/ai-go-hub/ai-go-admin/internal/infra/permission"
	"github.com/ai-go-hub/ai-go-admin/internal/model"
	repoAuth "github.com/ai-go-hub/ai-go-admin/internal/repository/admin/auth"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
)

// AdminLogService 管理员日志服务
type AdminLogService struct {
	service.IService[model.AdminLog]
	repo *repoAuth.AdminLogRepository
}

// NewAdminLogService 创建管理员日志服务实例
func NewAdminLogService(repo *repoAuth.AdminLogRepository) *AdminLogService {
	return &AdminLogService{
		IService: service.NewService(repo),
		repo:     repo,
	}
}

// List 覆写通用查询方法: 非超管仅可查看自己的日志
func (s *AdminLogService) List(ctx context.Context, opts service.Options) ([]model.AdminLog, error) {
	session, err := sessionFrom(opts)
	if err != nil {
		return nil, err
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	// 非超管则只查询自己的日志
	if !super {
		opts.Wheres = append(opts.Wheres, service.WhereGroup{
			Wheres: []service.Where{{
				Field:    "admin_id",
				Operator: "eq",
				Value:    session.ID,
			}},
		})
	}

	return s.IService.List(ctx, opts)
}

// Count 覆写通用统计方法: 与 List 使用相同的权限过滤条件
func (s *AdminLogService) Count(ctx context.Context, opts service.Options) (int64, error) {
	session, err := sessionFrom(opts)
	if err != nil {
		return 0, err
	}

	perm := permission.New()
	super, err := perm.IsSuperAdmin(ctx, session.ID)
	if err != nil {
		return 0, err
	}

	if !super {
		opts.Wheres = append(opts.Wheres, service.WhereGroup{
			Wheres: []service.Where{{
				Field:    "admin_id",
				Operator: "eq",
				Value:    session.ID,
			}},
		})
	}

	return s.IService.Count(ctx, opts)
}
