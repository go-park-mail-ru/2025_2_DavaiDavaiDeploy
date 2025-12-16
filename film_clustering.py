import pandas as pd
import numpy as np
from sklearn.cluster import KMeans
from sklearn.preprocessing import StandardScaler
from sklearn.decomposition import PCA
import matplotlib.pyplot as plt
from sqlalchemy import create_engine
from sqlalchemy import text
from collections import Counter
import math

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
        query = """
        SELECT f.id, f.title, f.year, f.genre_id, f.age_category, f.country_id, f.duration, 
               ROUND(COALESCE(AVG(ff.rating), 0), 1) as rating 
        FROM film f 
        LEFT JOIN film_feedback ff ON f.id = ff.film_id 
        GROUP BY f.id, f.title, f.year, f.genre_id, f.age_category, f.country_id, f.duration 
        ORDER BY rating DESC
        """
        films = pd.read_sql(query, engine)
        print(f"Загружено {len(films)} фильмов")
        return films
    except Exception as e:
        print(f"Ошибка при загрузке данных: {e}")
        return None

def create_genre_ordering():
    genre_order = {
        'нуар': 1,
        'исторические': 2,
        'документальные': 3,
        'биографии': 4,
        'драмы': 5,
        'мелодрамы': 6,
        'музыкальные': 7,
        'ромком': 8,
        'комедии': 9,
        'короткометражки': 10,
        'спортивные': 11,
        'семейные': 12,
        'мультфильмы': 13,
        'аниме': 14,
        'приключения': 15,
        'детективы': 16,
        'вестерны': 17,
        'боевики': 18,
        'криминал': 19,
        'триллеры': 20,
        'ужасы': 21,
        'мистика': 22,
        'фэнтези': 23,
        'фантастика': 24
    }
    return genre_order

def create_country_ordering():
    country_order = {
        'США': 1,
        'Великобритания': 2,
        'Канада': 3,
        'Франция': 4,
        'Германия': 5,
        'Япония': 6,
        'Индия': 7,
        'Россия': 8,
        'СССР': 9,
        'Новая Зеландия': 10
    }
    return country_order

def enforce_genre_proximity_constraint(labels, genre_ids, max_distance=2):
    unique_labels = np.unique(labels)
    
    for cluster in unique_labels:
        cluster_mask = labels == cluster
        cluster_genres = genre_ids[cluster_mask]
        
        if len(cluster_genres) == 0:
            continue
            
        max_genre = max(cluster_genres)
        min_genre = min(cluster_genres)
        
        if (max_genre - min_genre) > max_distance:
            genre_groups = {}
            for genre in cluster_genres:
                group_key = (genre // (max_distance + 1)) * (max_distance + 1)
                if group_key not in genre_groups:
                    genre_groups[group_key] = []
                genre_groups[group_key].append(genre)
            
            largest_group = max(genre_groups.values(), key=len)
            genre_to_keep = Counter(largest_group).most_common(1)[0][0]
            
            for i, genre in enumerate(cluster_genres):
                if genre != genre_to_keep:
                    cluster_indices = np.where(cluster_mask)[0]
                    point_idx = cluster_indices[i]
                    
                    best_cluster = None
                    best_score = float('inf')
                    
                    for other_cluster in unique_labels:
                        if other_cluster == cluster:
                            continue
                            
                        other_mask = labels == other_cluster
                        other_genres = genre_ids[other_mask]
                        
                        if len(other_genres) == 0:
                            continue
                            
                        other_max = max(other_genres)
                        other_min = min(other_genres)
                        
                        if abs(genre - other_max) <= max_distance and abs(genre - other_min) <= max_distance:
                            cluster_size = np.sum(labels == other_cluster)
                            score = cluster_size
                            
                            if score < best_score:
                                best_score = score
                                best_cluster = other_cluster
                    
                    if best_cluster is not None:
                        labels[point_idx] = best_cluster
    
    return labels

def constrained_kmeans(features, n_clusters, min_size, genre_ids):
    kmeans = KMeans(n_clusters=n_clusters, random_state=35, n_init=10)
    labels = kmeans.fit_predict(features)
    distances = kmeans.transform(features)
    
    labels = enforce_genre_proximity_constraint(labels, genre_ids)
    
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
                
            small_cluster_genres = genre_ids[labels == small_cluster]
            small_min_genre = min(small_cluster_genres) if len(small_cluster_genres) > 0 else 0
            small_max_genre = max(small_cluster_genres) if len(small_cluster_genres) > 0 else 0
                
            candidates = []
            for large_cluster in large_clusters:
                if cluster_sizes[large_cluster] <= min_size:
                    continue
                    
                points_in_large = np.where(labels == large_cluster)[0]
                for point_idx in points_in_large:
                    distance_diff = distances[point_idx, small_cluster] - distances[point_idx, large_cluster]
                    cluster_size_penalty = cluster_sizes[large_cluster] / max(cluster_sizes)
                    score = distance_diff * cluster_size_penalty
                    
                    point_genre = genre_ids[point_idx]
                    
                    new_min = min(small_min_genre, point_genre)
                    new_max = max(small_max_genre, point_genre)
                    
                    if (new_max - new_min) <= 2:
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
                
                labels = enforce_genre_proximity_constraint(labels, genre_ids)
                cluster_sizes = np.bincount(labels, minlength=n_clusters)
                
                small_cluster_genres = genre_ids[labels == small_cluster]
                small_min_genre = min(small_cluster_genres) if len(small_cluster_genres) > 0 else 0
                small_max_genre = max(small_cluster_genres) if len(small_cluster_genres) > 0 else 0
            
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
                
            small_cluster_genres = genre_ids[labels == small_cluster]
            small_min_genre = min(small_cluster_genres) if len(small_cluster_genres) > 0 else 0
            small_max_genre = max(small_cluster_genres) if len(small_cluster_genres) > 0 else 0
                
            points_in_largest = np.where(labels == largest_cluster)[0]
            
            if len(points_in_largest) > 0:
                point_distances = []
                for point_idx in points_in_largest:
                    dist = distances[point_idx, small_cluster]
                    point_genre = genre_ids[point_idx]
                    
                    new_min = min(small_min_genre, point_genre)
                    new_max = max(small_max_genre, point_genre)
                    
                    if (new_max - new_min) <= 2:
                        point_distances.append((dist, point_idx))
                
                point_distances.sort(key=lambda x: x[0])
                
                for i in range(min(needed, len(point_distances))):
                    _, point_idx = point_distances[i]
                    labels[point_idx] = small_cluster
                    final_sizes[small_cluster] += 1
                    final_sizes[largest_cluster] -= 1
    
    labels = enforce_genre_proximity_constraint(labels, genre_ids)
    
    return labels

films = loadFilmsData()
if films is not None:
    genre_order = create_genre_ordering()
    country_order = create_country_ordering()
    
    genre_mapping_query = "SELECT id, LOWER(title) as title_lower FROM genre"
    country_mapping_query = "SELECT id, name FROM country"
    
    engine = create_engine('postgresql://postgres:qwerty@localhost:5432/kinopoisk')
    genres_df = pd.read_sql(genre_mapping_query, engine)
    countries_df = pd.read_sql(country_mapping_query, engine)
    
    genre_id_to_title = dict(zip(genres_df['id'], genres_df['title_lower']))
    country_id_to_name = dict(zip(countries_df['id'], countries_df['name']))
    
    films['genre_title'] = films['genre_id'].map(genre_id_to_title)
    films['country_name'] = films['country_id'].map(country_id_to_name)
    
    films['genre_encoded'] = films['genre_title'].map(genre_order).fillna(99).astype(int)
    films['country_encoded'] = films['country_name'].map(country_order).fillna(99).astype(int)
    
    age_mapping = {age: i for i, age in enumerate(sorted(films['age_category'].dropna().unique()))}
    films['age_category_encoded'] = films['age_category'].map(age_mapping).fillna(0).astype(int)
    
    features = films[['year', 'duration', 'rating']].copy()
    
    features['genre_angle'] = films['genre_encoded'].apply(lambda x: 2 * math.pi * (x - 1) / 23 if x <= 24 else 0)
    features['genre_cos'] = np.cos(features['genre_angle'])
    features['genre_sin'] = np.sin(features['genre_angle'])
    
    features['country_angle'] = films['country_encoded'].apply(lambda x: 2 * math.pi * (x - 1) / 9 if x <= 10 else 0)
    features['country_cos'] = np.cos(features['country_angle'])
    features['country_sin'] = np.sin(features['country_angle'])
    
    features['age_encoded'] = films['age_category_encoded']
    
    scaler = StandardScaler()
    features_scaled = scaler.fit_transform(features)
    
    genre_ids_array = films['genre_encoded'].values
    
    cluster_labels = constrained_kmeans(features_scaled, n_clusters=35, min_size=7, genre_ids=genre_ids_array)
    
    films['cluster'] = cluster_labels
    
    print(f"Найдено {len(np.unique(cluster_labels))} кластеров")
    print("Распределение по кластерам:")
    print(films['cluster'].value_counts().sort_index())
    
    for cluster in np.unique(cluster_labels):
        cluster_films = films[films['cluster'] == cluster]
        unique_genres = cluster_films['genre_title'].unique()
        genre_nums = cluster_films['genre_encoded'].unique()
        genre_range = max(genre_nums) - min(genre_nums) if len(genre_nums) > 0 else 0
        print(f"Кластер {cluster}: {len(cluster_films)} фильмов, жанры: {unique_genres}, диапазон жанров: {genre_range}")
    
    pca = PCA(n_components=2)
    features_2d = pca.fit_transform(features_scaled)
    
    plt.figure(figsize=(12, 8))
    plt.subplot(1, 2, 1)
    scatter = plt.scatter(features_2d[:, 0], features_2d[:, 1], c=cluster_labels, cmap='viridis', alpha=0.7)
    plt.colorbar(scatter)
    plt.title('Кластеризация фильмов (2D проекция)')
    plt.xlabel('PCA Component 1')
    plt.ylabel('PCA Component 2')
    
    plt.subplot(1, 2, 2)
    cluster_counts = films['cluster'].value_counts().sort_index()
    plt.bar(cluster_counts.index, cluster_counts.values)
    plt.title('Количество фильмов в каждом кластере')
    plt.xlabel('Номер кластера')
    plt.ylabel('Количество фильмов')
    
    plt.tight_layout()
    plt.show()
    
    save_clusters_to_db(films[['title', 'year', 'cluster']])
