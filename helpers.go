package openid

import (
	"errors"
	"net/http"
	"net/url"

	uquery "github.com/google/go-querystring/query"
	"github.com/openbolt/openid/utils"
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

// Serialize response serializes an struct to an url query or fragment
func serializeResponse(redirect_uri url.URL, response_mode string, data interface{}) (url.URL, error) {
	var query url.Values
	if response_mode == "query" {
		query, _ = url.ParseQuery(redirect_uri.RawQuery)
	} else {
		query, _ = url.ParseQuery(redirect_uri.Fragment)
	}

	vals, err := uquery.Values(data)
	if err != nil {
		utils.ELog(errors.New("Cannot serialize to url"))
		return url.URL{}, err
	}

	// Merge the two maps
	for k, _ := range vals {
		query[k] = append(query[k], vals[k]...)
	}

	if response_mode == "query" {
		redirect_uri.RawQuery = query.Encode()
	} else {
		redirect_uri.Fragment = query.Encode()
	}
	return redirect_uri, nil
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
