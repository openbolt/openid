package openid

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/openbolt/openid/utils"
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
func (api *httpAPI) httpAuthorize(w http.ResponseWriter, r *http.Request) {
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

	// BUG return according to response_mode [fragment, query]
	if len(err.Error) > 0 {
		// If redirect_uri is not valid, show error as JSON
		u, e := url.Parse(parms.Get("redirect_uri"))
		if e != nil || u.String() == "" {
			utils.EDebug(e)
			r, _ := json.Marshal(err)
			w.Write(r)
			return
		}
		// If query argument is not valid for redirect_uri
		query, e := url.ParseQuery(u.RawQuery)
		if e != nil || u.String() == "" {
			utils.EDebug(e)
			r, _ := json.Marshal(err)
			w.Write(r)
			return
		}

		// Add error to query
		query.Add("error", err.Error)
		query.Add("error_description", err.ErrorDescription)
		if state := parms.Get("state"); state != "" {
			query.Add("state", state)
		}
		u.RawQuery = query.Encode()

		// Do 302 Redirect
		http.Redirect(w, r, u.String(), http.StatusFound)
	} else if resp.ok {
		// BUG(djboris) Proper success response handling
		r, _ := json.Marshal(resp)
		w.Write(r)
	} else {
		// Simply return if nor success, nor error
		return
	}
}

// /token
func (api *httpAPI) http_token(w http.ResponseWriter, r *http.Request) {
}

// /userinfo
func (api *httpAPI) http_userinfo(w http.ResponseWriter, r *http.Request) {
}
