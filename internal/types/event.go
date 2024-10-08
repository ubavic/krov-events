package types

import "time"

type EventId int

type Event struct {
	Id               EventId
	OrganizationCode OrganizationCode
	OrganizationName string
	Name             string
	EventType        EventType
	EventTypeName    string
	Description      string
	Website          string
	StartsAt         time.Time
	StartsAtStr      string
	StartAtIso       string
	EndsAt           *time.Time
	EndsAtStr        string
	EndsAtIso        string
	CityCode         CityCode
	CityName         string
	Address          string
	EntryPrice       int
	Language         int
	Canceled         bool
	CanceledAt       *time.Time
	CreatedAt        time.Time
	ModifiedAt       *time.Time
	DeletedAt        *time.Time
}

func (e *Event) FormatDates() {
	e.StartsAtStr = e.StartsAt.Format("15:04 02.01.2006")
	e.StartAtIso = e.StartsAt.Format(time.RFC3339)
	if e.EndsAt != nil {
		e.EndsAtStr = e.EndsAt.Format("15:04 02.01.2006")
		e.EndsAtIso = e.EndsAt.Format("2006-01-02T15:04:05Z")
	}
}
