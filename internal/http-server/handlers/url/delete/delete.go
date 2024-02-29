package delete

import (
	resp "URLShortenePetPrpoject/internal/lib/api/response"
	"URLShortenePetPrpoject/internal/lib/logger/sl"
	"URLShortenePetPrpoject/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.delete.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		err := urlDeleter.DeleteURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlNotFound) {
				log.Info("url not found", slog.String("alias", alias))

				render.JSON(w, r, resp.Error("not found"))

				return
			}

			log.Error("failed to delete url", sl.Err(err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		render.JSON(w, r, resp.Response{
			Status: resp.StatusOk,
		})
	}
}
