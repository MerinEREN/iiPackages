package rank

// beginner, master...
type Rank struct {
	ID     string            `datastore:"-"`
	Values map[string]string `json:"values"`
}

type Ranks []*Rank
