// Use this as an test program to generate a message on the kafka topic.
package main

import (
	"fmt"
	"go-micro-framework/microservice/configuration"
	"go-micro-framework/microservice/server"
	"log"
)

func main() {
	rpcCfg := configuration.MyRPCConfig()
	err := rpcCfg.LoadConfig("rpc-microservice/rpc.json")
	if err != nil {
		log.Fatal("Failed to LoadConfig - %#v\n", err)
		return
	}

	rpc, err := server.MyRPCClient(rpcCfg)
	if err != nil {
		log.Fatal("Failed to MyRPCClient")
		return
	}

	msg := &server.Message{
		Service:        server.Service("hello"),
		ServiceMessage: "Nitish Malhotra",
	}

	resp, err := rpc.Send("greeter", msg)
	if err != nil {
		log.Fatal("Failed to Send - %#v", err)
		return
	}
	fmt.Printf("Response Data %#v\n", resp)
}
