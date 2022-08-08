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

// type testJSON struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

var (
	shutDown bool = false
	// exampleSchemaDef      = "{\"type\":\"record\",\"name\":\"Example\",\"namespace\":\"test\"," +
	// 	"\"fields\":[{\"name\":\"ID\",\"type\":\"int\"},{\"name\":\"Width\",\"type\":\"int\"},{\"name\":\"Height\",\"type\":\"int\"}]}"
	// exampleSchemaDef = "{\"type\":\"record\",\"name\":\"Example\",\"namespace\":\"test\"," +
	// 	"\"fields\":[{\"name\":\"ID\",\"type\":\"int\"},{\"name\":\"Name\",\"type\":\"string\"}]}"
	avroSchemaDef = `
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

func produceFrames(client pulsar.Client) {
	fmt.Println("Started producer")
	properties := make(map[string]string)
	// schema := pulsar.NewJSONSchema(exampleSchemaDef, properties)
	codec, err := goavro.NewCodec(avroSchemaDef)
	if err != nil {
		log.Printf("Error creating new codec, %v\n", err)
		return
	}
	schema := pulsar.NewAvroSchema(codec.CanonicalSchema(), properties)
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic:           "input-10",
		Schema:          schema,
		BatchingMaxSize: 5000000,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()
	frameCount := 0
	// TODO: Read frame from webcam and push to stream
	webcam, _ := gocv.VideoCaptureDevice(0)
	defer webcam.Close()
	img := gocv.NewMat()
	for {
		frameCount++
		if shutDown {
			log.Println("Shutting down producer")
			break
		}
		webcam.Read(&img)
		// log.Printf("img.Cols(): %v img.Rows(): %v img.Type(): %v", img.Cols(), img.Rows(), img.Type())
		// gocv.Resize(img, &img, image.Point{1280, 720}, 0, 0, gocv.InterpolationDefault)
		// gocv.IMWrite("tmp.jpg", img)
		log.Printf("img.ToBytes(): %v", len(img.ToBytes()))
		log.Printf("img.Cols(): %v img.Rows(): %v img.Type(): %v", img.Cols(), img.Rows(), img.Type())
		log.Println("Before sending")
		_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
			Value: &FrameData{
				ID:      frameCount,
				Width:   img.Cols(),
				Height:  img.Rows(),
				MatType: int(img.Type()),
				Data:    img.ToBytes(),
			},
		})
		if err != nil {
			log.Println("Failed to publish message, ", err)
		}
		log.Printf("Processed %d frames\n", frameCount)
		// time.Sleep(80 * time.Millisecond)
		time.Sleep(1 * time.Second)
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
