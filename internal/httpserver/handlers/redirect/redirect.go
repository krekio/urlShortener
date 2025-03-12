package redirect

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/krekio/urlShortener.git/internal/lib/api/response"
	"github.com/krekio/urlShortener.git/internal/storage"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"
		log = log.With(
			slog.String("operation", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("not found"))
			return
		}
		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found")
			render.JSON(w, r, response.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", "error", err)
			render.JSON(w, r, response.Error("internal error"))
			return
		}
		log.Info("got url", "url", resURL)
		http.Redirect(w, r, resURL, http.StatusTemporaryRedirect)
	}
}
