package address

import (
	"google.golang.org/appengine"
)

type Address struct {
	Description string             `datastore:",noindex" json:"description"`
	Borough     string             `datastore:",noindex" json:"borough"`
	City        string             `datastore:"" json:"city"`
	Country     string             `datastore:"" json:"country"`
	Postcode    string             `datastore:",noindex" json:"postcode"`
	GeoPoint    appengine.GeoPoint `datastore:"",noindex json:"geoPoint"`
}

type Addresses []Address
