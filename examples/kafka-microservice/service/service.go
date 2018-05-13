package service

import (
	"go-micro-framework/examples/kafka-microservice/service/endpoint"
	server "go-micro-framework/microservice/server"

	log "github.com/sirupsen/logrus"
)

func printError(err error) (errString string) {
	if err != nil {
		return err.Error()
	}
	return
}

// ServiceMgr is the main structure of the microservice.
// It must implement the microservice interface.
// The ServiceMgr encapsulates all the servers used by the microservice.
type ServiceMgr struct {
	servers map[server.ServerType]server.ServerInterface
}

// MyServiceMgr returns an instance of the ServiceMgr type
func MyServiceMgr() (res *ServiceMgr, err error) {
	res = &ServiceMgr{
		servers: make(map[server.ServerType]server.ServerInterface),
	}
	return
}

// RegisterServer is an implementation of the microserver RegisterServer method.
// The main function calls this with all the instantiated servers, which are then stored in a map
// identifier by their Service type name.
func (ss *ServiceMgr) RegisterServer(serverType server.ServerType, server server.ServerInterface) (err error) {
	ss.servers[serverType] = server
	return
}

// Start is an implementation of the microservice Start method.
// It walks a map of all registered Servers and invokes each server's
// Start() methods, in-order.
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

// Stop is an implementation of the microservice Stop method.
// The Stop method walks a map of all the Started servers and gracefully
// shut them down in the order in which they were started.s
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

// Do all other initialization work, specific to the microservice here.
func (ss *ServiceMgr) init() (err error) {
	sv := ss.servers[server.ServerType("kafka")]
	// Register a namespace
	sv.RegisterNamespace("greeter")

	// Register your endpoints against the namespace with the server
	sv.RegisterService("greeter", server.Service("hello"), endpoint.MyHelloEndpoint())
	return
}
