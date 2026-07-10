# 每日数据报表

step: 拉取数据
  cap: ws-db-query
  sql: SELECT * FROM daily_metrics WHERE date = CURRENT_DATE
  -> output: raw_data

step: 汇总
  cap: ws-calc
  input: $raw_data
  operations: [sum, avg, max, min]
  group_by: category
  -> output: summary

step: 生成图表
  cap: ws-chart
  input: $summary
  chart_type: bar
  title: "日报 - {date}"
  -> output: chart_image

step: 发送邮件
  cap: ws-email
  input: [$summary, $chart_image]
  to: team@company.com
  subject: "日报 - {date}"
  template: daily_report
