# ws-lang 规范 v1.0

## 1. 概述

ws-lang 是一门领域特定编排语言，用于描述由多个能力（capability）组成的执行流程。
语言设计核心原则：

1. **LLM 原生** — 语法让 LLM 可以近乎零错误生成
2. **多语言** — 关键词可用任意人类语言书写
3. **DAG 原生** — 步骤依赖自动形成有向无环图
4. **确定性** — 相同输入永远得到相同输出

## 2. KV 行格式

### 2.1 基本语法

```
key: value
```

- 键值对用冒号+空格分隔
- 值不需要引号包裹
- 缩进（2 空格）控制层级

### 2.2 块结构

```
task: 名称
  step: 步骤名
    cap: 能力名
    key: value
```

### 2.3 内联块

```
step: 步骤名 { cap: ws-voice, input: file.mp3 }
```

花括号内用逗号分隔多个键值对。

### 2.4 注释

```
# 整行注释
step: 测试  # 行尾注释
```

### 2.5 引用

```
input: $步骤名       # 引用上游输出
input: [$步骤1, $步骤2]  # 合并多个来源
```

### 2.6 输出端口

```
-> output: 端口名
```

## 3. 关键词表

### 3.1 结构关键词

| 英文 | 中文 | 日文 | 用途 |
|------|------|------|------|
| task | 任务 | タスク | 定义一个任务 |
| step | 步骤 | ステップ | 定义一个步骤 |
| cap | 能力 | 能力/呼出 | 指定能力 |
| input | 输入 | 入力 | 步骤输入 |
| output | 输出 | 出力 | 输出端口 |

### 3.2 控制关键词

| 英文 | 中文 | 日文 | 用途 |
|------|------|------|------|
| on_error | 出错时 | エラー時 | 错误处理 |
| retry | 重试 | リトライ | 重试策略 |
| ignore | 忽略 | 無視 | 忽略错误 |
| stop | 停止 | 停止 | 停止执行 |
| timeout | 超时 | タイムアウト | 超时设置 |
| merge | 合并 | 結合 | 合并多路输入 |
| from | 来自 | 取得元 | 引用来源 |

## 4. 内部 YAML 格式

用户编写的 KV 格式在解析后，编译为规范的 YAML 格式供 ws-tasks 执行。

### 4.1 顶层结构

```yaml
name: "任务名称"
status: "pending"
task_id: "t-20260711-001"
steps:
  - id: s-1
    name: "步骤名"
    cap: "ws-voice"
    input:
      type: file
      path: "/path/to/file.mp3"
    outputs:
      - name: "transcript"
        from: "stdout.text"
    on_error:
      action: retry
      max_retries: 3
      interval: 5
```

### 4.2 输入类型

| 类型 | 说明 | 示例 |
|------|------|------|
| literal | 直接值 | `{ type: literal, value: "hello" }` |
| file | 文件路径 | `{ type: file, path: "/data/file.mp3" }` |
| ref | 引用上游 | `{ type: ref, from: "s-1" }` |
| merge | 合并多路 | `{ type: merge, sources: [...] }` |

### 4.3 错误处理

| 动作 | 说明 |
|------|------|
| stop | 停止整个任务（默认） |
| retry | 重试当前步骤 |
| ignore | 忽略错误，继续执行 |

## 5. 能力系统

每个能力是一个自描述的二进制程序，附带 `.cap.yaml` 清单。

### 5.1 能力清单格式

```yaml
name: ws-voice
version: 1.0.0
description: "语音识别与合成"
binary: /usr/local/bin/ws-voice
args:
  model:
    type: string
    default: whisper-large-v3
inputs:
  - type: file
    description: "音频文件"
outputs:
  transcript:
    from: stdout.text
    description: "识别文本"
```

### 5.2 能力命名规范

- 格式：`ws-<功能名>`
- 示例：`ws-voice`、`ws-image`、`ws-db-query`

## 6. 执行模型

### 6.1 DAG 调度

1. 解析所有步骤，识别依赖关系
2. 无依赖的步骤从第 0 层开始执行
3. 每层内的步骤**并行**执行
4. 依赖满足后，下一层开始执行

### 6.2 数据传递

- 上游步骤的输出通过输出端口暴露
- 下游步骤通过 `$步骤名` 或 `$输出端口名` 引用
- 合并多路输入使用 `[$a, $b]` 语法

### 6.3 状态流转

```
pending → running → done
                  → fail → (retry → running | stop | ignore)
```

## 7. 完整示例

### 7.1 会议录音处理

```ws
task: 会议处理
  step: 转录
    cap: ws-voice
    input: /recordings/meeting.mp3
    model: whisper-large-v3
    language: zh
    -> output: transcript
  step: 提取待办
    cap: ws-extract
    input: $transcript
    -> output: todos
  step: 总结
    cap: ws-summarize
    input: $transcript
    max_length: 500
    -> output: summary
  step: 通知
    cap: ws-wechat-send
    input: [$todos, $summary]
    target: 项目群
    on_error: ignore
```

### 7.2 编译后 YAML

```yaml
---
name: 会议处理
status: pending
steps:
  - id: s-1
    name: 转录
    cap: ws-voice
    input:
      type: file
      path: /recordings/meeting.mp3
    args:
      model: whisper-large-v3
      language: zh
    outputs:
      - name: transcript
        from: stdout.text

  - id: s-2
    name: 提取待办
    cap: ws-extract
    input:
      type: ref
      from: s-1
      output: transcript
    outputs:
      - name: todos
        from: stdout.text

  - id: s-3
    name: 总结
    cap: ws-summarize
    input:
      type: ref
      from: s-1
      output: transcript
    args:
      max_length: "500"
    outputs:
      - name: summary
        from: stdout.text

  - id: s-4
    name: 通知
    cap: ws-wechat-send
    input:
      type: merge
      sources:
        - type: ref
          from: s-2
        - type: ref
          from: s-3
    args:
      target: 项目群
    on_error:
      action: ignore
```

## 8. 翻译层规范

### 8.1 翻译表

翻译表是一个 JSON 文件，将自然语言关键词映射到规范英文关键词：

```json
{
  "任务": "task",
  "步骤": "step",
  "能力": "cap"
}
```

### 8.2 约束

1. **单射** — 每个关键词最多映射到一个规范词
2. **确定性** — 给定输入永远得到相同输出
3. **可扩展** — 任何人都可以添加新语言翻译表

### 8.3 语言声明

文件开头或参数中声明语言：

```
@lang: zh
```

## 9. JSON Schema

见 `schema/instruction-set.schema.json`。

## 10. 版本历史

| 版本 | 日期 | 变更 |
|------|------|------|
| 1.0 | 2026-07-11 | 初始版本 |
