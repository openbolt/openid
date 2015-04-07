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
	claims  Claimsource
	auth    Authsource
	client  Clientsource
	enduser EnduserIf

	// True, if server is fully started
	serving bool

	// Config... going on
}

// NewProvider returns an blank OpenID Provider instance
func NewProvider() *OpenID {
	op := new(OpenID)
	op.serving = false
	return op
}

// Serve starts the OpenID Provider
func (op *OpenID) Serve() error {
	if op.claims == nil {
		return errors.New("No Claimsource defined")
	}
	if op.auth == nil {
		return errors.New("No Authsource defined")
	}
	if op.client == nil {
		return errors.New("No Clientsource defined")
	}
	if op.enduser == nil {
		return errors.New("No EnduserIf defined")
	}
	op.serving = true

	return nil
}

// /userinfo
func (op *OpenID) UserInfo(tk jwt.Token) jwt.Token {
	if !op.serving {
		return jwt.Token{}
	}
	// TODO: Implement
	return jwt.Token{}
}

// /revoke
func (op *OpenID) Revoke(tk jwt.Token) jwt.Token {
	if !op.serving {
		return jwt.Token{}
	}
	// TODO: Implement
	return jwt.Token{}
}

// /token
func (op *OpenID) TokenEndpoint(tk jwt.Token) jwt.Token {
	if !op.serving {
		return jwt.Token{}
	}
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

	mux.HandleFunc("/authorize", api.httpAuthorize)
	mux.HandleFunc("/token", api.http_token)
	mux.HandleFunc("/userinfo", api.http_userinfo)
	mux.HandleFunc("/revoke", api.http_revoke)
	return nil
}

func (op *OpenID) SetClaimsource(src Claimsource) {
	op.claims = src
}

func (op *OpenID) SetAuthsource(src Authsource) {
	op.auth = src
}

func (op *OpenID) SetClientsource(src Clientsource) {
	op.client = src
}
func (op *OpenID) SetEnduserIf(src EnduserIf) {
	op.enduser = src
}
