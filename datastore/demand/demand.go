/*
Package demand "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package demand

import (
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Get returns limited entities from the start of the kind with given filter in an order,
// and also first and last cursors of the query as string with an error.
// Add location filter here and use accounts addres info.
func Get(ctx context.Context, tKey *datastore.Key) (Demands, string, string, error) {
	ds := make(Demands)
	var crsrEnd datastore.Cursor
	q := datastore.NewQuery("Demand")
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
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			crsrEnd, err = it.Cursor()
			return ds, crsrStart.String(), crsrEnd.String(), err
		}
		if err != nil {
			return nil, crsrStart.String(), crsrEnd.String(), err
		}
		d.ID = k.Encode()
		ds[d.ID] = d
	}
}

// GetNewest returns all the results from the begining to the previous start point
// via given filter and order.
// Add location filter here and use accounts addres info.
func GetNewest(ctx context.Context, crsrStartAsString string, tKey *datastore.Key) (
	Demands, string, error) {
	ds := make(Demands)
	crsrStart, err := datastore.DecodeCursor(crsrStartAsString)
	if err != nil {
		return nil, crsrStartAsString, err
	}
	q := datastore.NewQuery("Demand")
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
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			return ds, crsrStartAsString, err
		}
		if err != nil {
			return nil, crsrStartAsString, err
		}
		d.ID = k.Encode()
		ds[d.ID] = d
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
	q := datastore.NewQuery("Demand")
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
	Demands, string, error) {
	ds := make(Demands)
	crsrEnd, err := datastore.DecodeCursor(crsrEndAsString)
	if err != nil {
		return nil, crsrEndAsString, err
	}
	q := datastore.NewQuery("Demand")
	q = q.
		Filter("TagIDs =", tKey.Encode()).
		Filter("Status =", "active").
		Order("-LastModified").
		Start(crsrEnd).
		Limit(10)
	for it := q.Run(ctx); ; {
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			crsrEnd, err = it.Cursor()
			return ds, crsrEnd.String(), err
		}
		if err != nil {
			return nil, crsrEndAsString, err
		}
		d.ID = k.Encode()
		ds[d.ID] = d
	}
}

// GetByParent returns limited entities via demand key as parent from the previous cursor
// with given filters and order.
func GetByParent(ctx context.Context, crsrAsString string, pk *datastore.Key) (
	Demands, string, error) {
	var err error
	var crsr datastore.Cursor
	ds := make(Demands)
	q := datastore.NewQuery("Demand")
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
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return ds, crsr.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		d.ID = k.Encode()
		ds[d.ID] = d
	}
}

// GetAllLimited returns limited entities from the begining of the kind with given limit
// and an error.
func GetAllLimited(ctx context.Context, lim int) ([]Demand, error) {
	var dx []Demand
	q := datastore.NewQuery("Demand")
	q = q.Order("-LastModified").
		Project("TagIDs").
		Limit(lim)
	_, err := q.GetAll(ctx, &dx)
	return dx, err
}

// Put adds to or modifies an entity in the kind according to request method
// and returns the entity and an error.
func Put(s *session.Session, d *Demand, k *datastore.Key) (*Demand, error) {
	var err error
	if s.R.Method == "POST" {
		d.Created = time.Now()
	} else {
		// Updating an existing entity.
		tempD := new(Demand)
		if err = datastore.Get(s.Ctx, k, tempD); err != nil {
			return nil, err
		}
		d.Created = tempD.Created
	}
	d.LastModified = time.Now()
	k, err = datastore.Put(s.Ctx, k, d)
	return d, err
}

// UpdateStatus set's demand status to given value "v" by given encoded demand key "ek"
// and returns an error.
func UpdateStatus(ctx context.Context, ek, v string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	d := new(Demand)
	if err = datastore.Get(ctx, k, d); err != nil {
		return err
	}
	d.Status = v
	_, err = datastore.Put(ctx, k, d)
	return err
}
