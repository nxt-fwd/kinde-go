package applications

import "github.com/axatol/kinde-go/internal/enum"

// https://kinde.com/api/docs/#kinde-management-api-applications
type Application struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Type         Type     `json:"type"`
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	LoginURI     string   `json:"login_uri"`
	HomepageURI  string   `json:"homepage_uri"`
	LogoutURIs   []string `json:"logout_uris"`
	RedirectURIs []string `json:"redirect_uris"`
}

var _ enum.Enum[Type] = (*Type)(nil)

type Type string

const (
	TypeRegular               Type = "reg"
	TypeSinglePageApplication Type = "spa"
	TypeMachineToMachine      Type = "m2m"
)

func (t Type) Options() []Type {
	return []Type{
		TypeRegular,
		TypeSinglePageApplication,
		TypeMachineToMachine,
	}
}

func (t Type) Valid() error {
	return enum.Valid(t.Options(), t)
}
