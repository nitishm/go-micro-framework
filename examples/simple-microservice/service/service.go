package service

import (
	"go-micro-framework/examples/simple-microservice/service/endpoint"
	server "go-micro-framework/microservice/server"

	log "github.com/sirupsen/logrus"
)

func printError(err error) (errString string) {
	if err != nil {
		return err.Error()
	}
	return
}

type ServiceMgr struct {
	servers map[server.ServerType]server.ServerInterface
}

func MyServiceMgr() (res *ServiceMgr, err error) {
	res = &ServiceMgr{
		servers: make(map[server.ServerType]server.ServerInterface),
	}
	return
}

func (ss *ServiceMgr) RegisterServer(serverType server.ServerType, server server.ServerInterface) (err error) {
	ss.servers[serverType] = server
	return
}

func (ss *ServiceMgr) Start() (err error) {
	err = ss.init()
	if err != nil {
		return
	}
	for t, s := range ss.servers {
		log.Warnf("Starting server %v\n", t)
		err = s.Start()
		if err != nil {
			return err
		}
	}
	return
}

func (ss *ServiceMgr) Stop() (err error) {
	for t, s := range ss.servers {
		log.Warnf("Stopping server %v\n", t)
		err = s.Stop()
		if err != nil {
			return err
		}
	}
	return

}

func (ss *ServiceMgr) init() (err error) {
	helloEndpoint := endpoint.MyHelloEndpoint()
	sv := ss.servers[server.ServerType("kafka")]
	sv.RegisterNamespace("greeter")
	sv.RegisterService("greeter", server.Service("hello"), helloEndpoint)
	return
}
