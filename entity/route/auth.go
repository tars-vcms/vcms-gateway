package route

type AuthType uint8

const (
	// TOKEN 仅使用jwt payload字段中的role进行认证，用于安全要求较低的接口
	TOKEN AuthType = 1
	// API 每次请求调用下游RBAC服务，用以安全度要求较高的接口
	API AuthType = 2
)

type ServantAuth struct {
	Type AuthType
	// RoleName Token验证模式时所比较的用户组名
	RoleName []string
	// RbacCode 调用Rbac服务需要携带的Code
	RbacCode string
}
