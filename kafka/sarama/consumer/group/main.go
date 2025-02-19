// ConsumerTest
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"
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

var brokers string
var group string
var topic string
var cert string

// Reference : https://pkg.go.dev/github.com/shopify/sarama
func main() {
	topics := []string{topic}

	config := sarama.NewConfig()

	if cert != "" {
		tlsConfig, err := createTlsConfiguration(cert)
		if err != nil {
			panic(err)
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
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
