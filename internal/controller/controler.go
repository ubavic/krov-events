package controller

import (
	"context"
	"errors"
	"io/fs"
	"net/http"

	"decentrala.org/events/internal/model"
	"decentrala.org/events/internal/types"
	"decentrala.org/events/internal/view"
)

type Controller struct {
	model    model.Model
	view     view.View
	staticFS fs.FS
}

func NewController(Model model.Model, View view.View, StaticFS fs.FS) Controller {
	return Controller{
		model:    Model,
		view:     View,
		staticFS: StaticFS,
	}
}

func (controller *Controller) Mux() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServerFS(controller.staticFS)))

	mux.HandleFunc("GET /organizations/new/{$}", controller.newOrganization)
	mux.HandleFunc("POST /organizations/new/{$}", controller.postOrganization)

	mux.HandleFunc("GET /organizations/{id}/edit", controller.editOrganization)
	mux.HandleFunc("POST /organizations/{id}/edit", controller.editOrganizationPost)

	mux.HandleFunc("GET /organizations/{id}/", controller.getOrganization)
	mux.HandleFunc("GET /organizations/{$}", controller.getOrganizations)

	mux.HandleFunc("GET /events/new/{$}", controller.newEvent)
	mux.HandleFunc("POST /events/new/{$}", controller.postEvent)

	mux.HandleFunc("GET /events/{id}/edit", controller.editEvent)
	mux.HandleFunc("POST /events/{id}/edit", controller.editEventPost)

	mux.HandleFunc("GET /events/{id}/", controller.getEvent)
	mux.HandleFunc("DELETE /events/{id}/", controller.deleteEvent)
	mux.HandleFunc("GET /events/{$}", controller.getEvents)

	mux.HandleFunc("GET /login/{$}", controller.login)
	mux.HandleFunc("GET /logout/{$}", controller.logout)

	mux.HandleFunc("GET /about/", controller.aboutPage)
	mux.HandleFunc("GET /{$}", controller.getEvents)
	mux.HandleFunc("/", controller.serverError(404, errors.New("")))

	return controller.userMiddleware(mux)
}

func (controller *Controller) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := types.WebsiteUser{
			Admin:            false,
			OrganizationCode: "",
			Organization:     "",
			LoggedIn:         false,
		}

		ctx := context.WithValue(r.Context(), "user", user)
		newReq := r.WithContext(ctx)

		next.ServeHTTP(w, newReq)
	})
}

func getUser(r *http.Request) types.WebsiteUser {
	user, ok := r.Context().Value("user").(types.WebsiteUser)
	if !ok {
		return types.WebsiteUser{
			LoggedIn: false,
			Admin:    false,
		}
	}

	return user
}