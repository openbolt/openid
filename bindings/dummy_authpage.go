package bindings

import (
	"errors"
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/openbolt/openid"
)

// BUG Implement 5.5.1.1.  Requesting the "acr" Claim
func (ds *DummySource) Authpage(w http.ResponseWriter, r *http.Request) openid.AuthState {
	var warn string

	// Was submit button pressed
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

/*
 * Cacher
 */
// BUG TODO: Locking
var sessions map[string]openid.Session

func (ds *DummySource) Cache(c openid.Session) error {
	if sessions == nil {
		sessions = make(map[string]openid.Session)
	}
	sessions[c.Code] = c
	return nil
}

func (ds *DummySource) GetSession(code string) (openid.Session, error) {
	s, ok := sessions[code]
	if !ok {
		return openid.Session{}, errors.New("Invalid code")
	} else {
		return s, nil
	}
}

func (ds *DummySource) Retire(code string) {
	// BUG Implement!
}
