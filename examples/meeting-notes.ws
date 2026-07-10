# 会议录音处理

# 步骤一：语音转文字
step: 转录
  cap: ws-voice
  input: /recordings/meeting-0710.mp3
  model: whisper-large-v3
  language: zh
  diarize: true
  on_error: retry(3, 5s)
  -> output: transcript

# 步骤二：提取待办事项（与总结并行执行）
step: 提取待办
  cap: ws-extract
  input: $transcript
  type: todos
  -> output: todos_result

# 步骤三：生成会议总结（与提取待办并行）
step: 总结
  cap: ws-summarize
  input: $transcript
  max_length: 500
  -> output: summary

# 步骤四：合并结果并通知
step: 通知
  cap: ws-wechat-send
  input: [$todos_result, $summary]
  target: 项目群
  format: report
  on_error: ignore
