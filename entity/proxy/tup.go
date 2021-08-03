package proxy

import (
	"context"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/codec"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/requestf"
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"io/ioutil"
	"net/http"
)

func NewTupProxy(proxy *tars.ServantProxy) GatewayProxy {
	return &TupProxy{
		Proxy:    proxy,
		callback: nil,
	}
}

type TupProxy struct {
	Proxy    *tars.ServantProxy
	callback *GatewayProxyCallback
}

func (t *TupProxy) SetCallback(cb *GatewayProxyCallback) {
	t.callback = cb
}

func (t TupProxy) DoRequest(routeHttp *route.HttpRoute, w http.ResponseWriter, r *http.Request) (int, error) {
	reqContext, err := t.createContext(r)
	if err != nil {
		return 500, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 400, err
	}

	var status map[string]string
	resp := &requestf.ResponsePacket{}
	ctx := context.Background()
	if err := t.Proxy.Tars_invoke(ctx, 0, routeHttp.FuncName, body, status, reqContext, resp); err != nil {
		return 502, err
	}
	if t.callback != nil && t.callback.HandleResponseHeader != nil {
		if err := t.callback.HandleResponseHeader(w.Header()); err != nil {
			return 502, err
		}
	}
	if _, err = w.Write(codec.FromInt8(resp.SBuffer)); err != nil {
		return 502, err
	}
	return 200, nil
}

func (t TupProxy) createContext(r *http.Request) (map[string]string, error) {
	c := make(map[string]string)
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	c["query"] = r.Form.Encode()
	return c, nil
}
