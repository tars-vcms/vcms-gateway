package proxy

import (
	"context"
	"encoding/json"
	"github.com/TarsCloud/TarsGo/tars"
	"github.com/TarsCloud/TarsGo/tars/protocol/codec"
	"github.com/TarsCloud/TarsGo/tars/protocol/res/requestf"
	"github.com/tars-vcms/vcms-common/errs"
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

func (t *TupProxy) DoRequest(routeHttp *route.HttpRoute, w http.ResponseWriter, r *http.Request) (int, error) {
	reqContext, err := t.createContext(r, routeHttp.TransparentHeaders)
	if err != nil {
		return 500, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 400, err
	}

	input, err := t.marshalInput(routeHttp.InputName, body)
	if err != nil {
		return 400, err
	}

	var status map[string]string
	resp := &requestf.ResponsePacket{}
	ctx := context.Background()
	if err := t.Proxy.Tars_invoke(ctx, 0, routeHttp.FuncName, input, status, reqContext, resp); err != nil {
		return 502, err
	}
	if t.callback != nil && t.callback.HandleResponseHeader != nil {
		if err := t.callback.HandleResponseHeader(w.Header()); err != nil {
			return 502, err
		}
	}
	if err = errs.CatchError(resp.Context); err != nil {
		return 200, err
	}
	tarsResp := codec.FromInt8(resp.SBuffer)
	body, err = t.unmarshalOutput(routeHttp.OutputName, tarsResp)
	if err != nil {
		return 502, err
	}
	if _, err = w.Write(body); err != nil {
		return 502, err
	}
	return 200, nil
}

func (t *TupProxy) createContext(r *http.Request, tHeaders []string) (map[string]string, error) {
	c := make(map[string]string)
	if err := r.ParseForm(); err != nil {
		return nil, err
	}
	// 注入http query参数
	c["query"] = r.Form.Encode()
	header := r.Header
	// 注入需要透传的header
	for _, h := range tHeaders {
		c[h] = header.Get(h)
	}
	return c, nil
}

func (t *TupProxy) marshalInput(inputName string, httpBody []byte) ([]byte, error) {
	var body interface{}
	if err := json.Unmarshal(httpBody, &body); err != nil {
		return nil, err
	}

	input := make(map[string]interface{})
	input[inputName] = body

	return json.Marshal(input)
}

func (t *TupProxy) unmarshalOutput(outputName string, tarsResp []byte) ([]byte, error) {
	var resp map[string]interface{}

	if err := json.Unmarshal(tarsResp, &resp); err != nil {
		return tarsResp, err
	}
	output, ok := resp[outputName]
	if !ok {
		return tarsResp, nil
	}
	return json.Marshal(output)
}
