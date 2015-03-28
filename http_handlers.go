package openid

import (
	"net/http"
)

// /authorize
// ref 3.1.2.1
func http_authorize(w http.ResponseWriter, r *http.Request) {
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

	srv.Authorize(parms)
	// TODO: Response handling
}

// /token
func http_token(w http.ResponseWriter, r *http.Request) {
}

// /userinfo
func http_userinfo(w http.ResponseWriter, r *http.Request) {
}

// /revoke
func http_revoke(w http.ResponseWriter, r *http.Request) {
}
