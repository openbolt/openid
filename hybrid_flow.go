package openid

import "net/http"

func (op *OpenID) hybrid_flow(r *http.Request) (AuthSuccessResp, AuthErrResp) {
	// BUG(djboris) implement
	return AuthSuccessResp{}, AuthErrResp{}
}
