package model

type ZitadelUserinfo struct {
	Email             string                 `json:"email"`
	EmailVerified     bool                   `json:"email_verified"`
	FamilyName        string                 `json:"family_name"`
	Gender            string                 `json:"gender"`
	GivenName         string                 `json:"given_name"`
	Locale            string                 `json:"locale"`
	Name              string                 `json:"name"`
	Nickname          string                 `json:"nickname"`
	PreferredUsername string                 `json:"preferred_username"`
	Sub               string                 `json:"sub"`
	UpdatedAt         int                    `json:"updated_at"`
	ProjectRoles      map[string]interface{} `json:"urn:zitadel:iam:org:project:roles"`
}
