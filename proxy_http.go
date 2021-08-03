package main

import (
	"github.com/TarsCloud/TarsGo/tars"
	"net/http"
)

type ProxyHttp interface {
	handleRequest(w http.ResponseWriter, r *http.Request)
	handleError(w http.ResponseWriter, code int, err error)
	GetTarsHttpMux() *tars.TarsHttpMux
}

func NewProxyHttp() ProxyHttp {
	return newProxyHttpImp()
}
