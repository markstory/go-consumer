package main

import (
	"consumer"
	"flag"
	"log"
	"time"
)

func main() {
	var ini = flag.String("config", "./simple.ini", "The configuration file to use.")

	c, err := consumer.Create(*ini)
	if err != nil {
		log.Fatalf("Unable to create consumer. Error: %v", err)
	}
	err = c.Consume(worker)
	if err != nil {
		log.Fatalf("Unable to consume messages. Error: %v", err)
	}
}

// Example worker function
// Just prints to the logs and sleeps between receiving messages.
func worker(msg *consumer.Message) {
	log.Printf(
		"Message received Exchange: %s, RoutingKey: %s, Body: %s",
		msg.Exchange,
		msg.RoutingKey,
		msg.Body)

	msg.Ack(true)
	time.Sleep(1000 * time.Millisecond)
	log.Print("sleep complete")
}
