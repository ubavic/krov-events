package controller

import (
	"net/http"
)

func (controller *Controller) aboutPage(w http.ResponseWriter, r *http.Request) {
	pe := controller.view.NewExecutor(w, r.Context())
	err := pe.AboutPage()
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}
