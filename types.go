package openid

import (
	"net/http"
	"net/url"
	"time"
)

/* START stolen from net/url, TODO: Find better method */
type Values url.Values

// Get gets the first value associated with the given key.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (v Values) Get(key string) string {
	if v == nil {
		return ""
	}
	vs, ok := v[key]
	if !ok || len(vs) == 0 {
		return ""
	}
	return vs[0]
}

// Set sets the key to value. It replaces any existing
// values.
func (v Values) Set(key, value string) {
	v[key] = []string{value}
}

// Add adds the value to key. It appends to any existing
// values associated with key.
func (v Values) Add(key, value string) {
	v[key] = append(v[key], value)
}

// Del deletes the values associated with key.
func (v Values) Del(key string) {
	delete(v, key)
}

/* END stolen */

// TODO: implement
// Ref 3.1.3.3.  Successful Token Response
type AuthSuccessResp struct {
	// Flag if this response is valid, MUST NOT be exported
	ok     bool
	Params Values
}

// Ref 3.1.2.6.  Authentication Error Response
type AuthErrResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`
}

/*
 * Authsource, Claimsource and Clientsource are compatible with each together
 */
// Authsource implements the authenticating part
type Authsource interface {
	Auth(id string, params Values) AuthErrResp
	IsAuthenticated(params, header Values) (string, error)
	Revoke(id string)
	Register(id string, params Values) error
	Unregister(id string) error
}

// Claimsource returns claims according to `id`
type Claimsource interface {
	// returns value, ok?
	Get(id, claim, def string) (string, bool)
	Set(id, claim, value string) error
	Delete(id, claim string) error
	DeleteRef(id string) error
}

// Clientsource
type Clientsource interface {
	IsClient(id string) bool
	GetApplType(id string) string
	// returns "web", "native"
	ValidateRedirectUri(id, uri string) bool
}

// EnduserIf is used for rendering enduser dialogs
// Authpage must comply with `3.1.2.1.  Authentication Request`
// In short, it should implement the following:
//   - Claims: display, prompt, max_age, ui_locales, id_token_hint, login_hint, acr_values
type EnduserIf interface {
	// Q: Is user already authenticated/able to authN automaticaly?
	//   Y: Return AuthState:AuthOk=true
	//   N: Prompt for creds, set session
	// The redirect will be handled outside
	Authpage(w http.ResponseWriter, r *http.Request, params Values) AuthState
}
type AuthState struct {
	AuthOk        bool
	AuthAbort     bool
	AuthFailed    bool
	AuthPrompting bool
	Iss           string
	Sub           string
	AuthTime      time.Time
}
