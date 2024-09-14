package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"decentrala.org/events/internal/types"
)

func (controller *Controller) getEvents(w http.ResponseWriter, r *http.Request) {
	cityParam := types.CityCode(r.URL.Query().Get("city"))
	fromParam := r.URL.Query().Get("from")
	toParam := r.URL.Query().Get("to")

	from := parseDate(fromParam)
	to := parseDate(toParam)

	events, err := controller.model.GetEvents(cityParam, "", from, to)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	if from == nil {
		now := time.Now()
		from = &now
	}

	pe := controller.view.NewExecutor(w, r.Context())
	err = pe.EventsList(events, cityParam, formatDate(from), formatDate(to))
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) getEvent(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")

	eventId, err := strconv.Atoi(idValue)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	event, err := controller.model.GetEvent(types.EventId(eventId))
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())
	err = pe.EventPage(event)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) postEvent(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.LoggedIn {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	event, validationErrors := parseEventFromForm(r)
	if len(validationErrors) != 0 {
		pe := controller.view.NewExecutor(w, r.Context())
		err := pe.NewEvent(event, validationErrors)
		if err != nil {
			controller.serverError(http.StatusInternalServerError, err)(w, r)
		}
		return
	}

	eventId, err := controller.model.CreateEvent(event)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/events/%d", eventId), http.StatusSeeOther)
}

func (controller *Controller) deleteEvent(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.LoggedIn {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	id := r.URL.Query().Get("to")

	err := controller.model.DeleteEvent(id)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	http.Redirect(w, r, "/events/", http.StatusPermanentRedirect)
}

func (controller *Controller) newEvent(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.LoggedIn {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())

	err := pe.NewEvent(types.Event{}, nil)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) editEvent(w http.ResponseWriter, r *http.Request) {
	idValue := r.PathValue("id")

	eventId, err := strconv.Atoi(idValue)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	user := getUser(r)
	if !user.LoggedIn {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	event, err := controller.model.GetEvent(types.EventId(eventId))
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())
	err = pe.EditEvent(event, nil)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func (controller *Controller) editEventPost(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	if !user.LoggedIn {
		controller.serverError(http.StatusMethodNotAllowed, nil)(w, r)
		return
	}

	pe := controller.view.NewExecutor(w, r.Context())

	err := pe.NewEvent(types.Event{}, nil)
	if err != nil {
		controller.serverError(http.StatusInternalServerError, err)(w, r)
		return
	}
}

func parseEventFromForm(r *http.Request) (types.Event, []string) {
	validationErrors := []string{}
	var err error

	nameValue := r.FormValue("name")
	descriptionValue := r.FormValue("description")
	eventTypeValue := r.FormValue("eventType")
	fromValue := r.FormValue("from")
	toValue := r.FormValue("from")
	addressValue := r.FormValue("address")
	cityValue := r.FormValue("city")
	languageValue := r.FormValue("language")
	websiteValue := r.FormValue("website")
	entryValue := r.FormValue("entry")

	fmt.Println(fromValue)

	from := parsDateTime(fromValue)
	to := parsDateTime(toValue)
	if from == nil {
		validationErrors = append(validationErrors, "Format datuma početka nije dobar")
	}

	var eventType, entry, language int
	if eventTypeValue != "" {
		eventType, err = strconv.Atoi(eventTypeValue)
		if err != nil || eventType < 0 || eventType > 9 {
			validationErrors = append(validationErrors, "Vrednost tipa događaja nije dobra")
		}
	}

	if languageValue != "" {
		language, err = strconv.Atoi(languageValue)
		if err != nil || language < 0 || language > 3 {
			validationErrors = append(validationErrors, "Vrednost jezika nije dobra")
		}
	}

	if entryValue != "" {
		entry, err = strconv.Atoi(entryValue)
		if err != nil {
			validationErrors = append(validationErrors, "Vrednost cene ulaza nije dobra")
		}
	}

	if len(addressValue) > 100 {
		validationErrors = append(validationErrors, "Adresa može sadržati najviše 100 karaktera")
	}

	if len(descriptionValue) > 200 {
		validationErrors = append(validationErrors, "Opis može sadržati najviše 2000 karaktera")
	}

	if len(nameValue) < 5 || len(nameValue) > 100 {
		validationErrors = append(validationErrors, "Ime mora sadržati između 5 i 100 karaktera")
	}

	if websiteValue != "" && !isUrl(websiteValue) {
		validationErrors = append(validationErrors, "Format sajta nije dobar")
	}

	user := getUser(r)

	event := types.Event{
		OrganizationCode: user.OrganizationCode,
		Name:             nameValue,
		Description:      descriptionValue,
		EventType:        types.EventType(eventType),
		Address:          addressValue,
		Website:          websiteValue,
		EndsAt:           to,
		EntryPrice:       entry,
		Language:         language,
		CityCode:         types.CityCode(cityValue),
	}

	return event, validationErrors
}
