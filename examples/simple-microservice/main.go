package main

import (
	"flag"
	"fmt"
	"go-micro-framework/examples/simple-microservice/server"
	"go-micro-framework/examples/simple-microservice/service"
	"os"

	log "github.com/sirupsen/logrus"
)

var debug = flag.String("debug", "debug", "Debug levels : panic,fatal,error,warn,info,debug")

func blockForever() {
	c := make(chan struct{})
	<-c
}

func init() {
	//--- Setup the logrus logger ---
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func setLogLevel() {
	flag.Parse()
	switch *debug {
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	default:
		fmt.Printf("Level %s not supported.", *debug)
	}
	fmt.Printf("=============================\n")
	fmt.Printf("CURRENT LOG LEVEL : %s\n", log.GetLevel().String())
	fmt.Printf("=============================\n")
}

func main() {
	setLogLevel()

	// Instantiate the ServiceMgr
	ServiceMgr, err := service.MyServiceMgr()
	if err != nil {
		log.Fatalf("Failed MyServiceMgr - %v", err.Error())
		return
	}

	// Register your server with the microservice
	err = server.Register(ServiceMgr)
	if err != nil {
		log.Fatalf("Failed server Registration - %v", err.Error())
		return
	}

	// Start all your servers registered with the microservice
	err = ServiceMgr.Start()
	if err != nil {
		log.Fatalf("Failed Start - %v", err.Error())
		return
	}

	blockForever()
}
