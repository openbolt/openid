package openid

import (
	"errors"
	"net/http"
	"strings"

	"github.com/openbolt/openid/utils"
)

// 3.1.3.2.  Token Request Validation
func (op *OpenID) Token(w http.ResponseWriter, r *http.Request) (AuthSuccessResp, AuthErrResp) {
	if !op.serving {
		return AuthSuccessResp{}, AuthErrResp{}
	}

	clientID := GetParam(r, "client_id")
	if clientID == "" {
		err := AuthErrResp{
			Error:            "invalid_client",
			ErrorDescription: "No client_id given",
		}
		utils.EDebug(errors.New("returning invalid_client"), r)
		return AuthSuccessResp{}, err
	}

	// Authenticate the Client if it was issued Client Credentials or if it uses another Client Authentication method, per Section 9.
	if authok, autherr := op.AuthenticateClient(clientID, r); !authok {
		err := AuthErrResp{Error: "Undefined"}
		switch autherr {
		case CLIENT_NOT_ALLOWED:
			err = AuthErrResp{
				Error:            "invalid_client",
				ErrorDescription: "Client not allowed or cannot authenticate",
			}
		case REQUIRE_401:
			hdrs := new(http.Header)
			// TODO: Check for right value (rfc6749)
			hdrs.Add("WWW-Authenticate", "Basic realm=openid")
			err = AuthErrResp{
				Error:            "invalid_client",
				ErrorDescription: "Authentification needed",
				Headers:          *hdrs,
				StatusCode:       401,
			}
		}
		utils.EDebug(errors.New("returning invalid_client"), r)
		return AuthSuccessResp{}, err
	}

	// Ensure the Authorization Code was issued to the authenticated Client.
	// Verify that the Authorization Code is valid.
	// If possible, verify that the Authorization Code has not been previously used. => On exchange, the code will be retired
	session, err := op.Cache.GetSession(GetParam(r, "code"))
	if err != nil || session.ClientID != clientID {
		err := AuthErrResp{
			Error:            "invalid_grant",
			ErrorDescription: "Authorization Code is invalid",
		}
		utils.EDebug(errors.New("returning invalid_grant"), r)
		return AuthSuccessResp{}, err
	}

	// Ensure that the redirect_uri parameter value is identical to the redirect_uri parameter value that was included in the initial Authorization Request. If the redirect_uri parameter value is not present when there is only one registered redirect_uri value, the Authorization Server MAY return an error (since the Client should have included the parameter) or MAY proceed without an error (since OAuth 2.0 permits the parameter to be omitted in this case).
	if !op.Clientsrc.ValidateRedirectURI(clientID, GetParam(r, "redirect_uri")) {
		err := AuthErrResp{
			Error:            "invalid_grant",
			ErrorDescription: "Redirection URI is invalid",
		}
		utils.EDebug(errors.New("returning invalid_grant"), r)
		return AuthSuccessResp{}, err
	}

	// Verify that the Authorization Code used was issued in response to an OpenID Connect Authentication Request (so that an ID Token will be returned from the Token Endpoint).
	if !strings.Contains(session.Scope, "openid") {
		err := AuthErrResp{
			Error:            "invalid_request",
			ErrorDescription: "Code was not issued to an OIDC Auth Request",
		}
		utils.EDebug(errors.New("returning invalid_request"), r)
		return AuthSuccessResp{}, err
	}

	// Issue token according to variable `session`
	idToken, err := NewIDToken(session, op.accessTokenSignKey)
	atok := AccessToken{}
	atok.Load(session, op.accessTokenSignKey)
	if err == nil {
		op.Cache.Retire(GetParam(r, "code"))
		utils.EDebug(errors.New("returning ok"), r)
		return AuthSuccessResp{
			ok:          true,
			IDToken:     idToken,
			AccessToken: atok.Token,
			TokenType:   atok.TokenType,
			ExpiresIn:   atok.ExpiresIn,
		}, AuthErrResp{}
	} else {
		utils.EInfo(errors.New("Cannot generate IDToken: "+err.Error()), r)
		err := AuthErrResp{
			Error:            "invalid_request",
			ErrorDescription: "IDToken not avaiable",
		}
		utils.EDebug(errors.New("returning invalid_request"), r)
		return AuthSuccessResp{}, err
	}
}
