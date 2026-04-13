package server

import (
	"io/fs"
	"net/http"

	"github.com/channinghe/waken/internal/config"
	"github.com/channinghe/waken/internal/handler"
	"github.com/channinghe/waken/internal/middleware"
	"github.com/channinghe/waken/internal/repository"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(cfg config.Config, repo *repository.DeviceRepository, frontendFS fs.FS) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))

	deviceHandler := handler.NewDeviceHandler(repo, cfg)
	wakeHandler := handler.NewWakeHandler(repo, cfg)

	r.Get("/api/health", handler.Health)

	r.Group(func(r chi.Router) {
		r.Use(middleware.BearerAuth(cfg.AuthToken))

		r.Get("/api/devices", deviceHandler.List)
		r.Post("/api/devices", deviceHandler.Create)
		r.Put("/api/devices/{id}", deviceHandler.Update)
		r.Delete("/api/devices/{id}", deviceHandler.Delete)

		r.Post("/api/wake/{id}", wakeHandler.WakeByID)
		r.Post("/api/wake/name/{name}", wakeHandler.WakeByName)
		r.Post("/api/wake", wakeHandler.WakeByMAC)

		r.Get("/api/scan", handler.Scan)
	})

	if frontendFS != nil {
		spaHandler := spaFileServer(frontendFS)
		r.Handle("/*", spaHandler)
	}

	return r
}

func spaFileServer(frontendFS fs.FS) http.Handler {
	fileServer := http.FileServer(http.FS(frontendFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the file directly
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else if path[0] == '/' {
			path = path[1:]
		}

		if _, err := fs.Stat(frontendFS, path); err != nil {
			// File not found — serve index.html for SPA routing
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
