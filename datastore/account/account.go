/*
Package account has account struct type and datastore query functions.
*/
package account

import (
	"errors"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
)

// Errors
var (
	ErrPutAccount = errors.New("error while putting account into the datastore")
)

// Get returns entity via key or key name and an error.
func Get(s *session.Session, k interface{}) (*Account, error) {
	acc := new(Account)
	var err error
	switch v := k.(type) {
	case string:
		q := datastore.NewQuery("Account").
			Filter("KeyName =", v)
		it := q.Run(s.Ctx)
		_, err = it.Next(acc)
		if err == nil {
			acc.ID = v
		}
	case *datastore.Key:
		err = datastore.Get(s.Ctx, v, acc)
		if err == nil {
			acc.ID = v.StringID()
		}
	}
	return acc, err
}

// Delete deletes an entity via key and returns an error.
func Delete(s *session.Session, k *datastore.Key) error {
	err := datastore.Delete(s.Ctx, k)
	return err
}

/* func AddTags(s ...string) bool {
	return
} */
