# dify-cookbook
Dify API 的实用示例集合，包含命令行工具、Web应用和集成演示等多种实现方式。

## 项目介绍
这是一个 Dify API 的示例集合项目，旨在展示如何在不同场景下使用 Dify API。目前包含：

- 命令行聊天客户端（Go 实现）
- 更多示例开发中...

## Go 聊天客户端

一个基于命令行的 Dify API 聊天客户端，支持历史记录和打字机效果。

### 特性
- 流式输出
- 打字机效果显示
- 本地历史记录
- 支持隐藏思考过程
- 环境变量配置

### 使用方法

1. 配置环境变量：
```bash
export DIFY_API_KEY="your-api-key"
export DIFY_ENDPOINT="https://api.dify.ai/v1/chat-messages"
export USER_ID="your-user-id"
```
2. 运行程序：
```bash
go run main.go "你的问题"
```

### 配置说明
- DIFY_API_KEY : Dify API 密钥
- DIFY_ENDPOINT : Dify API 端点
- USER_ID : 用户标识
- ShowThink : 是否显示思考过程（默认：true）
- TypewriterDelay : 打字机效果延迟（默认：50ms）

## 贡献指南
欢迎提交 Pull Request 来贡献代码。在提交之前，请确保：
1. 代码风格保持一致
2. 添加必要的注释和文档
3. 更新 README 中的相关说明