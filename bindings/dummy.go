package bindings

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/openbolt/openid"
)

//go:generate go-bindata -o bindata.go -pkg bindings assets/

// DummySource defines an empty Claimsource and Clientsource
// which stores its data in RAM.
// This can be used for testcases or demo applications
type DummySource struct {
}

/*
 * Claimsource
 */
// Get
func (ds DummySource) Get(id, claim, def string) (string, bool) {
	// BUG(djboris) Implement
	return "foo", true
}

/*
 * Clientsource
 */
// IsClient returns true if id begins with `clt`
func (ds DummySource) IsClient(id string) bool {
	return strings.HasPrefix(id, "clt")
}
func (ds DummySource) GetClientType(id string) string {
	// BUG(djboris) Implement
	return "confidential"
}
func (ds DummySource) GetApplType(id string) string {
	// BUG(djboris) Implement
	return "web"
}

// Returns true if host equivalent to id[3::]. Example: cltlocalhost => localhost
func (ds DummySource) ValidateRedirectURI(id, uri string) bool {
	// BUG(djboris) Implement
	u, err := url.Parse(uri)
	if err != nil || len(id) < 3 {
		return false
	}
	return u.Host == id[3:len(id)]
}

/*
 * EnduserIf
 */
func (ds *DummySource) Authpage(w http.ResponseWriter, r *http.Request) openid.AuthState {
	var warn string

	if openid.GetParam(r, "_login") != "" {
		sub, iss := dummyAuth(openid.GetParam(r, "_username"), openid.GetParam(r, "_password"))
		if sub == "" {
			warn = "Wrong credentials"
			// Continues to loginpage...
		} else {
			res := openid.AuthState{}
			res.AuthOk = true
			res.Sub = sub
			res.Iss = iss
			res.AuthTime = time.Now()
			res.Acr = "0"
			res.Amr = "none"
			return res
		}

	} else if openid.GetParam(r, "_fail") != "" {
		return openid.AuthState{AuthFailed: true}
	} else if openid.GetParam(r, "_abort") != "" {
		return openid.AuthState{AuthAbort: true}
	}

	//
	// Display login form
	//
	tpl, _ := Asset("assets/pwlogin.html")
	t := template.New("pwauth.html")
	t.Parse(string(tpl))

	// Prepare form/URI values
	var query url.Values
	if r.Method == "GET" {
		query = r.URL.Query()
	} else {
		r.ParseForm()
		query = r.PostForm
	}

	// Prepare struct for template renderer
	vals := make(map[string]string)
	for k, v := range query {
		if k[0] != '_' {
			vals[k] = v[0]
		}
	}
	x := struct {
		Method  string
		Baseurl string
		Values  map[string]string
		Warn    string
	}{
		"post", // Use post instead of r.Method (c&p security)
		r.URL.Path,
		vals,
		warn,
	}

	// Execute
	t.Execute(w, x)
	return openid.AuthState{AuthPrompting: true}
}

// Username: sub dot iss (Example: djboris.myIssuer)
// PW: must be not empty
func dummyAuth(user, pw string) (sub, iss string) {
	if user == "" || pw == "" {
		return "", ""
	}

	v := strings.Split(user, ".")
	if len(v) < 2 {
		return "", ""
	} else {
		return v[0], v[1]
	}
}
