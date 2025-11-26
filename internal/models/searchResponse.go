package models

type SearchResponse struct {
	Films        []MainPageFilm  `json:"films"`
	Actors       []MainPageActor `json:"actors"`
	SearchString string          `json:"search_string"`
}
