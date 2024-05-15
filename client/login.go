package client

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	eeRgx   = regexp.MustCompile(`ee="(\d+)`)
	nnRgx   = regexp.MustCompile(`var nn="([0-9A-F]+)`)
	seqRgx  = regexp.MustCompile(`var seq="(\d+)`)
	homeRgx = regexp.MustCompile(`var token="([a-f0-9]+)"`)

	ErrorInternal      = errors.New("internal server error")
	ErrorSessionNotSet = errors.New("session not set")
)

func getAuthParams(host string) (*rsa.PublicKey, int, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/cgi/getParm", host), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("make request: %w", err)
	}

	req.Header.Add("Referer", host)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("do request: %w", err)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read body: %w", err)
	}

	if strings.Contains(string(bodyBytes), "500 Internal Server Error") {
		return nil, 0, ErrorInternal
	}

	eeRF := eeRgx.Find(bodyBytes)
	ee := strings.TrimPrefix(string(eeRF), "ee=\"")
	eeInt, err := strconv.ParseInt(ee, 16, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("parse ee to int: %w", err)
	}

	nnRF := nnRgx.Find(bodyBytes)
	nn := strings.TrimPrefix(string(nnRF), "var nn=\"")

	nnBigInt := new(big.Int)
	nnBigInt.SetString(nn, 16)

	seqRF := seqRgx.Find(bodyBytes)

	seq := strings.TrimPrefix(string(seqRF), "var seq=\"")

	seqInt, err := strconv.Atoi(seq)
	if err != nil {
		return nil, 0, fmt.Errorf("parse seq to int: %w", err)
	}

	return &rsa.PublicKey{
		N: nnBigInt,
		E: int(eeInt),
	}, seqInt, nil
}

func (c *Client) login(username string, password string) (string, string, error) {
	encryptedData, err := c.aes.Encrypt(fmt.Sprintf("%s\n%s", username, password))
	if err != nil {
		return "", "", fmt.Errorf("encrypt login data: %w", err)
	}

	encryptedSign := c.rsa.Encrypt(fmt.Sprintf(
		"key=%s&iv=%s&h=%s&s=%d",
		c.aes.Key, c.aes.IV, "undefined", c.seq+len(encryptedData)),
	)

	url := fmt.Sprintf("%s/cgi/login?data=%s&sign=%s&Action=1&LoginStatus=0", c.host, encryptedData, encryptedSign)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return "", "", fmt.Errorf("make login request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("do login request: %w", err)
	}

	var sessionID string
	for _, c := range res.Cookies() {
		if c.Name == "JSESSIONID" {
			sessionID = c.Value
		}
	}

	if sessionID == "deleted" || sessionID == "" {
		return "", "", ErrorSessionNotSet
	}

	log.Println("- logged in successfully")

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/", c.host), nil)
	if err != nil {
		return "", "", fmt.Errorf("make home request: %w", err)
	}

	req.Header.Add("Cookie", fmt.Sprintf("loginErrorShow=1; JSESSIONID=%s", sessionID))

	res, err = c.httpClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("do home request: %w", err)
	}

	homeB, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("read home request body: %w", err)
	}

	token := string(homeRgx.Find(homeB))
	token = strings.TrimPrefix(string(token), "var token=\"")

	if len(token) == 0 {
		return "", "", fmt.Errorf("token not found in home body: %w", err)
	}

	token = token[0 : len(token)-1]

	log.Println("- fetched token successfully")

	return sessionID, token, nil

}
