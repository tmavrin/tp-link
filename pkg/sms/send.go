package sms

import (
	"github.com/tmavrin/tp-link/client"
)

func (c *c) SendSMS(sms SMS) error {
	_, err := c.client.MakeRequest([]client.TPRequestWithArgs{
		{
			Req: smsSend,
			Args: map[string]string{
				"index":       "1",
				"to":          sms.To,
				"textContent": sms.Content,
			},
		},
	})
	return err
}
