package connections

// Strategy represents the type of connection
type Strategy string

const (
	// Built-in authentication strategies
	StrategyEmailPassword    Strategy = "email:password"    // Standard email/password authentication
	StrategyEmailOTP         Strategy = "email:otp"         // Email with one-time password
	StrategyPhoneOTP         Strategy = "phone:otp"         // Phone number with one-time password
	StrategyUsernamePassword Strategy = "username:password" // Username/password authentication
	StrategyUsernameOTP      Strategy = "username:otp"      // Username with one-time password

	// OAuth2 strategies
	StrategyOAuth2Apple     Strategy = "oauth2:apple"
	StrategyOAuth2AzureAD   Strategy = "oauth2:azure_ad"
	StrategyOAuth2Bitbucket Strategy = "oauth2:bitbucket"
	StrategyOAuth2Discord   Strategy = "oauth2:discord"
	StrategyOAuth2Facebook  Strategy = "oauth2:facebook"
	StrategyOAuth2Github    Strategy = "oauth2:github"
	StrategyOAuth2Gitlab    Strategy = "oauth2:gitlab"
	StrategyOAuth2Google    Strategy = "oauth2:google"
	StrategyOAuth2LinkedIn  Strategy = "oauth2:linkedin"
	StrategyOAuth2Microsoft Strategy = "oauth2:microsoft"
	StrategyOAuth2Patreon   Strategy = "oauth2:patreon"
	StrategyOAuth2Slack     Strategy = "oauth2:slack"
	StrategyOAuth2Stripe    Strategy = "oauth2:stripe"
	StrategyOAuth2Twitch    Strategy = "oauth2:twitch"
	StrategyOAuth2Twitter   Strategy = "oauth2:twitter"
	StrategyOAuth2Xero      Strategy = "oauth2:xero"

	// SAML strategy
	StrategySAMLCustom Strategy = "saml:custom"

	// WS-Federation strategy
	StrategyWSFedAzureAD Strategy = "wsfed:azure_ad"
)

// Options for different connection types

// SocialConnectionOptions represents options for OAuth2 social connections
type SocialConnectionOptions struct {
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	IsUseCustomDomain bool   `json:"is_use_custom_domain"`
}

// AzureADConnectionOptions represents options for Azure AD connections
type AzureADConnectionOptions struct {
	ClientID                     string   `json:"client_id"`
	ClientSecret                 string   `json:"client_secret"`
	HomeRealmDomains             []string `json:"home_realm_domains"`
	EntraIDDomain                string   `json:"entra_id_domain"`
	IsUseCommonEndpoint          bool     `json:"is_use_common_endpoint"`
	IsSyncUserProfileOnLogin     bool     `json:"is_sync_user_profile_on_login"`
	IsRetrieveProviderUserGroups bool     `json:"is_retrieve_provider_user_groups"`
	IsExtendedAttributesRequired bool     `json:"is_extended_attributes_required"`
}

// SAMLConnectionOptions represents options for SAML connections
type SAMLConnectionOptions struct {
	HomeRealmDomains      []string `json:"home_realm_domains"`
	SAMLEntityID          string   `json:"saml_entity_id"`
	SAMLASSURL            string   `json:"saml_acs_url"`
	SAMLIdpMetadataURL    string   `json:"saml_idp_metadata_url"`
	SAMLEmailKeyAttr      string   `json:"saml_email_key_attr"`
	SAMLFirstNameKeyAttr  string   `json:"saml_first_name_key_attr"`
	SAMLLastNameKeyAttr   string   `json:"saml_last_name_key_attr"`
	IsCreateMissingUser   bool     `json:"is_create_missing_user"`
	SAMLSigningCert       string   `json:"saml_signing_certificate"`
	SAMLSigningPrivateKey string   `json:"saml_signing_private_key"`
}
