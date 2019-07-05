package account

import (
	"github.com/MerinEREN/iiPackages/datastore/address"
	"github.com/MerinEREN/iiPackages/datastore/rank"
	"github.com/MerinEREN/iiPackages/datastore/score"
	"google.golang.org/appengine/datastore"
	"time"
)

// type Accounts []Account

// Account is the struct for register accounts.
// Hide name when sending.
// Type values are "individual" and "corporate".
type Account struct {
	ID           string            `datastore:"-"`
	Name         string            `json:"name"`
	Type         string            `json:"type"`
	Addresses    address.Addresses `datastore:"-" json:"addresses"`
	Status       string            `json:"status"`
	About        string            `json:"about"`
	Registered   time.Time         `json:"registered"`
	LastModified time.Time         `json:"lastModified"`
	Score        score.Score       `datastore:"-" json:"score"`
	RankIDs      []*datastore.Key  `json:"rankIDs"`
	Ranks        rank.Ranks        `datastore:"-" json:"ranks"`
	BankAccounts []BankAccount     `json:"bankAccount" valid:"bankAccount"`
}

// Accounts is a map[string]*Account.
type Accounts map[string]*Account

// BankAccount is the struct for store accounts bank account infos..
type BankAccount struct {
	IMEI string `json:"IMEI"`
}

// Entity USE THIS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Entity interface {
	// Use this for all structs
	// Update()
	// Upsert()
	// Delete()
}
