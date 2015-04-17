package openid

import (
	"net/http"
	"time"
)

// TODO: implement
// Ref 3.1.3.3.  Successful Token Response
type AuthSuccessResp struct {
	// Flag if this response is valid, MUST NOT be exported
	ok    bool   `url:"-"`
	Code  string `url:"code,omitempty"`
	State string `url:"state,omitempty"`
}

// Ref 3.1.2.6.  Authentication Error Response
type AuthErrResp struct {
	Error            string `url:"error"`
	ErrorDescription string `url:"error_description,omitempty"`
	ErrorUri         string `url:"error_uri,omitempty"`
	State            string `url:"state,omitempty"`
}

/*
 * Claimsource and Clientsource are compatible with each together
 */

// Claimsource returns claims according to `id`
type Claimsource interface {
	// returns value, ok?
	Get(id, claim, def string) (string, bool)
	//	Set(id, claim, value string) error
	//	Delete(id, claim string) error
	//	DeleteRef(id string) error
}

// Clientsource
type Clientsource interface {
	IsClient(id string) bool
	//	GetClientType(id string) string
	// returns "confidential", "public"
	GetApplType(id string) string
	// returns "web", "user-agent-based", "native"
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
	Authpage(w http.ResponseWriter, r *http.Request) AuthState
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
