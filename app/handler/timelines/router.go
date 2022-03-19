package timelines

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

// Create Handler for `/v1/timelines/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()
	h := &handler{app: app}

	r.Get("/public", h.Public)

	r.Route("/home", func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Get("/", h.Home)
	})
	return r
}
