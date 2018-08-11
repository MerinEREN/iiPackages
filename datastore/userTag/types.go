package userTag

import (
	"google.golang.org/appengine/datastore"
)

// UserTag datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type UserTag struct {
	UserKey *datastore.Key
	TagKey  *datastore.Key
}

// UserTags is a []*UserTag
type UserTags []*UserTag
