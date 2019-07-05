/*
Package role menages datastore's role kind.
*/
package role

import (
	"errors"
	"github.com/MerinEREN/iiPackages/datastore/roleTypeRole"
	"github.com/MerinEREN/iiPackages/datastore/roleUser"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindRole = errors.New("Error while getting role")
)

/*
GetMulti returns corresponding entities if keys provided.
Otherwise returns limited entitities from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(ctx context.Context, kx interface{}) (Roles, error) {
	rs := make(Roles)
	if kx, ok := kx.([]*datastore.Key); ok {
		rx := make([]*Role, len(kx))
		// RETURNED ENTITY LIMIT COULD BE A PROBLEM HERE !!!!!!!!!!!!!!!!!!!!!!!!!!
		err := datastore.GetMulti(ctx, kx, rx)
		if err != nil {
			return nil, err
		}
		for i, v := range rx {
			v.ContentID = kx[i].StringID()
			v.ID = kx[i].Encode()
			rs[v.ID] = v
		}
		return rs, err
	}
	q := datastore.NewQuery("Role").
		Order("-Created")
	for it := q.Run(ctx); ; {
		r := new(Role)
		k, err := it.Next(r)
		if err == datastore.Done {
			return rs, err
		}
		if err != nil {
			err = ErrFindRole
			return rs, err
		}
		r.ContentID = k.StringID()
		r.ID = k.Encode()
		rs[r.ID] = r
	}
}

// Put puts and returns an entity, and also returns an error.
func Put(ctx context.Context, r *Role) (*Role, error) {
	k := datastore.NewKey(ctx, "Role", r.ContentID, 0, nil)
	var err error
	r.Created = time.Now()
	k, err = datastore.Put(ctx, k, r)
	r.ID = k.Encode()
	return r, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities by the given limit.
func PutAndGetMulti(s *session.Session, r *Role) (Roles, error) {
	rs := make(Roles)
	rNew := new(Role)
	// USAGE "s.Ctx" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		rNew, err1 = Put(s.Ctx, r)
		if err1 != nil {
			return
		}
		rs, err1 = GetMulti(s.Ctx, nil)
		return
	}, nil)
	rs[rNew.ID] = rNew
	return rs, err
}

// Delete removes the entity and all the corresponding "roleTypeRole" and "roleUser"
// entities in a transaction by the provided encoded entity key
// and returns an error.
func Delete(ctx context.Context, ek string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	krtrx, err := roleTypeRole.GetKeys(ctx, k)
	if err != nil {
		return err
	}
	kurx, err := roleUser.GetKeysByUserOrRoleKey(ctx, k)
	if err != nil {
		return err
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	err = datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		err1 = roleTypeRole.DeleteMulti(ctx, krtrx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kurx)
		if err1 != nil {
			return
		}
		err1 = datastore.Delete(ctx, k)
		return
	}, opts)
	return err
}
