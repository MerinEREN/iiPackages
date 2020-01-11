/*
Package context "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package context

import (
	"crypto/md5"
	"github.com/MerinEREN/iiPackages/datastore/pageContext"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"io"
	"time"
)

/*
GetNextLimited returns limited entitities within an order after the given cursor.
If limit is 0 or greater than 40, default limit will be used.
*/
func GetNextLimited(ctx context.Context, crsrAsString string, lim int) (
	Contexts, string, error) {
	var err error
	var after datastore.Cursor
	cs := make(Contexts)
	q := datastore.NewQuery("Context").
		Order("-LastModified")
	if crsrAsString != "" {
		if after, err = datastore.DecodeCursor(crsrAsString); err != nil {
			return nil, crsrAsString, err
		}
		q = q.Start(after)
	}
	// if lim > 0 && lim < 40 {
	if lim != 0 {
		q = q.Limit(lim)
	} else {
		q = q.Limit(12)
	}
	for it := q.Run(ctx); ; {
		c := new(Context)
		k, err := it.Next(c)
		if err == datastore.Done {
			after, err = it.Cursor()
			return cs, after.String(), err
		}
		if err != nil {
			return nil, crsrAsString, err
		}
		c.ID = k.Encode()
		cs[c.ID] = c
	}
}

/*
UpdateMulti is a transaction that delets all PageContext entities for corresponding context
and then creates and puts new ones.
Finally, puts modified Contexts and returns nil map and an error.
*/
func UpdateMulti(ctx context.Context, cx []*Context) (Contexts, error) {
	var kx []*datastore.Key
	var kpcxDelete []*datastore.Key
	var kpcxPut []*datastore.Key
	var pcx pageContext.PageContexts
	var err error
	for _, v := range cx {
		k := new(datastore.Key)
		k, err = datastore.DecodeKey(v.ID)
		if err != nil {
			return nil, err
		}
		kpcx, err := pageContext.GetKeys(ctx, k)
		if err != nil {
			return nil, err
		}
		for _, v2 := range kpcx {
			kpcxDelete = append(kpcxDelete, v2)
		}
		kx = append(kx, k)
		v.LastModified = time.Now()
		for _, v2 := range v.PageIDs {
			kp := new(datastore.Key)
			kp, err = datastore.DecodeKey(v2)
			if err != nil {
				return nil, err
			}
			kpc := datastore.NewIncompleteKey(ctx, "PageContext", kp)
			kpcxPut = append(kpcxPut, kpc)
			pc := new(pageContext.PageContext)
			pc.ContextKey = k
			pcx = append(pcx, pc)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if err1 = datastore.DeleteMulti(ctx, kpcxDelete); err1 != nil {
			return
		}
		if _, err1 = datastore.PutMulti(ctx, kpcxPut, pcx); err1 != nil {
			return
		}
		_, err1 = datastore.PutMulti(ctx, kx, cx)
		return
	}, opts)
	return nil, err
}

/*
AddMulti creates and puts corresponding PageContexts and then puts newly created Contexts.
Also, returns new Contexts and an error.
*/
func AddMulti(ctx context.Context, cx []*Context) (Contexts, error) {
	var kx []*datastore.Key
	var kpcx []*datastore.Key
	var pcx pageContext.PageContexts
	var err error
	for _, v := range cx {
		k := new(datastore.Key)
		// Max stringID lenght for a datastore key is 500 acording to link
		// https://stackoverflow.com/questions/2557632/how-long-max-
		// characters-can-a-datastore-entity-key-name-be-is-it-bad-to-haver
		var stringID string
		if len(v.Values["en-US"]) > 100 {
			h := md5.New()
			io.WriteString(h, v.Values["en-US"])
			stringID = string(h.Sum(nil))
		} else {
			stringID = v.Values["en-US"]
		}
		k = datastore.NewKey(ctx, "Context", stringID, 0, nil)
		kx = append(kx, k)
		v.Created = time.Now()
		v.LastModified = time.Now()
		for _, v2 := range v.PageIDs {
			kp := new(datastore.Key)
			kp, err = datastore.DecodeKey(v2)
			if err != nil {
				return nil, err
			}
			kpc := datastore.NewIncompleteKey(ctx, "PageContext", kp)
			kpcx = append(kpcx, kpc)
			pc := new(pageContext.PageContext)
			pc.ContextKey = k
			pcx = append(pcx, pc)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if _, err1 = datastore.PutMulti(ctx, kpcx, pcx); err1 != nil {
			return
		}
		_, err1 = datastore.PutMulti(ctx, kx, cx)
		return
	}, opts)
	if err != nil {
		return nil, err
	}
	cs := make(Contexts)
	for i, v := range cx {
		v.ID = kx[i].Encode()
		cs[v.ID] = v
	}
	return cs, err
}

/*
AddMultiAndGetNextLimited is a transaction which puts the posted entities first
and then gets entities from the reseted cursor by the given limit.
Finally returnes received entities with posted entities added to them
as a map.
*/
func AddMultiAndGetNextLimited(ctx context.Context, crsrAsString string, cx []*Context) (
	Contexts, string, error) {
	csPut := make(Contexts)
	csGet := make(Contexts)
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		if csPut, err1 = AddMulti(ctx, cx); err1 != nil {
			return
		}
		csGet, crsrAsString, err1 = GetNextLimited(ctx, crsrAsString, 2222)
		return
	}, nil)
	if err == nil {
		for i, v := range csGet {
			csPut[i] = v
		}
	}
	return csPut, crsrAsString, err
}

/*
DeleteMulti removes the entities
and all the corresponding pageContext entities by the provided encoded keys
also returns an error.
*/
func DeleteMulti(ctx context.Context, ekx []string) error {
	var kx []*datastore.Key
	var kpcx []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
		kpcx2, err := pageContext.GetKeys(ctx, k)
		if err != nil {
			return err
		}
		for _, v2 := range kpcx2 {
			kpcx = append(kpcx, v2)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, kpcx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kx)
		return
	}, opts)
}
