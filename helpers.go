package openid

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	uquery "github.com/google/go-querystring/query"
	"github.com/openbolt/openid/utils"
)

// GetParam extracts the OAuth parameters from an http.Request according to the
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
func serializeResponse(redirectURI url.URL, responseMode string, data interface{}) (url.URL, error) {
	var query url.Values
	if responseMode == "query" {
		query, _ = url.ParseQuery(redirectURI.RawQuery)
	} else {
		query, _ = url.ParseQuery(redirectURI.Fragment)
	}

	vals, err := uquery.Values(data)
	if err != nil {
		utils.ELog(errors.New("Cannot serialize to url"))
		return url.URL{}, err
	}

	// Merge the two maps
	for k := range vals {
		query[k] = append(query[k], vals[k]...)
	}

	if responseMode == "query" {
		redirectURI.RawQuery = query.Encode()
	} else {
		redirectURI.Fragment = query.Encode()
	}
	return redirectURI, nil
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

// GetRandomString returns an random string with size `size`
func GetRandomString(size int) (string, error) {
	// Generate `code` for response
	sec := make([]byte, size)
	_, err := rand.Read(sec)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sec), nil
}

func ReadClaimsRequest(data string) (ClaimsRequest, error) {
	type csreq struct {
		Userinfo map[string]csreq `json:"userinfo"`
		IDToken  map[string]csreq `json:"id_token"`
	}

	type creq struct {
		Essential bool     `json:"essential"`
		Value     string   `json:"value"`
		Values    []string `json:"values"`
	}
	raw := csreq{}
	err := json.Unmarshal([]byte(data), raw)
	if err != nil {
		utils.ELog(err)
		return ClaimsRequest{}, err
	}

	// BUG Not implemented yes
	// TODO: Convert csreq to ClaimsRequest
	return ClaimsRequest{}, errors.New("Not implemented yet")
}
