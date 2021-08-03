package proxy

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewReserveProxy(target *url.URL) GatewayProxy {
	p := httputil.NewSingleHostReverseProxy(target)
	r := &ReserveProxy{
		Proxy: p,
	}
	targetQuery := target.RawQuery
	p.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path
		req.URL.RawPath = req.URL.EscapedPath()
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}

	p.ModifyResponse = func(resp *http.Response) error {
		if r.callback != nil && r.callback.HandleResponseHeader != nil {
			if err := r.callback.HandleResponseHeader(resp.Header); err != nil {
				return err
			}
		}
		return nil
	}
	return r
}

type ReserveProxy struct {
	Proxy    *httputil.ReverseProxy
	callback *GatewayProxyCallback
}

func (r *ReserveProxy) SetCallback(cb *GatewayProxyCallback) {
	r.callback = cb
}

func (r ReserveProxy) DoRequest(routeHttp *route.HttpRoute, rw http.ResponseWriter, req *http.Request) (int, error) {

	r.Proxy.ServeHTTP(rw, req)
	return 200, nil
}
