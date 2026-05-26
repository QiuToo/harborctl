# harborctl

[English](README.md) | [中文](README_ZH.md)

<p align="center">
  <img src="https://img.shields.io/badge/version-v1.1.1-blue?style=flat-square">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go">
  <img src="https://img.shields.io/badge/Harbor-v2.0+-4CAF50?style=flat-square">
  <img src="https://img.shields.io/license/MIT-yellow?style=flat-square">
</p>

> A command-line tool for managing Harbor Registry, written in Go with Cobra.

## ✨ Features

- **Project Management** - Create, list, delete, and inspect Harbor projects
- **Image/Repository Management** - Browse, delete, and clean up container images
- **User Management** - Create, list, and manage Harbor users
- **Member Management** - Add and remove project members
- **Garbage Collection** - Trigger and schedule GC jobs
- **System Information** - View Harbor status, health, and statistics
- **Search** - Search projects and repositories

## 📦 Installation

### From Source

```bash
git clone https://github.com/QiuToo/harborctl.git
cd harborctl
go build -o harborctl .
sudo mv harborctl /usr/local/bin/
```

### Pre-built Binaries

Download from [Releases](https://github.com/QiuToo/harborctl/releases)

## 🚀 Quick Start

### Configuration File

Create config at `/etc/harbor/harbor.yaml`:

```yaml
address: "192.168.2.222:80"
username: "admin"
password: "Harbor123456"
# scheme: "http"
# insecure: false
```

### Or Use Command Line Arguments

```bash
harborctl -addr http://192.168.2.222:80 -u admin -p Harbor123456 info
```

## 📖 Commands

### Project Management

```bash
harborctl project list                      # List all projects
harborctl project create my-project        # Create a private project
harborctl project create my-project --public  # Create a public project
harborctl project delete my-project     # Delete a project
harborctl project inspect my-project   # Show project details
```

### Image Management

```bash
harborctl image list                       # List all images
harborctl image list my-project          # List images in project
harborctl image tags my-project/repo     # List image tags
harborctl image delete my-project/repo:tag  # Delete image tag
harborctl image clean my-project         # Clean up old tags
```

### User Management

```bash
harborctl user list                       # List all users
harborctl user create john --email john@example.com --password 'Pass123!'
harborctl user delete john                # Delete a user
```

### Garbage Collection

```bash
harborctl gc list                        # List GC history
harborctl gc run                         # Trigger GC
harborctl gc run --dry-run              # Dry run mode
harborctl gc schedule                    # Show GC schedule
harborctl gc schedule update --cron "0 2 * *"  # Update schedule
```

### System Info

```bash
harborctl info               # System information
harborctl info --details    # With health and statistics
harborctl health            # Component health status
harborctl stat              # System statistics
harborctl search nginx      # Search projects/repos
```

### Member Management

```bash
harborctl member list my-project        # List members
harborctl member add my-project john   # Add member
```

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      harborctl                             │
├────────────────��────────────────────────────────────────────┤
│  CLI Layer (Cobra)                                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐          │
│  │ project │ │  image  │ │   user  │ │   gc    │ ...      │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘      │
├───────┼───────────┼───────────┼───────────┼──────────────┤
│  API Client                                          │
│  ┌───────────────────────────────────────────────┐  │
│  │  HTTP + Basic Auth                            │  │
│  │  /api/v2.0/*                               │  │
│  └───────────────────┬───────────────────────┘  │
├─────────────────────┼──────────────────────────────┤
│  Harbor Registry   │                                │
│  ┌────────────────▼───────────────────────┐    │
│  │  Projects | Repositories | Users | GC  │    │
│  └──────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## 📋 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `HARBOR_ADDR` | Harbor server address | - |
| `HARBOR_USER` | Username | admin |
| `HARBOR_PASS` | Password | - |
| `HARBOR_SCHEME` | HTTP or HTTPS | http |

## 🔧 Build

```bash
# Build for current platform
go build -o harborctl .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o harborctl-linux .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o harborctl-macos .
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- [Harbor](https://goharbor.io/) - Cloud Native Registry
- [Cobra](https://github.com/spf13/cobra) - CLI framework