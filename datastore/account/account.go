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

// Get returns the entity via key or key name and an error.
func Get(s *session.Session, i interface{}) (*Account, error) {
	acc := new(Account)
	var err error
	switch v := i.(type) {
	case string:
		q := datastore.NewQuery("Account").
			Filter("KeyName =", v)
		it := q.Run(s.Ctx)
		k := new(datastore.Key)
		k, err = it.Next(acc)
		if err == nil {
			acc.ID = k.Encode()
		}
	case *datastore.Key:
		err = datastore.Get(s.Ctx, v, acc)
		if err == nil {
			acc.ID = v.Encode()
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
