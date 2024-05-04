package client

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

type TPRequestWithArgs struct {
	Req  TPRequest
	Args map[string]string
}

func (c *Client) MakeRequest(reqs []TPRequestWithArgs) (string, error) {
	var data string

	for _, req := range reqs {
		if data != "" {
			data += "&"
		}
		data += strconv.Itoa(int(req.Req.Method))
	}

	data += "\r\n"

	for _, req := range reqs {
		command := fmt.Sprintf("[%s#0,0,0,0,0,0#0,0,0,0,0,0]%d,%d", req.Req.Controller, req.Req.Stack, len(req.Req.Attributes))

		data += fmt.Sprintf("%s\r\n", command)

		for _, v := range req.Req.Attributes {
			if req.Args[v] == "" {
				data += fmt.Sprintf("%s\r\n", v)
			} else {
				data += fmt.Sprintf("%s=%s\r\n", v, req.Args[v])
			}
		}
	}

	encryptedData, err := c.aes.Encrypt(data)
	if err != nil {
		return "", fmt.Errorf("encrypting data: %w", err)
	}

	signedData := c.rsa.Encrypt(fmt.Sprintf("h=undefined&s=%d", c.seq+len(encryptedData)))

	requestBody := fmt.Sprintf("sign=%s\r\ndata=%s\r\n", signedData, encryptedData)

	response, err := c.httpClient.Post(fmt.Sprintf("%s/cgi_gdpr", c.host), "text/plain", strings.NewReader(requestBody))
	if err != nil {
		return "", fmt.Errorf("making post request: %w", err)
	}

	encryptedResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("reading response body: %w", err)
	}

	body, err := c.aes.Decrypt(string(encryptedResponse))
	if err != nil {
		return "", fmt.Errorf("decrypting response body: %w", err)
	}

	if strings.Contains(body, "[error]0") {
		return body, nil
	}

	return body, fmt.Errorf("error response: %s", body)
}
