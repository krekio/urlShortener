package save

import (
	"errors"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"github.com/krekio/urlShortener.git/internal/lib/api/response"
	"github.com/krekio/urlShortener.git/internal/lib/random"
	"github.com/krekio/urlShortener.git/internal/storage"
	"log/slog"
	"net/http"
)

const (
	aliasLength = 6
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}
type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type urlSaver interface {
	SaveURL(urlToSave string, alias string) error
}

func New(log *slog.Logger, saver urlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to parse request", "error", err)
			render.JSON(w, r, response.Error("failed to parse request"))
			return
		}
		log.Info("request received", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("failed to validate request", "error", err)
			render.JSON(w, r, response.Error("failed to validate request"))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = random.RandomAlias(aliasLength)
		}
		err = saver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save url", "error", err)
			render.JSON(w, r, response.Error("failed to save url"))
		}
		log.Info("url saved", slog.String("url", req.URL))
		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}

}
