package types

type WebsiteUser struct {
	loggedIn         bool
	admin            bool
	organizationCode OrganizationCode
	organization     string
}

func NewWebsiteUser(admin bool, orgCode OrganizationCode, orgName string) WebsiteUser {
	return WebsiteUser{
		loggedIn:         true,
		admin:            admin,
		organizationCode: orgCode,
		organization:     orgName,
	}
}

func NewRegularVisitor() WebsiteUser {
	return WebsiteUser{
		loggedIn:         false,
		admin:            false,
		organizationCode: "",
		organization:     "",
	}
}

func (we *WebsiteUser) Admin() bool {
	return we.admin
}

func (we *WebsiteUser) LoggedIn() bool {
	return we.loggedIn
}

func (we *WebsiteUser) OrganizationCode() OrganizationCode {
	return we.organizationCode
}

func (we *WebsiteUser) Organization() string {
	return we.organization
}
