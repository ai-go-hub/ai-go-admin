package router

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	assets "github.com/ai-go-hub/ai-go-admin"
	"github.com/ai-go-hub/ai-go-admin/internal/infra/config"
	"github.com/ai-go-hub/ai-go-admin/internal/router/registry"

	"github.com/gin-gonic/gin"
)

func init() {
	registry.Register(func(r *gin.Engine) {
		// 静态资源合并托管: 嵌入资源优先，其余回退磁盘
		g := r.Group("/static")
		// 禁止浏览器对响应做 MIME 嗅探，避免上传的伪装文件被当可执行内容渲染
		g.Use(func(c *gin.Context) {
			c.Header("X-Content-Type-Options", "nosniff")
			c.Next()
		})

		// 禁止目录列表的静态文件服务器
		fileServer := http.StripPrefix("/static", http.FileServer(noListFS{http.FS(staticFS{})}))
		handler := func(c *gin.Context) {
			// 直接全部标记 404
			c.Status(http.StatusNotFound)

			// 存在的文件会被 http.FileServer 用 200 覆盖
			fileServer.ServeHTTP(c.Writer, c.Request)
		}
		g.GET("/*filepath", handler)
		g.HEAD("/*filepath", handler)

		// 自定义上传文件的存储路径后，或许不能再为其提供静态服务
		if prefix := config.Get().Upload.URLPrefix; prefix != "" && !strings.HasPrefix(prefix, "/static") {
			log.Printf("[warn] upload.url_prefix=%q 不在 /static 目录下，请自建服务并配置 cdn.url，否则该上传文件无法访问", prefix)
		}
	})
}

// staticFS 合并文件系统
//
// static/images 走嵌入资源
// 上传文件按 upload.url_prefix 配置映射
// 其余（如 static/dist）走磁盘（若也需嵌入写到 embed.go 文件即可）
type staticFS struct{}

func (staticFS) Open(name string) (fs.File, error) {
	// 1. 嵌入资源（static/images 等）
	if name != "." {
		if f, err := assets.FS.Open("static/" + name); err == nil {
			return f, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
	}

	// 2. 上传文件
	if rest, ok := uploadFilePath(name); ok {
		return os.DirFS(config.Get().Upload.Dir).Open(rest)
	}

	// 3. 磁盘兜底
	return os.DirFS("static").Open(name)
}

// uploadFilePath 判断 name 是否命中上传 URL 前缀，命中则返回去掉前缀后的相对路径
func uploadFilePath(name string) (string, bool) {
	prefix := config.Get().Upload.URLPrefix
	if !strings.HasPrefix(prefix, "/static") {
		return "", false
	}
	tail := strings.TrimPrefix(prefix, "/static/")
	if tail == "" || name == tail {
		return "", false
	}
	if rest, ok := strings.CutPrefix(name, tail+"/"); ok {
		return rest, true
	}
	return "", false
}

// noListFS 包装 http.FileSystem，
// 使目录的 Readdir 返回空，从而禁止 http.FileServer 输出目录列表
type noListFS struct {
	fs http.FileSystem
}

// Open 打开文件，目录返回禁止列表的文件
func (n noListFS) Open(name string) (http.File, error) {
	f, err := n.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return noListFile{f}, nil
}

// noListFile 包装 http.File，覆盖 Readdir 使目录列表为空
type noListFile struct {
	http.File
}

// Readdir 返回空列表，禁止目录列举
func (f noListFile) Readdir(int) ([]fs.FileInfo, error) {
	return nil, nil
}
