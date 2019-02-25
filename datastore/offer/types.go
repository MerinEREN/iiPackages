package offer

import (
	// "github.com/MerinEREN/iiPackages/datastore/price"
	"time"
)

// Offer is the struct for user offers.
// INFORM DEMAND OWNER WHEN AN OFFER MODIFIED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Status is available, accepted, notAccepted, changed, removed, successful,
// unsuccessful.
// backup (ONLY AUTHORIZED ACCOUNTS WHO ACCEPTED TO BE BACKUP) !!!!!!!!!!!!!!!!!!!!
// Demand key is the ancestor key.
// AccountID is used to link the account.
type Offer struct {
	ID           string    `datastore:"-"`
	UserID       string    `json:"userID"`
	AccountID    string    `datastore:"-" json:"accountID"`
	Explanation  string    `datastore:",noindex" json:"explanation"`
	Amount       float64   `datastore:",noindex" json:"amount"`
	Created      time.Time `datastore:",noindex" json:"created"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
	/* Price        price.Price `datastore:",noindex" json:"price"`
	StartTime      time.Time        `datastore:",noindex" json:"startTime"`
	Duration       string           `datastore:",noindex" json:"duration"`
	CustomerReview string           `json:"customerreview"`
	Score          score.Score      `json:"score"` */
}

// Offers is map[string]*Offer.
type Offers map[string]*Offer
