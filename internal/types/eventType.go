package types

type EventType int

const (
	EvenTypeNone EventType = iota
	EvenTypeOther
	EvenTypeMusic
	EvenTypeExhibition
	EvenTypeTheater
	EventTypeMovie
	EventTypeLecture
	EventTypeDiscussion
	EventTypeWorkshop
	EvenTypeHackathon
)

type EventTypeName struct {
	Value EventType
	Name  string
}

var EventTypes = []EventTypeName{
	{
		Value: EvenTypeNone,
		Name:  "",
	},
	{
		Value: EvenTypeMusic,
		Name:  "Muzika",
	},
	{
		Value: EvenTypeExhibition,
		Name:  "Izložba",
	},
	{
		Value: EvenTypeTheater,
		Name:  "Teatar",
	},
	{
		Value: EventTypeMovie,
		Name:  "Film",
	},
	{
		Value: EventTypeLecture,
		Name:  "Predavanje",
	},
	{
		Value: EventTypeDiscussion,
		Name:  "Diskusija",
	},
	{
		Value: EventTypeWorkshop,
		Name:  "Radionica",
	},
	{
		Value: EvenTypeHackathon,
		Name:  "Hakaton",
	},
	{
		Value: EvenTypeOther,
		Name:  "Ostalo",
	},
}

func GetEventTypeName(id EventType) string {
	for _, t := range EventTypes {
		if t.Value == id {
			return t.Name
		}
	}

	return ""
}
