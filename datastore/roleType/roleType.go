/*
Package roleType "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package roleType

import (
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

/*
GetMulti returns corresponding entities if keys provided.
Otherwise returns limited entitities from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(ctx context.Context, kx interface{}) (RoleTypes, error) {
	rts := make(RoleTypes)
	if kx, ok := kx.([]*datastore.Key); ok {
		rtx := make([]*RoleType, len(kx))
		// RETURNED ENTITY LIMIT COULD BE A PROBLEM HERE !!!!!!!!!!!!!!!!!!!!!!!!!!
		err := datastore.GetMulti(ctx, kx, rtx)
		if err != nil {
			return nil, err
		}
		for i, v := range rtx {
			v.ID = kx[i].StringID()
			rts[v.ID] = v
		}
		return rts, err
	}
	q := datastore.NewQuery("RoleType").
		Order("-Created")
	for it := q.Run(ctx); ; {
		rt := new(RoleType)
		k, err := it.Next(rt)
		if err == datastore.Done {
			return rts, err
		}
		if err != nil {
			return rts, err
		}
		rt.ID = k.StringID()
		rts[rt.ID] = rt
	}
}

// Put puts and returns an entity, and also returns an error.
func Put(ctx context.Context, rt *RoleType) (*RoleType, error) {
	k := datastore.NewKey(ctx, "RoleType", rt.ID, 0, nil)
	var err error
	rt.Created = time.Now()
	k, err = datastore.Put(ctx, k, rt)
	return rt, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities by the given limit.
func PutAndGetMulti(s *session.Session, rt *RoleType) (RoleTypes, error) {
	rts := make(RoleTypes)
	rtNew := new(RoleType)
	// USAGE "s.Ctx" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		rtNew, err1 = Put(s.Ctx, rt)
		if err1 != nil {
			return
		}
		rts, err1 = GetMulti(s.Ctx, nil)
		return
	}, nil)
	rts[rtNew.ID] = rtNew
	return rts, err
}

// Delete removes the entity by the provided stringID of a key and returns an error.
func Delete(ctx context.Context, ID string) error {
	k := datastore.NewKey(ctx, "RoleType", ID, 0, nil)
	return datastore.Delete(ctx, k)
}
