package photo

import (
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Errors
var (
	ErrFindPhoto = errors.New("error while checking photo existincy")
)

// Get returns an entity via key and an error.
func Get(ctx context.Context, k *datastore.Key) (*Photo, error) {
	p := new(Photo)
	err := datastore.Get(ctx, k, p)
	return p, err
}

// GetViaParent returns projected entity and an error via parent key.
func GetViaParent(ctx context.Context, pk *datastore.Key) (*Photo, error) {
	p := new(Photo)
	q := datastore.NewQuery("Photo")
	q = q.Ancestor(pk)
	q = q.Project("Link")
	it := q.Run(ctx)
	// BUG !!!!! If i made this function as naked return "it.Next" fails because of "p"
	_, err := it.Next(p)
	if err == datastore.Done {
		return nil, err
	}
	if err != nil {
		err = ErrFindPhoto
		return nil, err
	}
	return p, nil
}
