package openid

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"github.com/openbolt/openid/utils"
)

const (
	// AuthzCodeOctetsRand has the number of random bytes used in authz_code `code` generation
	AuthzCodeOctetsRand = 32
)

// Ref 3.1.  Authentication using the Authorization Code Flow
// The Authorization Code Flow returns an Authorization Code to the Client,
// which can then exchange it for an ID Token and an Access Token directly.
func (op *OpenID) authzCodeFlow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// Generate `code` for response
	sec := make([]byte, AUTHZ_CODE_OCTETS_RAND)
	_, err := rand.Read(sec)
	if err != nil {
		utils.ELog(err)

		resp := AuthErrResp{}
		resp.Error = "server_error"
		resp.ErrorDescription = "Server isn't able to fullfill your request"
		resp.State = GetParam(r, "state")
		return AuthSuccessResp{}, resp
	}

	// Generate response value
	suc := AuthSuccessResp{ok: true}
	suc.State = GetParam(r, "state")
	suc.Code = base64.StdEncoding.EncodeToString(sec)

	// Cache request to be able to respond with token
	ses := AuthzCodeSession{}
	ses.Code = suc.Code
	ses.ClientID = GetParam(r, "client_id")
	ses.Nonce = GetParam(r, "state")
	ses.AuthTime = time.Now()
	ses.MaxAge, _ = time.ParseDuration(GetParam(r, "max_age"))
	ses.Acr = state.Acr
	ses.ClaimsLocales = GetParam(r, "claim_locales")
	ses.Claims = GetParam(r, "claims") // TODO: Unmarshal
	// BUG Now it's needed to cache this (ses)

	return suc, AuthErrResp{}
}

func (op *OpenID) implicitFlow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// BUG(djboris) implement
	return AuthSuccessResp{}, AuthErrResp{}
}

func (op *OpenID) hybridFlow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// BUG(djboris) implement
	return AuthSuccessResp{}, AuthErrResp{}
}
