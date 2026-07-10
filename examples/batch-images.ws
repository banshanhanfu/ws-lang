# 批量图片处理

step: 下载
  cap: ws-storage
  input: s3://uploads/raw-images/20260711/
  -> output: images

step: 压缩
  cap: ws-image
  input: $images
  quality: 80
  format: webp
  max_width: 1920
  -> output: compressed

step: 缩略图
  cap: ws-image
  input: $images
  width: 200
  height: 200
  format: webp
  -> output: thumbnails

step: 水印
  cap: ws-image
  input: $compressed
  watermark: /assets/watermark.png
  position: bottom-right
  -> output: watermarked

step: 上传
  cap: ws-storage
  input: [$watermarked, $thumbnails]
  target: s3://uploads/processed/20260711/
  -> output: urls

step: 通知
  cap: ws-wechat-send
  input: $urls
  target: 设计组
  text: "图片处理完成，共 {count} 张"
