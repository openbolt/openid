package openid_test

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/openbolt/openid"
	"github.com/openbolt/openid/bindings"
)

var client *http.Client
var op *openid.OpenID
var ts *httptest.Server

func TestMain(m *testing.M) {
	// Setup HTTP(S) client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{Transport: tr}

	// Setup openid provider
	op = openid.NewProvider()

	// Set bindings
	src := new(bindings.DummySource)
	op.SetAuthsource(src)
	op.SetClaimsource(src)
	op.SetClientsource(src)
	op.SetEnduserIf(src)

	// Set TLS Server
	mux := http.NewServeMux()
	mux.HandleFunc("/backlink", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("backlink!"))
	})
	op.AddServer(mux)
	ts = httptest.NewUnstartedServer(mux)
	defer ts.Close()
	ts.StartTLS()

	// Start
	op.Serve()

	os.Exit(m.Run())
}

type getResp struct {
	Url    url.URL
	Body   string
	Status int
	Header http.Header
}

// sGet (simple Get) does an GET request to the OpenID server
func sGet(path string) getResp {
	resp, err := client.Get(ts.URL + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	r := getResp{}
	r.Url = *resp.Request.URL
	r.Status = resp.StatusCode
	bbody, _ := ioutil.ReadAll(resp.Body)
	r.Body = string(bbody)
	r.Header = resp.Header

	return r
}

// Test Authorization
func TestHttpAuthorize_negative(t *testing.T) {
	redirect_url := ts.URL + "/backlink"

	resp := sGet("/authorize")
	if !strings.Contains(resp.Body, "invalid_request") {
		t.Error("Wrong response when request with no params")
	}

	t.Log("Next")
	resp = sGet("/authorize?scope=openid&redirect_uri=" + redirect_url)
	if resp.Url.Path != "/backlink" {
		t.Error("Redirect on error not working")
	}
}

func TestHttpAuthorize_positive(t *testing.T) {
	redirect_url := ts.URL + "/backlink"

	resp := sGet("/authorize?" +
		"scope=openid" +
		"&client_id=cltTest1" +
		"&redirect_uri=" + redirect_url +
		"&response_type=code")
	t.Error(resp)
	if !strings.Contains(resp.Url.String(), "invalid_request") {
		t.Error("")
	}
}
