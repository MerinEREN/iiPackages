package user

import (
	// "github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/datastore/phone"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/tag"
	"google.golang.org/appengine/datastore"
	"time"
)

/*
Go's declaration syntax allows grouping of declarations. A single doc comment can introduce
a group of related constants or variables. Since the whole declaration is presented, such
a comment can often be perfunctory.
*/

// Account key is Ancestor
// UUID is Key Name
type User struct {
	ID                string           `datastore:"-"`
	Email             string           `json:"email"`
	Photo             photo.Photo      `datastore:"-" json:"photo"`
	Name              Name             `datastore: ",noindex" json:"name"`
	Gender            string           `json:"gender"`
	Status            string           `json:"status"`
	Type              string           `json:"type"`
	Roles             []string         `json:"roles"`
	TagIDs            []*datastore.Key `json:"-"`
	Tags              tag.Tags         `datastore: "-" json:"tags"`
	BirthDate         time.Time        `datastore: ",noindex" json:"birthDate"`
	Registered        time.Time        `datastore: ",noindex" json:"registered"`
	LastModified      time.Time        `datastore: ",noindex" json:"lastModified"`
	IsActive          bool             `json:"isactive"`
	PurchasedServices []*datastore.Key `datastore: "-" json:"purchasedSrvices"`
	Phones            phone.Phones     `datastore:"-" json:"phones"`
	// Password string `json:"password"`
	// Online, offline, frozen
	// User could be deactivated by superiors !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// Demands Demands `datastore: "-" json:"demands""`
	// Offers Offers `datastore: "-" json:"offers""`
	// ServicePacks ServicePacks `datastore: "-" json:"servicepacks""`
}

type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

type Users []*User

type Entity interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}
