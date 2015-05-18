package openid

import (
	_ "testing"
)

func Example_ReadClaimsRequest() {
	data := "blabla"
	ReadClaimsRequest(data)
	// Output: bla
}
