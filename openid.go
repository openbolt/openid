// Package OpenID implements the OpenID Core 1.0 provider

package openid

import (
	"errors"
	"net/http"

	"github.com/openbolt/openid/utils"
)

// OpenID implements the OpenID Provider (OP)
type OpenID struct {
	// Datasources
	Claimsrc  Claimsource
	Clientsrc Clientsource
	Enduser   EnduserIf

	// True, if server is fully started
	serving bool
}

// NewProvider returns an blank OpenID Provider instance
func NewProvider() *OpenID {
	op := new(OpenID)
	op.serving = false
	return op
}

// Serve starts the OpenID Provider
func (op *OpenID) Serve() error {
	if op.Claimsrc == nil {
		return errors.New("No Claimsource defined")
	}
	if op.Clientsrc == nil {
		return errors.New("No Clientsource defined")
	}
	if op.Enduser == nil {
		return errors.New("No EnduserIf defined")
	}
	op.serving = true

	return nil
}

// AddServer takes an mux and adds basic http endpoints for OpenID Connect
func (op *OpenID) AddServer(mux *http.ServeMux) error {
	api, err := newAPI(op)
	if err != nil {
		utils.ELog(err)
		return err
	}

	mux.HandleFunc("/authorize", api.Authorize)
	return nil
}
