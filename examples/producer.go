package main

import (
	"go-micro-framework/microservice/configuration"
	"go-micro-framework/microservice/server"
	"log"
)

func main() {
	kafkacfg := configuration.MyKafkaConfig()
	err := kafkacfg.LoadConfig("kafka.config")
	if err != nil {
		log.Fatal("Failed to LoadConfig - %#v\n", err)
		return
	}

	kp, err := server.MyKafkaProducer(kafkacfg)
	if err != nil {
		log.Fatal("Failed to MyKafkaProducer")
		return
	}

	err = kp.Start()
	if err != nil {
		log.Fatal("Failed to Start Kafka Producer")
		return
	}
	defer kp.Stop()

	msg := server.Message{
		Service:        server.Service("hello"),
		ServiceMessage: "Nitish Malhotra",
	}

	err = kp.Produce("greeter", msg)
	if err != nil {
		log.Fatal("Failed to Produce")
		return
	}
}
