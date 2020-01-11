/*
Package offer "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package offer

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

/*
GetNextByParentLimited returns limited and filtered entities within an order
after the given cursor.
If limit is 0 or greater than 40, default limit will be used.
*/
func GetNextByParentLimited(ctx context.Context, crsrAsString string, pk *datastore.Key, lim int) (
	Offers, string, error) {
	var err error
	var after datastore.Cursor
	os := make(Offers)
	q := datastore.NewQuery("Offer")
	q = q.
		Ancestor(pk).
		Filter("Status =", "active").
		Order("-Created")
	if crsrAsString != "" {
		if after, err = datastore.DecodeCursor(crsrAsString); err != nil {
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
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			after, err = it.Cursor()
			return os, after.String(), err
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
