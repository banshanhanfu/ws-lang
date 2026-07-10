# ws-lang en français
tâche: traitement-images
  étape: télécharger
    capacité: ws-storage
    entrée: s3://uploads/raw/
    -> sortie: images
  étape: compresser
    capacité: ws-image
    entrée: $images
    arguments: { format: webp, qualité: 80 }
    en_erreur: réessayer(3, 5)
  étape: notifier
    capacité: ws-wechat
    entrée: $compresser
    arguments: { destination: groupe-design }
