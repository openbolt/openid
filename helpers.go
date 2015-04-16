package openid

import (
	"net/http"
)

// GetParam extracts the OAuth parameters from an http.Request according to the
// OIDC spec
func GetParam(r *http.Request, param string) string {
	if r.Method == "GET" {
		// MUST URI Query String Serialization
		return r.URL.Query().Get(param)

	} else if r.Method == "POST" {
		// MUST Form Serialization
		r.ParseForm()
		return r.PostForm.Get(param)
	} else {
		return ""
	}
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
