package main

import (
	"fmt"
	"github.com/openbolt/openid"
	"github.com/openbolt/openid/bindings"
	"log"
	"net/http"
)

func main() {
	op := openid.NewProvider()
	src := bindings.DummySource{}

	// Add datasources
	op.SetAuthsource(src)
	op.SetClaimsource(src)
	op.SetClientsource(src)

	// Configure http api
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello_world)

	op.AddServer(mux)

	// Start OpenID Provider
	if err := op.Serve(); err != nil {
		log.Fatal(err)
	}

	// Add http listener to it
	fmt.Println("OpenID Connect 1.0 Core Provider demo started")
	fmt.Println("Go to https://localhost:8443")
	err := http.ListenAndServeTLS("localhost:8443", "demo.crt", "demo.pem", mux)
	log.Fatal(err)
}

func hello_world(w http.ResponseWriter, r *http.Request) {
	data, _ := Asset("hello.html")
	w.Write(data)
}

//go:generate go-bindata -o bindata.go hello.html
