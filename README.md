<br />
<div align="center">
    <img src="https://thankphp.com/ai-go/logo.png" alt="LOGO" />
    <h1 style="font-size: 36px;color: #18202e;font-weight: 600;margin: 0 0 6px 0;">AI GO ADMIN</h1>
    <p style="font-size: 17px;color: #6a8bad;margin-bottom: 10px;">AI 执笔快速构建，人掌分寸雕琢商用级后台管理系统</p>
    <a href="https://qm.qq.com/q/GE0qcmaTo6" target="_blank">加 QQ 群</a> |
    <a href="https://juejin.cn/column/7654319240119369737" target="_blank">开发手记（掘金）</a> |
    <a href="https://blog.csdn.net/weixin_48453209/category_13179276.html" target="_blank">开发手记（CSDN）</a> |
    <a href="https://gitee.com/ai-go-hub/ai-go-admin" target="_blank">Gitee 仓库</a> |
    <a href="https://github.com/ai-go-hub/ai-go-admin" target="_blank">GitHub 仓库</a>
</div>
<br />
<p align="center">
    <a href="https://go.dev/" target="_blank">
        <img src="https://img.shields.io/badge/Golang-%3E 1.25.0-brightgreen?color=91aac3&labelColor=439EFD" alt="vue">
    </a>
    <a href="https://gorm.io/" target="_blank">
        <img src="https://img.shields.io/badge/GORM-%3E 1.31.1-blue?color=91aac3&labelColor=439EFD" alt="vite">
    </a>
    <a href="https://v3.vuejs.org/" target="_blank">
        <img src="https://img.shields.io/badge/Vue-%3E 3.5-brightgreen?color=91aac3&labelColor=439EFD" alt="vue">
    </a>
    <a href="https://cn.vitejs.dev/" target="_blank">
        <img src="https://img.shields.io/badge/Vite-%3E 8.0-blue?color=91aac3&labelColor=439EFD" alt="vite">
    </a>
    <a href="https://element-plus.org/zh-CN/guide/changelog.html" target="_blank">
        <img src="https://img.shields.io/badge/Element--Plus-%3E 2.13-brightgreen?color=91aac3&labelColor=439EFD" alt="element plus">
    </a>
    <a href="https://www.tslang.cn/" target="_blank">
        <img src="https://img.shields.io/badge/TypeScript-%3E 6.0-blue?color=91aac3&labelColor=439EFD" alt="typescript">
    </a>
    <a href="https://github.com/ai-go-hub/ai-go-admin/blob/master/LICENSE" target="_blank">
        <img src="https://img.shields.io/badge/Apache2.0-license-blue?color=91aac3&labelColor=439EFD" alt="license">
    </a>
</p>

<br>
<div align="center">
  <img src="https://thankphp.com/ai-go/dashboard-radius-admin.png" />
</div>
<br>

### 介绍

🌈 基于 Golang + GORM + TypeScript + Vite + Pinia + Element Plus 等流行技术栈的后台管理系统，目前开发工作以【AI 执笔 + 人工逐行雕琢】的模式积极进行中，作者指挥左手 AI，右手同步手搓 Blog，全程记录开发细节，Blog 开源：[Github](https://github.com/ai-go-hub/ai-go-blog) | [Gitee](https://gitee.com/ai-go-hub/ai-go-blog)。

### 安装使用

##### 一、下载项目

```bash
git clone https://gitee.com/ai-go-hub/ai-go-admin

# 或

git clone https://github.com/ai-go-hub/ai-go-admin
```

##### 二、配置数据库

您可以使用环境配置，即：将 `.env.yaml.example` 重命名为 `.env.yaml`，其中填写好数据库连接信息。

或者不使用环境配置，直接于 `config/config.yaml` 文件内配置数据库连接信息。

#### 三、启动前后端服务

**🌐 服务端**

```bash
# 安装依赖
go mod tidy

# 启动服务（带热更新）
air

# 或使用 go 原生 run 命令启动服务
go run cmd/api/main.go
```

**🧑‍💻 前端**

```bash
# 安装依赖
pnpm i

# 启动开发服务（带热更新）
pnpm dev
```

### 联系我们

- [加 QQ 群：860406643](https://qm.qq.com/q/GE0qcmaTo6)
- [Gitee 仓库](https://gitee.com/ai-go-hub/ai-go-admin)、[GitHub 仓库](https://github.com/ai-go-hub/ai-go-admin)
- [开发手记（掘金）](https://juejin.cn/column/7654319240119369737)
- [开发手记（CSDN）](https://blog.csdn.net/weixin_48453209/category_13179276.html)
- [官方邮箱 1094963513@qq.com](mailto:1094963513@qq.com)

### 项目预览

|                                                             |                                                                     |
| ----------------------------------------------------------- | ------------------------------------------------------------------- |
| ![登录](https://thankphp.com/ai-go/preview/login-admin.png) | ![点选验证码](https://thankphp.com/ai-go/preview/captcha-admin.png) |

### 特别鸣谢

💕 感谢巨人提供肩膀，排名不分先后

- [Golang](https://github.com/golang/go)
- [Gin](https://github.com/gin-gonic/gin)
- [Vue](https://github.com/vuejs/core)
- [Element Plus](https://github.com/element-plus/element-plus)
- [TypeScript](https://github.com/microsoft/TypeScript)
- [Vite](https://github.com/vitejs/vite)
- [Axios](https://github.com/axios/axios)
- [login-animation](https://gitee.com/niumg9527/login-animation)

### 版权信息

🔐 AI GO ADMIN 遵循 `Apache2.0` 开源协议发布，提供无需授权的免费商用。

- ✅ 允许个人和企业直接商用、修改、二次开发、外包、私有化部署。
- ✅ 允许自由复制、修改、增删源码，衍生作品可以闭源分发、售卖，且无需开源修改部分。
- ✅ 约束：请保留原项目 `LICENSE` 文件、若需开源衍生作品，请于 `README` 的显著位置添加本仓库链接。
- ❌ 禁止冒充自有开源项目对外分发售卖。
- ❌ 禁止擅自使用本项目商标、产品名、品牌名。
