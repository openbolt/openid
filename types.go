package openid

import (
	"net/url"
)

// TODO: Better way
type Values url.Values

// TODO: Get ref and implement
type AuthSuccessResp struct {
	// Flag if this response is valid, MUST NOT be exported
	ok bool
}

// ref 3.1.2.6
type AuthErrResp struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorUri         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`
}

/*
 * Loginsource and Claimsource are compatible with each together
 */
// Loginsource implements the authenticating part
type Loginsource interface {
	Login(id string, params Values) AuthErrResp
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
