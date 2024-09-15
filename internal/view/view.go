package view

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"io/fs"

	"decentrala.org/events/internal/types"
)

type View struct {
	templates map[string]*template.Template
}

type PageExecutor struct {
	view   *View
	writer io.Writer
	reqCtx context.Context
}

func NewView(templateFs fs.FS) View {
	templates := make(map[string]*template.Template)

	files, err := fs.ReadDir(templateFs, ".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		pt, err := template.ParseFS(templateFs, file.Name(), "layout.html", "footer.html", "eventEntry.html")
		if err != nil {
			panic(err)
		}

		templates[file.Name()] = pt
	}

	return View{
		templates: templates,
	}
}

func (v *View) NewExecutor(writer io.Writer, reqCtx context.Context) PageExecutor {
	return PageExecutor{
		view:   v,
		writer: writer,
		reqCtx: reqCtx,
	}
}

func (pageExecutor *PageExecutor) executePage(templateName string, pageName string, data map[string]any) error {
	template, ok := pageExecutor.view.templates[templateName]
	if !ok {
		return fmt.Errorf("template %s not found", templateName)
	}

	wu, ok := pageExecutor.reqCtx.Value("user").(types.WebsiteUser)
	type userType struct {
		LoggedIn         bool
		Admin            bool
		Organization     string
		OrganizationCode string
	}

	user := userType{
		LoggedIn:         true,
		Admin:            true,
		Organization:     wu.Organization(),
		OrganizationCode: string(wu.OrganizationCode()),
	}

	if data == nil {
		data = make(map[string]any)
	}

	data["user"] = user

	pageData := struct {
		Title string
		Data  any
	}{
		Title: pageName,
		Data:  data,
	}

	err := template.Execute(pageExecutor.writer, pageData)
	return err
}

func (pe *PageExecutor) ErrorPage(w io.Writer, errorCode int) error {
	errorCodeString := fmt.Sprintf("%d", errorCode)

	data := make(map[string]any)
	data["ErrorCode"] = errorCodeString

	return pe.executePage("error.html", "Greška "+errorCodeString, data)
}
