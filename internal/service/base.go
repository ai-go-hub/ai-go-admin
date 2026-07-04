package service

import (
	"github.com/ai-go-hub/ai-go-admin/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Filter 筛选条件
type Filter struct {
	Field    string
	Operator string
	Value    any
}

// IService 通用服务接口
type IService[T any] interface {
	Create(c *gin.Context, entity *T) error
	Get(c *gin.Context, id uint) (*T, error)
	List(c *gin.Context, filters ...Filter) ([]T, error)
	Update(c *gin.Context, id uint, entity T) error
	Delete(c *gin.Context, id uint) error
}

// Service 通用服务实现
type Service[T any] struct {
	repo repository.IRepository[T]
}

// NewService 创建通用服务实例
func NewService[T any](repo repository.IRepository[T]) IService[T] {
	return &Service[T]{repo: repo}
}

// Create 创建记录
func (s *Service[T]) Create(c *gin.Context, entity *T) error {
	return s.repo.Create(c, entity)
}

// Get 根据主键 ID 查询单条记录
func (s *Service[T]) Get(c *gin.Context, id uint) (*T, error) {
	return s.repo.Get(c, id)
}

// List 查询全部记录
func (s *Service[T]) List(c *gin.Context, filters ...Filter) ([]T, error) {
	return s.repo.List(c, buildScopes(filters)...)
}

// Update 根据主键 ID 部分更新记录，子类可覆写添加业务校验
func (s *Service[T]) Update(c *gin.Context, id uint, entity T) error {
	return s.repo.Update(c, id, entity)
}

// Delete 根据主键 ID 删除记录，子类可覆写添加业务校验
func (s *Service[T]) Delete(c *gin.Context, id uint) error {
	return s.repo.Delete(c, id)
}

// buildScopes 将 Filter 数组转换为 GORM Scopes
func buildScopes(filters []Filter) []func(*gorm.Statement) {
	if len(filters) == 0 {
		return nil
	}
	scopes := make([]func(*gorm.Statement), 0, len(filters))
	for _, f := range filters {
		op := getOperatorByAlias(f.Operator)
		switch f.Operator {
		case "IS NULL", "IS NOT NULL":
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(f.Field + " IS NULL")
			})
		case "BETWEEN":
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(f.Field+" BETWEEN ? AND ?", f.Value)
			})
		default:
			scopes = append(scopes, func(stmt *gorm.Statement) {
				stmt.Where(f.Field+" "+op+" ?", f.Value)
			})
		}
	}
	return scopes
}

// getOperatorByAlias 符号类运算符别名 → SQL 运算符，单词类直接透传
func getOperatorByAlias(op string) string {
	switch op {
	case "eq", "":
		return "="
	case "ne":
		return "!="
	case "gt":
		return ">"
	case "gte":
		return ">="
	case "lt":
		return "<"
	case "lte":
		return "<="
	default:
		return op // LIKE、IN、NOT IN 等单词运算符直接透传
	}
}
