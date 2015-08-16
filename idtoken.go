package openid

import (
	"crypto/ecdsa"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// IDToken represents an OpenID Connect id_token
type IDToken struct {
	Token             *jwt.Token
	TokenSignedString string
}

// NewIDToken returns an IDToken according to parameters from Session
func NewIDToken(ses Session, signKey *ecdsa.PrivateKey) (*IDToken, error) {
	// BUG Implement
	tok := new(IDToken)
	//tok.Token = jwt.New(jwt.SigningMethodHS256) // HMAC
	tok.Token = jwt.New(jwt.SigningMethodES256) // ECDSA
	tok.Token.Claims["foo"] = "bar"
	tok.Token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	var err error
	tok.TokenSignedString, err = tok.Token.SignedString(signKey)
	return tok, err
}

// MarshalText is used to satisfy the encoding.TextMarshaler interface.
// It returns the IDToken as a byte slice. This way it is posible to easily serialize an IDToken
func (t *IDToken) MarshalText() (text []byte, err error) {
	return []byte(t.TokenSignedString), nil
}
