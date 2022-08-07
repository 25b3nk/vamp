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
)

var shutDown bool = false

func produceFrames(client pulsar.Client) {
	fmt.Println("Started producer")
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "input-source",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
    frameCount := 0
	// TODO: Read frame from webcam and push to stream
	// webcam, _ := gocv.VideoCaptureDevice(0)
	// img := gocv.NewMat()
	for {
        frameCount++
		if shutDown {
			break
		}
		// webcam.Read(&img)
        // TODO: Use avro for serdes
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: []byte("hello"),
		})
        log.Printf("Processed %d frames\n", frameCount)
		time.Sleep(80 * time.Millisecond)
		if err != nil {
			log.Println("Failed to publish message, ", err)
		}
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
	fmt.Println("Setting up pulsar client")
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               "pulsar://localhost:6650",
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate pulsar client: %v", err)
	}

	defer client.Close()
	produceFrames(client)
}
