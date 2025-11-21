package http

import (
	"html"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/search/delivery/grpc/gen"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type SearchHandler struct {
	client gen.SearchClient
}

func NewSearchHandler(client gen.SearchClient) *SearchHandler {
	return &SearchHandler{client: client}
}

// GetFilmsAndActorsFromSearch godoc
// @Summary Search films and actors
// @Tags search
// @Produce json
// @Param q query string true "Search string"
// @Param        films_count   query     int  false  "Number of films" default(10)
// @Param        films_offset  query     int  false  "Offset" default(0)
// @Param        actors_count   query     int  false  "Number of actors" default(10)
// @Param        actors_offset  query     int  false  "Offset" default(0)
// @Success 200 {object} models.SearchResponse
// @Failure 500
// @Router /search [get]
func (s *SearchHandler) GetFilmsAndActorsFromSearch(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))

	searchString := helpers.GetStringParameter(r, "q", "")
	searchString = html.EscapeString(searchString)

	filmsPager := models.Pager{
		Count:  helpers.GetParameter(r, "films_count", 10),
		Offset: helpers.GetParameter(r, "films_offset", 0),
	}
	actorsPager := models.Pager{
		Count:  helpers.GetParameter(r, "actors_count", 10),
		Offset: helpers.GetParameter(r, "actors_offset", 0),
	}

	result, err := s.client.SearchFilmsAndActors(r.Context(), &gen.SearchFilmsAndActorsRequest{
		SearchString: searchString,
		FilmsPager:   &gen.Pager{Count: int32(filmsPager.Count), Offset: int32(filmsPager.Offset)},
		ActorsPager:  &gen.Pager{Count: int32(actorsPager.Count), Offset: int32(actorsPager.Offset)},
	})

	if err != nil {
		log.LogHandlerError(logger, err, http.StatusInternalServerError)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	response := models.SearchResponse{}
	if len(result.Films) == 0 {
		response.Films = []models.MainPageFilm{}
	}
	if len(result.Actors) == 0 {
		response.Actors = []models.MainPageActor{}
	}

	for i := range result.Films {
		var film models.MainPageFilm
		film.ID = uuid.FromStringOrNil(result.Films[i].ID)
		film.Cover = result.Films[i].Cover
		film.Title = result.Films[i].Title
		film.Rating = result.Films[i].Rating
		film.Genre = result.Films[i].Genre
		film.Year = int(result.Films[i].Year)
		response.Films = append(response.Films, film)
	}

	for i := range result.Actors {
		var actor models.MainPageActor
		actor.ID = uuid.FromStringOrNil(result.Actors[i].ID)
		actor.RussianName = result.Actors[i].RussianName
		actor.Photo = result.Actors[i].Photo
		response.Actors = append(response.Actors, actor)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
