package route

import "github.com/tars-vcms/vcms-protocol/rbac_server"

type AuthType uint8

const (
	// TOKEN 仅使用jwt payload字段中的role进行认证，用于安全要求较低的接口
	TOKEN AuthType = 1
	// API 每次请求调用下游RBAC服务，用以安全度要求较高的接口
	API AuthType = 2
)

type ServantAuth struct {
	Type AuthType
	// RolesName Token验证模式时所比较的用户组名
	RolesName []string
	// RbacCode 调用Rbac服务需要携带的Code
	RbacID int64
}

var TARS_AUTH_TYPE_NAME = map[rbac_server.AUTH_TYPE]AuthType{
	rbac_server.AUTH_TYPE_TOKEN: TOKEN,
	rbac_server.AUTH_TYPE_API:   API,
}

func NewServantAuth(auth rbac_server.RouteAuth) *ServantAuth {
	return &ServantAuth{

		Type: TARS_AUTH_TYPE_NAME[auth.Type],

		RolesName: auth.RolesName,

		RbacID: auth.RbacID,
	}
}
