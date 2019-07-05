/*
Package account has account struct type and datastore query functions.
*/
package account

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Errors
var (
	ErrPutAccount = errors.New("error while putting account into the datastore")
)

// Get returns the entity via key or key name and an error.
func Get(ctx context.Context, i interface{}) (*Account, error) {
	acc := new(Account)
	var err error
	switch v := i.(type) {
	case string:
		q := datastore.NewQuery("Account").
			Filter("KeyName =", v)
		it := q.Run(ctx)
		k := new(datastore.Key)
		k, err = it.Next(acc)
		if err == nil {
			acc.ID = k.Encode()
		}
	case *datastore.Key:
		err = datastore.Get(ctx, v, acc)
		if err == nil {
			acc.ID = v.Encode()
		}
	}
	return acc, err
}

/* func AddTags(s ...string) bool {
	return
} */
