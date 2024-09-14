package types

import (
	"strings"
	"time"
)

type OrganizationCode string

type Organization struct {
	Code        OrganizationCode
	Name        string
	Description string
	Website     string
	Email       string
	token       *string
	apiAllowed  *bool
	admin       *bool
	Address     string
	CityCode    CityCode
	CityName    string
	OsmUrl      string
	CreatedAt   time.Time
	ModifiedAt  *time.Time
}

func (org *Organization) SetApiAllowed(apiAllowed bool) {
	if org.apiAllowed == nil {
		allowed := apiAllowed
		org.apiAllowed = &allowed
	}
}

func (org *Organization) IsApiAllowed() bool {
	if org.apiAllowed != nil {
		return *org.apiAllowed
	}

	return false
}

func (org *Organization) SetToken(token string) {
	if org.token == nil {
		newToken := strings.Clone(token)
		org.token = &newToken
	}
}

func (org *Organization) ValidateToken(token string) bool {
	if org.token == nil {
		return false
	}

	return strings.Compare(*org.token, token) == 0
}

func (org *Organization) SetAdmin(admin bool) {
	if org.admin == nil {
		a := admin
		org.admin = &a
	}
}

func (org *Organization) IsAdmin() bool {
	if org.admin != nil {
		return *org.admin
	}

	return false
}
