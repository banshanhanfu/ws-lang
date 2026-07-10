# ws-lang en español
tarea: procesar-imágenes
  paso: descargar
    capacidad: ws-storage
    entrada: s3://uploads/raw/
    -> salida: imágenes
  paso: comprimir
    capacidad: ws-image
    entrada: $imágenes
    argumentos: { formato: webp, calidad: 80 }
    en_error: reintentar(3, 5)
  paso: notificar
    capacidad: ws-wechat
    entrada: $comprimir
    argumentos: { destino: grupo-diseño }
