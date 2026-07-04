# AGENTS.md

本文件为 AGENTS 在当前代码库中工作时提供指导。

## 技术栈

后端：Go1.25 + Gin + GORM + PostgreSQL
前端：Vue3 + Element Plus + TypeScript + Vite + Pinia + Axios

## 回答偏好

- 回复使用中文
- 遇到有多种实现方案时，列出选项让我选择，而不是直接选一种

## 关键约定

- **包名**：全小写、单数、无下划线（`handler`、`service`）
- **文件名**：全小写、下划线分隔（`user_service.go`）
- **只使用 GET 和 POST 请求方式**：大多数 CDN/全站加速 服务对 PUT、DELETE 兼容性差
- **GORM**：使用 `GORM` 的 `Generics API（gorm.G[Model](db)....）`，而不是 `Traditional API`，且在使用 `Generics API` 时，一般应直接使用全局 db 实例（`internal\infra\database\database.go` 中的 `DB` 可获取），调用操作方法时再传递合适的 ctx 即可
- **避免 stutter**：包名已经表达的含义，结构体 / 函数不要再重复，如 `admin.AdminService` 改用 `admin.Service`

## 后端

### 分层结构

核心应用架构为：泛型驱动的四层架构模式（`internal/`），Handler（控制器）→ Service（业务逻辑）→ Repository（数据访问）→ Model（数据模型）。前三层均有泛型基类（对应目录的 `base.go` 文件），新增子模块时先嵌入再追加自定义方法。

### 配置自动加载

`internal/infra/config/config.go` 使用 viper，按以下顺序合并（后者覆盖前者）：

1. `config/*.yaml`（glob 全部 yaml，按文件名顺序 MergeInConfig）
2. 根目录 `.env.yaml`（若存在，覆盖上述配置；**已 gitignore，不入库**）
3. 环境变量（`AutomaticEnv`）

支持配置热重载（`OnConfigChange` + `WatchConfig`），通过 `config.Get()` 取带读写锁的副本。`.env.yaml.example` 是环境配置模板；`config/config.yaml` 是默认配置。新增配置项需同步在 `Config` 结构体加字段。

### 路由自动注册

`internal/router/registry/registry.go` 维护 `Routes []func(*gin.Engine)`。每个业务路由包在 `init()` 中调用 `registry.Register(...)` 注册自己；`router/index.go` 通过**空白导入** `_ ".../router/admin"`、`_ ".../router/common"` 触发 init，`router.Setup(engine)` 遍历执行。新增路由模块后，记得在 `router/index.go` 加空白导入。

### 模型自动迁移

`internal/model/model.go` 维护 `registered []any`。各模型文件在 `init()` 中调用 `Register(&Admin{})`；`database.Init()` 调用 `db.AutoMigrate(model.All()...)`。新增模型只需在 `model` 包内写 `init() { Register(&X{}) }`，并创建 `TableName()` 函数。

### 身份令牌（Token）

`internal/infra/token/token.go` 定义 `Driver` 接口与全局单例 `Manager()`（`sync.Once` 懒初始化，依据 `token.driver` 配置选驱动，当前仅 `database` 驱动）。Token 入库前做 **SHA256**，校验/删除时同样哈希后比对。认证中间件 `middleware/admin_auth.go`：`AdminAuth`（管理员强制登录认证）、`AdminAuthOptional`（管理员可选登录认证），从 `Authorization: Bearer <token>` 提取，校验后注入 `AdminSession` 到 context。

### 点选验证码（Click Captcha）

`internal/infra/captcha/click.go` 是核心特色模块：在背景图上绘制中文/大写字母/图标元素，带碰撞检测与混淆元素；用户按顺序点击，服务端比对 元素列表 + 坐标容差。`captcha.yaml` 控制元素类型、数量、过期、字体路径（`asset/font/SourceHanSansCN-Normal.ttf`）、背景图目录（`asset/captcha/click/background`）、图标目录（`asset/captcha/click/icon`）及图标中英文名映射。`bootstrap()` 首次调用加载中文字符池与图标元数据，每次 `Create` 顺带清理过期记录。`Check(req, deleteOnSuccess)` 的第二参控制验证成功是否删记录——预验场景传 `false`，登录场景传 `true`。

### 统一响应

`internal/response/response.go`：响应体 `{code, message, time, data}`，`code=0` 成功、`code=1` 失败（均 HTTP 200）。两种用法：函数式选项 `response.Success(c, response.WithData(x))` / 链式 `response.New(c).Code(...).Message(...).Send()`，日常使用函数式选项用法。

## 前端（web/）

### 路径别名与入口

- `/@/` → `src/`（tsconfig `paths` 与 vite `resolve.alias` 一致）。
- 入口 `src/main.ts`：注册 pinia（带 `pinia-plugin-persistedstate`）、router、element-plus、i18n、全局 icon。
- 路由用 **hash 模式**（`createWebHashHistory`）。

### 路由自动加载

`src/router/static.ts` 用 `import.meta.glob('./static/**/*.ts', { eager: true })` 自动收集 `src/router/static/` 下所有 `.ts` 默认导出（`RouteRecordRaw` 或 `RouteRecordRaw[]`）合并进 `staticRoutes`。`src/router/static/adminBase.ts` 定义后台基础路由 `/admin`。新增静态路由只需在该目录加文件。

### 状态管理（Pinia）

`src/stores/`：`config`（站点/语言/布局，持久化）、`adminInfo`（管理员信息含 token，持久化）、`menu`、`navTab`、`ref`。持久化 key 集中在 `src/stores/constant/cacheKey.ts`。

### 请求封装（`src/utils/request.ts`）

axios 实例 baseURL 取 `VITE_AXIOS_BASE_URL`。拦截器特性：

- **大小写转换**：请求 camelCase → snake_case，响应 snake_case → camelCase（`convertCase !== false` 时生效，FormData/URLSearchParams/Blob 不转）。
- **重复请求取消**：默认开启（`cancelDuplicate !== false`），按 method+url+params+data 生成 key，新请求 abort 旧请求。
- **loading**：`loading=true` 显示全屏 loading（计数式，支持并发）。
- **Bearer token**：自动从 `useAdminInfo().token` 注入 `Authorization`。
- 新增 API 请求函数放 `src/api/`，参考 `api/admin/index.ts`、`api/common.ts`。
