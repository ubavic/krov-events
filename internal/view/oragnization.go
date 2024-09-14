package view

import (
	"decentrala.org/events/internal/types"
)

func (pe *PageExecutor) OrganizationList(organizations []types.Organization, selectedCity types.CityCode) error {
	data := make(map[string]any)
	data["cities"] = types.Cities
	data["organizations"] = organizations
	data["selectedCity"] = selectedCity

	return pe.executePage("organizationList.html", "Organizacije", data)
}

func (pe *PageExecutor) OrganizationPage(organization types.Organization, organizationEvents []types.Event) error {
	data := make(map[string]any)
	data["organization"] = organization
	data["events"] = organizationEvents

	return pe.executePage("organizationPage.html", organization.Name, data)
}

func (pe *PageExecutor) NewOrganization(organization types.Organization) error {
	data := make(map[string]any)
	data["cities"] = types.Cities
	data["organization"] = organization
	data["edit"] = false
	data["allowApi"] = organization.IsApiAllowed()

	return pe.executePage("organizationForm.html", "Nova organizacija", data)
}

func (pe *PageExecutor) EditOrganization(organization types.Organization, validationErrors []string) error {
	data := make(map[string]any)
	data["cities"] = types.Cities
	data["organization"] = organization
	data["edit"] = true
	data["allowApi"] = organization.IsApiAllowed()

	return pe.executePage("organizationForm.html", "Izmeni organizaciju "+organization.Name, data)
}
