package openid

import (
	"net/http"
	"time"
)

// AuthSuccessResp holds all parameters which can be returned to the user if nothing fails
// Ref 3.1.3.3.  Successful Token Response
type AuthSuccessResp struct {
	// Flag if this response is valid, MUST NOT be exported
	ok    bool   `url:"-"`
	Code  string `url:"code,omitempty"`
	State string `url:"state,omitempty"`
}

// AuthErrResp holds all parameters which can be returned to the user in error case
// Ref 3.1.2.6.  Authentication Error Response
type AuthErrResp struct {
	Error            string `url:"error"`
	ErrorDescription string `url:"error_description,omitempty"`
	ErrorURI         string `url:"error_uri,omitempty"`
	State            string `url:"state,omitempty"`
}

/*
 * Cacher, Claimsource and Clientsource are compatible with each together
 */

// Cacher is used to cache sessions between code request and id_token retrieval
type Cacher interface {
	Cache(val AuthzCodeSession) error
	Retire(code string)
}

// Claimsource returns claims according to `id`
type Claimsource interface {
	// returns value, ok?
	Get(id, claim, def string) (string, bool)
	//	Set(id, claim, value string) error
	//	Delete(id, claim string) error
	//	DeleteRef(id string) error
}

// Clientsource is the databinding for OAuth 2.0 clients
type Clientsource interface {
	IsClient(id string) bool
	//	GetClientType(id string) string
	// returns "confidential", "public"
	GetApplType(id string) string
	// returns "web", "user-agent-based", "native"
	ValidateRedirectURI(id, uri string) bool
}

// EnduserIf is used for rendering enduser dialogs
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

// AuthState is used as the return value of EnduserIf
type AuthState struct {
	AuthOk        bool
	AuthAbort     bool
	AuthFailed    bool
	AuthPrompting bool
	Iss           string
	Sub           string
	AuthTime      time.Time
	Acr           string
	Amr           string
}

// AuthzCodeSession stores information about pending "code" requests
type AuthzCodeSession struct {
	Code     string
	ClientID string
	Nonce    string
	AuthTime time.Time

	// When max_age is used, the ID Token returned MUST
	// include an auth_time Claim Value.
	MaxAge time.Duration

	Acr           string
	ClaimsLocales string
	Claims        interface{} // Ref 5.5
}
