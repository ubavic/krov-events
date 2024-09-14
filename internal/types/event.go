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
	EndsAt           *time.Time
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