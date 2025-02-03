package roles

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	IsDefaultRole bool   `json:"is_default_role"`
}

type ListResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Roles   []Role `json:"roles"`
}

type CreateParams struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type CreateResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Role    Role   `json:"role"`
}

type UpdateParams struct {
	Name        string   `json:"name,omitempty"`
	Key         string   `json:"key,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type UpdateResponse CreateResponse

type DeleteResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type UpdatePermissionItem struct {
	ID        string `json:"id"`
	Operation string `json:"operation,omitempty"`
}

type UpdatePermissionsParams struct {
	Permissions []UpdatePermissionItem `json:"permissions"`
}

type UpdatePermissionsResponse struct {
	Code               string   `json:"code"`
	Message            string   `json:"message"`
	PermissionsAdded   []string `json:"permissions_added,omitempty"`
	PermissionsRemoved []string `json:"permissions_removed,omitempty"`
}

type RemovePermissionResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Key         string `json:"key"`
}

type ListPermissionsResponse struct {
	Code        string       `json:"code"`
	Message     string       `json:"message"`
	NextToken   string       `json:"next_token"`
	Permissions []Permission `json:"permissions,omitempty"`
} 
