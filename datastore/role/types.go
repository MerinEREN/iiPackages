package role

// Role represents user roles.
type Role struct {
	ID    string `datastore:"-"`
	Value string `json:"value"`
}

// Roles is a map of role pointers with role IDs as their keys.
type Roles map[string]*Role
