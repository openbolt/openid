package openid

// IDToken represents an OpenID Connect id_token
type IDToken struct {
}

// NewIDToken returns an IDToken according to parameters from Session
func NewIDToken(ses Session) *IDToken {
	return nil
}

// MarshalText is used to satisfy the encoding.TextMarshaler interface.
// It returns the IDToken as a byte slice. This way it is posible to easily serialize an IDToken
func (t *IDToken) MarshalText() (text []byte, err error) {
	return []byte("Not Implemented"), nil
}
