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

	resp, err := api.srv.Authorize(w, r)

	// Get default response_type for flow and override it if another is set
	var response_type string = "fragment"
	if getFlow(GetParam(r, "code")) == "authorization_code" {
		response_type = "query"
	}
	if tmp := GetParam(r, "response_type"); tmp != "" {
		response_type = tmp
	}

	if err.Error != "" {
		// If redirect_uri is not valid, show error as JSON
		redirect_uri := GetParam(r, "redirect_uri")
		client_id := GetParam(r, "client_id")
		flow := GetParam(r, "code")
		t := checkRedirectUri(redirect_uri, client_id, flow, api.srv.Clientsrc)
		u, e := url.Parse(redirect_uri)
		if e != nil || !t {
			utils.EDebug(e)
			r, _ := json.Marshal(err)
			w.Write(r)
			return
		}

		/*
		 * Add error to query or fragment
		 */
		var query url.Values
		if response_type == "query" {
			query, _ = url.ParseQuery(u.RawQuery)
		} else {
			query, _ = url.ParseQuery(u.Fragment)
		}

		// Add error to query
		query.Add("error", err.Error)
		query.Add("error_description", err.ErrorDescription)
		if state := GetParam(r, "state"); state != "" {
			query.Add("state", state)
		}

		if response_type == "query" {
			u.RawQuery = query.Encode()
		} else {
			u.Fragment = query.Encode()
		}

		/*
		 * Finish, Do 302 Redirect
		 */
		http.Redirect(w, r, u.String(), http.StatusFound)
	} else if resp.ok {
		// BUG(djboris) Proper success response handling
		r, _ := json.Marshal(resp)
		w.Write(r)
	}
}
