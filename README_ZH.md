# harborctl

[English](README.md) | [中文](README_ZH.md)

<p align="center">
  <img src="https://img.shields.io/badge/version-v1.1.1-blue?style=flat-square">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go">
  <img src="https://img.shields.io/badge/Harbor-v2.0+-4CAF50?style=flat-square">
  <img src="https://img.shields.io/license/MIT-yellow?style=flat-square">
</p>

> 使用 Go 和 Cobra 编写的 Harbor 镜像仓库命令行管理工具。

## ✨ 功能特性

- **项目管理** - 创建、查看、删除 Harbor 项目
- **镜像管理** - 浏览、删除和清理容器镜像
- **用户管理** - 创建、查看和管理 Harbor 用户
- **成员管理** - 添加和移除项目成员
- **垃圾回收** - 触发和调度 GC 任务
- **系统信息** - 查看状态、健康情况和统计
- **搜索** - 搜索项目和镜像仓库

## 📦 安装

### 从源码安装

```bash
git clone https://github.com/QiuToo/harborctl.git
cd harborctl
go build -o harborctl .
sudo mv harborctl /usr/local/bin/
```

### 下载预编译二进制

从 [Releases](https://github.com/QiuToo/harborctl/releases) 下载

## 🚀 快速开始

### 配置文件

创建配置文件 `/etc/harbor/harbor.yaml`:

```yaml
address: "192.168.2.222:80"
username: "admin"
password: "Harbor123456"
# scheme: "http"
# insecure: false
```

### 或使用命令行参数

```bash
harborctl -addr http://192.168.2.222:80 -u admin -p Harbor123456 info
```

## 📖 命令详解

### 项目管理

```bash
harborctl project list                      # 列出所有项目
harborctl project create my-project        # 创建私有项目
harborctl project create my-project --public  # 创建公有项目
harborctl project delete my-project     # 删除项目
harborctl project inspect my-project   # 查看项目详情
```

### 镜像管理

```bash
harborctl image list                       # 列出所有镜像
harborctl image list my-project          # 列出项目中的镜像
harborctl image tags my-project/repo     # 列出镜像标签
harborctl image delete my-project/repo:tag  # 删除镜像标签
harborctl image clean my-project         # 清理旧标签
```

### 用户管理

```bash
harborctl user list                       # 列出所有用户
harborctl user create john --email john@example.com --password 'Pass123!'
harborctl user delete john                # 删除用户
```

### 垃圾回收

```bash
harborctl gc list                        # 查看 GC 历史
harborctl gc run                         # 触发 GC
harborctl gc run --dry-run              # 试运行模式
harborctl gc schedule                    # 查看 GC 调度
harborctl gc schedule update --cron "0 2 * *"  # 更新调度
```

### 系统信息

```bash
harborctl info               # 系统信息
harborctl info --details    # 含健康和统计
harborctl health            # 组件健康状态
harborctl stat              # 系统统计
harborctl search nginx      # 搜索项目/镜像
```

### 成员管理

```bash
harborctl member list my-project        # 列出成员
harborctl member add my-project john   # 添加成员
```

## 🏗️ 架构图

```
┌─────────────────────────────────────────────────────────────┐
│                      harborctl                             │
├─────────────────────────────────────────────────────────────┤
│  CLI 层 (Cobra)                                          │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │ project │ │  image  │ │   user  │ │   gc    │ ...      │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘      │
├───────┼───────────┼───────────┼───────────┼──────────────┤
│  API 客户端                                        │
│  ┌───────────────────────────────────────────────┐  │
│  │  HTTP + Basic Auth                            │  │
│  │  /api/v2.0/*                               │  │
│  └───────────────────┬───────────────────────┘  │
├─────────────────────┼──────────────────────────────┤
│  Harbor 镜像仓库                                    │
│  ┌────────────────────────────────────────────┐    │
│  │  Projects | Repositories | Users | GC      │    │
│  └────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## 📋 环境变量

| 变量 | 说明 | 默认值 |
|----------|-------------|---------|
| `HARBOR_ADDR` | Harbor 服务器地址 | - |
| `HARBOR_USER` | 用户名 | admin |
| `HARBOR_PASS` | 密码 | - |
| `HARBOR_SCHEME` | HTTP 或 HTTPS | http |

## 🔧 编译

```bash
# 编译当前平台
go build -o harborctl .

# 编译 Linux 版本
GOOS=linux GOARCH=amd64 go build -o harborctl-linux .

# 编译 macOS 版本
GOOS=darwin GOARCH=amd64 go build -o harborctl-macos .
```

## 🤝 贡献

欢迎提交 Pull Request！

## 📄 许可证

MIT 许可证 - 查看 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- [Harbor](https://goharbor.io/) - 云原生镜像仓库
- [Cobra](https://github.com/spf13/cobra) - CLI 框架