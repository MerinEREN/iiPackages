/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package accountSettings

import (
	"github.com/MerinEREN/iiPackages/session"
)

func Handler(s *session.Session) {
	/* u, uKey, err := usr.Exist(ctx, ug.Email)
	if err == usr.FindUserError {
		log.Printf("Error while login user: %v\n", err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if u.Status == "suspended" {
		log.Printf("Suspended user %s trying to see "+
			"%s page !!!", u.Email, s)
		http.Error(w, "You are suspended", http.StatusForbidden)
		// fmt.Fprintf(w, "Permission denied !!!")
		return
	} else if !u.IsAdmin() {
		log.Printf("Unauthorized user %s trying to see "+
			"account settings page !!!", u.Email)
		// fmt.Fprintf(w, "Permission denied !!!")
		http.Error(w, "some error message", http.StatusUnauthorized)
		return
	}
	aKey := uKey.Parent()
	acc, err := account.Get(ctx, aKey)
	if err != nil {
		log.Printf("Error while getting user's account data: %v\n",
			err)
		// ALSO LOG THIS WITH DATASTORE LOG !!!!!!!!!!!!!!!!!!!!!!
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(acc)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Always send corresponding header values instead of defaults !!!!
	w.Write(b) */
}
