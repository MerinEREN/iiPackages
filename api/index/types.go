package index

import (
	"github.com/MerinEREN/iiPackages/datastore/account"
	"github.com/MerinEREN/iiPackages/datastore/user"
)

// Properties has to be kapitalized
// Otherwise they they can't be accessable at the client side.
type userAccount struct {
	User    map[string]*user.User       `json:"user"`
	Account map[string]*account.Account `json:"account"`
}
