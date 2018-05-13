package server

import (
	"encoding/json"
	"go-micro-framework/microservice/configuration"
	"go-micro-framework/microservice/endpoint"
	"log"
	"net"

	context "golang.org/x/net/context"

	grpc "google.golang.org/grpc"
)

type RPCNamespace struct {
	Name               string
	ServiceEndpointMap ServiceEndpointMap
}

func MyRPCNamespace(name string) (namespace RPCNamespace) {
	return RPCNamespace{
		Name:               name,
		ServiceEndpointMap: make(ServiceEndpointMap),
	}
}

type RPCServer struct {
	addr       string
	namespaces map[string]RPCNamespace
	server     *grpc.Server
}

func MyRPCServer(cfg *configuration.RPCConfig) (rpc *RPCServer, err error) {
	rpc = &RPCServer{
		addr:       cfg.Address,
		namespaces: make(map[string]RPCNamespace),
	}
	return
}

func (rpc *RPCServer) Start() (err error) {
	lis, err := net.Listen("tcp", rpc.addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rpc.server = grpc.NewServer()
	RegisterMessageServiceServer(rpc.server, rpc)
	go rpc.server.Serve(lis)

	return
}

func (rpc *RPCServer) Stop() (err error) {
	rpc.server.GracefulStop()
	return
}

func (rpc *RPCServer) RegisterNamespace(name string) {
	rpc.namespaces[name] = MyRPCNamespace(name)
	return
}

func (rpc *RPCServer) RegisterService(namespace string, service Service, ep endpoint.Endpoint) {
	rpc.namespaces[namespace].ServiceEndpointMap[service] = ep
	return
}

func (rpc *RPCServer) Command(c context.Context, msg *RPCMessage) (resp *RPCResponse, err error) {
	out, err := ByteHandler(msg.Data, rpc.namespaces[msg.Namespace].ServiceEndpointMap)
	if err != nil {
		resp = &RPCResponse{
			Error: err.Error(),
		}
		return
	}
	resp = &RPCResponse{
		Data: out,
	}
	return
}

type RPCClient struct {
	addr   string
	client MessageServiceClient
}

func MyRPCClient(cfg *configuration.RPCConfig) (rpc *RPCClient, err error) {
	conn, err := grpc.Dial(cfg.Address, grpc.WithInsecure())
	if err != nil {
		return
	}
	rpc = &RPCClient{
		addr:   cfg.Address,
		client: NewMessageServiceClient(conn),
	}
	return
}

func (rpc *RPCClient) Send(namespace string, msg interface{}) (resp string, err error) {
	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	respObj, err := rpc.client.Command(
		context.Background(),
		&RPCMessage{
			Namespace: namespace,
			Data:      b,
		},
	)
	if err != nil {
		return
	}

	resp = string(respObj.Data)
	return
}
