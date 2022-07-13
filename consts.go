package main

import (
	"fmt"
)

const url = "https://api.kitkot.io"

// TODO: Implement client certificate check

func GetURL(path string) string {
	return fmt.Sprintf("%s/%s",
		url, path)
}

// mTLS certificate and private key required to talk to our backend
var cloudflareCertificate = []byte("-----BEGIN CERTIFICATE-----\nMIIDSjCCAjKgAwIBAgIUNlTwS4Q9gDSeL6UGTApnCmUO05YwDQYJKoZIhvcNAQEL\nBQAwgagxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQH\nEw1TYW4gRnJhbmNpc2NvMRkwFwYDVQQKExBDbG91ZGZsYXJlLCBJbmMuMRswGQYD\nVQQLExJ3d3cuY2xvdWRmbGFyZS5jb20xNDAyBgNVBAMTK01hbmFnZWQgQ0EgMWJi\nMDZjYTVlMmQ3ZmIxMWVjMjNmODIwMDViYmE0MDgwHhcNMjIwNDI5MjI0NjAwWhcN\nMzcwNDI1MjI0NjAwWjAiMQswCQYDVQQGEwJVUzETMBEGA1UEAxMKQ2xvdWRmbGFy\nZTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABM38RhouS7slpuwTpQI0KN09HAeS\nww2n5tSNJN4TofDpevfJ+KvpB/5lmE0LQXp/Ua+SGceTA3OLb/5htFrE3b+jgbsw\ngbgwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYDVR0TAQH/BAIwADAdBgNVHQ4EFgQU\n5+kFzeJZ+I/0Lyw69++ZcPZa26kwHwYDVR0jBBgwFoAUFOXV1FaAfLluojP83Ajf\nRVIHSF0wUwYDVR0fBEwwSjBIoEagRIZCaHR0cDovL2NybC5jbG91ZGZsYXJlLmNv\nbS9hNWZmYzU3Yi05ZTQwLTRhZjAtYTQ1Yy0wZjdhZDhhMzBmNDMuY3JsMA0GCSqG\nSIb3DQEBCwUAA4IBAQAEv8RcL0/0LIvFEax+pfp8Jyfzc3NNwxyRSwbiSRH/Eozp\n+2QuBXbAV8FRO7r4AYhdeWqOmHTC3aywzTkLW3HV3bhxSBx+bbpTfqj9Ev0evDHV\nZGhln53Xs/YDOJDuoFvKb6thb0YgnO7vjy2JPYMIda649sMq/ZZl4RgiVnDV5cCn\nkbFmro1i8tPWuABEFBn2KGWBFLTuir+P5dHvsfao27cghvQUh6yTN2EqRrXObldp\npTKlHf+AKY5J8mqd0xJDJO65eYtw68uGIXFmvl9qRoGVn4EtWYeVJtq17qe5b4V1\nPbH/jmQtv82FctuJmMGHf+/c/TGRn4vC7v4bexFU\n-----END CERTIFICATE-----\n")
var cloudflarePrivateKey = []byte("-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgwebG9S4q4uDWVYFk\nPXC1t4w8eNXNXdLvOSBscNjqpKuhRANCAATN/EYaLku7JabsE6UCNCjdPRwHksMN\np+bUjSTeE6Hw6Xr3yfir6Qf+ZZhNC0F6f1GvkhnHkwNzi2/+YbRaxN2/\n-----END PRIVATE KEY-----\n")
