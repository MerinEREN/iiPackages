package tag

type Tag struct {
	ID     string            `datastore:"-"`
	Values map[string]string `json:"values"`
}

type Tags []*Tag
