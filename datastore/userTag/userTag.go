/*
Package userTag "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
7ne will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package userTag

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysProjected returns the user keys as a slice if the tag key is provided
// or returns the tag keys as a slice if user key is provided and also an error.
func GetKeysProjected(ctx context.Context, key *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("UserTag")
	kind := key.Kind()
	switch kind {
	case "Tag":
		q = q.Filter("TagKey =", key).
			Project("UserKey")
	case "User":
		q = q.Filter("UserKey =", key).
			Project("TagKey")
	}
	for it := q.Run(ctx); ; {
		ut := new(UserTag)
		_, err := it.Next(ut)
		if err == datastore.Done {
			return kx, err
		}
		if err != nil {
			return nil, err
		}
		if key.Kind() == "Tag" {
			kx = append(kx, ut.UserKey)
		} else {
			kx = append(kx, ut.TagKey)
		}
	}
}

// PutMulti puts entities with corresponding entity keys and returns an error.
func PutMulti(ctx context.Context, kx []*datastore.Key, utx UserTags) error {
	_, err := datastore.PutMulti(ctx, kx, utx)
	return err
}

// Delete deletes an entity by provided key and returns an error.
func Delete(ctx context.Context, k *datastore.Key) error {
	return datastore.Delete(ctx, k)
}

// GetCount returns the count of the entities that has provided key and an error.
/* func GetCount(s *session.Session, k *datastore.Key) (c int, err error) {
	q := datastore.NewQuery("UserTag")
	if k.Kind() == "User" {
		q = q.Filter("UserKey =", k)
	} else {
		q = q.Filter("TagKey =", k)
	}
	c, err = q.Count(s.Ctx)
	return
} */
