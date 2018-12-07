package user

import (
	"encoding/json"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
	"log"
)

// Get tries to return logged user from the memcache first,
// if fails, tries to return logged user's datastore key from the memcache and then
// gets the logged user via that key from the datastore.
// Finally, if both atemps above fails, tries to get logged user via user's email
// which stores in the session struct.
// And also returns an error.
// In adition to those, adds one or both of the logged user and the logged user's key
// to the memcache if they are apsent.
func Get(s *session.Session) (*user.User, error) {
	u := new(user.User)
	item, err := memcache.Get(s.Ctx, "u")
	if err == nil {
		err = json.Unmarshal(item.Value, u)
	} else {
		var bs []byte
		k := new(datastore.Key)
		item, err = memcache.Get(s.Ctx, "uKey")
		if err == nil {
			err = json.Unmarshal(item.Value, k)
			if err != nil {
				return nil, err
			}
			u, err = user.Get(s.Ctx, k)
		} else {
			u, k, err = user.GetViaEmail(s)
			if err != nil {
				return nil, err
			}
			bs, err = json.Marshal(k)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
			item = &memcache.Item{
				Key:   "uKey",
				Value: bs,
			}
			err = memcache.Add(s.Ctx, item)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			}
		}
		bs, err = json.Marshal(u)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		}
		item = &memcache.Item{
			Key:   "u",
			Value: bs,
		}
		err = memcache.Add(s.Ctx, item)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
		}
	}
	return u, err
}
