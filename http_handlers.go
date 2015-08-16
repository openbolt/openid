package openid

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/context"
	"github.com/openbolt/openid/utils"
	"github.com/pborman/uuid"
)

type httpAPI struct {
	srv *OpenID
}

func newAPI(srv *OpenID) (*httpAPI, error) {
	api := new(httpAPI)
	api.srv = srv

	return api, nil
}

// /authorize
// ref 3.1.2.1
func (api *httpAPI) Authorize(w http.ResponseWriter, r *http.Request) {
	context.Set(r, REQUEST_UUID, string(uuid.NewUUID().String()))

	// Return if Method not GET or POST
	if r.Method != "GET" && r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Method must be GET or POST"))
		return
	}

	// Run the authorization
	resp, err := api.srv.Authorize(w, r)

	// Get default response_mode for flow and override it if another is set
	var responseMode string
	if getFlow(GetParam(r, "response_type")) == "authorization_code" {
		responseMode = "query"
	} else {
		responseMode = "fragment"
	}
	utils.EDebug(errors.New("Using response_mode "+responseMode), r)

	if err.Error != "" {
		utils.ELog(errors.New("Auth failed: "+err.Error), r)

		// If redirect_uri is not valid, show error as JSON
		redirectURI := GetParam(r, "redirect_uri")
		clientID := GetParam(r, "client_id")
		flow := GetParam(r, "code")
		t := checkRedirectURI(redirectURI, clientID, flow, api.srv.Clientsrc)
		u, e := url.Parse(redirectURI)
		if e != nil || !t {
			utils.EDebug(e, r)
			r, _ := json.Marshal(err)
			w.Write(r)
			return
		}

		/*
		 * Add error to query or fragment
		 */
		err.State = GetParam(r, "state")
		*u, e = serializeResponse(*u, responseMode, err)
		if e != nil {
			utils.EDebug(e, r)
			r, _ := json.Marshal(err)
			w.Write(r)
			return
		}

		/*
		 * Finish, Do 302 Redirect
		 */
		utils.EDebug(errors.New("Redirecting to "+u.String()), r)
		http.Redirect(w, r, u.String(), http.StatusFound)
	} else if resp.ok {
		utils.EDebug(errors.New("Auth succeeded"), r)

		// Return success
		redirectURI := GetParam(r, "redirect_uri")
		u, _ := url.Parse(redirectURI)
		*u, _ = serializeResponse(*u, responseMode, resp)

		utils.EDebug(errors.New("Redirecting to "+u.String()), r)
		http.Redirect(w, r, u.String(), http.StatusFound)
	}
}

// 3.1.3.  Token Endpoint
// Must use TLS
func (api *httpAPI) Token(w http.ResponseWriter, r *http.Request) {
	context.Set(r, REQUEST_UUID, string(uuid.NewUUID().String()))

	// Return if Method not POST
	if r.Method != "POST" {
		err := AuthErrResp{
			Error:            "invalid_request",
			ErrorDescription: "Method must be POST",
		}
		w.WriteHeader(http.StatusBadRequest)
		r, _ := json.Marshal(err)
		w.Write(r)
		return
	}

	// Return if not HTTPS BUG: This doesn't work :(
	if false { //|| r.URL.Scheme != "https" {
		fmt.Println(r.URL.Scheme)
		err := AuthErrResp{
			Error:            "invalid_request",
			ErrorDescription: "TLS is required",
		}
		w.WriteHeader(http.StatusBadRequest)
		r, _ := json.Marshal(err)
		w.Write(r)
		return
	}

	resp, err := api.srv.Token(w, r)
	if err.Error != "" && err.StatusCode == 0 {
		data, _ := json.Marshal(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		utils.EDebug(errors.New(string(data)), r)
	} else if err.Error != "" {
		// TODO: Return StatusCode and Headers from struct
		data, _ := json.Marshal(resp)
		utils.EDebug(errors.New("X"+string(data)), r)
	} else if resp.ok {
		// 3.1.3.3.  Successful Token Response
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Pragma", "no-cache")
		data, _ := json.Marshal(resp)
		w.Write(data)
		utils.EDebug(errors.New(string(data)), r)
	} else {
		//BUG: Return 500
	}
}
