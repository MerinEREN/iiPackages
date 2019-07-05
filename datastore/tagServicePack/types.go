package tagServicePack

import (
	"google.golang.org/appengine/datastore"
	"time"
)

// TagServicePack datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Encoded tag key is key's stringID and servicePack key is the parent key.
type TagServicePack struct {
	Created time.Time
	TagKey  *datastore.Key
}

// TagServicePacks is a []*TagServicePack
type TagServicePacks []*TagServicePack
