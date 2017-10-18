package address

import (
	"google.golang.org/appengine"
)

// Address struct for account addresses.
type Address struct {
	Description string             `datastore:",noindex" json:"description"`
	Borough     string             `datastore:",noindex" json:"borough"`
	City        string             `json:"city"`
	Country     string             `json:"country"`
	Postcode    string             `datastore:",noindex" json:"postcode"`
	GeoPoint    appengine.GeoPoint `datastore:",noindex" json:"geoPoint"`
}

// Addresses is map[string]*Address.
type Addresses map[string]*Address
