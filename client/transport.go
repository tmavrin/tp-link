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
	Base http.RoundTripper
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	reqBodyClosed := false
	if req.Body != nil {
		defer func() {
			if !reqBodyClosed {
				req.Body.Close()
			}
		}()
	}

	req2 := cloneRequest(req)

	req2.Header.Add("Referer", t.Auth.Host)
	req2.Header.Add("Origin", t.Auth.Host+"/")
	req2.Header.Add("Cookie", fmt.Sprintf("loginErrorShow=1; JSESSIONID=%s", t.Auth.SessionID))
	req2.Header.Add("TokenID", t.Auth.TokenID)

	reqBodyClosed = true
	return t.base().RoundTrip(req2)
}

func (t *Transport) base() http.RoundTripper {
	if t.Base != nil {
		return t.Base
	}
	return http.DefaultTransport
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
