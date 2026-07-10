# ws-lang — Universal Capability Orchestration Language

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/banshanhanfu/ws-lang)](https://goreportcard.com/report/github.com/banshanhanfu/ws-lang)

**ws-lang** is a domain-specific orchestration language for composing capabilities (voice processing, image processing, search, notification, etc.) into executable instruction sets. It is designed for **LLM-native generation** — the language is optimized so that LLMs can output correct syntax with near-zero error rates.

### 🌍 Multi-Language by Design

ws-lang is **not** an English-only language. Keywords can be written in any human language, and a deterministic Translation Layer maps them to a canonical internal form.

**10 languages currently supported:**

| Language | Code | Example |
|:---------|:----:|:--------|
| 🇬🇧 English | `en` (default) | `task`, `step`, `cap`, `retry`, `merge` |
| 🇨🇳 Chinese | `zh` | `任务`, `步骤`, `能力`, `重试`, `合并` |
| 🇯🇵 Japanese | `ja` | `タスク`, `ステップ`, `呼出`, `リトライ`, `結合` |
| 🇪🇸 Spanish | `es` | `tarea`, `paso`, `capacidad`, `reintentar`, `combinar` |
| 🇫🇷 French | `fr` | `tâche`, `étape`, `capacité`, `réessayer`, `fusionner` |
| 🇵🇹 Portuguese | `pt` | `tarefa`, `etapa`, `capacidade`, `retentar`, `mesclar` |
| 🇷🇺 Russian | `ru` | `задача`, `шаг`, `возможность`, `повторить`, `объединить` |
| 🇮🇳 Hindi | `hi` | `कार्य`, `चरण`, `क्षमता`, `पुनःप्रयास`, `विलय` |
| 🇧🇩 Bengali | `bn` | `কাজ`, `ধাপ`, `ক্ষমতা`, `পুনরায়_চেষ্টা`, `একত্রিত` |
| 🇸🇦 Arabic | `ar` | `مهمة`, `خطوة`, `قدرة`, `إعادة_محاولة`, `دمج` |

```ws
# English
task "image-pipeline" {
    step download { cap: ws-storage }
}

# 中文
任务 "图像处理" {
    步骤 下载 { 能力: ws-storage }
}

# 日本語
タスク "画像処理" {
    ステップ ダウンロード { 呼出: ws-storage }
}

# Español
tarea "procesar-imágenes" {
    paso descargar { capacidad: ws-storage }
}

# Français
tâche "traitement-images" {
    étape télécharger { capacité: ws-storage }
}

# العربية
مهمة "معالجة-الصور" {
    خطوة تحميل { قدرة: ws-storage }
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
├── translations/    # Multi-language keyword mappings (*.json + generated translations.go)
├── examples/        # Example instruction sets (10 languages)
│   ├── meeting-notes.ws       # English
│   ├── meeting-notes.ws.zh    # Chinese
│   ├── spanish-test.ws        # Spanish
│   ├── arabic-test.ws         # Arabic
│   ├── french-test.ws         # French
│   ├── portuguese-test.ws     # Portuguese
│   ├── russian-test.ws        # Russian
│   ├── hindi-test.ws          # Hindi
│   ├── bengali-test.ws        # Bengali
│   └── japanese-test.ws       # Japanese
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

# With translation (Spanish keywords → canonical)
go run compiler/compiler.go -lang es examples/spanish-test.ws

# With translation (Arabic keywords → canonical)
go run compiler/compiler.go -lang ar examples/arabic-test.ws

# List all available languages
ls translations/*.json
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
