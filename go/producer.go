package c3alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func producer(m Message,producer *kafka.Producer) {

	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "cloud-glb-peve-monitor-alert"

	alert := Alert{}
	alert.Attributes.class = "dynatrace"
	alert.Attributes.McObjectURI = "dynatrace"
	alert.Attributes.Severity = "CRITICAL"
	alert.Attributes.Msg = m.Text
	alert.Attributes.McHost = "InternalTestAM5"
	alert.Attributes.McSmcAlias = "Prod application SRV1"
	alert.Attributes.McSmcID = "Model1_10000_S0109"
	alert.Attributes.McOwner = "Administrator"
	alert.Attributes.McPriority = "PRIORITY_4"

	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(alert)
	reqBodyBytes.Bytes()

	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value: reqBodyBytes.Bytes(),
	}, nil)


	// Wait for message deliveries before shutting down
	producer.Flush(15 * 1000)
}