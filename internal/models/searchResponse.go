package models

type SearchResponse struct {
	Films  []MainPageFilm  `json:"films"`
	Actors []MainPageActor `json:"actors"`
}
