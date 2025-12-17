package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"kinopoisk/internal/models"
	"kinopoisk/internal/pkg/helpers"
	"kinopoisk/internal/pkg/search/delivery/grpc/gen"
	"kinopoisk/internal/pkg/utils/log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	uuid "github.com/satori/go.uuid"
)

type SearchHandler struct {
	client     gen.SearchClient
	voiceToken string
	voiceVKURL string
}

func NewSearchHandler(client gen.SearchClient) *SearchHandler {
	voiceToken := os.Getenv("VOICE_TOKEN")
	voiceVKURL := os.Getenv("VK_VOICE_URL")
	return &SearchHandler{client: client, voiceToken: voiceToken, voiceVKURL: voiceVKURL}
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
	unescapedString, err := url.QueryUnescape(searchString)
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to unescape search string"), http.StatusBadRequest)
	} else {
		searchString = unescapedString
	}

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

	response.SearchString = searchString

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}

func (s *SearchHandler) VoiceSearch(w http.ResponseWriter, r *http.Request) {
	logger := log.GetLoggerFromContext(r.Context()).With(slog.String("func", log.GetFuncName()))
	var twoChannels bool

	voiceData, err := io.ReadAll(r.Body)
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to read voice data"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if len(voiceData) == 0 {
		log.LogHandlerError(logger, errors.New("failed to read voice data"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	reader := bytes.NewReader(voiceData)
	decoder := wav.NewDecoder(reader)

	ddbuff, err := decoder.FullPCMBuffer()
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to read voice data"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusBadRequest)
		return
	}

	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	tempDir := filepath.Join(currentDir, "tmp")

	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.LogHandlerError(logger,
			fmt.Errorf("failed to create temp dir %s: %w", tempDir, err),
			http.StatusInternalServerError,
		)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	convertedFile, err := os.CreateTemp(tempDir, "voice-*.wav")
	if err != nil {
		log.LogHandlerError(logger, err, http.StatusInternalServerError)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}
	defer os.Remove(convertedFile.Name())
	defer convertedFile.Close()

	if ddbuff.Format.NumChannels != 1 {
		twoChannels = true
		numSamples := len(ddbuff.Data) / 2
		monoData := make([]int, numSamples)

		for i := 0; i < numSamples; i++ {
			left := ddbuff.Data[i*2]
			right := ddbuff.Data[i*2+1]
			monoData[i] = (left + right) / 2
		}

		encoder := wav.NewEncoder(convertedFile, int(ddbuff.Format.SampleRate), 16, 1, 1)

		monobuff := &audio.IntBuffer{
			Format: &audio.Format{
				SampleRate:  ddbuff.Format.SampleRate,
				NumChannels: 1,
			},
			Data:           monoData,
			SourceBitDepth: 16,
		}

		err = encoder.Write(monobuff)
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to read voice data"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
			return
		}

		err = encoder.Close()
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to read voice data"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusBadRequest)
			return
		}
	}

	var voiceRequest *http.Request
	if twoChannels == false {
		voiceRequest, err = http.NewRequest("POST", s.voiceVKURL, bytes.NewReader(voiceData))
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to create request"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusInternalServerError)
			return
		}
	} else {
		_, err = convertedFile.Seek(0, 0)
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to seek to beginning of file"), http.StatusInternalServerError)
			helpers.WriteError(w, http.StatusInternalServerError)
			return
		}
		voiceRequest, err = http.NewRequest("POST", s.voiceVKURL, convertedFile)
		if err != nil {
			log.LogHandlerError(logger, errors.New("failed to create request"), http.StatusBadRequest)
			helpers.WriteError(w, http.StatusInternalServerError)
			return
		}
	}

	voiceRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.voiceToken))
	voiceRequest.Header.Set("Content-Type", "audio/wav")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: 30 * time.Second,
	}
	voiceResponse, err := client.Do(voiceRequest)
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to send request to VK"), http.StatusBadRequest)
		log.LogHandlerError(logger, err, http.StatusBadRequest)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	voiceResponseData, err := io.ReadAll(voiceResponse.Body)
	if err != nil {
		log.LogHandlerError(logger, err, http.StatusBadRequest)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}
	defer voiceResponse.Body.Close()

	if voiceResponse.StatusCode != http.StatusOK {
		log.LogHandlerError(logger, fmt.Errorf("failed to get voice data, status: %d, response: %s", voiceResponse.StatusCode, string(voiceResponseData)), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	var voiceResult models.VoiceResult
	err = json.Unmarshal(voiceResponseData, &voiceResult)
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to parse response"), http.StatusBadRequest)
		helpers.WriteError(w, http.StatusInternalServerError)
		return
	}

	if len(voiceResult.Result.Texts) == 0 {
		log.LogHandlerInfo(logger, "no text given", http.StatusOK)
		helpers.WriteJSON(w, models.SearchResponse{})
		return
	}

	bestConfidence := 0.0
	textResult := ""

	for _, text := range voiceResult.Result.Texts {
		if text.Confidence > bestConfidence {
			textResult = text.Text
			bestConfidence = text.Confidence
		}
	}

	if textResult == "" {
		log.LogHandlerInfo(logger, "no text given", http.StatusOK)
		helpers.WriteJSON(w, models.SearchResponse{})
		return
	}

	filmsPager := models.Pager{
		Count:  helpers.GetParameter(r, "films_count", 10),
		Offset: helpers.GetParameter(r, "films_offset", 0),
	}
	actorsPager := models.Pager{
		Count:  helpers.GetParameter(r, "actors_count", 10),
		Offset: helpers.GetParameter(r, "actors_offset", 0),
	}

	unescapedText, err := url.QueryUnescape(textResult)
	if err != nil {
		log.LogHandlerError(logger, errors.New("failed to unescape search string"), http.StatusBadRequest)
	} else {
		textResult = unescapedText
	}

	result, err := s.client.SearchFilmsAndActors(r.Context(), &gen.SearchFilmsAndActorsRequest{
		SearchString: textResult,
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

	response.SearchString = textResult

	helpers.WriteJSON(w, response)
	log.LogHandlerInfo(logger, "success", http.StatusOK)
}
