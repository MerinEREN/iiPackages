package demand

import (
	/* "github.com/MerinEREN/iiPackages/datastore/address"
	"github.com/MerinEREN/iiPackages/datastore/price" */
	"time"
)

// Demand attributes are;
// Type is remote or inPlace
// IF THERE IS AT LEAST ONE OFFER DO NOT LET USER TO CHANGE DEMAND !!!!!!!!!!!!!!!!
// Status is underConsideration, active, rejected, changed, removed, finished,
// disaproved.
// PIC is Person In Charge whom aprove this.
// User key is the ancestor.
// UserID and AccountID are used to if controles at frontend to be able to modify a demand
// and create an offer.
type Demand struct {
	ID           string    `datastore:"-"`
	UserID       string    `datastore:"-" json:"userID"`
	AccountID    string    `datastore:"-" json:"accountID"`
	Description  string    `datastore:",noindex" json:"description"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
	/* Type         string            `json:"type"`
	StartTime    time.Time         `json:"startTime"`
	EndTime      time.Time         `json:"endTime"`
	Price        price.Price       `json:"price"`
	Addresses    address.Addresses `json:"addresses"`
	LinksVideo   []string          `json:"linksVideo"`
	PIC          string            `json:"pic"`*/
}

// Demands is a *Demand map.
type Demands map[string]*Demand
