/*
apis Package contains response struct that API returns to client requests,
and also has sub packages which have handlers for all URLs.
*/
package apis

// response represends api response body.
type ResponseBody struct {
	Result      data   `json:"result"`
	NextPageUrl string `json:"nextPageUrl"`
}

type data interface{}
