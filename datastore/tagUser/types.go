package tagUser

import (
	"google.golang.org/appengine/datastore"
)

// TagUser datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Encoded tag key is key's stringID and user key is the parent key.
type TagUser struct {
	TagKey *datastore.Key
}

// TagUsers is a []*TagUser
type TagUsers []*TagUser
