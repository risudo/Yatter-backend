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

	r.Use(auth.Middleware(app))
	h := &handler{app: app}
	r.Post("/", h.Post)
	return r
}
