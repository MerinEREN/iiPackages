package demand

import (
	"github.com/MerinEREN/iiPackages/datastore/address"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/price"
	"github.com/MerinEREN/iiPackages/datastore/video"
	"google.golang.org/appengine/datastore"
	"time"
)

// Type is remote or inPlace
// IF THERE IS A WAY TO CREATE AN ENTITY WITH TWO ANCESTOR REMOVE Offers PROPERTY !
// IF THERE IS AT LEAST ONE OFFER DO NOT LET USER TO CHANGE DEMAND !!!!!!!!!!!!!!!!
// Status is underConsideration, active, rejected, changed, removed, finished,
// disaproved.
// Pic is Person In Charge whom aprove this
// User key is Ancestor
type Demand struct {
	ID           string            `datastore:"-"`
	Type         string            `json:"type"`
	Title        string            `datastore: ",noindex" json:"title"`
	Description  string            `datastore: ",noindex" json:"description"`
	StartTime    time.Time         `json:"startTime"`
	EndTime      time.Time         `json:"endTime"`
	Price        price.Price       `json:"price"`
	Addresses    address.Addresses `json:"addresses"`
	Created      time.Time         `json:"created"`
	LastModified time.Time         `datastore: ",noindex" json:"lastModified"`
	Status       string            `json:"status"`
	Pic          string            `json:"pic"`
	TagIDs       []*datastore.Key  `json:"tagIDs"`
	Photos       photo.Photos      `datastore: "-" json:"photos"`
	Videos       video.Videos      `datastore: "-" json:"videos"`
	// Tags           Tags              `datastore: "-" json:"tags"`
}

type Demands map[string]*Demand
