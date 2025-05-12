# Créer un dataframe avec les données fournies
data <- data.frame(
  algo = rep(c("md5", "cityhas64", "xxhash", "murmur", "clickhouse64", "sha256"), each = 3),
  value = c(240.19, 221.0, 214.0, 56.78, 55.46, 54.022, 109.38, 108.27, 131.179, 122.93, 114.43, 121.99, 57.24, 51.089, 58.22, 339.65, 341.02, 310.0)
)

# Calculer les médianes des valeurs pour chaque algorithme
medians <- aggregate(value ~ algo, data = data, median)

# Trier les algorithmes par ordre croissant des médianes
medians <- medians[order(medians$value), ]

# Obtenir l'ordre des algorithmes triés
sorted_algos <- medians$algo

# Convertir les algorithmes en facteur avec les niveaux triés
data$algo <- factor(data$algo, levels = sorted_algos)

# Charger la bibliothèque ggplot2 pour créer des graphiques
library(ggplot2)

# Créer un boxplot pour chaque catégorie "algo" triée
p <- ggplot(data, aes(x = algo, y = value, fill = algo)) +
  geom_boxplot() +
  labs(title = "Algorithm execution time : 2017 files - 34.5 GB",
       subtitle = "2x Intel(R) Xeon(R) Gold 6140 CPU @ 2.30 GHz",
       x = "Algorithm",
       y = "Seconds") +
  theme_minimal()

# Sauvegarder le boxplot dans un fichier PNG
ggsave("boxplot_sorted.png", plot = p, device = "png", width = 8, height = 6, dpi = 300)

