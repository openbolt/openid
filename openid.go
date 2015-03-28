// OpenID Connect Core 1.0 provider implementation
package openid

import (
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

// OAuthErrors
// ref 18.3
const (
	interaction_required       = 0
	login_required             = 1
	account_selection_required = 2
	consent_required           = 4
	invalid_request_uri        = 8
	invalid_request_object     = 16
	request_not_supported      = 32
	request_uri_not_supported  = 64
	registration_not_supported = 128
)

// Runtimestore is an interface for getting and setting runtime data
type Runtimestore interface {
	// TODO: implement
}

type OpenID struct {
	// Config
	offline_access bool
	token_expire   int // sec
}

// NewProvider returns an blank OpenID Provider instance
func NewProvider() OpenID {
	return OpenID{}
}

var srv *OpenID

// Serve starts the OpenID Provider
func (op *OpenID) Serve() error {
	srv = op
	return nil
}

// TODO: Continue!!!!!!
// /authorize
func (op *OpenID) Authorize(parms Values) (AuthSuccessResp, AuthErrResp) {
	// TODO: Switch to Hybrid, Implicit

	// Authorization Code Flow
	// ref 3.1.2.2
	err1 := validate_oauth_params(parms)                                                                // ref Rule 1
	err2 := validate_scope_param(parms)                                                                 // ref Rule 2
	err3 := validate_req_params(parms, []string{"scope", "response_type", "client_id", "redirect_uri"}) // ref Rule 3
	err4 := validate_sub_param(parms)                                                                   // ref Rule 4

	// Check first part of validation
	if len(err1.Error) != 0 {
		return AuthSuccessResp{}, err1
	}
	if len(err2.Error) != 0 {
		return AuthSuccessResp{}, err2
	}
	if len(err3.Error) != 0 {
		return AuthSuccessResp{}, err3
	}
	if len(err4.Error) != 0 {
		return AuthSuccessResp{}, err4
	}

	// Run through flow
	// ref 3
	switch parms["response_type"][0] {
	case "code":
		return op.authz_code_flow(parms)
	case "id_token", "id_token token":
		return op.implizit_flow(parms)
	case "code id_token", "code token", "code id_token token":
		return op.hybrid_flow(parms)
	default:
		// TODO: invalid_request response
		return AuthSuccessResp{}, AuthErrResp{}
	}
}

func (op *OpenID) authz_code_flow(parms Values) (AuthSuccessResp, AuthErrResp) {
	return AuthSuccessResp{}, AuthErrResp{}
}
func (op *OpenID) implizit_flow(parms Values) (AuthSuccessResp, AuthErrResp) {
	return AuthSuccessResp{}, AuthErrResp{}
}
func (op *OpenID) hybrid_flow(parms Values) (AuthSuccessResp, AuthErrResp) {
	return AuthSuccessResp{}, AuthErrResp{}
}

// /userinfo
func (op *OpenID) UserInfo(tk jwt.Token) jwt.Token {
	return jwt.Token{}
}

// /revoke
func (op *OpenID) Revoke(tk jwt.Token) jwt.Token {
	return jwt.Token{}
}

// /token
func (op *OpenID) TokenEndpoint(tk jwt.Token) jwt.Token {
	return jwt.Token{}
}

// AddServer takes an mux and adds basic http endpoints for OpenID Connect
func (op *OpenID) AddServer(mux *http.ServeMux) error {
	mux.HandleFunc("/authorize", http_authorize)
	mux.HandleFunc("/token", http_token)
	mux.HandleFunc("/userinfo", http_userinfo)
	mux.HandleFunc("/revoke", http_revoke)
	return nil
}

func (op *OpenID) AddClaimsource(ds Claimsource) {
}
