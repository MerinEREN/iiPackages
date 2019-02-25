/*
Package offers "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package offers

import (
	"github.com/MerinEREN/iiPackages/api"
	"github.com/MerinEREN/iiPackages/datastore/offer"
	"github.com/MerinEREN/iiPackages/datastore/user"
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
	"log"
	"net/http"
	"strconv"
)

// Handler returns account's offers via user ID if the user is admin
// otherwise returns only logged user's offers.
// Or returns demand's offers via demand ID.
// If the request method is POST, puts the offer to the datastore.
func Handler(s *session.Session) {
	switch s.R.Method {
	case "POST":
		dID := s.R.FormValue("dID")
		if dID == "" {
			log.Printf("Path: %s, Error: no demand ID\n", s.R.URL.Path)
			http.Error(s.W, "No demand ID", http.StatusBadRequest)
			return
		}
		uID := s.R.FormValue("uID")
		if uID == "" {
			log.Printf("Path: %s, Error: no user ID\n", s.R.URL.Path)
			http.Error(s.W, "No user ID", http.StatusBadRequest)
			return
		}
		explanation := s.R.FormValue("explanation")
		if explanation == "" {
			log.Printf("Path: %s, Error: no explanation value\n", s.R.URL.Path)
			http.Error(s.W, "No explanation value", http.StatusBadRequest)
			return
		}
		amountString := s.R.FormValue("amount")
		if amountString == "" {
			log.Printf("Path: %s, Error: no amount value\n", s.R.URL.Path)
			http.Error(s.W, "No amount value", http.StatusBadRequest)
			return
		}
		amount, err := strconv.ParseFloat(amountString, 64)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		o := &offer.Offer{
			UserID:      uID,
			Explanation: explanation,
			Amount:      amount,
			Status:      "active",
		}
		pk, err := datastore.DecodeKey(dID)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		k := datastore.NewIncompleteKey(s.Ctx, "Offer", pk)
		_, err = offer.Put(s, o, k)
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		s.W.WriteHeader(http.StatusNoContent)
	default:
		err := s.R.ParseForm()
		if err != nil {
			log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
			http.Error(s.W, err.Error(), http.StatusInternalServerError)
			return
		}
		accID := s.R.Form.Get("aID")
		dID := s.R.Form.Get("dID")
		var crsrAsStringx []string
		os := make(offer.Offers)
		URL := s.R.URL
		q := URL.Query()
		rb := new(api.ResponseBody)
		if accID != "" {
			// DUMMY BLOCK !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
			aKey, err := datastore.DecodeKey(accID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			uKeyx, err := user.GetKeysByParentOrdered(s.Ctx, aKey)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			if len(s.R.Form["cs"]) == 0 {
				crsrAsStringx = make([]string, len(uKeyx))
			} else {
				crsrAsStringx = s.R.Form["cs"]
			}
			for i, v := range crsrAsStringx {
				os2, crsrAsString, err := offer.GetByUserID(s.Ctx, v, uKeyx[i].Encode())
				if err != nil && err != datastore.Done {
					log.Printf("Path: %s, Request: get account offers via users keys, Error: %v\n", s.R.URL.Path, err)
					http.Error(s.W, err.Error(),
						http.StatusInternalServerError)
					return
				}
				for i2, v2 := range os2 {
					os[i2] = v2
				}
				if i == 1 {
					q.Set("cs", crsrAsString)
				} else {
					q.Add("cs", crsrAsString)
				}
			}
			URL.RawQuery = q.Encode()
			rb.PrevPageURL = URL.String()
		} else if dID != "" {
			dKey, err := datastore.DecodeKey(dID)
			if err != nil {
				log.Printf("Path: %s, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			crsrAsString := s.R.Form.Get("c")
			os, crsrAsString, err = offer.GetByParent(s.Ctx, crsrAsString, dKey)
			if err != nil && err != datastore.Done {
				log.Printf("Path: %s, Request: get demand offers via demand key, Error: %v\n", s.R.URL.Path, err)
				http.Error(s.W, err.Error(),
					http.StatusInternalServerError)
				return
			}
			q.Set("c", crsrAsString)
			URL.RawQuery = q.Encode()
			rb.PrevPageURL = URL.String()
		} else {
			// For timeline
		}
		if len(os) == 0 {
			s.W.WriteHeader(http.StatusNoContent)
			return
		}
		rb.Result = os
		api.WriteResponseJSON(s, rb)
	}
}
