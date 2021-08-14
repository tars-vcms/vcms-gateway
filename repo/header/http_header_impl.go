package header

import (
	"github.com/tars-vcms/vcms-gateway/entity/config"
	"github.com/tars-vcms/vcms-gateway/repo/rcfgs"
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
	gatewayCfg := &config.GatewayConfig{
		HeaderMap: make(map[string]string),
	}
	if err := rcfgs.GetInstance().GetConfig(config.GATEWAY_FILE_NAME, rcfgs.STRUCT, gatewayCfg); err != nil {
		panic(err.Error())
	}
	headerMap := gatewayCfg.HeaderMap
	for k, v := range DefaultHeaders {
		headerMap[k] = v
	}
	return &HttpHeaderImpl{
		headers: headerMap,
	}
}
