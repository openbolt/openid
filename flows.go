package openid

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/openbolt/openid/utils"
)

const AUTHZ_CODE_OCTETS_RAND = 32

// Ref 3.1.  Authentication using the Authorization Code Flow
// The Authorization Code Flow returns an Authorization Code to the Client,
// which can then exchange it for an ID Token and an Access Token directly.
func (op *OpenID) authz_code_flow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
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

	suc := AuthSuccessResp{ok: true}
	suc.State = GetParam(r, "state")
	suc.Code = base64.StdEncoding.EncodeToString(sec)
	return suc, AuthErrResp{}
}

func (op *OpenID) implicit_flow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// BUG(djboris) implement
	return AuthSuccessResp{}, AuthErrResp{}
}

func (op *OpenID) hybrid_flow(r *http.Request, state AuthState) (AuthSuccessResp, AuthErrResp) {
	// BUG(djboris) implement
	return AuthSuccessResp{}, AuthErrResp{}
}
