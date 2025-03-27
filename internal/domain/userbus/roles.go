package userbus

import "fmt"

var (
	RoleUser  = newRole("USER")
	RoleAdmin = newRole("ADMIN")
)

// set of known roles to this app.
var roles = make(map[string]Role)

// Role represents a role in our application,since we require some validation over it we created a type.
type Role struct {
	value string
}

func newRole(value string) Role {
	r := Role{value: value}
	roles[r.value] = r
	return r
}

func (r Role) String() string {
	return r.value
}

func (r Role) Equal(role Role) bool {
	return r.value == role.value
}

func (r Role) MarshalText() (text []byte, err error) {
	return []byte(r.value), nil
}

func ParseRole(value string) (Role, error) {
	r, ok := roles[value]
	if !ok {
		return Role{}, fmt.Errorf("invalid role: %q", value)
	}
	return r, nil
}

func ParseSliceOfRoles(roles []string) ([]Role, error) {
	parsed := make([]Role, len(roles))

	for i, str := range roles {
		role, err := ParseRole(str)
		if err != nil {
			return nil, err
		}
		parsed[i] = role
	}
	return parsed, nil
}

func EncodeRoles(roles []Role) []string {
	encoded := make([]string, len(roles))
	for i, role := range roles {
		encoded[i] = role.String()
	}

	return encoded
}
