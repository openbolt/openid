# Access Token
This describes the access token format for OAuth 2.0 use.

## Format
The access token is returned as two key-value pairs, separated by semicolon.
`TKN_SIG=TokenSignature; TKN_DATA=TokenData`
`TokenSignature` and `TokenData` are base64 strings.

### TokenData
TokenData is an JSON Object, which is serialized to an UTF-8 string and then encoded in base64.
The contents are the following:
```
{
	"exp": // Expiration time [seconds after 1970]
	"iss": // Issuer
	"sub": // Subject
	"iat": // Issuing time [seconds after 1970]
	"amr": // Authentication Methods References
	"acr": // Authentication Context Class Reference
	"nonce": // Nonce
}
```

### TokenSignature
RSA Signature of TokenData, encoded as base64

## Verification
