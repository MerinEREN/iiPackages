/*
Package tag "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package tag

import (
	"errors"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindTag = errors.New("Error while getting tag")
)

/*
GetMulti returns corresponding entities if keys provided.
Otherwise returns limited entitities from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(ctx context.Context, kx interface{}) (Tags, error) {
	ts := make(Tags)
	if kx, ok := kx.([]*datastore.Key); ok {
		tx := make([]*Tag, len(kx))
		// RETURNED ENTITY LIMIT COULD BE A PROBLEM HERE !!!!!!!!!!!!!!!!!!!!!!!!!!
		err := datastore.GetMulti(ctx, kx, tx)
		if err != nil {
			return nil, err
		}
		for i, v := range tx {
			v.ContentID = kx[i].StringID()
			v.ID = kx[i].Encode()
			ts[v.ID] = v
		}
		return ts, err
	}
	q := datastore.NewQuery("Tag").
		Order("-Created")
	for it := q.Run(ctx); ; {
		t := new(Tag)
		k, err := it.Next(t)
		if err == datastore.Done {
			return ts, err
		}
		if err != nil {
			err = ErrFindTag
			return ts, err
		}
		t.ContentID = k.StringID()
		t.ID = k.Encode()
		ts[t.ID] = t
	}
}

// Put puts and returns an entity, and also returns an error.
func Put(ctx context.Context, t *Tag) (*Tag, error) {
	k := datastore.NewKey(ctx, "Tag", t.ContentID, 0, nil)
	var err error
	t.Created = time.Now()
	k, err = datastore.Put(ctx, k, t)
	t.ID = k.Encode()
	return t, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities by the given limit.
func PutAndGetMulti(s *session.Session, t *Tag) (Tags, error) {
	ts := make(Tags)
	tNew := new(Tag)
	// USAGE "s.Ctx" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		tNew, err1 = Put(s.Ctx, t)
		if err1 != nil {
			return
		}
		ts, err1 = GetMulti(s.Ctx, nil)
		return
	}, nil)
	ts[tNew.ID] = tNew
	return ts, err
}

// Delete removes the entity by the provided encoded entity key and returns an error.
func Delete(ctx context.Context, ek string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	return datastore.Delete(ctx, k)
}

// DeleteMulti removes the entitys by the provided encoded entity key slice
// and returns an error.
/* func DeleteMulti(s *session.Session, ekx []string) error {
	var kx []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
	}
	return datastore.DeleteMulti(s.Ctx, kx)
}*/
