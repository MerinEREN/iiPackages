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

// GetFilteredByAncestorLimited returns the limited entities filtered by ancestor.
func GetFilteredByAncestorLimited(ctx context.Context, pk *datastore.Key, lim int) (
	Photos, error) {
	var px []*Photo
	q := datastore.NewQuery("Photo")
	q = q.
		Ancestor(pk).
		Project("Link")
	if lim > 0 && lim < 40 {
		q = q.Limit(lim)
	} else {
		q = q.Limit(12)
	}
	kx, err := q.GetAll(ctx, &px)
	if err != nil {
		return nil, err
	}
	ps := make(Photos)
	for i, v := range kx {
		px[i].ID = v.Encode()
		ps[v.Encode()] = px[i]
	}
	return ps, err
}
