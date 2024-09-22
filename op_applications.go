package kinde

// https://kinde.com/api/docs/#kinde-management-api-applications
type Application struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
