/*
Package servicePack "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package servicePack

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

/*
GetMulti returns the corresponding entities as a map with their "ID"s assigned
by their keys and an error.
*/
func GetMulti(ctx context.Context, kx []*datastore.Key) (ServicePacks, error) {
	var spx []*ServicePack
	err := datastore.GetMulti(ctx, kx, spx)
	if err != nil {
		return nil, err
	}
	sps := make(map[string]*ServicePack)
	for i, v := range spx {
		v.ID = kx[i].Encode()
		sps[v.ID] = v
	}
	return sps, err
}

/*
GetNextByParentLimited returns limited entities filtered by user key as parent
from the previous cursor in an order.
Also, returns the next cursor as string and an error.
*/
func GetNextByParentLimited(ctx context.Context, crsrAsString string, pk *datastore.Key,
	lim int) (ServicePacks, string, error) {
	var err error
	var after datastore.Cursor
	sps := make(ServicePacks)
	q := datastore.NewQuery("ServicePack")
	q = q.Ancestor(pk).
		Order("-LastModified")
	if crsrAsString != "" {
		after, err = datastore.DecodeCursor(crsrAsString)
		if err != nil {
			return nil, crsrAsString, err
		}
		q = q.Start(after)
	}
	if lim > 0 && lim < 40 {
		q = q.Limit(lim)
	} else {
		q = q.Limit(12)
	}
	for it := q.Run(ctx); ; {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			after, err = it.Cursor()
			return sps, after.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		sp.ID = k.Encode()
		sps[sp.ID] = sp
	}
}

// GetAllLimited returns limited entities from the begining of the kind with given limit
// and an error.
func GetAllLimited(ctx context.Context, lim int) ([]ServicePack, error) {
	var spx []ServicePack
	q := datastore.NewQuery("ServicePack")
	q = q.Order("-LastModified").
		Project("TagIDs").
		Limit(lim)
	_, err := q.GetAll(ctx, &spx)
	return spx, err
}
