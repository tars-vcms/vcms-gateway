package route

type ServantType uint8

const (
	TarsServant         = 1
	ReverseProxyServant = 2
)

type HttpRoute struct {
	Path        string
	ServantName string
	FuncName    string
	Type        ServantType
	Auth        *ServantAuth
	Children    []*HttpRoute
}
