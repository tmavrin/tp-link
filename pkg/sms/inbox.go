package sms

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tmavrin/tp-link/client"
)

func (c *c) GetInbox() ([]SMS, error) {
	response, err := c.client.MakeRequest([]client.TPRequestWithArgs{
		{
			Req: inboxList,
			Args: map[string]string{
				"PageNumber": "1",
			},
		},
		{
			Req: inboxListEntries,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed making get inbox request: %w", err)
	}

	smsList, err := parseInbox(response)
	if err != nil {
		return nil, fmt.Errorf("failed parsing inbox: %w", err)
	}

	return smsList, nil
}

func parseInbox(response string) ([]SMS, error) {
	var (
		smsList []SMS
		err     error
	)

	msgs := strings.Split(response, ",0,0,0,0,0]1\n")

	for _, m := range msgs[1:] {
		var sms SMS

		props := strings.Split(m, "\n")

		if strings.HasPrefix(props[0], "index=") {
			sms.Index, err = strconv.Atoi(props[0][6:])
			if err != nil {
				return nil, fmt.Errorf("failed parsing msg index: %w", err)
			}
		}
		if strings.HasPrefix(props[1], "from=") {
			sms.From = props[1][5:]
		}
		if strings.HasPrefix(props[2], "content=") {
			sms.Content = props[2][8:]
		}

		if strings.HasPrefix(props[3], "receivedTime=") {
			sms.ReceiveTime, err = time.Parse(time.DateTime, props[3][13:])
			if err != nil {
				return nil, fmt.Errorf("failed parsing msg received time: %w", err)
			}
		}

		if strings.HasPrefix(props[4], "unread=") {
			sms.Unread = props[4][7:] == "1"
		}

		smsList = append(smsList, sms)
	}

	return smsList, nil
}
