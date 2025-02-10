// ConsumerTest
package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM/sarama"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// Reference : https://pkg.go.dev/github.com/shopify/sarama
func main() {
	const brokers = "10.254.1.51:29291,10.254.1.51:29292,10.254.1.51:29293"
	// topics := []string{"call_topic"}
	topics := []string{"test"}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), "2", config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		fmt.Println("consumes start")
		handler := exampleConsumerGroupHandler{}
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
