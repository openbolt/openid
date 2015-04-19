package openid

import (
	"errors"
	"net/http"

	"github.com/openbolt/openid/utils"
)

// Ref 3.1.2.1. Authentication Request
// An Authentication Request is an OAuth 2.0 Authorization Request that requests
// that the End-User be authenticated by the Authorization Server.
func (op *OpenID) Authorize(w http.ResponseWriter, r *http.Request) (AuthSuccessResp, AuthErrResp) {
	if !op.serving {
		return AuthSuccessResp{}, AuthErrResp{}
	}

	// ref 3.1.2.2
	err1 := validate_oauth_params(r)             // ref Rule 1
	err2 := validate_scope_param(r)              // ref Rule 2
	err3 := validate_req_params(r, op.Clientsrc) // ref Rule 3

	// Check first part of validation
	if len(err1.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 1"))
		return AuthSuccessResp{}, err1
	}
	if len(err2.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 2"))
		return AuthSuccessResp{}, err2
	}
	if len(err3.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 3"))
		return AuthSuccessResp{}, err3
	}

	// Ref 3.1.2.3.  Authorization Server Authenticates End-User
	state := op.Enduser.Authpage(w, r)

	// Respond to enduser if not successfully authenticated
	if state.AuthAbort {
		utils.EDebug(errors.New("Auth aborted"))
		err := AuthErrResp{}
		err.Error = "login_required"
		err.ErrorDescription = "Authentication aborted"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	} else if state.AuthFailed {
		utils.EDebug(errors.New("Auth failed"))
		err := AuthErrResp{}
		err.Error = "access_denied"
		err.ErrorDescription = "Authentication failed"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	} else if state.AuthPrompting {
		utils.EDebug(errors.New("Auth prompting"))
		// Simply return if a prompt is presented
		return AuthSuccessResp{}, AuthErrResp{}
	} else if !state.AuthOk {
		utils.EDebug(errors.New("Auth requests reload"))
		// Reload page using the same method
		w.Header().Set("Location", r.RequestURI)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if state.AuthOk {
		utils.EDebug(errors.New("Auth ok"))
		// Do really nothing?
	}

	// Can only be checked after authentification
	// (compare "sub" with requested `claims`->`sub`)
	err4 := validate_sub_param(r, state.Sub) // ref Rule 4
	if len(err4.Error) != 0 {
		utils.EDebug(errors.New("Failed Rule 4"))
		return AuthSuccessResp{}, err4
	}

	// Run through flow
	// ref 3
	switch getFlow(GetParam(r, "response_type")) {
	case "authorization_code":
		return op.authz_code_flow(r, state)
	case "implicit":
		return op.implicit_flow(r, state)
	case "hybrid":
		return op.hybrid_flow(r, state)
	default:
		utils.EDebug(errors.New("invalid response_type"))
		err := AuthErrResp{}
		err.Error = "invalid_request"
		err.ErrorDescription = "Invalid `code` request sent"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	}
}
