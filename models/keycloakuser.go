package models

type KeycloakUser struct {
	Sub string `json:"sub"`
	EmailVerified bool `json:"email_verified"`
	Name string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Email string `json:"email"`
}