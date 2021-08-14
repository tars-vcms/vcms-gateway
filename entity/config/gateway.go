package config

const (
	GATEWAY_FILE_NAME = "gateway.conf"
)

type GatewayConfig struct {
	Name      string            `tars:"/gateway/Name"`
	HeaderMap map[string]string `tars:"/gateway<header>"`
}
