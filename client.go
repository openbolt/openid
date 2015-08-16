package openid

import "net/http"

// AuthenticateClient returns true when the client is successfuly authenticated as defined in OAuth 2.0 Spec
// 9.  Client Authentication
func (op *OpenID) AuthenticateClient(client_id string, req *http.Request) (bool, int) {
	// BUG Implement
	return true, 0

	// false, REQUIRE_401
	// false, CLIENT_NOT_ALLOWED
}
