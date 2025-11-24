package http

import (
	"errors"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/films/delivery/grpc/gen"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CompilationHandler struct {
	client gen.FilmsClient
}

func NewCompilationHandler(client gen.FilmsClient) *CompilationHandler {
	return &CompilationHandler{client: client}
}

// GetCompilation godoc
// @Summary Получить подборку по ID
// @Description Возвращает информацию о конкретной подборке по её идентификатору
// @Tags compilations
// @Accept json
// @Produce json
// @Param id path string true "UUID подборки"
// @Success 200 {object} models.Compilation "Успешный ответ с данными подборки"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 404 {object} object "Подборка не найдена"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /compilations/{id} [get]
func (g *CompilationHandler) GetCompilation(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of compilation"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	compilation, err := g.client.GetCompilation(r.Context(), &gen.GetCompilationRequest{CompilationId: id.String()})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := models.Compilation{
		ID:          uuid.FromStringOrNil(compilation.Compilation.Id),
		Title:       compilation.Compilation.Name,
		Description: compilation.Compilation.Description,
		Icon:        compilation.Compilation.Icon,
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetCompilations godoc
// @Summary Получить список подборок
// @Description Возвращает список всех подборок с пагинацией
// @Tags compilations
// @Accept json
// @Produce json
// @Param count query int false "Количество элементов (по умолчанию 10)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {array} models.Compilation "Успешный ответ со списком подборок"
// @Failure 404 {object} object "Подборки не найдены"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /compilations/ [get]
func (g *CompilationHandler) GetCompilations(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	pager := helpers.GetPagerFromRequest(r)

	compilations, err := g.client.GetCompilations(r.Context(), &gen.GetCompilationsRequest{
		Pager: &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := []models.Compilation{}
	for i := range compilations.Compilations {
		compilation := models.Compilation{
			ID:          uuid.FromStringOrNil(compilations.Compilations[i].Id),
			Title:       compilations.Compilations[i].Name,
			Description: compilations.Compilations[i].Description,
			Icon:        compilations.Compilations[i].Icon,
		}
		response = append(response, compilation)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

// GetFilmsByCompilation godoc
// @Summary Получить фильмы из подборки
// @Description Возвращает список фильмов, входящих в конкретную подборку, с пагинацией
// @Tags compilations
// @Accept json
// @Produce json
// @Param id path string true "UUID подборки"
// @Param count query int false "Количество элементов (по умолчанию 10)"
// @Param offset query int false "Смещение (по умолчанию 0)"
// @Success 200 {array} models.FavFilm "Успешный ответ со списком фильмов"
// @Failure 400 {object} object "Неверный формат UUID"
// @Failure 404 {object} object "Подборка или фильмы не найдены"
// @Failure 500 {object} object "Внутренняя ошибка сервера"
// @Router /compilations/{id}/films [get]
func (g *CompilationHandler) GetFilmsByCompilation(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])
	if err != nil {
		log.LogHandlerError(logger, errors.New("invalid id of compilation"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	pager := helpers.GetPagerFromRequest(r)

	films, err := g.client.GetFilmsByCompilation(r.Context(), &gen.GetFilmsByCompilationRequest{
		CompilationId: id.String(),
		Pager:         &gen.Pager{Count: int32(pager.Count), Offset: int32(pager.Offset)},
	})
	if err != nil {
		st, _ := status.FromError(err)
		switch st.Code() {
		case codes.NotFound:
			helpers.WriteError(w, http.StatusNotFound)
		case codes.InvalidArgument:
			helpers.WriteError(w, http.StatusBadRequest)
		default:
			helpers.WriteError(w, http.StatusInternalServerError)
		}
		return
	}

	response := []models.FavFilm{}
	for i := range films.Films {
		film := models.FavFilm{
			ID:               uuid.FromStringOrNil(films.Films[i].Id),
			Image:            films.Films[i].Image,
			Title:            films.Films[i].Title,
			Rating:           films.Films[i].Rating,
			Year:             int(films.Films[i].Year),
			Genre:            films.Films[i].Genre,
			ShortDescription: films.Films[i].ShortDescription,
			Duration:         int(films.Films[i].Duration),
		}
		response = append(response, film)
	}

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
