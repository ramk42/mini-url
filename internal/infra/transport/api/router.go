package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/ramk42/mini-url/internal/infra/database"
	"github.com/ramk42/mini-url/internal/infra/env"
	"github.com/ramk42/mini-url/internal/infra/repository"
	"github.com/ramk42/mini-url/internal/infra/transport/api/handler"
	apimiddleware "github.com/ramk42/mini-url/internal/infra/transport/api/middleware"
	"github.com/ramk42/mini-url/internal/usecase"
	"net/http"
)

func createRouter() *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.RequestID)
	r.Use(apimiddleware.RequestLogger)
	r.NotFound(JSONNotFoundHandler)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Recoverer)

	// Routes
	r.Get("/health", HealthProbe)

	urlRepository := repository.NewURL(database.GetInstance())

	baseURL := env.GetEnvAsString("BASE-URL", "http://localhost:8080")
	shortenerUsecase := usecase.NewURLShortener(urlRepository, baseURL)
	urlHandler := handler.NewURL(shortenerUsecase)

	r.Post("/shorten", urlHandler.Shorten)
	r.Get("/{shortCode}", urlHandler.Resolve)

	return r
}

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code"`
}

func JSONNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render.JSON(w, r, ErrorResponse{
		Error:      "Resource not found",
		StatusCode: http.StatusNotFound,
	})
}
