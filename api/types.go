/*
Package api contains response struct that API returns to client requests,
and also has sub packages which have handlers for all URLs.
*/
package api

// ResponseBody represends api response body.
// "Reset" resets corresponding client datas if the method is "GET" and has filters for
// purposes like serch.
type ResponseBody struct {
	Result      data   `json:"result"`
	Reset       bool   `json:"reset"`
	NextPageURL string `json:"nextPageURL"`
	PrevPageURL string `json:"prevPageURL"`
}

type data interface{}
