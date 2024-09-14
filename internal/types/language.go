package types

type Language int

const (
	None Language = iota
	Serbian
	English
)

var Languages = []struct {
	Language Language
	Name     string
}{
	{Language: None, Name: ""},
	{Language: Serbian, Name: "Srpski"},
	{Language: English, Name: "Engleski"},
}
