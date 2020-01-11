/*
Package language "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package language

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

/*
GetAll returns all the entities in an order with some projections
from the begining of the kind.
*/
func GetAll(ctx context.Context) (Languages, error) {
	var lx []*Language
	q := datastore.NewQuery("Language")
	q = q.
		Project("ContextID").
		Order("-Created")
	kx, err := q.GetAll(ctx, &lx)
	if err != nil {
		return nil, err
	}
	ls := make(Languages)
	for i, v := range kx {
		lx[i].ID = v.StringID()
		ls[v.StringID()] = lx[i]
	}
	return ls, err
}
