package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tmavrin/tp-link/encryption"
)

type Client struct {
	host       string
	httpClient *http.Client
	rsa        *encryption.RSA
	aes        *encryption.AES
	seq        int
}

func Authenticate(host string, username string, password string) (*Client, error) {
	rsaKey, seq, err := getAuthParams(host)
	if errors.Is(err, ErrorInternal) {
		rsaKey, seq, err = getAuthParams(host)
	}
	if err != nil {
		return nil, fmt.Errorf("failed getting auth params: %w", err)
	}

	// TODO: generate this randomly
	key := "1234567890123456"
	iv := "somepass12345672"

	auth := &Auth{
		Host: host,
	}

	c := Client{
		host: host,
		rsa:  encryption.NewRSA(rsaKey),
		aes:  encryption.NewAES([]byte(key), []byte(iv)),
		seq:  seq,
		httpClient: &http.Client{
			Transport: &Transport{
				Base: http.DefaultTransport,
				Auth: auth,
			},
		},
	}

	auth.SessionID, auth.TokenID, err = c.login(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed calling login: %w", err)
	}

	return &c, nil
}

func (c *Client) Close() error {
	return c.MakeRequest(Logout, nil)
}
