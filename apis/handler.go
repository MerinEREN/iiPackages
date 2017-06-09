package apis

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, r *http.Request, rb *ResponseBody) {
	bs, err := json.Marshal(rb)
	if err != nil {
		log.Printf("Path: %s, Error: %v\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(bs)
}
