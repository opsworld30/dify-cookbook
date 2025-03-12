# dify-cookbook
A collection of practical examples and implementations using Dify API, including CLI tools, web applications, and integration demos.

## Introduction
This is a collection of examples showcasing how to use the Dify API in different scenarios. Currently includes:

- Command-line chat client (Go implementation)
- More examples coming soon...

## Go Chat Client

A command-line Dify API chat client with history recording and typewriter effect.

### Features
- Streaming output
- Typewriter effect display
- Local history storage
- Thinking process toggle
- Environment variable configuration

### Usage

1. Set environment variables:
```bash
export DIFY_API_KEY="your-api-key"
export DIFY_ENDPOINT="https://api.dify.ai/v1/chat-messages"
export USER_ID="your-user-id"
```
2. Run program:
```bash
go run main.go "你的问题"
```

### Configuration
- DIFY_API_KEY : Dify API key
- DIFY_ENDPOINT : Dify API endpoint
- USER_ID : User ID
- ShowThink : Whether to show thinking process (default: true)
- TypewriterDelay : Typewriter effect delay (default: 50ms)

## Contribution Guide
Welcome to contribute pull request. Make sure:

1. Code style is consistent
2. Add necessary comments and documentation
3. Update README.md with relevant information