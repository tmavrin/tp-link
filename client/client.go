package client

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
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

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Authenticate(host string, username string, password string) (*Client, error) {
	var (
		c   *Client
		err error
	)

	// sometimes we get random error on auth so we can retry
	for i := 0; i < 3; i++ {
		c, err = authenticate(host, username, password)
		if err != nil && !errors.Is(err, ErrorSessionNotSet) {
			return nil, err
		}
		if err == nil {
			return c, nil
		}
		fmt.Printf("Failed authenticating. Retry %d\n", i+1)
	}
	return nil, fmt.Errorf("failed to auth after 3 retries, %s", err)
}

func authenticate(host string, username string, password string) (*Client, error) {
	rsaKey, seq, err := getAuthParams(host)
	if errors.Is(err, ErrorInternal) {
		rsaKey, seq, err = getAuthParams(host)
	}
	if err != nil {
		return nil, fmt.Errorf("failed getting auth params: %w", err)
	}

	log.Println("- fetched auth params")

	key := make([]byte, 16)
	for i := range key {
		key[i] = letters[rand.Intn(len(letters))]
	}

	iv := make([]byte, 16)
	for i := range iv {
		iv[i] = letters[rand.Intn(len(letters))]
	}

	auth := &Auth{
		Host: host,
	}

	c := Client{
		host: host,
		rsa:  encryption.NewRSA(rsaKey),
		aes:  encryption.NewAES(key, iv),
		seq:  seq,
		httpClient: &http.Client{
			Transport: &Transport{
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
	_, err := c.MakeRequest([]TPRequestWithArgs{{Req: Logout, Args: nil}})
	if err != nil {
		return err
	}

	log.Println("- logged out successfully")

	return nil
}
