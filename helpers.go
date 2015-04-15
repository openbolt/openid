package openid

import (
	"net/http"
)

func getParam(r http.Request, param string) string {
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

// TODO: Ev. Validate(param) func?
