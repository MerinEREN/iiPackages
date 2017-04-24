package cookie

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	// "github.com/MerinEREN/iiPackages/user"
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
func Set(w http.ResponseWriter, r *http.Request, name, value string) error {
	// COOKIE IS A PART OF THE HEADER, SO U SHOULD SET THE COOKIE BEFORE EXECUTING A
	// TEMPLATE OR WRITING SOMETHING TO THE BODY !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	c, err := r.Cookie(name)
	if err == http.ErrNoCookie {
		c, err = create(name, value)
		http.SetCookie(w, c)
	} else {
		if isUserDataChanged(c) {
			// DELETING CORRUPTED COOKIE AND CREATING NEW ONE !!!!!!!!!!!!!!!!!
			Delete(w, r, name)
			c, _ = create(name, value)
			http.SetCookie(w, c)
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
	// log.Printf("Cookie is: %v. Error is: %v.", c, err)
	return
}

func Delete(w http.ResponseWriter, r *http.Request, n string) error {
	c, err := r.Cookie(n)
	if err == http.ErrNoCookie {
		return err
	}
	c.MaxAge = -1
	// If path is different can't delete cookie without cookie's path.
	// Maybe should use cookie path even paths are same.
	c.Path = r.URL.Path
	http.SetCookie(w, c)
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
	// log.Printf("Marshalled cookie data is %s\n", string(bs))
	// log.Printf("Cookie value for "+c.Name+" is: %s\n", c.Value)
	c.Value += "|" + base64.StdEncoding.EncodeToString(bs)
	code := GetHmac(c.Value)
	c.Value += "|" + code
	// log.Printf("Cookie value for "+c.Name+" is: %s\n", c.Value)
	return
}

func isUserDataChanged(c *http.Cookie) bool {
	cvSlice := strings.Split(c.Value, "|")
	uuidData := cvSlice[0] + "|" + cvSlice[1]
	returnedCode := GetHmac(uuidData)
	if returnedCode != cvSlice[2] {
		log.Printf("%s cookie value is corrupted. Cookie HMAC is %s, "+
			"genereted HMAC is %s", c.Name, cvSlice[2], returnedCode)
		returnedCookieData := decodeThanUnmarshall(cvSlice[1])
		log.Printf("Returned cookie data is %v", returnedCookieData)
		return true
	}
	return false
}

// MAKE GENERIC RETURN TYPE !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func GetData(r *http.Request, n string) (*SessionData, error) {
	c, err := r.Cookie(n)
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

func GetHmac(i interface{}) string {
	h := hmac.New(sha256.New, []byte("someKey"))
	s, ok := i.(string)
	if ok {
		io.WriteString(h, s)
	}
	var r io.Reader
	r, ok = i.(io.Reader)
	if ok {
		io.Copy(h, r)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
