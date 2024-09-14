package controller

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func (controller *Controller) serverError(statusCode int, err error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(statusCode, err, r.URL)

		pe := controller.view.NewExecutor(w, r.Context())
		pe.ErrorPage(w, statusCode)
	}
}

func parseDate(str string) *time.Time {
	if str == "" {
		return nil
	}

	date, err := time.Parse(time.DateOnly, str)
	if err != nil {
		return nil
	}

	return &date
}

func parsDateTime(str string) *time.Time {
	if str == "" {
		return nil
	}

	date, err := time.Parse("2006-01-02T15:04", str)
	if err != nil {
		return nil
	}

	return &date
}

func formatDate(date *time.Time) string {
	if date == nil {
		return ""
	}

	return date.Format(time.DateOnly)
}

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
