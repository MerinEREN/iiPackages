package pageContent

import (
	"google.golang.org/appengine/datastore"
)

// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type PageContent struct {
	ContentKey *datastore.Key
	PageKey    *datastore.Key
}

// PageContents is a []*PageContent
type PageContents []*PageContent
