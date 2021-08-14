package main

import (
	"github.com/TarsCloud/TarsGo/tars"
)

func main() {
	// Get server routes
	cfg := tars.GetServerConfig()
	proxyHttp := NewProxyHttp()
	tars.AddHttpServant(proxyHttp.GetTarsHttpMux(), cfg.App+"."+cfg.Server+".ProxyHTTPObj")

	// Run application
	tars.Run()
}
