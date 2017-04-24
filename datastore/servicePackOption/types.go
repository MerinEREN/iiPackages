package servicePackOption

import (
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/price"
	"github.com/MerinEREN/iiPackages/datastore/video"
	"time"
)

// Status is underConsideration, disaproved, active, passive, changed, removed
// Pic is Person In Charge whom aprove this
// ServicePack key is Ancestor
type ServicePackOption struct {
	ID           string        `datastore:"-"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Duration     time.Duration `json:"duration"`
	Price        price.Price   `json:"price"`
	Created      time.Time     `json:"created"`
	LastModified time.Time     `json:"lastModified"`
	Status       string        `json:"status"`
	Pic          string        `json:"pic"`
	Photos       photo.Photos  `datastore:"-" json:"photos"`
	Videos       video.Videos  `datastore:"-" json:"videos"`
	// Evaluation     Evaluation `json:"evaluation"`
	// CustomerReview string     `json:"customerreview"`
}

type ServicePackOptions []ServicePackOption
