package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/linkedin/goavro/v2"
	"gocv.io/x/gocv"
)

// FIXME: Move the struct to separate package so that all the modules can use it
type FrameData struct {
	ID      int    `json:"id"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	MatType int    `json:"matType"`
	Data    []byte `json:"data"`
}

var (
	shutDown      bool = false
	avroSchemaDef      = `
		{
			"type" : "record",
			"namespace" : "VAMP",
			"name" : "framedata",
			"fields" : [
				{ "name" : "id" , "type" : "int" },
				{ "name" : "width" , "type" : "int" },
				{ "name" : "height" , "type" : "int" },
				{ "name" : "matType" , "type" : "int" },
				{ "name" : "data" , "type" : "bytes" }
			]
		}
	`
)

// func producer(client pulsar.Client) {
// 	producer, err := client.CreateProducer(pulsar.ProducerOptions{
// 		Topic: "my-topic",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
// 		Payload: []byte("hello"),
// 	})

// 	defer producer.Close()

// 	if err != nil {
// 		log.Println("Failed to publish message, ", err)
// 	}
// 	log.Println("Published message")
// }

func consumer(client pulsar.Client) {
	properties := make(map[string]string)
	codec, err := goavro.NewCodec(avroSchemaDef)
	if err != nil {
		return
	}
	schema := pulsar.NewAvroSchema(codec.CanonicalSchema(), properties)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond) // Timeout while receiving the message
	defer cancel()
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       "input-10",
		SubscriptionName:            "my-sub",
		Type:                        pulsar.Shared,
		Schema:                      schema,
		SubscriptionInitialPosition: pulsar.SubscriptionPositionEarliest,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer consumer.Close()

	var s FrameData

	for {
		if shutDown {
			log.Println("Stopping consumer")
			break
		}
		msg, err := consumer.Receive(ctx)
		if err != nil {
			continue
			// log.Fatal(err)
		}
		err = msg.GetSchemaValue(&s)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Received message msgId: %#v\n", msg.ID())
		img, err := gocv.NewMatFromBytes(s.Height, s.Width, gocv.MatType(s.MatType), s.Data)
		if err != nil {
			log.Printf("Error reading image from bytes, %v", err)
		}
		log.Printf("img.Cols(): %v img.Rows(): %v img.Type(): %v", img.Cols(), img.Rows(), img.Type())
		consumer.Ack(msg)
	}
	if err := consumer.Unsubscribe(); err != nil {
		log.Fatal(err)
	}
}

func gracfulExit() {
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	fmt.Println("Setting up graceful exit")

	go func() {
		sig := <-gracefulStop
		shutDown = true
		log.Printf("caught sig: %+v", sig)
		log.Println("Gracefully exiting")
	}()
}

func main() {
	go gracfulExit()
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate pulsar client: %v", err)
	}

	defer client.Close()

	// producer(client)
	consumer(client)
}
