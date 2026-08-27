package dto

// CRUDFields CRUD 字段数据
type CRUDFields struct {
	Title       string         `json:"title"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	DataType    string         `json:"dataType,omitempty"`
	Length      int            `json:"length"`
	Precision   int            `json:"precision"`
	Default     string         `json:"default,omitempty"`
	DefaultType string         `json:"defaultType"`
	Null        bool           `json:"null"`
	PrimaryKey  bool           `json:"primaryKey"`
	Unique      bool           `json:"unique"`
	Generated   string         `json:"generated,omitempty"`
	Comment     string         `json:"comment"`
	DesignType  string         `json:"designType"`
	Table       map[string]any `json:"table"`
	Form        map[string]any `json:"form"`
	UUID        string         `json:"uuid"`
}

// CRUDTable CRUD 表数据
type CRUDTable struct {
	Name                 string              `json:"name"`
	Comment              string              `json:"comment"`
	QuickSearchField     []string            `json:"quickSearchField"`
	DefaultSortField     string              `json:"defaultSortField"`
	FormFields           []string            `json:"formFields"`
	ColumnFields         []string            `json:"columnFields"`
	DefaultSortType      string              `json:"defaultSortType"`
	GenerateRelativePath string              `json:"generateRelativePath"`
	IsRepositoryModel    int                 `json:"isRepositoryModel"`
	ModelFile            string              `json:"modelFile"`
	HandlerFile          string              `json:"handlerFile"`
	RepositoryFile       string              `json:"repositoryFile"`
	RouterFile           string              `json:"routerFile"`
	ServiceFile          string              `json:"serviceFile"`
	WebViewsDir          string              `json:"webViewsDir"`
	RoutePath            string              `json:"routePath"`
	DesignChange         []TableDesignChange `json:"designChange"`
	Rebuild              string              `json:"rebuild"`
}

// GenerateFileBasicDataInfo 生成文件的基本信息
type GenerateFileBasicDataInfo struct {
	Type     string `json:"type"`               // 生成文件类型: model/handler/service/repository/router/views/lang
	Path     string `json:"path"`               // 相对路径
	App      string `json:"app"`                // handler/service/repository/router/views 的一级子目录，如 admin、common
	Dir      string `json:"dir"`                // 生成目录
	File     string `json:"file,omitempty"`     // 文件完整路径
	Package  string `json:"package,omitempty"`  // go 文件 package 值
	LastName string `json:"last_name"`          // 路径的最后一级，如 user_log 的 log
	Name     string `json:"name"`               // 大驼峰路径，如 user_log > UserLog
	CnFile   string `json:"cn_file,omitempty"`  // 语言包 zh-cn 文件路径
	EnFile   string `json:"en_file,omitempty"`  // 语言包 en 文件路径
	LangKey  string `json:"lang_key,omitempty"` // 语言翻译 key 前缀
}

// TableDesignChange CRUD 表设计变更记录
type TableDesignChange struct {
	Type    string `json:"type"`
	Index   int    `json:"index,omitempty"`
	Name    string `json:"name"`
	NewName string `json:"newName"`
	Sync    bool   `json:"sync,omitempty"`
}

// CrudLogData CRUD 日志数据
type CrudLogData struct {
	ID     *uint        // 日志 ID
	Table  CRUDTable    // 表数据
	Fields []CRUDFields // 字段数据
	Sql    *string      // 执行的 SQL
	Status string       // 状态: generating/succeeded/failed/deleted
}

// ChatMessage AI 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest AI 对话请求
type ChatRequest struct {
	Model       string        `json:"model" binding:"required"` // 模型名
	Messages    []ChatMessage `json:"messages"`                 // 对话上下文
	Temperature *float64      `json:"temperature"`              // 温度参数
	TopP        *float64      `json:"top_p"`                    // 核采样参数
}
