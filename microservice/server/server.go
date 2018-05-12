package server

import "go-micro-framework/microservice/endpoint"

type ServerType string

type Service string

type ServiceEndpointMap map[Service]endpoint.Endpoint

type Handler func(in []byte, serviceEndpointMap ServiceEndpointMap) (err error)

type StartStopInterface interface {
	Start() (err error)
	Stop() (err error)
}

type ServerInterface interface {
	StartStopInterface
	RegisterNamespace(namespace string)
	RegisterService(namespace string, service Service, endpoint endpoint.Endpoint)
}
