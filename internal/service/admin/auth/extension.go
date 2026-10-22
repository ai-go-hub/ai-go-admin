package auth

import (
	"errors"

	"github.com/ai-go-hub/ai-go-admin/internal/dto"
	"github.com/ai-go-hub/ai-go-admin/internal/service"
)

// AuthExtension 管理员操作扩展参数
type AuthExtension struct {
	AdminSession *dto.AdminSession
}

// sessionFrom 从扩展参数中取出管理员会话信息
func sessionFrom(opts service.Options) (*dto.AdminSession, error) {
	ext, ok := opts.Extension.(*AuthExtension)
	if !ok || ext.AdminSession == nil {
		return nil, errors.New("参数错误，缺少 AdminSession 扩展数据")
	}
	return ext.AdminSession, nil
}
