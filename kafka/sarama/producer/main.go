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

// Reference : https://pkg.go.dev/github.com/shopify/sarama

var brokers string
var topic string
var cert string

var asyncProducer sarama.AsyncProducer
var syncProducer sarama.SyncProducer

func init() {
	brokerList := strings.Split(brokers, ",")
	asyncProducer = AsyncWriter(brokerList)
	syncProducer = SyncWriter(brokerList)
}

func SyncWriter(brokerList []string) sarama.SyncProducer {
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()

	if cert != "" {
		tlsConfig, err := createTlsConfiguration(cert)
		if err != nil {
			panic(err)
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	// config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message

	// config.Producer.Compression = sarama.CompressionNone
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages

	config.Producer.Retry.Max = 10 // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		fmt.Println("Failed to start Sarama producer:", err)
	}

	return producer
}

func AsyncWriter(brokerList []string) sarama.AsyncProducer {
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()

	if cert != "" {
		tlsConfig, err := createTlsConfiguration(cert)
		if err != nil {
			panic(err)
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}

	// config.Producer.Partitioner = sarama.NewManualPartitioner
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForLocal // Only wait for the leader to ack

	// config.Producer.Compression = sarama.CompressionNone
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages

	// config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		fmt.Println("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			fmt.Println("Failed to write access log entry:", err)
		}
	}()

	go func() {
		for msg := range producer.Successes() {
			fmt.Printf("async produce succeeded. topic: %s, partition: %d. offset : %d\n", msg.Topic, msg.Partition, msg.Offset)
		}
	}()

	return producer
}

func AsyncProducer() {
	asyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		// Key:   sarama.StringEncoder("CallId"),
		Value: sarama.StringEncoder("Hello Go!"),
		// Partition: int32(0),
	}
}

func SyncProducer() {
	// returned partition, offset, err
	partition, offset, err := syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		// Key:   sarama.StringEncoder("CallId"),
		Value: sarama.StringEncoder("Hello Go!"),
		// Partition: int32(0),
		// Offset
		// Metadata
		// Timestamp
		// Headers
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("sync produce succeeded. topic: %s, partition: %d. offset : %d\n", topic, partition, offset)
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

func main() {
	// AsyncProducer()
	// SyncProducer()
	SyncProducer()
}
