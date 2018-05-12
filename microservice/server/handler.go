package server

import (
	"context"
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func ByteHandler(in []byte, serviceEndpointMap ServiceEndpointMap) (err error) {
	msg := &Message{}
	err = json.Unmarshal(in, msg)
	if err != nil {
		return
	}

	if ep, ok := serviceEndpointMap[msg.Service]; !ok {
		err = fmt.Errorf("No handler found for SERVICE [%v]", msg.Service)
		return
	} else {
		resp, err := ep.Handle(context.Background(), []byte(msg.ServiceMessage))
		if err != nil {
			return err
		}
		log.Infof("Response ByteHandler %#v\n", resp)
	}
	return
}
