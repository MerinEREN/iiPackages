package phone

// User key is Anchestor
type Phone struct {
	ID          string `datastore:"-"`
	Type        string `json:"stype"`
	CountryCode string `json:countryCode`
	Number      string `json:number`
	Exstension  string `json:exstension`
}

type Phones map[string]*Phone
