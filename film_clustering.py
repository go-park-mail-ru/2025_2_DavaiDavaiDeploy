import pandas as pd
import numpy as np
from sklearn.cluster import AffinityPropagation
from sklearn.preprocessing import StandardScaler
import matplotlib.pyplot as plt
import seaborn as sns
from sqlalchemy import create_engine
import warnings
from sklearn.cluster import KMeans
from sqlalchemy import text


def save_clusters_to_db(films_df):
    try:
        engine = create_engine('postgresql://postgres:qwerty@localhost:5432/kinopoisk')

        with engine.connect() as conn:
            for _, film in films_df.iterrows():
                update_query = f"""
                UPDATE film 
                SET cluster_id = {film['cluster']} 
                WHERE title = '{film['title'].replace("'", "''")}' AND year = {film['year']}
                """
                conn.execute(text(update_query))

            conn.commit()

        print("Кластеры сохранены в БД")
    except Exception as e:
        print(f"Ошибка сохранения: {e}")


def loadFilmsData():
    try:
        engine = create_engine('postgresql://postgres:qwerty@localhost:5432/kinopoisk')
        query = "SELECT f.title, f.year, f.genre_id, f.age_category, f.country_id, f.duration, ROUND(COALESCE(AVG(ff.rating), 0), 1) as rating FROM film f LEFT JOIN film_feedback ff ON f.id = ff.film_id GROUP BY f.id, f.title, f.year, f.genre_id, f.age_category ORDER BY rating DESC"
        films = pd.read_sql(query, engine)

        print(f"Загружено {len(films)} фильмов")
        return films

    except Exception as e:
        print(f"Ошибка при загрузке данных: {e}")
        return None


def constrained_kmeans(features, n_clusters, min_size):
    kmeans = KMeans(n_clusters=n_clusters, random_state=42)
    labels = kmeans.fit_predict(features)

    for _ in range(100):  
        cluster_sizes = np.bincount(labels, minlength=n_clusters)

        if all(size >= min_size for size in cluster_sizes):
            break

        smallest_cluster = np.argmin(cluster_sizes)

        distances = kmeans.transform(features)

        closest_points = np.argsort(distances[:, smallest_cluster])
        for point_idx in closest_points:
            if cluster_sizes[smallest_cluster] >= min_size:
                break
            if labels[point_idx] != smallest_cluster:
                labels[point_idx] = smallest_cluster
                cluster_sizes = np.bincount(labels, minlength=n_clusters)

    return labels


films = loadFilmsData()

features = films[['rating', 'year', 'genre_id', 'age_category', 'country_id', 'duration']].copy()

genreMapping = {genre: i for i, genre in enumerate(features['genre_id'].unique())}
ageRatingMapping = {age: i for i, age in enumerate(features['age_category'].unique())}
countryMapping = {country: i for i, country in enumerate(features['country_id'].unique())}

featuresEncoded = features.copy()
featuresEncoded['genre_encoded'] = featuresEncoded['genre_id'].map(genreMapping)
featuresEncoded['age_category_encoded'] = featuresEncoded['age_category'].map(ageRatingMapping)
featuresEncoded['country_encoded'] = featuresEncoded['country_id'].map(countryMapping)

featuresFinal = featuresEncoded[
    ['rating', 'year', 'genre_encoded', 'age_category_encoded', 'duration', 'country_encoded']]

scaler = StandardScaler()
featuresScaled = scaler.fit_transform(featuresFinal)

af = KMeans(n_clusters=7, random_state=42)
clusterLabels = constrained_kmeans(featuresScaled, n_clusters=7, min_size=7)

films['cluster'] = clusterLabels

print(f"Найдено {len(np.unique(clusterLabels))} кластеров")
print("Распределение по кластерам:")
print(films['cluster'].value_counts().sort_index())

plt.figure(figsize=(12, 8))

from sklearn.decomposition import PCA

pca = PCA(n_components=2)
features2D = pca.fit_transform(featuresScaled)

plt.subplot(1, 2, 1)
scatter = plt.scatter(features2D[:, 0], features2D[:, 1], c=clusterLabels, cmap='viridis', alpha=0.7)
plt.colorbar(scatter)
plt.title('Кластеризация фильмов (2D проекция)')
plt.xlabel('PCA Component 1')
plt.ylabel('PCA Component 2')

plt.subplot(1, 2, 2)
clusterCounts = films['cluster'].value_counts().sort_index()
plt.bar(clusterCounts.index, clusterCounts.values)
plt.title('Количество фильмов в каждом кластере')
plt.xlabel('Номер кластера')
plt.ylabel('Количество фильмов')

plt.tight_layout()
plt.show()

cluster_1_films = films[films['cluster'] == 0]
cluster_3_films = films[films['cluster'] == 3]

print("0 кластер")
for index, film in cluster_1_films.iterrows():
    print(f"{film['title']} | {film['rating']} | {film['year']} | {film['genre_id']} | {film['age_category']}")
print("")
print("6 кластер")
for index, film in cluster_3_films.iterrows():
    print(f"{film['title']} | {film['rating']} | {film['year']} | {film['genre_id']} | {film['age_category']}")

save_clusters_to_db(films)