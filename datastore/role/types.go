package role

type Role struct {
	ID     string            `datastore: "" json:"id"`
	Values map[string]string `datastore: "" json:"values"`
}

type Roles map[string]*Role
