package openid

/*
 * Ref 3.1.  Authentication using the Authorization Code Flow
 * The Authorization Code Flow returns an Authorization Code to the Client,
 * which can then exchange it for an ID Token and an Access Token directly.
 *
 * Ref 3.1.1.  Authorization Code Flow Steps
 *
 * The Authorization Code Flow goes through the following steps.
 *
 *  1  Client prepares an Authentication Request containing the desired
 *      request parameters.
 *  2  Client sends the request to the Authorization Server.
 *  3  Authorization Server Authenticates the End-User.
 *  4  Authorization Server obtains End-User Consent/Authorization.
 *  5  Authorization Server sends the End-User back to the Client with an
 *      Authorization Code.
 *  6  Client requests a response using the Authorization Code at the
 *       Token Endpoint.
 *  7  Client receives a response that contains an ID Token and Access Token
 *       in the response body.
 *  8  Client validates the ID token and retrieves the End-User's
 *       Subject Identifier.
 */

func (op *OpenID) authz_code_flow(parms Values) (AuthSuccessResp, AuthErrResp) {
	// BUG(djboris) implement
	return AuthSuccessResp{}, AuthErrResp{}
}
