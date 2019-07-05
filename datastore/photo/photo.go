package photo

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetMainByAncestor returns the entity with main type and an error.
func GetMainByAncestor(ctx context.Context, pType string, pk *datastore.Key) (
	Photo, error) {
	q := datastore.NewQuery("Photo")
	q = q.
		Ancestor(pk).
		Filter("Type =", pType).
		Project("Link")
	it := q.Run(ctx)
	var p Photo
	k, err := it.Next(&p)
	p.ID = k.Encode()
	return p, err
}

/*
// GetViaParent returns the entities projected via parent key and an error.
func GetViaParent(ctx context.Context, pk *datastore.Key, after, lim string) (
	Photos, error) {
	ps := make(Photos)
	q := datastore.NewQuery("Photo")
	q = q.
		Ancestor(pk).
		Order("-LastModified")
		Project("Link")
	for it := q.Run(ctx); ; {
		p := new(Photo)
		// BUG !!!!! If i made this function as naked return "it.Next" fails because of "p"
		_, err := it.Next(p)
		if err == datastore.Done {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
*/
