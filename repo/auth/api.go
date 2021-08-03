package auth

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
)

type HttpAuth interface {

	// ChallengeHttpAuth  验证http中的jwt，若需要回调用RbacAPI
	ChallengeHttpAuth(request *http.Request, response http.ResponseWriter, route *route.HttpRoute) error

	ParseJwtPayload(request *http.Request) (string, error)
}

func NewHttpAuth() HttpAuth {
	return newHttpAuthImpl()
}
