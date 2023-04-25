package messageservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/apache/pulsar-client-go/pulsar"
	"log"
	"time"
)

func Producer(pulsarUrl, pTopic string) {
	fmt.Println("Pulsar Producer")
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               pulsarUrl,
		OperationTimeout:  30 * time.Second,
		ConnectionTimeout: 30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Could not instantiate Pulsar client: %v", err)
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: pTopic,
	})

	if err != nil {
		log.Fatal(err)
	}

	//_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
	//	Key: "Hello,This is a message from Pulsar Producer!\n",
	//})
	//msg := pulsar.ProducerMessage{
	//	Payload: []byte("Hello,This is a message from Pulsar Producer!\n"),
	//}
	aa := map[string]string{"pan": "zhong"}
	b, _ := json.Marshal(aa)
	msg := pulsar.ProducerMessage{
		Payload: []byte(string(b)),
		//Key: "Hello,This is a message from Pulsar Producer!\n",
	}
	if err, _ := producer.Send(context.Background(), &msg); err != nil {
		log.Fatalf("Producer could not send message:%v", err)
	}

	defer producer.Close()

	if err != nil {
		fmt.Println("Failed to publish message", err)
	}
	fmt.Println("Published message")
}
