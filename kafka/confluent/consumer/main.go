// Example channel-based high-level Apache Kafka consumer
package main

/**
 * Copyright 2016 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var brokers string
var group string
var topic string

func main() {
	topics := []string{topic}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               brokers,
		"group.id":                        group,
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": false,
		// Enable generation of PartitionEOF when the
		// end of a partition is reached.
		"enable.partition.eof": true,
		"auto.offset.reset":    "earliest"})

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Consumer %v\n", c)

	err = c.SubscribeTopics(topics, nil)

	run := true

	consumed := 0

	for run == true {
		select {
		case sig := <-sig:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false

		case ev := <-c.Events():
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Println("assigned partitions")
				_, _ = fmt.Fprintf(os.Stderr, "%% %v\n", e)
				_ = c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Println("revoked partitions")
				_, _ = fmt.Fprintf(os.Stderr, "%% %v\n", e)
				_ = c.Unassign()
			case *kafka.Message:
				fmt.Printf("Topic %s Consumed message offset %d , Partition %d\n", *e.TopicPartition.Topic, e.TopicPartition.Offset, e.TopicPartition.Partition)
				consumed++
				fmt.Printf("Consumed: %d\n", consumed)
				fmt.Println(string(e.Key))
				fmt.Println(string(e.Value))
				time.Sleep(time.Second * 1)
			case kafka.PartitionEOF:
				// fmt.Printf("%% Reached %v\n", e)
			case kafka.Error:
				// Errors should generally be considered as informational, the client will try to automatically recover
				_, _ = fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			}
		}
	}

	fmt.Printf("Closing consumer\n")
	_ = c.Close()
}
