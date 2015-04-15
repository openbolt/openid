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
	Claims  Claimsource
	Auth    Authsource
	Client  Clientsource
	Enduser EnduserIf

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
	if op.Claims == nil {
		return errors.New("No Claimsource defined")
	}
	if op.Auth == nil {
		return errors.New("No Authsource defined")
	}
	if op.Client == nil {
		return errors.New("No Clientsource defined")
	}
	if op.Enduser == nil {
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
	return nil
}
