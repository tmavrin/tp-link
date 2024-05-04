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
	Stack      int
	Attributes []string
}

var (
	Logout TPRequest = TPRequest{
		Method:     MethodCGI,
		Controller: "/cgi/logout",
		Attributes: []string{},
	}
)
