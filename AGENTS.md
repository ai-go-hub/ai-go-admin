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
- **源代码文件**：统一采用 UTF-8（无 BOM）编码；换行符统一使用 LF（\n）；文件末尾保留一个空行，行尾禁止空白空格

## 后端

### 分层结构

核心应用架构为：泛型驱动的四层架构模式（`internal/`），`handler`（控制器）→ `service`（业务逻辑）→ `repository`（数据访问）→ `model`（数据模型）。前三层均有泛型基类（对应目录的 `base.go` 文件），新增子模块时先嵌入再追加自定义方法。

### 配置自动加载

`internal/infra/config/config.go` 使用 viper，按以下顺序合并（后者覆盖前者）：

1. `config/*.yaml`（glob 全部 yaml，按文件名顺序 MergeInConfig）
2. 根目录 `.env.yaml`（若存在，覆盖上述配置；**已 gitignore，不入库**）
3. 环境变量（`AutomaticEnv`）

支持配置热重载（`OnConfigChange` + `WatchConfig`），通过 `config.Get()` 取带读写锁的副本。`.env.yaml.example` 是环境配置模板；`config/config.yaml` 是默认配置。新增配置项需同步在 `Config` 结构体加字段。

### 路由自动注册

`internal/router/registry/registry.go` 维护 `Routes` 和 `AdminRoutes`。每个业务路由包在 `init()` 中调用 `registry.Register(...)` 或 `registry.RegisterAdmin(...)` 注册自己；`router/index.go` 通过**空白导入** `_ ".../router/admin"`、`_ ".../router/common"` 触发 `init`，`router.Setup(engine)` 遍历执行。新增路由模块后，记得在 `router/index.go` 加空白导入。

### 数据库迁移

迁移文件位于 `cmd/migrate/migrations/`，遵循 golang-migrate 命名规范 `{版本号}_{名称}.up.sql` / `{版本号}_{名称}.down.sql`，经 `go:embed` 内嵌进二进制。

迁移按模型文件对应：`000001_admin` 对应 `internal/model/admin.go`，`000002_common` 对应 `internal/model/common.go`。种子数据保留单独迁移（`000003_seed`）。

通过 `aigo migrate up` 命令应用全部待执行迁移，`aigo migrate status` 查看状态，`aigo migrate create <name>` 创建新迁移文件。

新增模型需创建对应的 up/down 迁移 SQL 文件，并按版本号递增命名；SQL 内以 `__PREFIX__` 作为表前缀占位符。

### 身份令牌（Token）

`internal/infra/token/token.go` 定义 `Driver` 接口与全局单例 `Manager()`（`sync.Once` 懒初始化，依据 `token.driver` 配置选驱动，当前仅 `database` 驱动）。Token 入库前做 **SHA256**，校验/删除时同样哈希后比对。认证中间件 `middleware/admin_auth.go`：`AdminAuth`（管理员强制登录认证）、`AdminAuthOptional`（管理员可选登录认证），从 `Authorization: Bearer <token>` 提取，校验后注入 `AdminSession` 到 context。

### 中间件

除认证中间件 `AdminAuth、AdminAuthOptional` 以外，额外提供了 `AdminPermission`（管理员权限校验中间件，无需验权的路由则无需使用此中间件，使用此中间件必须先使用 `AdminAuth`中间件）、`AdminLog`（管理员日志记录中间件，已为 `/admin` 分组路由统一注册）

### 统一响应

`internal/kit/httpx/response.go`：响应体 `{code, message, time, data}`，`code=0` 成功、`code=1` 失败（均 HTTP 200）。两种用法：函数式选项 `httpx.Success(c, httpx.WithData(x))` / 链式 `httpx.New(c).Code(...).Message(...).Send()`，日常使用函数式选项用法。

### 其他主要目录

- pkg 公共库代码（与项目无关，不依赖配置、数据库等基础设施的公共代码，甚至可被其他项目引用）
- static 静态服务根目录
- internal/dto 数据传输对象定义
- internal/infra 基础设施目录，含：配置加载、配置实例暴露、数据库连接与实例暴露、token 工具箱、上传工具箱、验证码等
- internal/kit 成套工具箱，可以依赖基础设施，需要成套以避免成为 common，含：httpx、urlx

## 前端（web/）

### 路径别名与入口

- `@/` → `src/`（tsconfig `paths` 与 vite `resolve.alias` 一致）。
- 入口 `src/main.ts`：注册 pinia（带 `pinia-plugin-persistedstate`）、router、element-plus、i18n、全局 icon。
- 路由用 **hash 模式**（`createWebHashHistory`）。

### 路由自动加载注册

`src/router/static.ts` 用 `import.meta.glob('./static/**/*.ts', { eager: true })` 自动收集 `src/router/static/` 下所有 `.ts` 默认导出（`RouteRecordRaw` 或 `RouteRecordRaw[]`）合并进 `staticRoutes`。`src/router/static/adminBase.ts` 定义后台基础路由 `/admin`。新增静态路由只需在该目录加文件。

### 状态管理（Pinia）

`src/stores/`：`config`（站点/语言/布局，持久化）、`adminInfo`（管理员信息含 token，持久化）、`menu`、`navTab`、`ref`。持久化 key 集中在 `src/stores/constant/cacheKey.ts`。

### 请求封装（`src/utils/request.ts`）

axios 实例 baseURL 取 `VITE_AXIOS_BASE_URL`。拦截器特性：

- **大小写转换**：请求 camelCase → snake_case，响应 snake_case → camelCase（`opts.convertCase === true` 时生效，默认关闭，FormData/URLSearchParams/Blob 不转）。
- **重复请求取消**：默认开启（`cancelDuplicate !== false`），按 method+url+params+data 生成 key，新请求 abort 旧请求。
- **loading**：`loading=true` 显示全屏 loading（计数式，支持并发）。
- **Bearer token**：自动从 `useAdminInfo().token` 注入 `Authorization`。
- 新增 API 请求函数放 `src/api/`，参考 `api/admin/index.ts`、`api/common.ts`。

### 输入组件

除了 `element plus` 组件库提供的输入组件外，我们自行封装了以下输入组件：

|输入组件名|来源/位置|绑定值类型/示例|
|----|----|----|
|省份区域选择器|`src\components\agInput\components\areaSelect.vue`|`'1,2,3'`|
|图标选择器|`src\components\agInput\components\iconSelect.vue`|`'el-close'`|
|远程下拉|`src\components\agInput\components\remoteSelect.vue`|单选: `1`，多选: `[1, 2, 3]`|
|文件上传|`src\components\agInput\components\agUpload.vue`|单文件上传: `'a.png'`，多文件上传: `['a.png', 'b.png']`|
|数组|`src\components\agInput\components\array.vue`|`[{key: '', value: ''}]`|
|富文本编辑器|`src\components\agInput\components\editor.vue`|`'内容'`|

额外提供了一个统一入口组件 `src\components\agInput\index.vue`，只传递 `type` 即可渲染对应输入组件（包括 25 总输入组件支持），虽然写法简洁，但可维护性有限（属性传递、插槽方面），只应在 `动态表单` 场景下使用；在使用我们自行封装的输入组件时，应总是导入对应输入组件后使用，如：

```vue
<template>
    <el-form-item label="分组" prop="group">
        <RemoteSelect
            field="name"
            v-model="formItems.group"
            remote-url="/admin/auth/group/list"
        />
    </el-form-item>
</template>

<script setup lang="ts">
import RemoteSelect from '@/components/agInput/components/remoteSelect.vue'
</script>
```

### 表格

`src/components/table/` 基于 Element Plus `el-table` 封装，与 `src\hooks\useTableManager.ts` hook 配合使用。

- **`index.vue`**：表格主体，遍历 `manager.table.column` 渲染列，支持 `slot` 自定义渲染
- **`header/index.vue`**：表头工具栏，支持刷新、编辑、删除、公共搜索、快速搜索、列显示等按钮
- **`header/comSearch.vue`**：公共搜索组件，支持多字段组合条件搜索
- **`types\table.d.ts`**：表格类型定义（表格列可用属性、各种操作前后钩子、公共搜索可用操作符号、表格表单数据等）
- **`index.ts`**：工具函数（`getCellValue` 取值、`getDefaultOptButtons` 默认操作按钮、`findIndexRow` 树形查找）
- **`cellRenderer/`**：单元格渲染器，一个文件等于一个渲染器，使用列配置的 `render` 属性指定渲染器（`buttons,color,customRender,customTemplate,datetime,icon,image,images,switch,tag,tags,url,slot`）

`useTableManager` 列配置（`column` 数组）关键属性：

| 属性 | 说明 |
|------|------|
| `prop` | 字段名，支持 `.` 嵌套（如 `user.nickname`） |
| `label` | 列标题 |
| `render` | 渲染器类型（`image`、`tag`、`switch`、`slot` 等） |
| `operator` | 搜索操作符（`eq`、`LIKE`、`ILIKE`、`BETWEEN`、`IN` 等） |
| `formatter` | 自定义格式化函数 |
| `quickSearch` | 是否参与快速搜索 |
| `showOverflowTooltip` | 超出显示 `tooltip` |

### 多语言

`src/lang/index.ts` 基于 `vue-i18n` + `yaml`，实现懒加载。

- 使用 `import.meta.glob` 批量加载 `lang/{lang}/**/*.yaml`，按文件路径自动构建嵌套的 `messages` 结构
- 调用 `setLang(lang)` 切换语言时动态加载对应语言包
- 页面中使用 `vue-i18n` 的 `useI18n()` 或 `$t()` 获取翻译文本
- 新增语言文件：在 `lang/{lang}/` 下创建 `.yaml` 文件，自动被 `glob` 识别加载

### 样式

`src/styles/` 使用 SCSS 管理样式，入口 `index.scss`。

- **`var.scss`**：CSS 变量定义（`--ag*` 前缀）
- **`element.scss`**：Element Plus 组件样式覆盖
- **`dark.scss`**：暗黑模式样式
- **`loading.scss`**：全局加载样式
- **`app.scss`**：应用基础样式

总是优先使用 element plus 的 `--el` 开头和框架自定义的 `--ag` 开头的 CSS 变量。

### 工具

`src/utils/` 存放通用工具函数。

- **`request.ts`**：axios 封装
- **`common.ts`**：通用工具（`fullURL` 资源地址拼接、`auth` 鉴权等）
- **`router.ts`**：路由工具（`handleAdminRoute` 动态路由注册、`openMenu` 菜单打开、菜单处理、权限节点组装）
- **`layout.ts`**：布局工具（`mainHeight` 计算内容区高度、`setNavTabsWidth` 设置导航标签宽度）
- **`storage.ts`**：本地存储封装（`localStorage` 和 `sessionStorage`）
- **`validate.ts`**：表单验证规则
- **`mask.ts`**：遮罩层控制
- **`random.ts`**：随机数生成
- **`horizontalScroll.ts`**：横向滚动辅助

### 布局和页面

布局文件位于 `src/layouts/`，页面文件位于 `src/views/`。

**后台布局**（`layouts/admin/`）：

- **`index.vue`**：布局入口，根据配置选择容器模式
- **`container/`**：5 种布局容器（`classic`、`default`、`double`、`leftSplit`、`streamline`），通过 `config.layout.mode` 切换
- **`components/menu/`**：菜单组件（`menuVertical`、`menuHorizontal`、`menuLeftSplit` 等）
- **`components/navBar/`**：导航栏组件，每种布局模式对应一个导航栏
- **`components/aside/`**：侧边栏，含工具栏（`menuSearch` 菜单搜索）
- **`router-view/main.vue`**：路由视图容器，支持多标签页、keepalive

**页面结构**：

- 每个功能模块一个目录，包含 `index.vue`（主页面）+ `dialogForm.vue`（弹窗表单）
- 页面使用 `useTableManager` + `TableHeader` + `Table` + `DialogForm` 组合模式
- 表格的 API 请求通过 `api/table.ts` 的 `TableManagerAPI` 类封装，传入控制器 `URL` 即可自动生成增删改查接口，请求自定义 API 时，可于 `src\api\` 合适位置建立请求函数后导入使用