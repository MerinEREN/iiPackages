/*
Package offer "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package offer

import (
	"github.com/MerinEREN/iiPackages/datastore/demand"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// GetByParent returns limited entities via offer key as parent from the previous cursor
// with given filters and order.
func GetByParent(ctx context.Context, crsrAsString string, pk *datastore.Key) (
	Offers, string, error) {
	var err error
	var crsr datastore.Cursor
	os := make(Offers)
	q := datastore.NewQuery("Offer")
	q = q.
		Ancestor(pk).
		Filter("Status =", "active").
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
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return os, crsr.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		o.ID = k.Encode()
		ku, err := datastore.DecodeKey(o.UserID)
		if err != nil {
			return nil, crsrAsString, err
		}
		o.AccountID = ku.Parent().Encode()
		os[o.ID] = o
	}
}

// GetByUserID returns limited entities filtered by encoded user key as UserID
// from the previous cursor in an order.
// Also returns an updated cursor as string and the error.
func GetByUserID(ctx context.Context, crsrAsString, uID string) (
	Offers, string, error) {
	var err error
	var crsr datastore.Cursor
	os := make(Offers)
	q := datastore.NewQuery("Offer")
	q = q.
		Filter("UserID =", uID).
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
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return os, crsr.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		o.ID = k.Encode()
		os[o.ID] = o
	}
}

// Put adds to or modifies an entity in the kind according to request method
// and returns the entity and an error.
func Put(s *session.Session, o *Offer, k *datastore.Key) (*Offer, error) {
	var err error
	if s.R.Method == "POST" {
		o.Created = time.Now()
	} else {
		// Updating an egzisting entity.
		tempO := new(Offer)
		if err = datastore.Get(s.Ctx, k, tempO); err != nil {
			return nil, err
		}
		o.Created = tempO.Created
	}
	o.LastModified = time.Now()
	_, err = datastore.Put(s.Ctx, k, o)
	return o, err
}

// UpdateStatus set's offer status to given value "v" by given encoded offer key "ek"
// and returns an error.
func UpdateStatus(ctx context.Context, ek, v string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	o := new(Offer)
	return datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if err1 = datastore.Get(ctx, k, o); err1 != nil {
			return
		}
		o.Status = v
		if _, err1 = datastore.Put(ctx, k, o); err1 != nil {
			return
		}
		if v == "accepted" {
			err1 = demand.UpdateStatus(ctx, k.Parent().Encode(), v)
		}
		return
	}, nil)
}
