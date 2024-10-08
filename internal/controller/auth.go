package controller

import (
	"net/http"
	"time"

	"decentrala.org/events/internal/types"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Organization     string `json:"organization"`
	OrganizationCode string `json:"organizationCode"`
	IsAdmin          bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

func (controller *Controller) login(w http.ResponseWriter, r *http.Request) {
	orgParam := r.URL.Query().Get("org")
	tokenParam := r.URL.Query().Get("token")

	pe := controller.view.NewExecutor(w, r.Context())
	if len(orgParam) < 3 || len(tokenParam) < 30 {
		pe.LoginPage(false)
		return
	}

	organization, err := controller.model.GetOrganization(types.OrganizationCode(orgParam))
	if err != nil {
		pe.LoginPage(false)
		return
	}

	if !organization.ValidateToken(tokenParam) {
		pe.LoginPage(false)
		return
	}

	expirationTime := time.Now().Add(300 * time.Minute)
	claims := &Claims{
		Organization:     organization.Name,
		OrganizationCode: string(organization.Code),
		IsAdmin:          true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(controller.key)
	if err != nil {
		pe.LoginPage(false)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   3600,
	})

	pe.LoginPage(true)
}

func (controller *Controller) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})

	pe := controller.view.NewExecutor(w, r.Context())
	pe.LogoutPage()
}

func (controller *Controller) verifyToken(tokenString string) types.WebsiteUser {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return controller.key, nil
	})

	invalidUser := types.WebsiteUser{}

	if err != nil {
		return invalidUser
	}

	if !token.Valid {
		return invalidUser
	}

	var claims Claims

	return types.NewWebsiteUser(
		true,
		"dmz",
		claims.Organization,
	)
}
