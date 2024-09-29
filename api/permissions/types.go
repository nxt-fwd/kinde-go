package permissions

// https://kinde.com/api/docs/#kinde-management-api-roles
type Permission struct {
	ID          string `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
