package openid

import (
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
	code, err := GetRandomString(AuthzCodeOctetsRand)
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
	suc.Code = code

	// Cache request to be able to respond with token
	ses := Session{}
	ses.Code = suc.Code
	ses.ClientID = GetParam(r, "client_id")
	ses.Nonce = GetParam(r, "nonce")
	ses.Scope = GetParam(r, "scope")
	ses.AuthTime = state.AuthTime
	ses.MaxAge, _ = time.ParseDuration(GetParam(r, "max_age"))
	ses.Acr = state.Acr
	ses.ClaimsLocales = GetParam(r, "claim_locales")
	ses.Claims = GetParam(r, "claims") // TODO: Unmarshal

	// Now it's needed to cache this
	op.Cache.Cache(ses)

	return suc, AuthErrResp{}
}

func (op *OpenID) implicitFlow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// Generate an session, no need to save/cache
	ses := Session{}
	ses.ClientID = GetParam(r, "client_id")
	ses.Nonce = GetParam(r, "nonce")
	ses.Scope = GetParam(r, "scope")
	ses.AuthTime = state.AuthTime
	ses.MaxAge, _ = time.ParseDuration(GetParam(r, "max_age"))
	ses.Acr = state.Acr
	ses.ClaimsLocales = GetParam(r, "claim_locales")
	ses.Claims = GetParam(r, "claims") // TODO: Unmarshal

	suc := AuthSuccessResp{ok: true}
	suc.State = GetParam(r, "state")
	suc.IDToken = NewIDToken(ses)

	if GetParam(r, "response_type") != "id_token" {
		tok := NewAccessToken(ses)
		suc.AccessToken = tok.Token
		suc.TokenType = tok.TokenType
		suc.ExpiresIn = tok.ExpiresIn
	}

	return suc, AuthErrResp{}
}

func (op *OpenID) hybridFlow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// Generate `code` for response
	code, err := GetRandomString(AuthzCodeOctetsRand)
	if err != nil {
		utils.ELog(err)

		resp := AuthErrResp{}
		resp.Error = "server_error"
		resp.ErrorDescription = "Server isn't able to fullfill your request"
		resp.State = GetParam(r, "state")
		return AuthSuccessResp{}, resp
	}

	// Cache request to be able to respond with token
	ses := Session{}
	ses.Code = code
	ses.ClientID = GetParam(r, "client_id")
	ses.Nonce = GetParam(r, "nonce")
	ses.Scope = GetParam(r, "scope")
	ses.AuthTime = state.AuthTime
	ses.MaxAge, _ = time.ParseDuration(GetParam(r, "max_age"))
	ses.Acr = state.Acr
	ses.ClaimsLocales = GetParam(r, "claim_locales")
	ses.Claims = GetParam(r, "claims") // TODO: Unmarshal

	// Generate response value
	suc := AuthSuccessResp{ok: true}
	suc.State = GetParam(r, "state")
	suc.Code = code
	suc.IDToken = NewIDToken(ses)

	if GetParam(r, "response_type") != "id_token" {
		tok := NewAccessToken(ses)
		suc.AccessToken = tok.Token
		suc.TokenType = tok.TokenType
		suc.ExpiresIn = tok.ExpiresIn
	}

	// Now it's needed to cache this
	op.Cache.Cache(ses)

	return suc, AuthErrResp{}
}
