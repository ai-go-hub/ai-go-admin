// Package assets 内嵌运行时资源（验证码素材、静态图片、config）
package assets

import "embed"

//go:embed all:config all:asset all:static/images
var FS embed.FS
