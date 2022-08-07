package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func producer(client pulsar.Client) {
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "my-topic",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: []byte("hello"),
	})

	defer producer.Close()

	if err != nil {
		log.Println("Failed to publish message, ", err)
	}
	log.Println("Published message")
}

func consumer(client pulsar.Client) {
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "my-topic",
		SubscriptionName: "my-sub",
		Type:             pulsar.Shared,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	msg, err := consumer.Receive(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received message msgId: %#v -- content: '%s'\n",
		msg.ID(), string(msg.Payload()))
    consumer.Ack(msg)
    if err := consumer.Unsubscribe(); err != nil {
        log.Fatal(err)
    }
}

func main() {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:6650",
		OperationTimeout:  30 * time.Minute,
		ConnectionTimeout: 30 * time.Minute,
	})
	if err != nil {
		log.Fatalf("Could not instantiate pulsar client: %v", err)
	}

	defer client.Close()

	// producer(client)
	consumer(client)
}
