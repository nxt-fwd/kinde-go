package users

import "time"

type User struct {
	ID                string     `json:"id"`
	ProvidedID       string     `json:"provided_id,omitempty"`
	PreferredEmail    string     `json:"preferred_email"`
	LastName         string     `json:"last_name"`
	FirstName        string     `json:"first_name"`
	IsSuspended      bool       `json:"is_suspended"`
	Picture          string     `json:"picture,omitempty"`
	CreatedOn        time.Time  `json:"created_on"`
	LastSignedIn     *time.Time `json:"last_signed_in,omitempty"`
	UpdatedOn        time.Time  `json:"updated_on"`
}

type ListResponse struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Users    []User `json:"users"`
	NextToken string `json:"next_token,omitempty"`
}

type CreateParams struct {
	Profile     Profile     `json:"profile"`
	Identities  []Identity  `json:"identities,omitempty"`
	Password    string      `json:"password,omitempty"`
	OrgCode     string      `json:"org_code,omitempty"`
}

type Profile struct {
	GivenName   string `json:"given_name"`
	FamilyName  string `json:"family_name"`
	Email       string `json:"email"`
	ProvidedID  string `json:"provided_id,omitempty"`
}

type Identity struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type CreateResponse struct {
	ID         string     `json:"id"`
	Created    bool       `json:"created"`
	Identities []Identity `json:"identities"`
}

type ListParams struct {
	PageSize  int    `json:"page_size,omitempty"`
	NextToken string `json:"next_token,omitempty"`
	Sort      string `json:"sort,omitempty"`
}

type UpdateParams struct {
	GivenName   string `json:"given_name,omitempty"`
	FamilyName  string `json:"family_name,omitempty"`
	ProvidedID  string `json:"provided_id,omitempty"`
	IsSuspended *bool  `json:"is_suspended,omitempty"`
}

type UpdateResponse struct {
	ID                       string  `json:"id"`
	Email                    *string `json:"email"`
	Picture                  *string `json:"picture"`
	GivenName               string  `json:"given_name"`
	FamilyName              string  `json:"family_name"`
	IsSuspended             bool    `json:"is_suspended"`
	IsPasswordResetRequested bool    `json:"is_password_reset_requested"`
}

type DeleteResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
} 
