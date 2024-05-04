package sms

import (
	"fmt"
	"time"

	"github.com/tmavrin/tp-link/client"
)

type SMS struct {
	Index       int
	From        string
	To          string
	Content     string
	ReceiveTime time.Time
	Unread      bool
}

func (s *SMS) String() string {
	return fmt.Sprintf(`
index: %d
from: %s
to: %s
content: %s
receive_time: %s
unread: %v
`,
		s.Index,
		s.From,
		s.To,
		s.Content,
		s.ReceiveTime.String(),
		s.Unread,
	)
}

type c struct {
	client *client.Client
}

func New(client *client.Client) *c {
	return &c{client: client}
}

var (
	smsSend client.TPRequest = client.TPRequest{
		Method:     client.MethodSet,
		Controller: "LTE_SMS_SENDNEWMSG",
		Attributes: []string{"index", "to", "textContent"},
	}

	inboxList client.TPRequest = client.TPRequest{
		Method:     client.MethodSet,
		Controller: "LTE_SMS_RECVMSGBOX",
		Attributes: []string{"PageNumber"},
	}

	inboxListEntries client.TPRequest = client.TPRequest{
		Method:     client.MethodGL,
		Stack:      1,
		Controller: "LTE_SMS_RECVMSGENTRY",
		Attributes: []string{"index", "from", "content", "receivedTime", "unread"},
	}
)
