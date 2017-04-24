package content

import (
	"encoding/json"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	// "io/ioutil"
	// "log"
)

// Use Key List method instead of this one.
func Get(ctx context.Context, page string) (Contents, error) {
	// filename := title + ".html"
	//USE CURRENT WORKING DIRECTORY IN PATH !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	// BECAUSE ioutil.ReadFile USES CALLER PACKAGE'S DIRECTORY AS CURRENT WORKING
	// DIRECTORY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	/* body, err := ioutil.ReadFile("../page/templates/" + filename)
	if err != nil {
		return nil, err
	} */
	pc := new(PageContent)
	c := new(Content)
	var cs Contents
	csItem, err := memcache.Get(ctx, page)
	if err == memcache.ErrCacheMiss {
		qpc := datastore.NewQuery("PageContent").Filter("PageID =", page)
		for it := qpc.Run(ctx); ; {
			_, err = it.Next(pc)
			if err == datastore.Done {
				return cs, ErrDatastoreDoneContentID
			}
			if err != nil {
				return nil, err
			}
			qc := datastore.NewQuery("Content").Filter("ContentID =", pc.ContentID)
			it2 := qc.Run(ctx)
			_, err = it2.Next(c)
			if err != datastore.Done && err != nil {
				return nil, err
			}
			if err == datastore.Done {
				err = ErrDatastoreDoneContent
			}
			cs = append(cs, *c)
		}

		bs, err := json.Marshal(cs)
		if err != nil {
			return nil, err
		}
		csItem = &memcache.Item{
			Key:   page,
			Value: bs,
		}
		if err = memcache.Set(ctx, csItem); err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal(csItem.Value, &cs)
		if err != nil {
			return nil, err
		}
	}
	return cs, nil
}
