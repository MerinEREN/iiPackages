/*
Package servicePack "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package servicePack

import (
	dstore "github.com/MerinEREN/iiPackages/datastore"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

// Get returns limited entities from the start of the kind with given filters and order.
// Add location filter here and use accounts addres info.
func Get(ctx context.Context, uTagIDs []*datastore.Key) (
	ServicePacks, datastore.Cursor, datastore.Cursor, error) {
	sps := make(ServicePacks)
	var cEnd datastore.Cursor
	q := datastore.NewQuery("ServicePack")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Limit(10)
	it := q.Run(ctx)
	cStart, err := it.Cursor()
	if err != nil {
		log.Printf("ServicePack Get Cursor Error: %v\n", err)
	}
	for {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			cEnd, err = it.Cursor()
			return sps, cStart, cEnd, err
		}
		if err != nil {
			return nil, cStart, cEnd, err
		}
		sp.ID = k.StringID()
		sps[sp.ID] = sp
	}
}

// GetNewest returns all the results from the begining to the previous start point
// with given filters and order.
// Add location filter here and use accounts addres info.
func GetNewest(ctx context.Context, cStart datastore.Cursor, uTagIDs []*datastore.Key) (
	ServicePacks, datastore.Cursor, error) {
	sps := make(ServicePacks)
	q := datastore.NewQuery("ServicePack")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		End(cStart)
	it := q.Run(ctx)
	cStart, err := it.Cursor()
	if err != nil {
		log.Printf("ServicePack GetNewest Cursor Error: %v\n", err)
	}
	for {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			return sps, cStart, err
		}
		if err != nil {
			return nil, cStart, err
		}
		sp.ID = k.StringID()
		sps[sp.ID] = sp
	}
}

// GetNewestCount returns the results count from the begining to the previous start point
// with given filters and order.
// Add location filter here and use accounts addres info.
func GetNewestCount(ctx context.Context, cStart datastore.Cursor, uTagIDs []*datastore.Key) (
	cnt int, err error) {
	q := datastore.NewQuery("ServicePack")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		End(cStart)
	cnt, err = q.Count(ctx)
	return
}

// GetOldest returns limited entities from the previous cEnd cursor
// with given filters and order.
// Add location filter here and use accounts addres info.
func GetOldest(ctx context.Context, cEnd datastore.Cursor, uTagIDs []*datastore.Key) (
	ServicePacks, datastore.Cursor, error) {
	sps := make(ServicePacks)
	q := datastore.NewQuery("ServicePack")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Start(cEnd).
		Limit(10)
	for it := q.Run(ctx); ; {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			cEnd, err = it.Cursor()
			return sps, cEnd, err
		}
		if err != nil {
			return nil, cEnd, err
		}
		sp.ID = k.StringID()
		sps[sp.ID] = sp
	}
}

// GetViaParent returns limited entities from the previous cursor
// with given filters and order.
func GetViaParent(ctx context.Context, crsr datastore.Cursor, pk *datastore.Key) (
	ServicePacks, datastore.Cursor, error) {
	sps := make(ServicePacks)
	q := datastore.NewQuery("ServicePack")
	q = q.Ancestor(pk).
		Order("-LastModified").
		Start(crsr).
		Limit(10)
	for it := q.Run(ctx); ; {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return sps, crsr, err
		}
		if err != nil {
			return nil, crsr, err
		}
		sp.ID = k.Encode()
		sps[sp.ID] = sp
	}
}
