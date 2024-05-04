package client

import (
	"fmt"
	"net/http"
)

type Auth struct {
	Host      string
	TokenID   string
	SessionID string
}

type Transport struct {
	Auth *Auth
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Referer", t.Auth.Host)
	req.Header.Add("Origin", t.Auth.Host+"/")
	req.Header.Add("Cookie", fmt.Sprintf("loginErrorShow=1; JSESSIONID=%s", t.Auth.SessionID))
	req.Header.Add("TokenID", t.Auth.TokenID)
	return http.DefaultTransport.RoundTrip(req)
}
