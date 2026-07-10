# KV Line Format — ws-lang 的核心用户接口

## 为什么需要 KV Line Format？

我们发现 LLM 生成 JSON/YAML 时经常出错——漏括号、多逗号、字符串不闭合。
KV Line Format 专为 LLM 设计，每一行是一个键值对，结构由缩进控制：

```
key: value
key: value
  subkey: value
```

**LLM 的错误率接近 0%**，因为：
1. 没有括号需要匹配
2. 没有逗号需要分隔
3. 字符串不需要引号包裹
4. 缩进自动确定层级

## 三层格式体系

```
┌──────────────────────────────────────────┐
│    KV Line Format (用户/LLM 编写)         │  ← 最易写，最容易生成
├──────────────────────────────────────────┤
│    Parser + Validator                     │  ← 自动转换
├──────────────────────────────────────────┤
│    YAML (内部规范格式)                     │  ← 带注释，可审计
├──────────────────────────────────────────┤
│    JSON Schema (最终校验)                  │  ← 标准交换格式
└──────────────────────────────────────────┘
```

## 语法规则

### 基本行

```
key: value
```

- `key` — 字母、数字、下划线，区分大小写
- `value` — 自由文本，不需要引号
- 冒号后面必须跟一个空格

### 多层结构

缩进控制层级，**建议 2 空格**：

```
step: 转录
  cap: ws-voice
  input: /recordings/file.mp3
  on_error: retry(3, 5s)
```

等价于 YAML：
```yaml
step: 转录
  cap: ws-voice
  input: /recordings/file.mp3
  on_error: retry(3, 5s)
```

### 注释

```
# 这是注释（整行）
step: 转录
  cap: ws-voice  # 行尾注释（v2 支持）
```

### 引用

```
-> output: transcript     # 定义输出端口
input: $transcript        # 引用上游输出
input: [$a, $b]           # 合并多个输入
```

## 示例对照

### KV 格式（用户写）

```ws
task: 会议处理
  step: 转录
    cap: ws-voice
    input: /recordings/meeting.mp3
    -> output: transcript
  step: 总结
    cap: ws-summarize
    input: $transcript
```

### 编译后 YAML（内部存储）

```yaml
task: 会议处理
  steps:
    - id: s-1
      name: 转录
      cap: ws-voice
      input:
        type: file
        path: /recordings/meeting.mp3
      outputs:
        - name: transcript
          from: stdout.text
    - id: s-2
      name: 总结
      cap: ws-summarize
      input:
        type: ref
        from: s-1
        output: transcript
      depends_on:
        - s-1
```

## 解析器实现要点

1. 按行读取，忽略注释和空行
2. 根据缩进确定层级
3. 每行分割为 key/value
4. 检测特殊标记（`$`引用、`->`输出、`[]`数组）
5. 构建为中间 AST
6. AST → YAML 序列化
