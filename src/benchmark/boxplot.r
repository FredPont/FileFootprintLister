# Créer un dataframe avec les données fournies
data <- data.frame(
  algo = rep(c("md5", "cityhas64", "xxhash", "murmur", "clickhouse64", "sha256"), each = 3),
  value = c(240.19, 221.0, 214.0, 56.78, 55.46, 54.022, 109.38, 108.27, 131.179, 122.93, 114.43, 121.99, 57.24, 51.089, 58.22, 339.65, 341.02, 310.0)
)

# Charger la bibliothèque ggplot2 pour créer des graphiques
library(ggplot2)

# Créer un boxplot pour chaque catégorie "algo"
ggplot(data, aes(x = algo, y = value, fill = algo)) +
  geom_boxplot() +
  labs(title = "Algorithm execution time",
       x = "Algorithm",
       y = "Seconds") +
  theme_minimal()

