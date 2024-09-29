package apis

// https://kinde.com/api/docs/#kinde-management-api-apis
type API struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Audience        string           `json:"audience"`
	IsManagementAPI bool             `json:"is_management_api"`
	Applications    []APIApplication `json:"applications,omitempty"`
}

type APIApplication struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive *bool  `json:"is_active"`
}
