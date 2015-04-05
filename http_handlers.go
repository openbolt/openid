package openid

import (
	"encoding/json"
	"net/http"
)

type httpAPI struct {
	srv *OpenID
}

func newHttpAPI(srv *OpenID) (*httpAPI, error) {
	api := new(httpAPI)
	api.srv = srv

	return api, nil
}

// /authorize
// ref 3.1.2.1
func (api *httpAPI) http_authorize(w http.ResponseWriter, r *http.Request) {
	// Return if Method not GET or POST
	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Method must be GET or POST"))
		return
	}

	// Extract paramters
	var parms Values

	if r.Method == "GET" {
		// MUST URI Query String Serialization
		parms = Values(r.URL.Query())

	} else if r.Method == "POST" {
		// MUST Form Serialization
		r.ParseForm()
		parms = Values(r.PostForm)
	}

	resp, err := api.srv.Authorize(w, r, parms)

	// BUG(djboris) Proper response handling
	if len(err.Error) > 0 {
		r, _ := json.Marshal(err)
		w.Write(r)
	} else {
		r, _ := json.Marshal(resp)
		w.Write(r)
	}
	// ENDBUG
}

// /token
func (api *httpAPI) http_token(w http.ResponseWriter, r *http.Request) {
}

// /userinfo
func (api *httpAPI) http_userinfo(w http.ResponseWriter, r *http.Request) {
}

// /revoke
func (api *httpAPI) http_revoke(w http.ResponseWriter, r *http.Request) {
}
