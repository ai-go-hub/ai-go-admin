# CLAUDE.md

本文件为 Claude Code 在当前代码库中工作时提供指导。

## 项目概述

爱购商城（ai-go-mall），技术栈：Go 1.25 + Gin + GORM + PostgreSQL。

## 回答偏好

- 回复使用中文
- 遇到有多种实现方案时，列出选项让我选择，而不是直接选一种
- 当我询问某方案是否合理时，请先根据社区惯例判断是否合理，合理直接实现，不合理取消实现并解释原因（社区惯例指对应技术栈的社区，如 golang 开源社区，Gin 开源社区，开源高星仓库，官方文档，权威 blog 等）
- 基类虽是 OOP 中的称呼，但在本项目中可以 `基类` 一词代表泛型驱动的三层通用实现（Handler → Service → Repository），若某一层单独出现时，可称为 `基控制器、基服务、基仓储`

## 关键约定

- **只使用 GET 和 POST 请求方式**：大多数 CDN/全站加速 服务对 PUT、DELETE 兼容性差
- **包名**：全小写、单数、无下划线（`handler`、`service`）
- **文件名**：全小写、下划线分隔（`user_service.go`）
- **导出符号**：大驼峰，私有符号小驼峰
- **数据库字段**：蛇形命名（`user_name`、`created_at`）
- **统一响应**：`response.Success(c, opts...)` / `response.Fail(c, opts...)`，支持 functional options 和链式调用两种风格，优先使用 functional options
- **GORM**：尽量使用 `GORM` 的 `Generics API（gorm.G[Model](db)....）`，而不是 `Traditional API`；在使用 `Generics API` 时，一般应直接使用全局 db 实例（`internal\infra\database\database.go` 中的 `DB` 可获取），调用操作方法时再传递合适的 ctx 即可。
- **避免 stutter**：包名已经表达的含义，结构体 / 函数不要再重复，如 `admin.AdminService` 改用 `admin.Service`。

## 四层应用架构和泛型基类

核心应用架构为：泛型驱动的四层架构模式，请求 → Handler（控制器）→ Service（业务逻辑）→ Repository（数据访问）→ Model（数据模型）

前三层都有**泛型基础实现**，具体模块通过**组合（嵌入）**复用：

| 层级       | 泛型基类        | 接口             | 文件                          |
| ---------- | --------------- | ---------------- | ----------------------------- |
| Repository | `Repository[T]` | `IRepository[T]` | `internal/repository/base.go` |
| Service    | `Service[T]`    | `IService[T]`    | `internal/service/base.go`    |
| Handler    | `Handler[T]`    | 无               | `internal/handler/base.go`    |

**扩展模式**（以 User 为例）：

- `UserRepository` 嵌入 `*Repository[model.User]`，可访问 `DB()` 编写自定义查询
- `UserService` 嵌入 `IService[model.User]` 并持有 `*UserRepository`，可覆盖业务逻辑
- `UserHandler` 嵌入 `*Handler[model.User]` 并持有 `*UserService`，注册路由可调用 `RegisterBaseRoutes` 后再追加自定义路由

**目录结构规划**：

随着项目功能扩张，各层按以下规则组织子目录：

| 层级       | 子目录规则                                                                                             |
| ---------- | ------------------------------------------------------------------------------------------------------ |
| Handler    | 按用户身份设立子目录：`admin/`（后台用户/后台端）、`user/`（前台会员/会员端）、`business/`（供应商）等 |
| Service    | 按用户身份设立子目录，与 Handler 对应                                                                  |
| Repository | 按用户身份设立子目录，与 Handler 对应                                                                  |
| Model      | **不设子目录**，按业务模块组织文件（如 `admin.go`、`user.go`），一个文件可包含多个相关模型             |

**设计理由**：

1. **Handler/Service/Repository 按用户身份分目录**：不同身份的用户（后台管理员、前台会员、供应商）有各自独立的业务逻辑和 API 接口，分目录隔离更清晰
2. **脱离用户身份的通用服务**：可放入 `common/` 子目录或不设子目录（如验证码、配置、地区数据等跨身份服务）
3. **Model 不按身份分目录**：
   - 一个模型文件可包含多个相关模型（如 `admin.go` 内同时定义 `Admin`、`AdminGroup`、`AdminRule`）
   - 不同身份用户可能共用同一模型（如验证码、配置、省份数据等）

**结构示例**：

```
internal/
├── handler/
│   └── admin/           # 后台管理员相关接口
├── service/
│   └── admin/           # 后台管理员相关业务逻辑
├── repository/
│   └── admin/           # 后台管理员相关数据访问
└── model/
    ├── admin.go         # Admin, AdminGroup, AdminRule 等模型
    ├── user.go          # User, UserProfile 等模型
    └── common.go        # 验证码、配置等通用模型
```

## 路由自动发现

`internal/router/registry/` 提供 `Routes` 切片 + `Register(fn)` 函数。

子模块在 `init()` 中调用 `registry.Register(func(r *gin.Engine) { ... })` 自注册路由分组，`internal/router/router.go` 通过空白导入触发 `init()`：

```
router.go ──_ import──→ admin/ ──init()──→ registry.Register(...)
           ├─_ import──→ user/  ──init()──→ registry.Register(...)
           └─Setup() ──→ 遍历 registry.Routes
```

新增路由模块：

1. 新建子目录并于 `init()` 调 `registry.Register`
2. 在 `router.go` 加一行空白导入以触发 `init()`

## 启动流程

1. `config.Init()` — 合并 `config/*.yaml` + `.env.yaml`，环境变量覆盖，热加载
2. `database.Init()` — 连接 PostgreSQL，配置读写分离，`AutoMigrate` 所有 `model.Register()` 的模型
3. `engine.Use(database.Middleware())` — 注入 `*gorm.DB` 到 `gin.Context`
4. `router.Setup(engine)` — 遍历 `registry.Routes` 注册所有路由
5. `engine.Run()`
