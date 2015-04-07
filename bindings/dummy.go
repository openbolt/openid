package bindings

import (
	"net/http"
	"strings"

	"github.com/openbolt/openid"
)

// DummySource defines an empty Authsource, Claimsource and Clientsource
// which stores its data in RAM.
// This can be used for testcases or demo applications
type DummySource struct {
}

/*
 *Authsource
 */
// Auth
func (ds DummySource) Auth(id string, params openid.Values) openid.AuthErrResp {
	// BUG(djboris) Implement
	return openid.AuthErrResp{}
}
func (ds DummySource) Revoke(id string) {
	// BUG(djboris) Implement
}
func (ds DummySource) Register(id string, params openid.Values) error {
	// BUG(djboris) Implement
	return nil
}
func (ds DummySource) Unregister(id string) error {
	// BUG(djboris) Implement
	return nil
}
func (ds DummySource) IsAuthenticated(params, header openid.Values) (string, error) {
	return "", nil
}

/*
 * Claimsource
 */
// Get
func (ds DummySource) Get(id, claim, def string) (string, bool) {
	// BUG(djboris) Implement
	return "foo", true
}
func (ds DummySource) Set(id, claim, value string) error {
	// BUG(djboris) Implement
	return nil
}
func (ds DummySource) Delete(id, claim string) error {
	// BUG(djboris) Implement
	return nil
}
func (ds DummySource) DeleteRef(id string) error {
	// BUG(djboris) Implement
	return nil
}

/*
 * Clientsource
 */
// IsClient returns true if id begins with `clt`
func (ds DummySource) IsClient(id string) bool {
	return strings.HasPrefix(id, "clt")
}
func (ds DummySource) GetApplType(id string) string {
	// BUG(djboris) Implement
	return "web"
}
func (ds DummySource) ValidateRedirectUri(id, uri string) bool {
	// BUG(djboris) Implement
	return true
}

/*
 * EnduserIf
 */
func (ds *DummySource) Authpage(w http.ResponseWriter, r *http.Request, params openid.Values) openid.AuthState {
	w.Write([]byte("demopage"))
	return openid.AuthState{AuthPrompting: true}
}
