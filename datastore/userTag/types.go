package userTag

import (
	"google.golang.org/appengine/datastore"
)

// UserTag datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// User key is the parent key.
type UserTag struct {
	TagKey *datastore.Key
}

// UserTags is a []*UserTag
type UserTags []*UserTag
