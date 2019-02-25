/*
Package servicePack "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package servicePack

import (
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Get returns limited entities from the start of the kind with given filter in an order,
// and also first and last cursors of the query as string with an error.
// Add location filter here and use accounts addres info.
func Get(ctx context.Context, tKey *datastore.Key) (ServicePacks, string, string, error) {
	sps := make(ServicePacks)
	var crsrEnd datastore.Cursor
	q := datastore.NewQuery("ServicePack")
	q = q.
		Filter("TagIDs =", tKey.Encode()).
		Filter("Status =", "active").
		Order("-LastModified").
		Limit(10)
	it := q.Run(ctx)
	crsrStart, err := it.Cursor()
	if err != nil {
		return nil, "", "", err
	}
	for {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			crsrEnd, err = it.Cursor()
			return sps, crsrStart.String(), crsrEnd.String(), err
		}
		if err != nil {
			return nil, crsrStart.String(), crsrEnd.String(), err
		}
		sp.ID = k.Encode()
		sps[sp.ID] = sp
	}
}

// GetNewest returns all the results from the begining to the previous start point
// via given filter and order.
// Add location filter here and use accounts addres info.
func GetNewest(ctx context.Context, crsrStartAsString string, tKey *datastore.Key) (
	ServicePacks, string, error) {
	sps := make(ServicePacks)
	crsrStart, err := datastore.DecodeCursor(crsrStartAsString)
	if err != nil {
		return nil, crsrStartAsString, err
	}
	q := datastore.NewQuery("ServicePack")
	q = q.
		Filter("TagIDs =", tKey.Encode()).
		Filter("Status =", "active").
		Order("-LastModified").
		End(crsrStart)
	it := q.Run(ctx)
	crsrStart, err = it.Cursor()
	if err != nil {
		return nil, crsrStartAsString, err
	}
	crsrStartAsString = crsrStart.String()
	for {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			return sps, crsrStartAsString, err
		}
		if err != nil {
			return nil, crsrStartAsString, err
		}
		sp.ID = k.Encode()
		sps[sp.ID] = sp
	}
}

// GetNewestKeysFilteredByTag returns the keys from the begining to the previous
// start point with given filter and order.
// Add location filter here and use accounts addres info.
func GetNewestKeysFilteredByTag(ctx context.Context, crsrStartAsString string, tKey *datastore.Key) (
	[]*datastore.Key, error) {
	var kx []*datastore.Key
	crsrStart, err := datastore.DecodeCursor(crsrStartAsString)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery("ServicePack")
	q = q.
		Filter("TagIDs =", tKey.Encode()).
		Filter("Status =", "active").
		Order("-LastModified").
		End(crsrStart).
		KeysOnly()
	for it := q.Run(ctx); ; {
		k, err := it.Next(nil)
		if err == datastore.Done {
			return kx, err
		}
		if err != nil {
			return nil, err
		}
		kx = append(kx, k)
	}
}

// GetOldest returns limited(10) entities starting from the previously returned end cursor
// with given filter and order.
// Add location filter here and use accounts addres info.
func GetOldest(ctx context.Context, crsrEndAsString string, tKey *datastore.Key) (
	ServicePacks, string, error) {
	sps := make(ServicePacks)
	crsrEnd, err := datastore.DecodeCursor(crsrEndAsString)
	if err != nil {
		return nil, crsrEndAsString, err
	}
	q := datastore.NewQuery("ServicePack")
	q = q.
		Filter("TagIDs =", tKey.Encode()).
		Filter("Status =", "active").
		Order("-LastModified").
		Start(crsrEnd).
		Limit(10)
	for it := q.Run(ctx); ; {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			crsrEnd, err = it.Cursor()
			return sps, crsrEnd.String(), err
		}
		if err != nil {
			return nil, crsrEndAsString, err
		}
		sp.ID = k.Encode()
		sps[sp.ID] = sp
	}
}

// GetByParent returns limited entities via user key as parent from the previous cursor
// with given filters and order.
func GetByParent(ctx context.Context, crsrAsString string, pk *datastore.Key) (
	ServicePacks, string, error) {
	var err error
	var crsr datastore.Cursor
	sps := make(ServicePacks)
	q := datastore.NewQuery("ServicePack")
	q = q.Ancestor(pk).
		Order("-LastModified")
	if crsrAsString != "" {
		crsr, err = datastore.DecodeCursor(crsrAsString)
		if err != nil {
			return nil, crsrAsString, err
		}
		q = q.Start(crsr)
	}
	q = q.Limit(10)
	for it := q.Run(ctx); ; {
		sp := new(ServicePack)
		k, err := it.Next(sp)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return sps, crsr.String(), err
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

// Put adds to or modifies an entity in the kind according to request method
// and returns an error.
func Put(s *session.Session, sp *ServicePack, pID string) error {
	k := new(datastore.Key)
	var err error
	if s.R.Method == "POST" {
		pk, err := datastore.DecodeKey(pID)
		if err != nil {
			return err
		}
		k = datastore.NewIncompleteKey(s.Ctx, "ServicePack", pk)
		sp.Created = time.Now()
	} else {
		// Updating an egzisting entity.
		k, err = datastore.DecodeKey(sp.ID)
		if err != nil {
			return err
		}
		tempSP := new(ServicePack)
		if err = datastore.Get(s.Ctx, k, tempSP); err != nil {
			return err
		}
		sp.Created = tempSP.Created
	}
	sp.LastModified = time.Now()
	_, err = datastore.Put(s.Ctx, k, sp)
	return err
}
