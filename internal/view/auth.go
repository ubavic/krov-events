package view

func (pe *PageExecutor) LoginPage(successful bool) error {
	data := make(map[string]any)
	data["successful"] = successful

	return pe.executePage("login.html", "Login", data)
}

func (pe *PageExecutor) LogoutPage() error {
	return pe.executePage("logout.html", "Logout", nil)
}
