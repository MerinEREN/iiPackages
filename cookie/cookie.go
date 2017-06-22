/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package cookie

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/MerinEREN/iiPackages/crypto"
	"github.com/MerinEREN/iiPackages/session"
	"log"
	"net/http"
	"strings"
)

// Cookie error variables
var (
	ErrCorruptedCookie = errors.New("Cookie data corrupted")
)

// CHANGE THIS DUMMY COOKIE STRUCT !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type SessionData struct {
	Photo string
}

// Adding uuid and hash to the cookie and check hash code
func Set(s *session.Session, name, value string) error {
	// COOKIE IS A PART OF THE HEADER, SO U SHOULD SET THE COOKIE BEFORE EXECUTING A
	// TEMPLATE OR WRITING SOMETHING TO THE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	c, err := s.R.Cookie(name)
	if err == http.ErrNoCookie {
		c, err = create(name, value)
		http.SetCookie(s.W, c)
	} else {
		if isUserDataChanged(c) {
			// DELETING CORRUPTED COOKIE AND CREATING NEW ONE !!!!!!!!!!!!!!!!!
			Delete(s, name)
			c, _ = create(name, value)
			http.SetCookie(s.W, c)
			err = ErrCorruptedCookie
		}
	}
	return err
}

func create(n, v string) (c *http.Cookie, err error) {
	c = &http.Cookie{
		Name: n,
		// U CAN USE UUID AS VALUE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		Value: v,
		// NOT GOOD PRACTICE
		// ADDING USER DATA TO A COOKIE
		// WITH NO WAY OF KNOWING WHETER OR NOT THEY MIGHT HAVE ALTERED
		// THAT DATA !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// HMAC WOULD ALLOW US TO DETERMINE WHETHER OR NOT THE DATA IN THE
		// COOKIE WAS ALTERED !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// HOWEVER, BEST TO STORE USER DATA ON THE SERVER AND KEEP
		// BACKUPS !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// Value: "emil = merin@inceis.net" + "JSON data" + "whatever",
		// IF SECURE IS TRUE THIS COOKIE ONLY SEND WITH HTTP2 !!!!!!!!!!!!!!!!!!!!!
		// Secure: true,
		// HttpOnly: true MEANS JAVASCRIPT CAN NOT ACCESS THE COOKIE !!!!!!!!!!!!!!
		HttpOnly: false,
	}
	err = setValue(c)
	return
}

func Delete(s *session.Session, n string) error {
	c, err := s.R.Cookie(n)
	if err == http.ErrNoCookie {
		return err
	}
	c.MaxAge = -1
	// If path is different can't delete cookie without cookie's path.
	// Maybe should use cookie path even paths are same.
	c.Path = s.R.URL.Path
	http.SetCookie(s.W, c)
	return err
}

// Setting different kind of struct for different cookies
func setValue(c *http.Cookie) (err error) {
	var cd interface{}
	if strings.Contains(c.Name, "/") {
		cd = SessionData{
			Photo: "img/MKA.jpg",
		}
	}
	var bs []byte
	bs, err = json.Marshal(cd)
	if err != nil {
		return
	}
	c.Value += "|" + base64.StdEncoding.EncodeToString(bs)
	code, err := crypto.GetMAC(c.Value)
	if err != nil {
		return
	}
	c.Value += "|" + code
	return
}

func isUserDataChanged(c *http.Cookie) bool {
	cvSlice := strings.Split(c.Value, "|")
	uuidData := cvSlice[0] + "|" + cvSlice[1]
	if !crypto.CheckMAC(uuidData, cvSlice[2]) {
		log.Printf("%s cookie value is corrupted.", c.Name)
		return true
	}
	return false
}

// MAKE GENERIC RETURN TYPE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func GetData(s session.Session, n string) (*SessionData, error) {
	c, err := s.R.Cookie(n)
	if err == http.ErrNoCookie {
		return &SessionData{}, err
	}
	cvSlice := strings.Split(c.Value, "|")
	return decodeThanUnmarshall(cvSlice[1]), nil
}

// MAKE GENERIC RETURN TYPE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func decodeThanUnmarshall(cd string) *SessionData {
	decodedBase64, err := base64.StdEncoding.DecodeString(cd)
	if err != nil {
		log.Printf("Error while decoding cookie data. Error is %v\n", err)
	}
	var cookieData SessionData
	err = json.Unmarshal(decodedBase64, &cookieData)
	if err != nil {
		log.Printf("Cookie data unmarshaling error. %v\n", err)
	}
	return &cookieData
}
