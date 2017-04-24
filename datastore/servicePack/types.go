package servicePack

import (
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/price"
	"github.com/MerinEREN/iiPackages/datastore/score"
	"github.com/MerinEREN/iiPackages/datastore/video"
	"google.golang.org/appengine/datastore"
	"time"
)

// Type is remote or inPlace
// Status is underConsideration, disaproved, active, passive, changed, removed
// Pic is Person In Charge whom aprove this
// User key is Ancestor
type ServicePack struct {
	ID             string           `datastore:"-"`
	Type           string           `json:"type"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Duration       time.Duration    `json:"duration"`
	Price          price.Price      `json:"price"`
	Created        time.Time        `json:"created"`
	LastModified   time.Time        `json:"lastModified"`
	Status         string           `json:"status"`
	Pic            string           `json:"pic"`
	TagIDs         []*datastore.Key `json:"tagIDs"`
	Score          score.Score      `json:"score"`
	CustomerReview string           `json:"customerReview"`
	Photos         photo.Photos     `datastore: "-" json:"photos"`
	Videos         video.Videos     `datastore: "-" json:"videos"`
	// Extras         servicePackOption.ServicePackOptions `datastore: "-" json:"extras"`
}

type ServicePacks []ServicePack
