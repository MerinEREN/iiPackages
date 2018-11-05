package userTag

import (
	"google.golang.org/appengine/datastore"
)

// UserTag datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type UserTag struct {
	TagKey *datastore.Key
}

// UserTags is a []*UserTag
type UserTags []*UserTag
