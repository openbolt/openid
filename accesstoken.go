package openid

import "time"

// AccessToken represents an OAuth 2.0 access_token
// BUG Implement an signed access-token. So there is no need to cache/lookup this
// on the resource server
type AccessToken struct {
	Token     string
	TokenType string
	ExpiresIn time.Duration
}

func (t *AccessToken) Load(ses Session) *AccessToken {
	// BUG: Not implemented
	return &AccessToken{}
}
