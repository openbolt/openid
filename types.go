package openid

import (
	"net/url"
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
	ok bool
}

// Ref 3.1.2.6.  Authentication Error Response
type AuthErrResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`
}

/*
 * Authsource and Claimsource are compatible with each together
 */
// Authsource implements the authenticating part
type Authsource interface {
	Auth(id string, params Values) AuthErrResp
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
