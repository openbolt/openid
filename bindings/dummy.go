package bindings

import (
	"net/url"
	"strings"
)

//go:generate go-bindata -o bindata.go -pkg bindings assets/

// DummySource defines an empty Claimsource and Clientsource
// which stores its data in RAM.
// This can be used for testcases or demo applications
type DummySource struct {
}

/*
 * Claimsource
 */
// Get
func (ds DummySource) Get(id, claim, def string) (string, bool) {
	// BUG(djboris) Implement
	return "foo", true
}

/*
 * Clientsource
 */
// IsClient returns true if id begins with `clt`
func (ds DummySource) IsClient(id string) bool {
	return strings.HasPrefix(id, "clt")
}
func (ds DummySource) GetClientType(id string) string {
	// BUG(djboris) Implement
	return "confidential"
}
func (ds DummySource) GetApplType(id string) string {
	// BUG(djboris) Implement
	return "web"
}

// Returns true if host == localhost
func (ds DummySource) ValidateRedirectURI(id, uri string) bool {
	u, _ := url.Parse(uri)
	return u.Host == "localhost:8443" || u.Host == "localhost:8080"
}
