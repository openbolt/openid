// OpenID Connect Core 1.0 provider implementation
package openid

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/openbolt/openid/utils"
)

type OpenID struct {
	// Datasources
	claims *Claimsource
	auth   *Authsource
	client *Clientsource

	// Config... going on
}

// NewProvider returns an blank OpenID Provider instance
func NewProvider() OpenID {
	return OpenID{}
}

// Serve starts the OpenID Provider
// TODO: Only accept requests when this is called
func (op *OpenID) Serve() error {
	if op.claims == nil {
		return errors.New("No claimsource defined")
	}
	if op.auth == nil {
		return errors.New("No authsource defined")
	}
	if op.client == nil {
		return errors.New("No clientsource defined")
	}

	return nil
}

// Ref 3.1.2.1. Authentication Request
// An Authentication Request is an OAuth 2.0 Authorization Request that requests
// that the End-User be authenticated by the Authorization Server.
func (op *OpenID) Authorize(parms Values) (AuthSuccessResp, AuthErrResp) {
	// ref 3.1.2.2
	err1 := validate_oauth_params(parms)           // ref Rule 1
	err2 := validate_scope_param(parms)            // ref Rule 2
	err3 := validate_req_params(parms, *op.client) // ref Rule 3
	err4 := validate_sub_param(parms)              // ref Rule 4

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

// /userinfo
func (op *OpenID) UserInfo(tk jwt.Token) jwt.Token {
	// TODO: Implement
	return jwt.Token{}
}

// /revoke
func (op *OpenID) Revoke(tk jwt.Token) jwt.Token {
	// TODO: Implement
	return jwt.Token{}
}

// /token
func (op *OpenID) TokenEndpoint(tk jwt.Token) jwt.Token {
	// TODO: Implement
	return jwt.Token{}
}

// AddServer takes an mux and adds basic http endpoints for OpenID Connect
func (op *OpenID) AddServer(mux *http.ServeMux) error {
	api, err := newHttpAPI(op)
	if err != nil {
		utils.ELog(err)
		return err
	}

	mux.HandleFunc("/authorize", api.http_authorize)
	mux.HandleFunc("/token", api.http_token)
	mux.HandleFunc("/userinfo", api.http_userinfo)
	mux.HandleFunc("/revoke", api.http_revoke)
	return nil
}

func (op *OpenID) SetClaimsource(src *Claimsource) {
	op.claims = src
}

func (op *OpenID) SetAuthsource(src *Authsource) {
	op.auth = src
}

func (op *OpenID) SetClientsource(src *Clientsource) {
	op.client = src
}
