/*
Package demand "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package demand

import (
	"github.com/MerinEREN/iiPackages/datastore/photo"
	"github.com/MerinEREN/iiPackages/datastore/tagDemand"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

/*
Add puts the demand, it's tag relations and it's photos
also returns it's key and an error.
*/
func Add(ctx context.Context, k *datastore.Key, d *Demand, tIDs []string,
	px []*photo.Photo) (*datastore.Key, error) {
	tdx := make([]*tagDemand.TagDemand, 0, cap(tIDs))
	ktdx := make([]*datastore.Key, 0, cap(tIDs))
	kpx := make([]*datastore.Key, 0, cap(px))
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) (
		err1 error) {
		k, err1 = datastore.Put(ctx, k, d)
		if err1 != nil {
			return
		}
		for _, v := range tIDs {
			ktd := datastore.NewKey(ctx, "TagDemand", v, 0, k)
			ktdx = append(ktdx, ktd)
			kt := new(datastore.Key)
			kt, err1 = datastore.DecodeKey(v)
			if err1 != nil {
				return
			}
			td := &tagDemand.TagDemand{
				Created: time.Now(),
				TagKey:  kt,
			}
			tdx = append(tdx, td)
		}
		_, err1 = datastore.PutMulti(ctx, ktdx, tdx)
		if err1 != nil {
			return
		}
		for i := 0; i < len(px); i++ {
			kp := datastore.NewIncompleteKey(ctx, "Photo", k)
			kpx = append(kpx, kp)
		}
		_, err1 = datastore.PutMulti(ctx, kpx, px)
		return
	}, nil)
	return k, err
}

/*
GetMulti returns the corresponding entities as a map with their "ID"s assigned
by their keys and an error.
*/
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

/*
UpdateStatus set's demand status to given value "v" by given encoded demand key "ek"
and returns an error.
*/
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

/*
Update updates and returns (only with "ID", "LastModified" and "Status" fields) the entity
by the given encoded entity key "ek" and also returns an error.
*/
func Update(ctx context.Context, d *Demand, ek string) (*Demand, error) {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return nil, err
	}
	d2 := new(Demand)
	if err = datastore.Get(ctx, k, d2); err != nil {
		return nil, err
	}
	d.Created = d2.Created
	d.LastModified = time.Now()
	d.Status = "updated"
	_, err = datastore.Put(ctx, k, d)
	d3 := &Demand{
		ID:           ek,
		LastModified: d.LastModified,
		Status:       d.Status,
	}
	return d3, err
}
