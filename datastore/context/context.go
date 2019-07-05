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
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"io"
	"time"
)

/*
GetNextLimited returns limited entitities after the given cursor.
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
		after, err = datastore.DecodeCursor(crsrAsString)
		if err != nil {
			return nil, crsrAsString, err
		}
		q = q.Start(after)
	}
	// if lim > 0 && lim < 40 {
	if lim != 0 {
		q = q.Limit(lim)
	} else {
		q = q.Limit(20)
	}
	for it := q.Run(ctx); ; {
		c := new(Context)
		k, err := it.Next(c)
		if err == datastore.Done {
			after, err = it.Cursor()
			return cs, after.String(), err
		}
		if err != nil {
			return cs, crsrAsString, err
		}
		c.ID = k.Encode()
		cs[c.ID] = c
	}
}

/*
PutMulti is a transaction that delets all PageContext entities for corresponding context
if only the request method is "PUT".
And then creates and puts new ones.
Finally, puts modified or newly created Contexts and returns new Contexts if request method
is POST.
*/
func PutMulti(s *session.Session, cx []*Context) (Contexts, error) {
	var kx []*datastore.Key
	var kpcxDelete []*datastore.Key
	var kpcxPut []*datastore.Key
	var pcx pageContext.PageContexts
	var err error
	for _, v := range cx {
		k := new(datastore.Key)
		if s.R.Method == "PUT" {
			k, err = datastore.DecodeKey(v.ID)
			if err != nil {
				return nil, err
			}
			kpcx2, err := pageContext.GetKeysOnly(s.Ctx, k)
			if err != datastore.Done {
				return nil, err
			}
			for _, v2 := range kpcx2 {
				kpcxDelete = append(kpcxDelete, v2)
			}
		} else {
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
			k = datastore.NewKey(s.Ctx, "Context", stringID, 0, nil)
			v.Created = time.Now()
		}
		kx = append(kx, k)
		v.LastModified = time.Now()
		for _, v2 := range v.PageIDs {
			kp := new(datastore.Key)
			kp, err = datastore.DecodeKey(v2)
			if err != nil {
				return nil, err
			}
			kpc := datastore.NewIncompleteKey(s.Ctx, "PageContext", kp)
			kpcxPut = append(kpcxPut, kpc)
			pc := new(pageContext.PageContext)
			pc.ContextKey = k
			pcx = append(pcx, pc)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if s.R.Method == "PUT" {
			err1 = datastore.DeleteMulti(ctx, kpcxDelete)
			if err1 != nil {
				return
			}
		}
		_, err1 = datastore.PutMulti(ctx, kpcxPut, pcx)
		if err1 != nil {
			return
		}
		_, err1 = datastore.PutMulti(ctx, kx, cx)
		return
	}, opts)
	if err != nil {
		return nil, err
	}
	if s.R.Method == "POST" {
		cs := make(Contexts)
		for i, v := range cx {
			v.ID = kx[i].Encode()
			cs[v.ID] = v
		}
		return cs, err
	}
	return nil, err
}

// PutMultiAndGetNextLimited is a transaction which puts the posted entities first
// and then gets entities from the reseted cursor by the given limit.
// Finally returnes received entities with posted entities added to them
// as a map.
func PutMultiAndGetNextLimited(s *session.Session, crsr datastore.Cursor, cx []*Context) (
	Contexts, datastore.Cursor, error) {
	csPut := make(Contexts)
	csGet := make(Contexts)
	// USAGE "s" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!!!!!
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if csPut, err1 = PutMulti(s, cx); err1 != nil {
			return
		}
		if csGet, crsr, err1 = GetNextLimited(s.Ctx, crsr, 2222); err1 == nil ||
			err1 == datastore.Done {
			for i, v := range csGet {
				csPut[i] = v
			}
		}
		return
	}, nil)
	return csPut, crsr, err
}

// DeleteMulti removes the entities
// and all the corresponding pageContext entities by the provided encoded keys
// also returns an error.
func DeleteMulti(ctx context.Context, ekx []string) error {
	var kx []*datastore.Key
	var kpcx []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
		kpcx2, err := pageContext.GetKeysOnly(ctx, k)
		if err != datastore.Done {
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
