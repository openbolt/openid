package bindings

import (
	"github.com/openbolt/openid"
)

// DemoSource defines an empty LoginSource and ClaimSource
// which stores its data in RAM. This can be used for testcases or demo applications
type DemoSource interface {
	// Loginsource
	Login(id string, params openid.Values) openid.AuthErrResp
	Revoke(id string)
	Register(id string, params openid.Values) error
	Unregister(id string) error

	// Claimsource
	Get(id, claim, def string) (string, bool)
	Set(id, claim, value string) error
	Delete(id, claim string) error
	DeleteRef(id string) error
}
