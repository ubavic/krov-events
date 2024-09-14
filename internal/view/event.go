package view

import (
	"decentrala.org/events/internal/types"
)

func (pe *PageExecutor) EventsList(events []types.Event, selectedCity types.CityCode, fromDate, toDate string) error {
	data := make(map[string]any)
	data["cities"] = types.Cities
	data["events"] = events
	data["selectedCity"] = selectedCity
	data["from"] = fromDate
	data["to"] = toDate

	return pe.executePage("eventList.html", "Događaji", data)
}

func (pe *PageExecutor) EventPage(event types.Event) error {
	data := make(map[string]any)
	data["event"] = event

	return pe.executePage("eventPage.html", event.Name, data)
}

func (pe *PageExecutor) NewEvent(event types.Event, validationErrors []string) error {
	data := make(map[string]any)
	data["cities"] = types.Cities
	data["eventTypes"] = types.EventTypes
	data["languages"] = types.Languages
	data["event"] = event
	data["validationErrors"] = validationErrors
	data["edit"] = false

	return pe.executePage("eventForm.html", "Novi događaj", data)
}

func (pe *PageExecutor) EditEvent(event types.Event, validationErrors []string) error {
	data := make(map[string]any)
	data["cities"] = types.Cities
	data["eventTypes"] = types.EventTypes
	data["languages"] = types.Languages
	data["event"] = event
	data["validationErrors"] = validationErrors
	data["edit"] = true

	return pe.executePage("eventForm.html", "Izmeni događaj "+event.Name, data)
}
