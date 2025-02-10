// ConsumerTest
package main

import (
	"fmt"
	"strings"

	"github.com/IBM/sarama"
)

/*
//Consumer Group Handler

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ ConsumerGroupSession) error { return nil }
func (h exampleConsumerGroupHandler) ConsumeClaim(sess ConsumerGroupSession, claim ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
        sess.MarkMessage(msg, "")
    }
    return nil
}
*/

// Reference : https://pkg.go.dev/github.com/shopify/sarama
func main() {
	var nTntIdx int32 = 0 // Partition Index Set

	config := sarama.NewConfig()
	// config.Version = version

	// default 256
	config.ChannelBufferSize = 1000000
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	const brokers = "10.254.1.51:29291,10.254.1.51:29292,10.254.1.51:29293"
	brokerList := strings.Split(brokers, ",")
	client, err := sarama.NewClient(brokerList, config)

	if err != nil {
		panic(err)
	}

	// if you want completed message -> 0
	lastOffset, err := client.GetOffset("call_topic", nTntIdx, sarama.OffsetNewest)

	if err != nil {
		panic(err)
	}

	// if consumer group isn't exist , create it
	// group, err := sarama.NewConsumerGroupFromClient("consumer-group-name", client)
	consumer, err := sarama.NewConsumerFromClient(client)

	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Read Specific Partition From Topic
	partitionConsumer, err := consumer.ConsumePartition("call_topic", nTntIdx, lastOffset)
	// ctx := context.Background()
	// handler := exampleConsumerGroupHandler{}
	// err := group.Consume(ctx, []string{"call_topic"}, handler)

	if err != nil {
		panic(err)
	}

	// if ConsumerGroup Skip Underlines
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	consumed := 0

	fmt.Println("consumes start")
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			fmt.Printf("Topic %s Consumed message offset %d , Partition %d\n", msg.Topic, msg.Offset, msg.Partition)
			consumed++
			fmt.Printf("Consumed: %d\n", consumed)
			fmt.Println(string(msg.Key))
			fmt.Println(string(msg.Value))
			fmt.Println("")
		}
	}
}
