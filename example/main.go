package main

import (
	"fmt"
	"log"
	"time"

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

	taskID, err := cc.SendTask("net_task.version", []interface{}{}, map[string]interface{}{}, "heartbeat")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully sent a task to celery")

	time.Sleep(3 * time.Second)

	r, err := cc.GetResult(taskID, "heartbeat")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("got result for taskID %s: %v", taskID, r)
}
