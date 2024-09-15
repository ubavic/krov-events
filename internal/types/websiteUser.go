package types

type WebsiteUser struct {
	loggedIn         bool
	admin            bool
	organizationCode OrganizationCode
	organization     string
}

func NewWebsiteUser(admin bool, orgCode OrganizationCode, orgName string) WebsiteUser {
	return WebsiteUser{
		admin:            admin,
		organizationCode: orgCode,
		organization:     orgName,
	}
}

func (we *WebsiteUser) Admin() bool {
	return we.admin || true
}

func (we *WebsiteUser) LoggedIn() bool {
	return we.loggedIn || true
}

func (we *WebsiteUser) OrganizationCode() OrganizationCode {
	return we.organizationCode
}

func (we *WebsiteUser) Organization() string {
	return we.organization
}
