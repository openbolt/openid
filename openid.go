// Package OpenID implements the OpenID Core 1.0 provider

package openid

import (
	"crypto/ecdsa"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/openbolt/openid/utils"
)

// OpenID implements the OpenID Provider (OP)
type OpenID struct {
	// Datasources
	Claimsrc  Claimsource
	Clientsrc Clientsource
	Enduser   EnduserIf
	Cache     Cacher

	// Initialized on OpenID.Serve()
	AccessTokenSignKeyFile string
	accessTokenSignKey     *ecdsa.PrivateKey

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
	if op.Cache == nil {
		return errors.New("No Cache defined")
	}

	// Load AccessToken Sign Key
	raw, err := ioutil.ReadFile(op.AccessTokenSignKeyFile)
	if err != nil {
		return errors.New("Cannot read AccessTokenSignKeyFile: " + err.Error())
	}
	key, err := loadSigningKey(raw)
	if err != nil {
		return err
	}
	op.accessTokenSignKey = key

	// Activate
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
	mux.HandleFunc("/token", api.Token)
	return nil
}
