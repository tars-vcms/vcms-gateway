package auth

import (
	"github.com/tars-vcms/vcms-gateway/entity/route"
	"net/http"
)

type HttpAuthImpl struct {
}

func (h HttpAuthImpl) ParseJwtPayload(request *http.Request) (string, error) {
	return "", nil
}

func (h HttpAuthImpl) ChallengeHttpAuth(request *http.Request, response http.ResponseWriter, route *route.HttpRoute) error {
	return nil
}

func newHttpAuthImpl() *HttpAuthImpl {
	return &HttpAuthImpl{}
}
