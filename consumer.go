package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config-file-path> <target-url>\n",
			os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	targetUrl := os.Args[2]
	conf := ReadConfig(configFile)
	conf["group.id"] = "kafka-go-getting-started"
	conf["auto.offset.reset"] = "earliest"

	c, err := kafka.NewConsumer(&conf)

	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	topic := "metrics"
	err = c.SubscribeTopics([]string{topic}, nil)
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			req, err := http.NewRequest("POST", targetUrl, bytes.NewReader(ev.Value))
			if err != nil {
				// handle err
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-protobuf")
			req.Header.Set("User-Agent", "promkafka-0.0.2")
			req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				// handle err
				log.Fatal(err)
			}
			defer resp.Body.Close()
			fmt.Printf("Consumed event from topic %s.\n",
				*ev.TopicPartition.Topic)
		}
	}

	c.Close()
}
