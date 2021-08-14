package cache

type MQRouteCmdType string

const (
	ROUTE_VERSION_SUFFIX = "_route_version"

	MQ_SUFFIX = "_route"

	// CMD_ROUTE_UPDATE 路由更新KEY
	CMD_ROUTE_UPDATE MQRouteCmdType = "CMD_ROUTE_UPDATE"

	// CMD_ROUTE_REQUIRE 向route请求预热路由信息
	CMD_ROUTE_REQUIRE MQRouteCmdType = "CMD_ROUTE_REQUIRE"
)

type MQRouteCmd struct {
	CMD MQRouteCmdType `json:"cache"`

	Payload string `json:"payload"`

	Timestamp int64 `json:"timestamp"`
}
