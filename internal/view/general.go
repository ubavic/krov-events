package view

func (pe *PageExecutor) AboutPage() error {
	return pe.executePage("about.html", "O projektu", nil)
}
