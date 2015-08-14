package openid

import (
	"errors"
	"net/http"

	"github.com/openbolt/openid/utils"
)

func (op *OpenID) Token(w http.ResponseWriter, r *http.Request) (AuthSuccessResp, AuthErrResp) {
	if !op.serving {
		return AuthSuccessResp{}, AuthErrResp{}
	}
	// BUG Authenticate the Client if it was issued Client Credentials or if it uses another Client Authentication method, per Section 9.
	// BUG Ensure the Authorization Code was issued to the authenticated Client.
	// BUG Verify that the Authorization Code is valid.
	// BUG If possible, verify that the Authorization Code has not been previously used.
	// BUG Ensure that the redirect_uri parameter value is identical to the redirect_uri parameter value that was included in the initial Authorization Request. If the redirect_uri parameter value is not present when there is only one registered redirect_uri value, the Authorization Server MAY return an error (since the Client should have included the parameter) or MAY proceed without an error (since OAuth 2.0 permits the parameter to be omitted in this case).
	// BUG Verify that the Authorization Code used was issued in response to an OpenID Connect Authentication Request (so that an ID Token will be returned from the Token Endpoint).
	utils.ELog(errors.New("token_endpoint not implemented"))
	return AuthSuccessResp{}, AuthErrResp{}
}
