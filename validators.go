package openid

import (
	"errors"
	"net/http"
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
	rxResponseName = "[" + DIGIT + ALPHA + "_" + "]+"

	// grant-name = 1*name-char
	// name-char  = "-" / "." / "_" / DIGIT / ALPHA
	rxGrantName = "[\\-\\._" + DIGIT + ALPHA + "]+"

	// type-name  = 1*name-char
	// name-char  = "-" / "." / "_" / DIGIT / ALPHA
	rxTypeName = "[\\-\\._" + DIGIT + ALPHA + "]+"

	// Regex from rfc3986#appendix-B
	rxURIReference = "^(([^:/?#]+):)?(//([^/?#]*))?([^?#]*)(\\?([^#]*))?(#(.*))?"
)

var (
	// client-id     = *VSCHAR
	reClientID = regexp.MustCompile(VSCHAR + "*")

	// client-secret = *VSCHAR
	reClientSecret = regexp.MustCompile(VSCHAR + "*")

	// response-type = response-name *( SP response-name )
	reResponseType = regexp.MustCompile(
		rxResponseName + "(\\s" + rxResponseName + ")*")

	// scope       = scope-token *( SP scope-token )
	// scope-token = 1*NQCHAR
	reScope = regexp.MustCompile(NQCHAR + "+(\\s" + NQCHAR + "+)*")

	// state      = 1*VSCHAR
	reState = regexp.MustCompile(VSCHAR + "+")

	// redirect-uri      = URI-reference
	reRedirectURI = regexp.MustCompile(rxURIReference)

	// error             = 1*NQSCHAR
	reError = regexp.MustCompile(NQSCHAR + "+")

	// error-description = 1*NQSCHAR
	reErrorDescription = regexp.MustCompile(NQSCHAR + "+")

	// error-uri         = URI-reference
	reErrorURI = regexp.MustCompile(rxURIReference)

	// grant-type = grant-name / URI-reference
	reGrantType = regexp.MustCompile("(" + rxGrantName + "|" + rxURIReference + ")")

	// code       = 1*VSCHAR
	reCode = regexp.MustCompile(VSCHAR + "+")

	// access-token = 1*VSCHAR
	reAccessToken = regexp.MustCompile(VSCHAR + "+")

	// token-type = type-name / URI-reference
	reTokenType = regexp.MustCompile("(" + rxTypeName + "|" + rxURIReference + ")")

	// expires-in = 1*DIGIT
	reExpiresIn = regexp.MustCompile(DIGIT + "+")

	// username = *UNICODECHARNOCRLF
	reUsername = regexp.MustCompile(UNICODECHARNOCRLF + "*")

	// password = *UNICODECHARNOCRLF
	rePassword = regexp.MustCompile(UNICODECHARNOCRLF + "*")

	// refresh-token = 1*VSCHAR
	reRefreshToken = regexp.MustCompile(VSCHAR + "+")

	// The syntax for new endpoint parameters is defined in
	// OAuth 2.0 Section 8.2
	// param-name = 1*name-char
	// name-char  = "-" / "." / "_" / DIGIT / ALPHA
	reParamName = regexp.MustCompile("[\\-\\._" + DIGIT + ALPHA + "]+")
)

/*
 * 3.1.2.2.  Authentication Request Validation
 */

// Rule 1:
// The Authorization Server MUST validate all the OAuth 2.0 parameters according
// to the OAuth 2.0 specification.
// Ref rfc6749 appendix-A
func validateOAuthParams(r *http.Request) AuthErrResp {
	test := func(parm string, re *regexp.Regexp, cont *string) bool {
		if GetParam(r, parm) != "" && !re.MatchString(GetParam(r, parm)) {
			utils.EDebug(errors.New(parm+" malformed"), r)
			*cont = *cont + parm + ";"
			return false
		}
		return true
	}

	var errParm = new(string)
	var ok = true
	ok = ok && test("client_id", reClientID, errParm)
	ok = ok && test("client_secret", reClientSecret, errParm)
	ok = ok && test("response_type", reResponseType, errParm)
	ok = ok && test("scope", reScope, errParm)
	ok = ok && test("state", reState, errParm)
	ok = ok && test("redirect_uri", reRedirectURI, errParm)
	ok = ok && test("error", reError, errParm)
	ok = ok && test("error_description", reErrorDescription, errParm)
	ok = ok && test("error_uri", reErrorURI, errParm)
	ok = ok && test("grant_type", reGrantType, errParm)
	ok = ok && test("code", reCode, errParm)
	ok = ok && test("access_token", reAccessToken, errParm)
	ok = ok && test("token_type", reTokenType, errParm)
	ok = ok && test("expires_in", reExpiresIn, errParm)
	ok = ok && test("username", reUsername, errParm)
	ok = ok && test("password", rePassword, errParm)
	ok = ok && test("refresh_token", reRefreshToken, errParm)

	resp := AuthErrResp{}
	if ok {
		utils.EDebug(errors.New("returning ok"), r)
	} else {
		resp.Error = "invalid_request"
		resp.ErrorDescription = "One or more malformed request parameters: " + *errParm
		resp.State = GetParam(r, "state")

		utils.EDebug(errors.New("returning invalid_request"), r)
	}
	return resp
}

// Rule 2:
// Verify that a scope parameter is present and contains the openid scope value.
func validateScopeParam(r *http.Request) AuthErrResp {
	args := strings.Split(GetParam(r, "scope"), " ")

	// Check if openid value in scope
	var ok = false
	for _, v := range args {
		if v == "openid" {
			ok = true
		}
	}

	resp := AuthErrResp{}
	if ok {
		utils.EDebug(errors.New("returning ok"), r)
	} else {
		resp.Error = "invalid_request"
		resp.ErrorDescription = "Scope doesn't contain openid"
		resp.State = GetParam(r, "state")

		utils.EDebug(errors.New("returning invalid_request"), r)
	}
	return resp
}

// Rule 3:
// The Authorization Server MUST verify that all the REQUIRED parameters are
// present and their usage conforms to this specification.
//   Scope  will not be tested as it is already done in validate_scope_param
//   Required params: scope, response_type, client_id, redirect_uri
//   For Implicit: + nonce
func validateReqParams(r *http.Request, clt Clientsource) AuthErrResp {
	var ok = true
	var errs string

	// Check existence of parameters
	if GetParam(r, "scope") == "" {
		errs += "scope missing"
		ok = false
	}
	if GetParam(r, "response_type") == "" {
		errs += "response_type missing"
		ok = false
	}
	if GetParam(r, "client_id") == "" {
		errs += "client_id missing"
		ok = false
	}
	if GetParam(r, "redirect_uri") == "" {
		errs += "redirect_uri missing"
		ok = false
	}

	// response_type
	// OAuth 2.0 Response Type value that determines the authorization
	// processing flow to be used, including what parameters are returned
	// from the endpoints used.
	flow := getFlow(GetParam(r, "response_type"))
	if flow == "" {
		s := "invalid response_type"
		utils.EDebug(errors.New(s), r)
		errs += s + ";"
		ok = false
	}

	// client_id
	// OAuth 2.0 Client Identifier valid at the Authorization Server.
	t := clt.IsClient(GetParam(r, "client_id"))
	if !t {
		errs += "no client with this id;"
	}
	ok = ok && t

	// redirect_uri
	// This URI MUST exactly match one of the Redirection URI values for the
	// Client pre-registered at the OpenID Provider, with the matching
	// performed as described in Section 6.2.1 of [RFC3986] (Simple String Comparison).
	// When using this flow, the Redirection URI SHOULD use the https scheme;
	// however, it MAY use the http scheme, provided that the
	// Client Type is confidential, as defined in Section 2.1 of OAuth 2.0,
	// and provided the OP allows the use of http Redirection URIs in this case.
	// The Redirection URI MAY use an alternate scheme, such as one that is
	// intended to identify a callback into a native application.
	clientID := GetParam(r, "client_id")
	redirectURI := GetParam(r, "redirect_uri")
	t = checkRedirectURI(redirectURI, clientID, flow, clt)
	if !t {
		errs += "invalid or not allowed redirect_uri;"
	}
	ok = ok && t

	// nonce (required only for implicit flow)
	// String value used to associate a Client session with an ID Token,
	// and to mitigate replay attacks. The value is passed through unmodified
	// from the Authentication Request to the ID Token. Sufficient entropy
	// MUST be present in the nonce values used to prevent attackers from guessing values
	if flow == "implicit" && len(GetParam(r, "nonce")) == 0 {
		s := "nonce not present in implicit flow"
		utils.EDebug(errors.New(s), r)
		errs += s + ";"
		ok = false
	}

	// Return
	resp := AuthErrResp{}
	if ok {
		utils.EDebug(errors.New("returning ok"), r)
	} else {
		resp.Error = "invalid_request"
		resp.ErrorDescription = "One or more not valid parameters: " + errs
		resp.State = GetParam(r, "state")

		utils.EDebug(errors.New("returning invalid_request"), r)
	}
	return resp
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
func validateSubParam(r *http.Request, sub string) AuthErrResp {
	// BUG(djboris) Implement
	utils.EInfo(errors.New("returning ok, as Not Implemented"), r)
	return AuthErrResp{}
}

// checkRedirectURI validates an redirect_uri according to flow type
func checkRedirectURI(redirectURI, clientID, flow string, clt Clientsource) bool {
	if flow == "implicit" {
		// ...the Redirection URI MUST NOT use the http scheme unless
		// the Client is a native application, in which case it
		// MAY use the http: scheme with localhost as the hostname.
		uri, err := url.Parse(redirectURI)
		if err != nil {
			utils.EInfo(err, nil)
			return false
		}
		if uri.Scheme == "http" &&
			(clt.GetApplType(clientID) != "native" || uri.Host != "localhost") {
			utils.EDebug(errors.New("Not compatible redirect_uri"), nil)
			return false
		}

	} else {
		if !clt.ValidateRedirectURI(clientID, redirectURI) {
			utils.EDebug(errors.New("Client hasn't registered this redirect_uri"), nil)
			return false
		}
	}
	return true
}
