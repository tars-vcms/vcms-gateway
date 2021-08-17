package route

import "github.com/tars-vcms/vcms-protocol/route_manager"

type HttpServantType uint8

type HttpRouteType uint8

const (
	TARS_SERVANT          = 1
	RESERVE_PROXY_SERVANT = 2
	// DICTIONARY_SERVANT 标识本路由仅仅为一个文件夹，不具有路由
	DICTIONARY_SERVANT = 3
)

var TARS_SERVANT_TYPE_NAME = map[route_manager.SERVANT_TYPE]HttpServantType{
	route_manager.SERVANT_TYPE_TARS_SERVANT:          TARS_SERVANT,
	route_manager.SERVANT_TYPE_RESERVE_PROXY_SERVANT: RESERVE_PROXY_SERVANT,
	route_manager.SERVANT_TYPE_DICTIONARY_SERVANT:    DICTIONARY_SERVANT,
}

type HttpRoute struct {
	Path               string
	ServantName        string
	FuncName           string
	InputName          string
	OutputName         string
	ServantType        HttpServantType
	Auth               *ServantAuth
	TransparentHeaders []string
	Children           []*HttpRoute
}

func NewRoutes(routes []route_manager.RouteTable) []*HttpRoute {
	if routes == nil {
		return nil
	}
	formatRoutes := make([]*HttpRoute, len(routes))
	for _, r := range routes {
		formatRoutes = append(formatRoutes, &HttpRoute{
			Path:               r.Path,
			ServantName:        r.ServantName,
			FuncName:           r.FuncName,
			InputName:          r.InputName,
			OutputName:         r.OutputName,
			ServantType:        TARS_SERVANT_TYPE_NAME[r.Type],
			Auth:               NewServantAuth(r.Auth),
			TransparentHeaders: r.TransparentHeaders,
			Children:           NewRoutes(r.Children),
		})
	}
	return nil
}
