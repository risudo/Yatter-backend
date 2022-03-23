package status

import (
	"net/http"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/statuses/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()
	h := &handler{app: app}

	r.Route("/", func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Post("/", h.Post)
	})

	r.Route("/{id}", func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Delete("/", h.Delete)
	})

	r.Get("/{id}", h.Fetch)

	return r
}
