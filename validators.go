package openid

import (
	"errors"
	"net/url"
	"regexp"
	"strings"

	"github.com/openbolt/openid/utils"
)

// A-BNF common definitions. Ref OAuth 2.0 rfc6749#appendix-A
const (
	VSCHAR            = "[\\x20-\\x7E]"
	NQCHAR            = "[\\x21\\x23-\\x5B\\x5D-\\x7E]"
	NQSCHAR           = "[\\x20-\\x21\\x23-\\x5B\\x5D-\\x7E]"
	UNICODECHARNOCRLF = "[\\x09\\x20-\\x7E\\x80-\\x{D7FF}\\x{E000}-\\x{FFFD}\\x{10000}-\\x{10FFFF}]"

	ALPHA = "\\x41-\\x5A\\x61-\\x7A"
	DIGIT = "\\x30-\\x39"
	SP    = "\\x20"

	// response-name = 1*response-char
	// response-char = "_" / DIGIT / ALPHA
	rx_response_name = "[" + DIGIT + ALPHA + "_" + "]+"

	// grant-name = 1*name-char
	// name-char  = "-" / "." / "_" / DIGIT / ALPHA
	rx_grant_name = "[\\-\\._" + DIGIT + ALPHA + "]+"

	// type-name  = 1*name-char
	// name-char  = "-" / "." / "_" / DIGIT / ALPHA
	rx_type_name = "[\\-\\._" + DIGIT + ALPHA + "]+"

	// BUG(djboris) Implement it according to rfc3986#appendix-A
	// TODO: Do it right with BNF
	rx_uri_reference = "^https?:"
)

var (
	// client-id     = *VSCHAR
	re_client_id = regexp.MustCompile(VSCHAR + "*")

	// client-secret = *VSCHAR
	re_client_secret = regexp.MustCompile(VSCHAR + "*")

	// response-type = response-name *( SP response-name )
	re_response_type = regexp.MustCompile(
		rx_response_name + "(\\s" + rx_response_name + ")*")

	// scope       = scope-token *( SP scope-token )
	// scope-token = 1*NQCHAR
	re_scope = regexp.MustCompile(NQCHAR + "+(\\s" + NQCHAR + "+)*")

	// state      = 1*VSCHAR
	re_state = regexp.MustCompile(VSCHAR + "+")

	// redirect-uri      = URI-reference
	re_redirect_uri = regexp.MustCompile(rx_uri_reference)

	// error             = 1*NQSCHAR
	re_error = regexp.MustCompile(NQSCHAR + "+")

	// error-description = 1*NQSCHAR
	re_error_description = regexp.MustCompile(NQSCHAR + "+")

	// error-uri         = URI-reference
	re_error_uri = regexp.MustCompile(rx_uri_reference)

	// grant-type = grant-name / URI-reference
	re_grant_type = regexp.MustCompile("(" + rx_grant_name + "|" + rx_uri_reference + ")")

	// code       = 1*VSCHAR
	re_code = regexp.MustCompile(VSCHAR + "+")

	// access-token = 1*VSCHAR
	re_access_token = regexp.MustCompile(VSCHAR + "+")

	// token-type = type-name / URI-reference
	re_token_type = regexp.MustCompile("(" + rx_type_name + "|" + rx_uri_reference + ")")

	// expires-in = 1*DIGIT
	re_expires_in = regexp.MustCompile(DIGIT + "+")

	// username = *UNICODECHARNOCRLF
	re_username = regexp.MustCompile(UNICODECHARNOCRLF + "*")

	// password = *UNICODECHARNOCRLF
	re_password = regexp.MustCompile(UNICODECHARNOCRLF + "*")

	// refresh-token = 1*VSCHAR
	re_refresh_token = regexp.MustCompile(VSCHAR + "+")

	// The syntax for new endpoint parameters is defined in
	// OAuth 2.0 Section 8.2
	// param-name = 1*name-char
	// name-char  = "-" / "." / "_" / DIGIT / ALPHA
	re_param_name = regexp.MustCompile("[\\-\\._" + DIGIT + ALPHA + "]+")
)

/*
 * 3.1.2.2.  Authentication Request Validation
 */

// Rule 1:
// The Authorization Server MUST validate all the OAuth 2.0 parameters according
// to the OAuth 2.0 specification.
// Ref rfc6749 appendix-A
func validate_oauth_params(params Values) AuthErrResp {
	test := func(parm string, re *regexp.Regexp, cont *string) bool {
		if params.Get(parm) != "" && !re.MatchString(params.Get(parm)) {
			utils.EDebug(errors.New(parm + " malformed"))
			*cont = *cont + parm + ";"
			return false
		} else {
			return true
		}
	}

	var errParm *string = new(string)
	var ok bool = true
	ok = ok && test("client_id", re_client_id, errParm)
	ok = ok && test("client_secret", re_client_secret, errParm)
	ok = ok && test("response_type", re_response_type, errParm)
	ok = ok && test("scope", re_scope, errParm)
	ok = ok && test("state", re_state, errParm)
	ok = ok && test("redirect_uri", re_redirect_uri, errParm)
	ok = ok && test("error", re_error, errParm)
	ok = ok && test("error_description", re_error_description, errParm)
	ok = ok && test("error_uri", re_error_uri, errParm)
	ok = ok && test("grant_type", re_grant_type, errParm)
	ok = ok && test("code", re_code, errParm)
	ok = ok && test("access_token", re_access_token, errParm)
	ok = ok && test("token_type", re_token_type, errParm)
	ok = ok && test("expires_in", re_expires_in, errParm)
	ok = ok && test("username", re_username, errParm)
	ok = ok && test("password", re_password, errParm)
	ok = ok && test("refresh_token", re_refresh_token, errParm)

	resp := AuthErrResp{}
	if ok {
		utils.EDebug(errors.New("returning ok"))
		return resp
	} else {
		resp.Error = "invalid_request"
		resp.ErrorDescription = "One or more malformed request parameters: " + *errParm
		resp.State = params.Get("state")

		utils.EDebug(errors.New("returning invalid_request"))
		return resp
	}
}

// Rule 2:
// Verify that a scope parameter is present and contains the openid scope value.
func validate_scope_param(params Values) AuthErrResp {
	args := strings.Split(params.Get("scope"), " ")

	// Check if openid value in scope
	var ok bool = false
	for _, v := range args {
		if v == "openid" {
			ok = true
		}
	}

	resp := AuthErrResp{}
	if ok {
		utils.EDebug(errors.New("returning ok"))
		return resp
	} else {
		resp.Error = "invalid_request"
		resp.ErrorDescription = "Scope doesn't contain openid"
		resp.State = params.Get("state")

		utils.EDebug(errors.New("returning invalid_request"))
		return resp
	}
}

// Rule 3:
// The Authorization Server MUST verify that all the REQUIRED parameters are
// present and their usage conforms to this specification.
//   Scope  will not be tested as it is already done in validate_scope_param
//   Required params: scope, response_type, client_id, redirect_uri
//   For Implicit: + nonce
func validate_req_params(params Values, clt Clientsource) AuthErrResp {
	var ok bool = true
	var errs string

	// Check flow
	flow := getFlow(params.Get("response_type"))
	if flow == "" {
		s := "invalid response_type"
		utils.EDebug(errors.New(s))
		errs += s + ";"
		ok = false
	}

	// Check client_id
	t := clt.IsClient(params.Get("client_id"))
	if !t {
		errs += "no client with this id;"
	}
	ok = ok && t

	// Check redirect_uri
	client_id := params.Get("client_id")
	if flow == "implicit" {
		// ...the Redirection URI MUST NOT use the http scheme unless
		// the Client is a native application, in which case it
		// MAY use the http: scheme with localhost as the hostname.
		uri, err := url.Parse(params.Get("redirect_uri"))
		if err != nil {
			utils.EInfo(err)
			errs += "failed to parse redirect_uri: " + err.Error() + ";"
			ok = false
		}
		if uri.Scheme == "http" &&
			(clt.GetApplType(client_id) != "native" || uri.Host != "localhost") {
			s := "Not compatible redirect_uri"
			utils.EDebug(errors.New(s))
			errs += s + ";"
			ok = false
		}

	} else {
		t := clt.ValidateRedirectUri(client_id, params.Get("redirect_uri"))
		if !t {
			errs += "redirect_uri validation failed;"
		}
		ok = ok && t
	}

	// For implicit: nonce is required
	if flow == "implicit" && len(params.Get("nonce")) == 0 {
		s := "nonce not present in implicit flow"
		utils.EDebug(errors.New(s))
		errs += s + ";"
		ok = false
	}

	// Return
	resp := AuthErrResp{}
	if ok {
		utils.EDebug(errors.New("returning ok"))
		return resp
	} else {
		resp.Error = "invalid_request"
		resp.ErrorDescription = "One or more not valid parameters: " + errs
		resp.State = params.Get("state")

		utils.EDebug(errors.New("returning invalid_request"))
		return resp
	}
}

// Rule 4
// If the sub (subject) Claim is requested with a specific value for the
// ID Token, the Authorization Server MUST only send a positive response if
// the End-User identified by that sub value has an active session with the
// Authorization Server or has been Authenticated as a result of the request.
// The Authorization Server MUST NOT reply with an ID Token or Access Token
// for a different user, even if they have an active session with the
// Authorization Server. Such a request can be made either using an
// id_token_hint parameter or by requesting a specific Claim Value as described
// in Section 5.5.1, if the claims parameter is supported by the implementation.
func validate_sub_param(params Values) AuthErrResp {
	// BUG(djboris) Implement
	// This should be called after auth
	return AuthErrResp{}
}

// getFlow returns authorization_code, implicit or hybrid. If any error occours,
// "" will be returned
func getFlow(field string) string {
	// Returns according flow

	// Ref 3
	switch field {
	case "code":
		return "authorization_code"
	case "id_token", "id_token token":
		return "implicit"
	case "code id_token", "code token", "code id_token token":
		return "hybrid"
	default:
		return ""
	}
}
