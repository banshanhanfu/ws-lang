# ws-lang em português
tarefa: processar-imagens
  passo: baixar
    capacidade: ws-storage
    entrada: s3://uploads/raw/
    -> saída: imagens
  passo: comprimir
    capacidade: ws-image
    entrada: $imagens
    argumentos: { formato: webp, qualidade: 80 }
    ao_erro: repetir(3, 5)
  passo: notificar
    capacidade: ws-wechat
    entrada: $comprimir
    argumentos: { destino: grupo-design }
