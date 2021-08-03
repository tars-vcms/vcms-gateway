package header

import (
	"github.com/tars-vcms/vcms-gateway/repo/rcfg"
	"net/http"
)

type HttpHeaderImpl struct {
	headers map[string]string
}

func (h HttpHeaderImpl) CheckRequestHeader(req *http.Request) (httpCode int, err error) {
	return 200, err
}

func (h HttpHeaderImpl) InjectCommonHeader(header http.Header) {
	for k, v := range h.headers {
		header.Set(k, v)
	}
}

func (h HttpHeaderImpl) HandleErrorInfo(w http.ResponseWriter, err error) {
	panic("implement me")
}

func newHttpHeaderImpl() *HttpHeaderImpl {
	var content string
	if err := rcfg.GetInstance().GetConfig("gateway.conf", rcfg.TEXT, &content); err != nil {
		panic(err.Error())
	}
	c, err := rcfg.GetInstance().ParseConfig(content)
	if err != nil {
		panic(err.Error())
	}
	headers := map[string]string{}
	for k, v := range DefaultHeaders {
		headers[k] = v
	}
	for k, v := range c.GetMap("/gateway/headers") {
		headers[k] = v
	}
	return &HttpHeaderImpl{
		headers: headers,
	}
}
