package accounts

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

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()
	h := &handler{app: app}

	r.Route("/{username}", func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Post("/follow", h.Follow)
		r.Post("/unfollow", h.Unfollow)
	})

	r.Route("/relationships", func(r chi.Router) {
		r.Use(auth.Middleware(app))
		r.Get("/", h.Relationships)
	})

	r.Route("/udpate_credentials", func(r chi.Router) {
		r.Use(auth.Middleware(app))
	})

	r.Post("/", h.Create)
	r.Get("/{username}", h.Fetch)
	r.Get("/{username}/following", h.Following)
	r.Get("/{username}/followers", h.Followers)

	return r
}
