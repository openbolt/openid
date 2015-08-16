package openid

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"hash"
	"io"
	"math/big"
	"time"

	"github.com/openbolt/openid/utils"
)

// AccessToken represents an OAuth 2.0 access_token
// BUG Implement an signed access-token. So there is no need to cache/lookup this
// on the resource server
type AccessToken struct {
	Token     string
	TokenType string
	ExpiresIn time.Duration
}

type AccessTokenPayload struct {
	ClientID string
	Scope    string
	AuthTime time.Time
	Validity time.Time
}

func (t *AccessToken) Load(ses Session, signkey *ecdsa.PrivateKey) *AccessToken {
	t.TokenType = "bearer"
	t.ExpiresIn = 300 // TODO: Make this param

	// Generate payload
	payload := AccessTokenPayload{
		ClientID: ses.ClientID,
		Scope:    ses.Scope,
		AuthTime: ses.AuthTime,
		Validity: time.Now().Add(t.ExpiresIn * time.Second),
	}

	data, _ := json.Marshal(payload)
	t.Token = base64.StdEncoding.EncodeToString(data)

	// Sign
	var h hash.Hash
	h = sha256.New()
	r := big.NewInt(0)
	s := big.NewInt(0)
	io.WriteString(h, t.Token)

	signhash := h.Sum(nil)
	r, s, err := ecdsa.Sign(rand.Reader, signkey, signhash)
	if err != nil {
		utils.ELog(err, nil)
		return new(AccessToken)
	}

	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	sigtext := base64.StdEncoding.EncodeToString(signature)
	// BUG Need to sign ";ES256;" also. Security.
	t.Token = t.Token + ";ES256;" + sigtext
	return &AccessToken{}
}

// loadSigningKey reads the bytes of an PEM file to extract the ECDSA private key
func loadSigningKey(keydat []byte) (*ecdsa.PrivateKey, error) {
	var block *pem.Block
	block, _ = pem.Decode(keydat)
	if block == nil {
		return nil, errors.New("Cannot decode key")
	}

	var eckey interface{}
	var err error
	if eckey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		return nil, err
	}

	var pkey *ecdsa.PrivateKey
	var ok bool
	if pkey, ok = eckey.(*ecdsa.PrivateKey); !ok {
		return nil, errors.New("Cannot find ECDSA private key")
	}

	return pkey, nil
}
