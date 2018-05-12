package endpoint

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type HelloEndpoint struct {
}

func MyHelloEndpoint() *HelloEndpoint {
	return &HelloEndpoint{}
}

func (ep *HelloEndpoint) Handle(ctx context.Context, request []byte) (response []byte, err error) {
	log.Infof("\n=======\nHello %s\n=======\n", request)
	return
}
