package client

import (
	"fmt"
	"io"
	"strings"
)

func (c *Client) MakeRequest(req TPRequest, args map[string]string) error {
	for _, v := range req.Attributes {
		if args[v] == "" {
			return fmt.Errorf("'%s' not passed in as args", v)
		}
	}

	command := fmt.Sprintf("[%s#0,0,0,0,0,0#0,0,0,0,0,0]0,%d", req.Controller, len(args))

	data := fmt.Sprintf("%d\r\n%s\r\n", req.Method, command)

	for _, v := range req.Attributes {
		data += fmt.Sprintf("%s=%s\r\n", v, args[v])
	}

	encryptedData, err := c.aes.Encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypting data: %w", err)
	}

	signedData := c.rsa.Encrypt(fmt.Sprintf("h=undefined&s=%d", c.seq+len(encryptedData)))

	requestBody := fmt.Sprintf("sign=%s\r\ndata=%s\r\n", signedData, encryptedData)

	response, err := c.httpClient.Post(fmt.Sprintf("%s/cgi_gdpr", c.host), "text/plain", strings.NewReader(requestBody))
	if err != nil {
		return fmt.Errorf("making post request: %w", err)
	}

	encryptedResponse, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	body, err := c.aes.Decrypt(string(encryptedResponse))
	if err != nil {
		return fmt.Errorf("decrypting response body: %w", err)
	}

	if strings.Contains(body, "[error]0") {
		return nil
	}

	return fmt.Errorf("error response: %s", body)
}
