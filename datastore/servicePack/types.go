package servicePack

import (
	"time"
)

// ServicePack attributes are;
// Type is remote or inPlace
// Status is underConsideration, disaproved, active, passive, changed, removed
// PIC is Person In Charge whom aprove this.
// User key is the ancestor.
type ServicePack struct {
	ID           string    `datastore:"-"`
	UserID       string    `datastore:"-" json:"userID"`
	AccountID    string    `datastore:"-" json:"accountID"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
	/* LinksVideo   []string  `json:"linksVideo"`
	Duration       time.Duration    `json:"duration"`
	Price          price.Price      `json:"price"`
	PIC            string           `json:"pic"`
	Score          score.Score      `json:"score"`
	CustomerReview string           `json:"customerReview"` */
	// Extras         servicePackOption.ServicePackOptions `datastore: "-" json:"extras"`
}

// ServicePacks is a *ServicePack map.
type ServicePacks map[string]*ServicePack
