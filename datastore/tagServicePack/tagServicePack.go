/*
Package tagServicePack "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
7ne will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package tagServicePack

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysByServicePackOrTagKey returns the servicePack keys as a slice if the tag key is provided
// or returns the tag keys as a slice if servicePack key is provided and also an error.
func GetKeysByServicePackOrTagKey(ctx context.Context, key *datastore.Key) (
	[]*datastore.Key, error) {
	q := datastore.NewQuery("TagServicePack")
	kind := key.Kind()
	switch kind {
	case "Tag":
		q = q.
			Filter("TagKey =", key).
			KeysOnly()
		return q.GetAll(ctx, nil)
	default:
		// For "ServicePack" kind
		var ktx []*datastore.Key
		q = q.
			Ancestor(key).
			KeysOnly()
		kx, err := q.GetAll(ctx, nil)
		if err != nil {
			return nil, err
		}
		for _, v := range kx {
			kt, err := datastore.DecodeKey(v.StringID())
			if err != nil {
				return nil, err
			}
			ktx = append(ktx, kt)
		}
		return ktx, nil
	}
}

// GetDistinctLatestLimited returns distinct entities from limited results
// and also an error.
func GetDistinctLatestLimited(ctx context.Context, lim int) ([]TagServicePack, error) {
	var tdx []TagServicePack
	q := datastore.NewQuery("TagServicePack")
	q = q.Order("-Created").
		Project("TagKey").
		Limit(lim).
		Distinct()
	_, err := q.GetAll(ctx, &tdx)
	return tdx, err
}

/*
GetPrevKeysParentsFilteredByTagKey returns the key's parent key
from the begining to the previous start point with given filter and order.
And also an error.
*/
// ADD LOCATION FILTER HERE AND USE ACCOUNTS ADDRES INFO.
func GetPrevKeysParentsFilteredByTagKey(ctx context.Context, crsrAsString string,
	kt *datastore.Key) ([]*datastore.Key, string, error) {
	var kspx []*datastore.Key
	before, err := datastore.DecodeCursor(crsrAsString)
	if err != nil {
		return nil, crsrAsString, err
	}
	q := datastore.NewQuery("TagServicePack")
	q = q.
		Filter("TagKey =", kt).
		Order("-Created").
		End(before).
		KeysOnly()
	it := q.Run(ctx)
	before, err = it.Cursor()
	if err != nil {
		return nil, crsrAsString, err
	}
	for {
		k, err := it.Next(nil)
		if err == datastore.Done {
			return kspx, before.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		kspx = append(kspx, k.Parent())
	}
}

/*
GetNextKeysParentsFilteredByTagKeyLimited returns limited number of the key's parent key
from the previous end point with given filter and order.
And also an error.
*/
// ADD LOCATION FILTER HERE AND USE ACCOUNTS ADDRES INFO.
func GetNextKeysParentsFilteredByTagKeyLimited(ctx context.Context, crsrAsString string,
	kt *datastore.Key, lim int) ([]*datastore.Key, string, error) {
	var kspx []*datastore.Key
	after, err := datastore.DecodeCursor(crsrAsString)
	if err != nil {
		return nil, crsrAsString, err
	}
	q := datastore.NewQuery("TagServicePack")
	q = q.
		Filter("TagKey =", kt).
		Order("-Created").
		Start(after).
		KeysOnly()
	if lim > 0 && lim < 40 {
		q = q.Limit(lim)
	} else {
		q = q.Limit(20)
	}
	for it := q.Run(ctx); ; {
		k, err := it.Next(nil)
		if err == datastore.Done {
			after, err = it.Cursor()
			return kspx, after.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		kspx = append(kspx, k.Parent())
	}
}

// GetKeys returns the tagServicePack keys by servicePack or tag key and an error.
/*
func GetKeys(ctx context.Context, key *datastore.Key) ([]*datastore.Key, error) {
	q := datastore.NewQuery("TagServicePack")
	kind := key.Kind()
	switch kind {
	case "Tag":
		q = q.Filter("TagKey =", key)
	default:
		// For "ServicePack" kind
		q = q.Ancestor(key)
	}
	q = q.KeysOnly()
	return q.GetAll(ctx, nil)
}
*/

// GetCount returns the count of the entities that has the provided key and an error.
/* func GetCount(s *session.Session, k *datastore.Key) (c int, err error) {
	q := datastore.NewQuery("TagServicePack")
	if k.Kind() == "ServicePack" {
		q = q.Ancestor(k)
	} else {
		q = q.Filter("TagKey =", k)
	}
	c, err = q.Count(s.Ctx)
	return
} */
