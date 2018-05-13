package server

import "go-micro-framework/microservice/endpoint"

// ServerType is a typedef for server identifier
type ServerType string

// Service is a typedef for service identifier
type Service string

// ServiceEndpointMap is a typedef for a map of Service endpoints identified by their service type
type ServiceEndpointMap map[Service]endpoint.Endpoint

// Handler is a func signature implemented by the Endpoint Handle() method
type Handler func(in []byte, serviceEndpointMap ServiceEndpointMap) (out []byte, err error)

// StartStopInterface is composed of the Start() and Stop() methods to be implemented by the server
type StartStopInterface interface {
	// Start is used to start the server
	Start() (err error)
	// Stop is used to gracefully stop the server
	Stop() (err error)
}

// ServerInterface defines the methods that all servers registed with the microservice must implement
type ServerInterface interface {
	StartStopInterface
	// RegisterNamesapce is used to register a namespace (like a kafka channel/topic or grpc namespace)
	// with the server.
	RegisterNamespace(namespace string)
	// RegisterService is used to register a service and its endpoints with the server.
	RegisterService(namespace string, service Service, endpoint endpoint.Endpoint)
}
