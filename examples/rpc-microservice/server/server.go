package server

import (
	"go-micro-framework/microservice"
	"go-micro-framework/microservice/configuration"
	"go-micro-framework/microservice/server"
)

// Register your server here. Server options are available in SDS_common
// microservice/server package.
func Register(ms microservice.Microservice) (err error) {
	rpcCfg := configuration.MyRPCConfig()
	err = rpcCfg.LoadConfig("rpc.json")
	if err != nil {
		return
	}

	// Instantiate the rpc server
	rpcserver, err := server.MyRPCServer(rpcCfg)
	if err != nil {
		return
	}

	// Register the server with the microservice.
	// Microservice.Start() walks through all the servers registered with the microservice
	// and call their respective Start() routines.
	err = ms.RegisterServer(server.ServerType("rpc"), rpcserver)
	if err != nil {
		return
	}
	return
}
