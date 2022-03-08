package status

import (
	"net/http"
	"yatter-backend-go/app/app"
	"yatter-backend-go/app/handler/auth"

	"github.com/go-chi/chi"
)

type handler struct {
	app *app.App
}

func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Route("/", func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Post("/", h.Post)
	})
	r.Get("/{id}", h.Get)
	return r
}

func NewDeleteRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Route("/{id}", func (r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Delete("/", h.Delete)
	})

	return r
}
