package server

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

// Message defines the message structure to be used by the microservice framework.
// It provides a mechanism to demux messages coming in a namespace by using its Service
// identifier.
type Message struct {
	// Service is the service identifier. Set this to a unique string for a given server
	Service Service `json:"service,omitempty"`
	// ServiceMessage is the actual payload passed over the network.
	ServiceMessage string `json:"service_message,omitempty"`
}

// ByteHandler is an implementation of the ServerHandler.
// It takes in a slice of bytes, coming in from the network
// and demultiplexes the message using a map provided to it with all the
// service endpoints.
func ByteHandler(in []byte, serviceEndpointMap ServiceEndpointMap) (err error) {
	msg := &Message{}
	err = json.Unmarshal(in, msg)
	if err != nil {
		return
	}

	ep, ok := serviceEndpointMap[msg.Service]
	if !ok {
		err = fmt.Errorf("No handler found for SERVICE [%v]", msg.Service)
		return
	}

	resp, err := ep.Handle(context.Background(), []byte(msg.ServiceMessage))
	if err != nil {
		return err
	}
	log.Infof("Response ByteHandler %#v\n", resp)

	return
}
