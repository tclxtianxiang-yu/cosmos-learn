# Cosmos Learn

**Cosmos Learn** 是一个使用 Cosmos SDK 和 Tendermint 构建的教学示例区块链项目，通过 [Ignite CLI](https://ignite.com/cli) 创建。

## 项目特性

### 核心功能
- ✅ 基于 Cosmos SDK 的模块化区块链架构
- ✅ Tendermint 共识机制
- ✅ 自定义 Blog 模块，包含完整的 CRUD 操作
- ✅ IBC (Inter-Blockchain Communication) 跨链通信能力
- ✅ 跨链发送博客文章的功能

### Blog 模块功能

Blog 模块实现了以下功能：

1. **文章管理 (Post CRUD)**
   - 创建文章 (Create Post)
   - 读取文章 (Read Post)
   - 更新文章 (Update Post)
   - 删除文章 (Delete Post)
   - 列出所有文章 (List Posts)

2. **IBC 跨链通信**
   - 通过 IBC 协议向其他链发送文章
   - 接收来自其他链的文章
   - 跨链消息确认机制

## 快速开始

### 前置要求
- Go 1.22 或更高版本
- Ignite CLI v29.6.0 或更高版本

### 启动区块链

```bash
# 进入项目目录
cd cosmos-learn

# 启动开发环境
ignite chain serve
```

`serve` 命令会自动完成以下操作：
- 安装依赖
- 编译代码
- 初始化区块链
- 启动开发环境

## 使用示例

### 文章管理操作

#### 创建文章
```bash
# 使用 Alice 账户创建一篇文章
cosmos-learnd tx blog create-post "我的第一篇文章" "这是文章内容" --from alice --chain-id cosmoslearn
```

#### 查询所有文章
```bash
cosmos-learnd query blog list-post
```

#### 查询特定文章
```bash
cosmos-learnd query blog show-post [post-id]
```

#### 更新文章
```bash
cosmos-learnd tx blog update-post [post-id] "新标题" "新内容" --from alice --chain-id cosmoslearn
```

#### 删除文章
```bash
cosmos-learnd tx blog delete-post [post-id] --from alice --chain-id cosmoslearn
```

### IBC 跨链通信操作

#### 向其他链发送文章
```bash
# 通过 IBC 向其他链发送文章
cosmos-learnd tx blog send-post [channel-id] "跨链文章标题" "跨链文章内容" --from alice --chain-id cosmoslearn
```

注意：IBC 跨链通信需要先建立 IBC 连接和通道。

## 项目结构

```
cosmos-learn/
├── app/                    # 应用程序配置和初始化
├── cmd/                    # CLI 命令
├── proto/                  # Protobuf 定义
│   └── cosmoslearn/
│       └── blog/           # Blog 模块的 protobuf 定义
├── x/                      # 自定义模块
│   └── blog/               # Blog 模块实现
│       ├── keeper/         # 模块的状态管理
│       ├── types/          # 类型定义
│       ├── client/         # CLI 客户端
│       └── module/         # 模块接口实现
├── config.yml              # 开发环境配置
└── readme.md               # 项目文档
```

## 配置说明

区块链的开发环境可以通过 `config.yml` 文件进行配置。该文件包含：

- **账户配置**：预设的测试账户（alice, bob）及其初始代币
- **验证者配置**：区块链验证者节点的设置
- **Faucet 配置**：测试代币水龙头设置
- **默认代币**：链的默认代币单位

更多配置选项请参考 [Ignite CLI 文档](https://docs.ignite.com)。

### Web 前端

Ignite CLI 提供了前端脚手架功能（基于 Vue），可以快速构建区块链的 Web 前端：

```bash
ignite scaffold vue
```

此命令可以在区块链项目中运行。更多信息请参考 [Ignite 前端开发仓库](https://github.com/ignite/web)。

## 技术栈

- **Cosmos SDK**: 区块链应用框架
- **Tendermint Core**: 拜占庭容错共识引擎
- **IBC Protocol**: 跨链通信协议
- **Protobuf**: 数据序列化格式
- **gRPC**: RPC 框架

## 发布版本

要发布区块链的新版本，创建并推送带有 `v` 前缀的新标签：

```bash
git tag v0.1
git push origin v0.1
```

这将自动创建一个包含配置目标的草稿发布。

### 安装二进制文件

要安装区块链节点二进制文件的最新版本：

```bash
curl https://get.ignite.com/username/cosmos-learn@latest! | sudo bash
```

注意：`username/cosmos-learn` 应该与推送源代码的 GitHub 仓库的用户名和仓库名匹配。了解更多关于[安装过程](https://github.com/ignite/installer)。

## 学习资源

### Cosmos 生态系统
- [Cosmos SDK 官方文档](https://docs.cosmos.network)
- [Tendermint Core 文档](https://docs.tendermint.com)
- [IBC 协议规范](https://github.com/cosmos/ibc)

### Ignite CLI
- [Ignite CLI 官网](https://ignite.com/cli)
- [Ignite CLI 教程](https://docs.ignite.com/guide)
- [Ignite CLI 文档](https://docs.ignite.com)
- [开发者社区](https://discord.com/invite/ignitecli)

### 相关项目
- [Cosmos Hub](https://hub.cosmos.network)
- [Osmosis](https://osmosis.zone) - DEX 示例
- [Juno](https://junonetwork.io) - 智能合约平台

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

本项目是一个教学示例项目，用于学习 Cosmos SDK 和区块链开发。
