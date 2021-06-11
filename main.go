package main

import (
	sw "c3alert/go"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var (
	kafkaBrokerList        = "kafka:9092"
	kafkaTopic             = "metrics"
	topicTemplate          *template.Template
	kafkaCompression       = "none"
	kafkaBatchNumMessages  = "10000"
	kafkaSslClientCertFile = ""
	kafkaSslClientKeyFile  = ""
	kafkaSslClientKeyPass  = ""
	kafkaSslCACertFile     = ""
	kafkaSecurityProtocol  = ""
)

func main() {
	log.Printf("creating kafka producer")

	kafkaConfig := kafka.ConfigMap{
		"bootstrap.servers":   kafkaBrokerList,
		"compression.codec":   kafkaCompression,
		"batch.num.messages":  kafkaBatchNumMessages,
		"go.batch.producer":   true,  // Enable batch producer (for increased performance).
		"go.delivery.reports": false, // per-message delivery reports to the Events() channel
	}

	if kafkaSslClientCertFile != "" && kafkaSslClientKeyFile != "" && kafkaSslCACertFile != "" {
		if kafkaSecurityProtocol == "" {
			kafkaSecurityProtocol = "ssl"
		}

		if kafkaSecurityProtocol != "ssl" && kafkaSecurityProtocol != "sasl_ssl" {
			log.Printf("invalid config: kafka security protocol is not ssl based but ssl config is provided")
		}

		kafkaConfig["security.protocol"] = kafkaSecurityProtocol
		kafkaConfig["ssl.ca.location"] = kafkaSslCACertFile              // CA certificate file for verifying the broker's certificate.
		kafkaConfig["ssl.certificate.location"] = kafkaSslClientCertFile // Client's certificate
		kafkaConfig["ssl.key.location"] = kafkaSslClientKeyFile          // Client's key
		kafkaConfig["ssl.key.password"] = kafkaSslClientKeyPass          // Key password, if any.
	}

	producer, err := kafka.NewProducer(&kafkaConfig)
	if err != nil {
		log.Printf("couldn't create kafka producer")
	}

	log.Printf("Server started")
	router := sw.NewRouter(producer)
	if router != nil {
		log.Printf("%s", "Gollia HTTP Router Started")
	}

	log.Fatal(http.ListenAndServe(":8090", router))
}

func init() {
	log.Printf("Init Service")
	if value := os.Getenv("KAFKA_BROKER_LIST"); value != "" {
		kafkaBrokerList = value
	}

	if value := os.Getenv("KAFKA_TOPIC"); value != "" {
		kafkaTopic = value
	}

	if value := os.Getenv("KAFKA_COMPRESSION"); value != "" {
		kafkaCompression = value
	}

	if value := os.Getenv("KAFKA_BATCH_NUM_MESSAGES"); value != "" {
		kafkaBatchNumMessages = value
	}

	if value := os.Getenv("KAFKA_SSL_CLIENT_CERT_FILE"); value != "" {
		kafkaSslClientCertFile = value
	}

	if value := os.Getenv("KAFKA_SSL_CLIENT_KEY_FILE"); value != "" {
		kafkaSslClientKeyFile = value
	}

	if value := os.Getenv("KAFKA_SSL_CLIENT_KEY_PASS"); value != "" {
		kafkaSslClientKeyPass = value
	}

	if value := os.Getenv("KAFKA_SSL_CA_CERT_FILE"); value != "" {
		kafkaSslCACertFile = value
	}

	if value := os.Getenv("KAFKA_SECURITY_PROTOCOL"); value != "" {
		kafkaSecurityProtocol = strings.ToLower(value)
	}

	var err error
	topicTemplate, err = parseTopicTemplate(kafkaTopic)
	if err != nil {
		log.Printf("%s", "couldn't parse the topic template")
	}
}

func parseTopicTemplate(tpl string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"replace": func(old, new, src string) string {
			return strings.Replace(src, old, new, -1)
		},
		"substring": func(start, end int, s string) string {
			if start < 0 {
				start = 0
			}
			if end < 0 || end > len(s) {
				end = len(s)
			}
			if start >= end {
				panic("template function - substring: start is bigger (or equal) than end. That will produce an empty string.")
			}
			return s[start:end]
		},
	}
	return template.New("topic").Funcs(funcMap).Parse(tpl)
}
