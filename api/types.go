/*
Package api contains response struct that API returns to client requests,
and also has sub packages which have handlers for all URLs.
*/
package api

// ResponseBody represends api response body.
// Reset resets corresponding client datas if it is true.
// If POST always set the reset true, because the ID is created by server side.
// If PUT only reset if the modified data needed from the backand
// Like an image needed from the cloud storage.
type ResponseBody struct {
	Result      data   `json:"result"`
	Reset       bool   `json:"reset"`
	NextPageURL string `json:"nextPageURL"`
	PrevPageURL string `json:"prevPageURL"`
}

type data interface{}
