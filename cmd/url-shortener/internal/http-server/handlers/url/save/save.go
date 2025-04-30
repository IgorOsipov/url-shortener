package save

import (
	"errors"
	res "iosipoff/url-shortener/cmd/url-shortener/internal/lib/api/response"
	"iosipoff/url-shortener/cmd/url-shortener/internal/lib/random"
	"iosipoff/url-shortener/cmd/url-shortener/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	res.Response
	Alias string `json:"alias,omitempty"`
}

const aliasLength = 6

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (uint, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(slog.String("op", op), slog.String("request_id", middleware.GetReqID(r.Context())))

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request", slog.Any("error", err))
			render.JSON(w, r, res.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("error", err))
			render.JSON(w, r, res.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.Any("error", err))
			render.JSON(w, r, res.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", slog.Any("error", err))
			render.JSON(w, r, res.Error("failed to save url"))
			return
		}

		log.Info("url added", slog.Any("id", id))

		render.JSON(w, r, Response{
			Response: res.OK(),
			Alias:    alias,
		})
	}
}
