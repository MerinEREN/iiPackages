package photo

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Errors
var (
	ErrFindPhoto = errors.New("Error while checking photo existincy.")
)

func Get(ctx context.Context, k *datastore.Key) (*Photo, *datastore.Key, error) {
	p := new(Photo)
	q := datastore.NewQuery("Photo").Ancestor(k)
	it := q.Run(ctx)
	// BUG !!!!! If i made this function as naked return "it.Next" fails because of "u"
	k, err := it.Next(p)
	if err == datastore.Done {
		return nil, nil, err
	}
	if err != nil {
		err = ErrFindPhoto
		return nil, nil, err
	}
	return p, k, nil
}
