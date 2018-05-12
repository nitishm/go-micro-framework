package microservice

import "go-micro-framework/microservice/server"

type Microservice interface {
	RegisterServer(serverType server.ServerType, server server.ServerInterface) (err error)
	Start() (err error)
	Stop() (err error)
}
