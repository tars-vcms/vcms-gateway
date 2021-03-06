package header

import "net/http"

var DefaultHeaders = map[string]string{
	"Content-ServantType":          "application/json",
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "POST, GET",
	"Cache-Control":                "no-cache",
	"Tars-Ret":                     "0",
	"Tars-Msg":                     "success",
}

type HttpHeader interface {
	InjectCommonHeader(w http.Header)

	CheckRequestHeader(req *http.Request) (httpCode int, err error)

	HandleErrorInfo(w http.ResponseWriter, err error)
}

func NewHttpHeader() HttpHeader {
	return newHttpHeaderImpl()
}
