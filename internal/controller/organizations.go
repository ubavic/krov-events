package controller

import (
	"net/http"

	"decentrala.org/events/internal/types"
)

func (controller *Controller) getOrganizations(w http.ResponseWriter, r *http.Request) {
	cityParam := types.CityCode(r.URL.Query().Get("city"))

	organizations, err := controller.model.GetOrganizations(cityParam)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())
	err = pe.OrganizationList(organizations, cityParam)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) getOrganization(w http.ResponseWriter, r *http.Request) {
	orgId := r.PathValue("id")

	if orgId == "" {
		http.Redirect(w, r, "/organizations/", http.StatusSeeOther)
		return
	}

	organization, err := controller.model.GetOrganization(types.OrganizationCode(orgId))
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	events, err := controller.model.GetOrganizationNext10Events(organization.Code)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())

	err = pe.OrganizationPage(organization, events)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) newOrganization(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.Admin() {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())
	err := pe.NewOrganization(types.Organization{})
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) postOrganization(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.Admin() {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	organization, validationErrors := parseOrganizationFromForm(r)
	if len(validationErrors) > 0 {
		pe := controller.view.NewExecutor(w, r.Context())
		err := pe.NewOrganization(organization)
		if err != nil {
			controller.serverError(http.StatusInternalServerError, err)(w, r)
		}
		return
	}

	err := controller.model.CreateOrganization(organization)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	http.Redirect(w, r, "/organizations/"+string(organization.Code), http.StatusSeeOther)
}

func (controller *Controller) editOrganization(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.Admin() {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	orgId := r.PathValue("id")
	if orgId == "" {
		http.Redirect(w, r, "/organizations/", http.StatusSeeOther)
		return
	}

	organization, err := controller.model.GetOrganization(types.OrganizationCode(orgId))
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())
	err = pe.EditOrganization(organization, nil)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) editOrganizationPost(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.Admin() {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	organization, validationErrors := parseOrganizationFromForm(r)
	if (user.OrganizationCode() != organization.Code) && !user.Admin() {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	if len(validationErrors) > 0 {
		pe := controller.view.NewExecutor(w, r.Context())
		err := pe.EditOrganization(organization, validationErrors)
		if err != nil {
			controller.serverError(http.StatusInternalServerError, err)(w, r)
		}
		return
	}

	err := controller.model.ModifyOrganization(organization)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	http.Redirect(w, r, "/organizations/"+string(organization.Code), http.StatusSeeOther)
}

func parseOrganizationFromForm(r *http.Request) (types.Organization, []string) {
	validationErrors := make([]string, 0)

	organization := types.Organization{
		Name:        r.FormValue("name"),
		Code:        types.OrganizationCode(r.FormValue("code")),
		Description: r.FormValue("description"),
		Email:       r.FormValue("email"),
		Website:     r.FormValue("website"),
		Address:     r.FormValue("address"),
		CityCode:    types.CityCode(r.FormValue("city")),
		OsmUrl:      r.FormValue("osm"),
	}

	if r.FormValue("osm") == "on" {
		organization.SetApiAllowed(true)
	}

	if len := len(organization.Code); len < 3 || len > 20 {
		validationErrors = append(validationErrors, "Kod mora biti između 3 i 20 karaktera")
	}

	if len := len(organization.Name); len < 3 || len > 100 {
		validationErrors = append(validationErrors, "Ime mora biti između 3 i 100 karaktera")
	}

	if len(organization.Description) > 2000 {
		validationErrors = append(validationErrors, "Opis mora biti kraći od 2000 karaktera")
	}

	if organization.Website != "" && !isUrl(organization.Website) {
		validationErrors = append(validationErrors, "Format sajta nije dobar")
	}

	if organization.OsmUrl != "" && !isUrl(organization.OsmUrl) {
		validationErrors = append(validationErrors, "Format OSM adrese nije dobar")
	}

	return organization, validationErrors
}
