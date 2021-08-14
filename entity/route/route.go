package route

type ServantType uint8

const (
	TARS_SERVANT          = 1
	RESERVE_PROXY_SERVANT = 2
)

type HttpRoute struct {
	Path               string       `json:"path"`
	ServantName        string       `json:"servant_name"`
	FuncName           string       `json:"func_name"`
	InputName          string       `json:"input_name"`
	OutputName         string       `json:"output_name"`
	Type               ServantType  `json:"type"`
	Auth               *ServantAuth `json:"auth"`
	TransparentHeaders []string     `json:"transparent_headers"`
	Children           []*HttpRoute `json:"children"`
}
