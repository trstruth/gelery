package main

import (
	"fmt"
	"log"

	"github.com/trstruth/gelery"
)

func main() {
	brokerInfo := &gelery.CeleryBrokerInfo{
		Type: "redis",
		Host: "localhost",
		Port: "5544",
	}

	cc, err := gelery.NewCeleryClient(
		gelery.WithBroker(brokerInfo),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = cc.SendTask("net_task.version", []interface{}{}, map[string]interface{}{}, "heartbeat")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully sent a task to celery")
}
