/*
Package demand "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package demand

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// GetMulti returns the corresponding entities as a map with their "ID"s assigned
// by their keys and an error.
func GetMulti(ctx context.Context, kx []*datastore.Key) (Demands, error) {
	var dx []*Demand
	err := datastore.GetMulti(ctx, kx, dx)
	if err != nil {
		return nil, err
	}
	ds := make(map[string]*Demand)
	for i, v := range dx {
		v.ID = kx[i].Encode()
		ds[v.ID] = v
	}
	return ds, err
}

/*
GetNextByParentLimited returns limited entities filtered by user key as parent
from the previous cursor in an order.
Also, returns the next cursor as string and an error.
*/
func GetNextByParentLimited(ctx context.Context, crsrAsString string, pk *datastore.Key,
	lim int) (Demands, string, error) {
	var err error
	var after datastore.Cursor
	ds := make(Demands)
	q := datastore.NewQuery("Demand")
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
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			after, err = it.Cursor()
			return ds, after.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		d.ID = k.Encode()
		ds[d.ID] = d
	}
}

// GetLatestLimited returns limited entities from the begining of the kind with given limit
// and an error.
func GetLatestLimited(ctx context.Context, lim int) ([]Demand, error) {
	var dx []Demand
	q := datastore.NewQuery("Demand")
	q = q.Order("-LastModified").
		Project("TagIDs").
		Limit(lim)
	_, err := q.GetAll(ctx, &dx)
	return dx, err
}

// Update modifies an entity in the kind and returns the entity and an error.
func Update(ctx context.Context, k *datastore.Key, d *Demand) (*datastore.Key, *Demand, error) {
	var err error
	tempD := new(Demand)
	if err = datastore.Get(ctx, k, tempD); err != nil {
		return nil, nil, err
	}
	d.Created = tempD.Created
	d.LastModified = time.Now()
	k, err = datastore.Put(ctx, k, d)
	return k, d, err
}

// UpdateStatus set's demand status to given value "v" by given encoded demand key "ek"
// and returns an error.
func UpdateStatus(ctx context.Context, ek, v string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	d := new(Demand)
	/* return datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if err1 = datastore.Get(ctx, k, d); err1 != nil {
			return
		}
		d.Status = v
		_, err1 = datastore.Put(ctx, k, d)
		// UPDATE DEMAND'S OFFER'S STATUSES !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		return
	}, nil) */
	if err = datastore.Get(ctx, k, d); err != nil {
		return err
	}
	d.Status = v
	_, err = datastore.Put(ctx, k, d)
	return err
}
