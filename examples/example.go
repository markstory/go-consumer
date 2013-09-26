package main

import (
	"consumer"
	"flag"
	"log"
	"time"
)

func main() {
	var ini = flag.String("config", "./test.ini", "The configuration file to use.")

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
// Should eventually get type coming out of consumer
// instead of an amqp.Delivery value.
func worker (msg *consumer.Message) {
	log.Printf("msg received %s", msg.Body)
	msg.Ack(true)
	time.Sleep(2000 * time.Millisecond)
	log.Print("woke up")
}
