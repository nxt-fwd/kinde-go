package organizations

import "time"

type Color struct {
	Hex string `json:"hex"`
	Raw string `json:"raw"`
}

type Organization struct {
	Code                 string    `json:"code"`
	Logo                 *string   `json:"logo,omitempty"`
	Name                 string    `json:"name"`
	Handle               *string   `json:"handle,omitempty"`
	LogoDark             *string   `json:"logo_dark,omitempty"`
	CreatedOn            time.Time `json:"created_on,omitempty"`
	IsDefault            bool      `json:"is_default"`
	LinkColor            *Color    `json:"link_color,omitempty"`
	ThemeCode            *string   `json:"theme_code,omitempty"`
	ExternalID           *string   `json:"external_id,omitempty"`
	FaviconSvg          *string   `json:"favicon_svg,omitempty"`
	ButtonColor          *Color    `json:"button_color,omitempty"`
	ColorScheme          string    `json:"color_scheme"`
	LinkColorDark        *Color    `json:"link_color_dark,omitempty"`
	BackgroundColor      *Color    `json:"background_color,omitempty"`
	FaviconFallback      *string   `json:"favicon_fallback,omitempty"`
	ButtonColorDark      *Color    `json:"button_color_dark,omitempty"`
	ButtonTextColor      *Color    `json:"button_text_color,omitempty"`
	CardBorderRadius     *string   `json:"card_border_radius,omitempty"`
	InputBorderRadius    *string   `json:"input_border_radius,omitempty"`
	ButtonBorderRadius   *string   `json:"button_border_radius,omitempty"`
	BackgroundColorDark  *Color    `json:"background_color_dark,omitempty"`
	ButtonTextColorDark  *Color    `json:"button_text_color_dark,omitempty"`
	IsAllowRegistrations *bool     `json:"is_allow_registrations,omitempty"`
	IsAutoMembershipEnabled bool      `json:"is_auto_membership_enabled"`
}

type ListResponse struct {
	Code          string         `json:"code"`
	Message       string         `json:"message"`
	Organizations []Organization `json:"organizations"`
	NextToken     string         `json:"next_token"`
}

type CreateParams struct {
	Name       string `json:"name"`
	Code       string `json:"code,omitempty"`
	Handle     string `json:"handle,omitempty"`
	IsPersonal bool   `json:"is_personal,omitempty"`
}

type CreateResponse struct {
	Code         string       `json:"code"`
	Message      string       `json:"message"`
	Organization Organization `json:"organization"`
}

type GetResponse CreateResponse

type UpdateParams struct {
	Name                          string   `json:"name,omitempty"`
	ExternalID                    string   `json:"external_id,omitempty"`
	BackgroundColor              string   `json:"background_color,omitempty"`
	ButtonColor                  string   `json:"button_color,omitempty"`
	ButtonTextColor             string   `json:"button_text_color,omitempty"`
	LinkColor                   string   `json:"link_color,omitempty"`
	BackgroundColorDark         string   `json:"background_color_dark,omitempty"`
	ButtonColorDark             string   `json:"button_color_dark,omitempty"`
	ButtonTextColorDark        string   `json:"button_text_color_dark,omitempty"`
	LinkColorDark              string   `json:"link_color_dark,omitempty"`
	ThemeCode                  string   `json:"theme_code,omitempty"`
	Handle                     string   `json:"handle,omitempty"`
	IsCustomAuthConnectionsEnabled bool     `json:"is_custom_auth_connections_enabled,omitempty"`
	IsAutoJoinDomainList      bool       `json:"is_auto_join_domain_list,omitempty"`
	AllowedDomains            []string `json:"allowed_domains,omitempty"`
	IsEnableAdvancedOrgs      bool       `json:"is_enable_advanced_orgs,omitempty"`
	IsEnforceMfa              bool       `json:"is_enforce_mfa,omitempty"`
}

type UpdateResponse CreateResponse

type DeleteResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AddUser struct {
	ID          string   `json:"id"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type AddUsersParams struct {
	Users []AddUser `json:"users"`
}

type AddUsersResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Role represents a role in an organization
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Description string `json:"description"`
} 