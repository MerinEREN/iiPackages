package account

import (
	//"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/datastore/address"
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/rank"
	"github.com/MerinEREN/iiPackages/datastore/score"
	"google.golang.org/appengine/datastore"
	"time"
)

// type Accounts []Account

// Hide name when sending.
type Account struct {
	ID           string            `datastore:"-"`
	Photo        photo.Photo       `datastore:"-" json:"photo"`
	Name         string            `json:"name"`
	Addresses    address.Addresses `json:"addresses"`
	Status       string            `json:"status"`
	About        string            `json:"about"`
	Score        score.Score       `datastore:"-" json:"score"`
	Registered   time.Time         `json:"registered"`
	LastModified time.Time         `json:"lastModified"`
	RankIDs      []*datastore.Key  `json:"rankIDs"`
	Ranks        rank.Ranks        `datastore:"-" json:"ranks"`
	BankAccounts []BankAccount     `json:"bankAccount" valid:"bankAccount"`
}

type Accounts map[string]*Account

type BankAccount struct {
	IMEI string `json:"IMEI"`
}

type Entity interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}
