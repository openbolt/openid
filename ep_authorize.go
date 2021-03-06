package openid

import (
	"errors"
	"net/http"

	"github.com/openbolt/openid/utils"
)

// Authorize is called to process an authorization-request from an http api
// Ref 3.1.2.1. Authentication Request
// An Authentication Request is an OAuth 2.0 Authorization Request that requests
// that the End-User be authenticated by the Authorization Server.
func (op *OpenID) Authorize(w http.ResponseWriter, r *http.Request) (AuthSuccessResp, AuthErrResp) {
	if !op.serving {
		return AuthSuccessResp{}, AuthErrResp{}
	}

	// ref 3.1.2.2
	err1 := validateOAuthParams(r)             // ref Rule 1
	err2 := validateScopeParam(r)              // ref Rule 2
	err3 := validateReqParams(r, op.Clientsrc) // ref Rule 3

	// Check first part of validation
	if len(err1.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 1"), r)
		return AuthSuccessResp{}, err1
	}
	if len(err2.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 2"), r)
		return AuthSuccessResp{}, err2
	}
	if len(err3.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 3"), r)
		return AuthSuccessResp{}, err3
	}

	// Ref 3.1.2.3.  Authorization Server Authenticates End-User
	state := op.Enduser.Authpage(w, r)

	// Respond to enduser if not successfully authenticated
	if state.AuthAbort {
		utils.EDebug(errors.New("Auth aborted"), r)
		err := AuthErrResp{}
		err.Error = "login_required"
		err.ErrorDescription = "Authentication aborted"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	} else if state.AuthFailed {
		utils.EDebug(errors.New("Auth failed"), r)
		err := AuthErrResp{}
		err.Error = "access_denied"
		err.ErrorDescription = "Authentication failed"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	} else if state.AuthPrompting {
		utils.EDebug(errors.New("Auth prompting"), r)
		// Simply return if a prompt is presented
		return AuthSuccessResp{}, AuthErrResp{}
	} else if !state.AuthOk {
		utils.EDebug(errors.New("Auth requests reload"), r)
		// Reload page using the same method
		w.Header().Set("Location", r.RequestURI)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if state.AuthOk {
		utils.EDebug(errors.New("Authpage returned ok"), r)
	}

	// Can only be checked after authentification
	// (compare "sub" with requested `claims`->`sub`)
	err4 := validateSubParam(r, state.Sub) // ref Rule 4
	if len(err4.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 4"), r)
		return AuthSuccessResp{}, err4
	}

	// Run through flow
	// ref 3
	switch getFlow(GetParam(r, "response_type")) {
	case "authorization_code":
		utils.EDebug(errors.New("Using authzCodeFlow"), r)
		return op.authzCodeFlow(r, state)
	case "implicit":
		utils.EDebug(errors.New("Using implicit flow"), r)
		return op.implicitFlow(r, state)
	case "hybrid":
		utils.EDebug(errors.New("Using hybrid flow"), r)
		return op.hybridFlow(r, state)
	default:
		utils.EDebug(errors.New("invalid response_type, cannot find flow"), r)
		err := AuthErrResp{}
		err.Error = "invalid_request"
		err.ErrorDescription = "Invalid `code` request sent"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	}
}
