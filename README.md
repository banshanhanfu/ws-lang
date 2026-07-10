# ws-lang — Universal Capability Orchestration Language

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/banshanhanfu/ws-lang)](https://goreportcard.com/report/github.com/banshanhanfu/ws-lang)

**ws-lang** is a domain-specific orchestration language for composing capabilities (voice processing, image processing, search, notification, etc.) into executable instruction sets. It is designed for **LLM-native generation** — the language is optimized so that LLMs can output correct syntax with near-zero error rates.

### 🌍 Multi-Language by Design

ws-lang is **not** an English-only language. Keywords can be written in any human language, and a deterministic Translation Layer maps them to a canonical internal form:

```ws
# Chinese
任务 "会议处理" {
    步骤 转录 { 能力: ws-voice }
}

# Japanese
タスク "会議処理" {
    ステップ 転写 { 能力: ws-voice }
}

# English
task "meeting" {
    step transcribe { cap: ws-voice }
}
```

### 🧩 Architecture

```
User Intent
  → Frontend Program (LLM)
      → KV Line Format (LLM-native)
          → Parser + Validator
              → YAML Instruction Set (canonical)
                  → ws-tasks Engine (execution)
```

### ✨ Key Features

- **KV Line Format** — LLM-optimized syntax: no brackets, no commas, no quotes to mismatch
- **Deterministic Translation** — Multi-language support via lookup tables, not AI (zero ambiguity)
- **DAG-native** — Steps automatically form dependency graphs; parallel branches execute concurrently
- **Capability Abstraction** — Each capability (`ws-voice`, `ws-image`, etc.) is a self-describing binary with a `.cap.yaml` manifest
- **Self-documenting** — YAML internal format supports comments for execution audit trails

### 📦 Repository Structure

```
ws-lang/
├── parser/          # KV line format parser (Go)
├── compiler/        # KV → YAML compiler
├── schema/          # JSON Schema for instruction sets
├── translations/    # Multi-language keyword mappings
├── examples/        # Example instruction sets
└── docs/            # Documentation
```

### 🚀 Quick Start

```bash
# Parse a KV format instruction set
go run parser/parser.go examples/meeting-notes.ws

# Compile to canonical YAML
go run compiler/compiler.go examples/meeting-notes.ws

# With translation (Chinese keywords → canonical)
go run compiler/compiler.go -lang zh examples/meeting-notes.ws.zh
```

### 📖 Learn More

- [Specification](SPEC.md)
- [KV Line Format](docs/kv-format.zh.md)
- [Translation Layer](docs/translation-layer.md)
- [Architecture](docs/architecture.md)

### 🤝 Contributing

Contributions are welcome! Add a new language translation, improve the parser, or suggest capability manifests.

### 📄 License

MIT
