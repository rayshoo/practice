// Example channel-based Apache Kafka producer
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
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	const brokers = "10.254.1.51:29291,10.254.1.51:29292,10.254.1.51:29293"
	// topic := "call_topic"
	topic := "test"

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created Producer %v\n", p)

	doneChan := make(chan bool)

	go func() {
		defer close(doneChan)
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				m := ev
				if m.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
				} else {
					fmt.Printf("produce succeeded. topic: %s. partition: %d. offset : %d\n", *m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
				}
				// return
			default:
				fmt.Printf("Ignored event: %s\n", ev)
			}
		}
	}()

	value := "Hello Go!"

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: 1}, Value: []byte(value)}
			p.ProduceChannel() <- &kafka.Message{TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, Value: []byte(value)}
		}()
	}
	wg.Wait()

	// wait for delivery report goroutine to finish
	_ = <-doneChan

	p.Close()
}
