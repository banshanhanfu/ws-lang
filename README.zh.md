# ws-lang — 通用能力编排语言

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/banshanhanfu/ws-lang)](https://goreportcard.com/report/github.com/banshanhanfu/ws-lang)

**ws-lang** 是一门领域特定编排语言，用于将各种能力（语音处理、图像处理、搜索、通知等）组合成可执行的指令集。
它专门为 **LLM 原生生成** 而设计——语言本身经过优化，让 LLM 可以近乎零错误地输出正确语法。

### 🌍 天然多语言

ws-lang **不是**英语专属语言。关键词可以用任何人类语言编写，通过确定性的翻译层（Translation Layer）映射为规范的内部形式：

```ws
# 中文
任务 "会议处理" {
    步骤 转录 { 能力: ws-voice }
}

# 日本語
タスク "会議処理" {
    ステップ 転写 { 能力: ws-voice }
}

# English
task "meeting" {
    step transcribe { cap: ws-voice }
}
```

### 🧩 系统架构

```
用户意图
  → 前置程序（LLM）
      → KV 行格式（LLM 原生）
          → 解析器 + 校验器
              → YAML 指令集（规范格式）
                  → ws-tasks 引擎（执行）
```

### ✨ 核心特性

- **KV 行格式** — 专为 LLM 优化的语法：没有括号、没有逗号、没有需要配对的引号
- **确定性翻译** — 通过查找表实现多语言支持，不依赖 AI（零歧义）
- **DAG 原生** — 步骤自动构成依赖关系图，并行分支并发执行
- **能力抽象** — 每个能力（`ws-voice`、`ws-image` 等）是一个自描述二进制，附带 `.cap.yaml` 清单
- **自文档化** — YAML 内部格式支持注释，便于执行审计

### 📦 仓库结构

```
ws-lang/
├── parser/          # KV 行格式解析器（Go）
├── compiler/        # KV → YAML 编译器
├── schema/          # 指令集 JSON Schema
├── translations/    # 多语言关键词映射表
├── examples/        # 示例指令集
└── docs/            # 文档
```

### 🚀 快速开始

```bash
# 解析 KV 格式指令集
go run parser/parser.go examples/meeting-notes.ws

# 编译为标准 YAML
go run compiler/compiler.go examples/meeting-notes.ws

# 使用翻译（中文关键词 → 规范格式）
go run compiler/compiler.go -lang zh examples/meeting-notes.ws
```

### 📖 了解更多

- [完整规范](SPEC.md)
- [KV 行格式](docs/kv-format.zh.md)
- [翻译层](docs/translation-layer.md)
- [系统架构](docs/architecture.md)

### 🤝 参与贡献

欢迎贡献！添加新的语言翻译、改进解析器、或提交能力清单。

### 📄 许可证

MIT
