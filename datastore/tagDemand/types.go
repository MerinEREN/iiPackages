package tagDemand

import (
	"google.golang.org/appengine/datastore"
	"time"
)

// TagDemand datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Encoded tag key is key's stringID and demand key is the parent key.
type TagDemand struct {
	Created time.Time
	TagKey  *datastore.Key
}

// TagDemands is a []*TagDemand
type TagDemands []*TagDemand
