/*
Package role menages datastore's role kind.
*/
package role

import (
	"errors"
	"github.com/MerinEREN/iiPackages/datastore/roleTypeRole"
	"github.com/MerinEREN/iiPackages/datastore/roleUser"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Errors
var (
	ErrFindRole = errors.New("Error while getting role")
)

/*
GetAll returns all the entities in an order(-Created) and an error.
*/
func GetAll(ctx context.Context) (Roles, error) {
	var rx []*Role
	q := datastore.NewQuery("Role")
	q = q.
		Order("-Created")
	kx, err := q.GetAll(ctx, &rx)
	if err != nil {
		return nil, err
	}
	rs := make(Roles)
	for i, v := range kx {
		rx[i].ID = v.Encode()
		rx[i].ContextID = v.StringID()
		rs[v.Encode()] = rx[i]
	}
	return rs, err
}

/*
GetMulti returns corresponding entities if keys provided.
Otherwise returns limited entitities from the given cursor.
If limit is nil default limit will be used.
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
			v.ContextID = kx[i].StringID()
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
		r.ContextID = k.StringID()
		r.ID = k.Encode()
		rs[r.ID] = r
	}
}
*/

/*
Add puts the role first
and then creates and puts all the corresponding roleTypeRoles with role key
as the parent key.
Also returns an error.
*/
func Add(ctx context.Context, k *datastore.Key, r *Role) error {
	var krtrx []*datastore.Key
	var rtrx roleTypeRole.RoleTypeRoles
	for _, v := range r.RoleTypes {
		krtr := datastore.NewIncompleteKey(ctx, "RoleTypeRole", k)
		krtrx = append(krtrx, krtr)
		krt := datastore.NewKey(ctx, "RoleType", v, 0, nil)
		rtr := &roleTypeRole.RoleTypeRole{
			RoleTypeKey: krt,
		}
		rtrx = append(rtrx, rtr)
	}
	return datastore.RunInTransaction(ctx, func(ctx context.Context) (
		err1 error) {
		k, err1 = datastore.Put(ctx, k, r)
		if err1 != nil {
			return
		}
		_, err1 = datastore.PutMulti(ctx, krtrx, rtrx)
		return
	}, nil)
}

/*
Delete removes the entity and all the corresponding "roleTypeRole" and "roleUser"
entities in a transaction by the provided encoded entity key
and returns an error.
*/
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
	if err != datastore.Done {
		return err
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, krtrx)
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
}
