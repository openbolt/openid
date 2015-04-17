package openid

import "net/http"

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
		return AuthSuccessResp{}, err1
	}
	if len(err2.Error) != 0 {
		return AuthSuccessResp{}, err2
	}
	if len(err3.Error) != 0 {
		return AuthSuccessResp{}, err3
	}

	// Ref 3.1.2.3.  Authorization Server Authenticates End-User
	state := op.Enduser.Authpage(w, r)

	// Respond to enduser if not successfully authenticated
	if state.AuthAbort {
		err := AuthErrResp{}
		err.Error = "login_required"
		err.ErrorDescription = "Authentication aborted"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	} else if state.AuthFailed {
		err := AuthErrResp{}
		err.Error = "access_denied"
		err.ErrorDescription = "Authentication failed"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	} else if state.AuthPrompting {
		// Simply return if a prompt is presented
		return AuthSuccessResp{}, AuthErrResp{}
	} else if !state.AuthOk {
		// Reload page using the same method
		w.Header().Set("Location", r.RequestURI)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if state.AuthOk {
		// Do really nothing?
	}

	// BUG(djboris) Check additional Request Values
	err4 := validate_sub_param(r, state.Sub) // ref Rule 4
	if len(err4.Error) != 0 {
		return AuthSuccessResp{}, err4
	}

	// var idTokenVals Values
	//idTokenVals.Set("nonce", GetParam(r, "nonce")) //TODO: Ignore if empty
	// TODO: display, prompt, max_age, ui_locales, acr_values

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
		err := AuthErrResp{}
		err.Error = "invalid_request"
		err.ErrorDescription = "Invalid `code` request sent"
		err.State = GetParam(r, "state")
		return AuthSuccessResp{}, err
	}
}
