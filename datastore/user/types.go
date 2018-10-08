package user

import (
	"github.com/MerinEREN/iiPackages/datastore/phone"
	"google.golang.org/appengine/datastore"
	"time"
)

/*
Go's declaration syntax allows grouping of declarations. A single doc comment can introduce
a group of related constants or variables. Since the whole declaration is presented, such
a comment can often be perfunctory.
*/

// User account key is Ancestor
type User struct {
	ID                string           `datastore:"-"`
	Email             string           `json:"email"`
	Name              Name             `json:"name"`
	Link              string           `json:"link"`
	Gender            string           `json:"gender"`
	Status            string           `json:"status"`
	Type              string           `json:"type"`
	Roles             []string         `json:"roles"`
	BirthDate         time.Time        `datastore:",noindex" json:"birthDate"`
	Created           time.Time        `json:"created"`
	LastModified      time.Time        `datastore:",noindex" json:"lastModified"`
	IsActive          bool             `json:"isactive"`
	PurchasedServices []*datastore.Key `datastore:"-" json:"purchasedSrvices"`
	Phones            phone.Phones     `datastore:"-" json:"phones"`
	// Password string `json:"password"`
	// Online, offline, frozen
	// User could be deactivated by superiors !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// Demands Demands `datastore:"-" json:"demands""`
	// Offers Offers `datastore:"-" json:"offers""`
	// ServicePacks ServicePacks `datastore:"-" json:"servicepacks""`
}

// Name is user name struct with "First" and "Last" fields.
type Name struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// Users is a users map with encoded user key as map key.
type Users map[string]*User

// Entity interface to implement all structs i guess.
type Entity interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}
