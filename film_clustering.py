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
    kmeans = KMeans(n_clusters=n_clusters, random_state=42, n_init=10)
    labels = kmeans.fit_predict(features)
    distances = kmeans.transform(features)

    max_iterations = 100
    for iteration in range(max_iterations):
        cluster_sizes = np.bincount(labels, minlength=n_clusters)
        
        small_clusters = np.where(cluster_sizes < min_size)[0]
        large_clusters = np.where(cluster_sizes > min_size)[0]
        
        if len(small_clusters) == 0:
            break
        for small_cluster in small_clusters:
            needed = min_size - cluster_sizes[small_cluster]
            
            if needed <= 0:
                continue
            candidates = []
            
            for large_cluster in large_clusters:
                if cluster_sizes[large_cluster] <= min_size:
                    continue
                    
                points_in_large = np.where(labels == large_cluster)[0]
                for point_idx in points_in_large:
                    distance_diff = distances[point_idx, small_cluster] - distances[point_idx, large_cluster]
                    cluster_size_penalty = cluster_sizes[large_cluster] / max(cluster_sizes)
                    
                    score = distance_diff * cluster_size_penalty
                    candidates.append((score, point_idx, large_cluster, small_cluster))
            
            candidates.sort(key=lambda x: x[0])
            
            moved = 0
            for score, point_idx, large_cluster, small_cluster in candidates:
                if moved >= needed:
                    break
                labels[point_idx] = small_cluster
                cluster_sizes[small_cluster] += 1
                cluster_sizes[large_cluster] -= 1
                moved += 1
            
            large_clusters = np.where(cluster_sizes > min_size)[0]
            
            if len(large_clusters) == 0:
                break
    
    final_sizes = np.bincount(labels, minlength=n_clusters)
    small_clusters_final = np.where(final_sizes < min_size)[0]
    
    if len(small_clusters_final) > 0:
        
        for small_cluster in small_clusters_final:
            needed = min_size - final_sizes[small_cluster]
            if needed <= 0:
                continue
            
            largest_cluster = np.argmax(final_sizes)
            if final_sizes[largest_cluster] <= min_size:
                continue
            
            points_in_largest = np.where(labels == largest_cluster)[0]
            
            if len(points_in_largest) > 0:
                point_distances = []
                for point_idx in points_in_largest:
                    dist = distances[point_idx, small_cluster]
                    point_distances.append((dist, point_idx))
                
                point_distances.sort(key=lambda x: x[0])
                
                for i in range(min(needed, len(point_distances))):
                    _, point_idx = point_distances[i]
                    labels[point_idx] = small_cluster
                    final_sizes[small_cluster] += 1
                    final_sizes[largest_cluster] -= 1
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
    ['year', 'genre_encoded', 'genre_encoded', 'genre_encoded', 'age_category_encoded', 'age_category_encoded', 'country_encoded', 'country_encoded']]

scaler = StandardScaler()
featuresScaled = scaler.fit_transform(featuresFinal)

af = KMeans(n_clusters=35, random_state=42)
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
