package client

type TPMethod int

const (
	MethodGet TPMethod = 1
	MethodSet TPMethod = 2
	MethodDel TPMethod = 4
	MethodGL  TPMethod = 5
	MethodGS  TPMethod = 6
	MethodCGI TPMethod = 8
)

type TPRequest struct {
	Method     TPMethod
	Controller string
	Attributes []string
}

var (
	SMSSend TPRequest = TPRequest{
		Method:     MethodSet,
		Controller: "LTE_SMS_SENDNEWMSG",
		Attributes: []string{"index", "to", "textContent"},
	}

	Logout TPRequest = TPRequest{
		Method:     MethodCGI,
		Controller: "/cgi/logout",
		Attributes: []string{},
	}
)
