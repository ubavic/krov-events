package types

type CityCode string

type CityName struct {
	Code CityCode
	Name string
}

var Cities []CityName = []CityName{
	{Code: "", Name: ""},
	{Code: "BG", Name: "Beograd"},
	{Code: "BO", Name: "Bor"},
	{Code: "ČA", Name: "Čačak"},
	{Code: "JA", Name: "Jagodina"},
	{Code: "KG", Name: "Kragujevac"},
	{Code: "KI", Name: "Kikinda"},
	{Code: "KŠ", Name: "Kruševac"},
	{Code: "KV", Name: "Kraljevo"},
	{Code: "LE", Name: "Leskovac"},
	{Code: "LO", Name: "Loznica"},
	{Code: "NI", Name: "Niš"},
	{Code: "NP", Name: "Novi Pazar"},
	{Code: "NS", Name: "Novi Sad"},
	{Code: "PA", Name: "Pančevo"},
	{Code: "PI", Name: "Pirot"},
	{Code: "PK", Name: "Prokuplje"},
	{Code: "PO", Name: "Požarevac"},
	{Code: "ŠA", Name: "Šabac"},
	{Code: "SD", Name: "Smederevo"},
	{Code: "SM", Name: "Sremska Mitrovica"},
	{Code: "SO", Name: "Sombor"},
	{Code: "SU", Name: "Subotica"},
	{Code: "UE", Name: "Užice"},
	{Code: "VA", Name: "Valjevo"},
	{Code: "VR", Name: "Vranje"},
	{Code: "VŠ", Name: "Vršac"},
	{Code: "ZA", Name: "Zaječar"},
	{Code: "ZR", Name: "Zrenjanin"},
}

func GetCityName(id CityCode) string {
	for _, city := range Cities {
		if city.Code == id {
			return city.Name
		}
	}

	return ""
}
