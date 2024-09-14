package types

type WebsiteUser struct {
	LoggedIn         bool
	Admin            bool
	OrganizationCode OrganizationCode
	Organization     string
}
