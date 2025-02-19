// ConsumerTest
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
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

var brokers string
var topic string
var cert string

// Reference : https://pkg.go.dev/github.com/shopify/sarama
func main() {
	var nTntIdx int32 = 0 // Partition Index Set

	config := sarama.NewConfig()
	// config.Version = version

	if cert != "" {
		tlsConfig, err := createTlsConfiguration(cert)
		if err != nil {
			panic(err)
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	// default 256
	config.ChannelBufferSize = 1000000
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}

	brokerList := strings.Split(brokers, ",")
	client, err := sarama.NewClient(brokerList, config)

	if err != nil {
		panic(err)
	}

	// if you want completed message -> 0
	lastOffset, err := client.GetOffset(topic, nTntIdx, sarama.OffsetNewest)

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
	partitionConsumer, err := consumer.ConsumePartition(topic, nTntIdx, lastOffset)
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

func createTlsConfiguration(path string) (*tls.Config, error) {
	paths := strings.Split(path, "\\")
	path = filepath.Join(paths...)

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := x509.NewCertPool()
	p.AppendCertsFromPEM(f)

	return &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    p,
	}, nil
}
