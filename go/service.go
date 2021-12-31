package c3alert

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
)

func serviceHandler(p *kafka.Producer) func(w http.ResponseWriter, r *http.Request) {
	if p == nil {
		panic("nil Kafka producer!")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var m Message
		err := json.NewDecoder(r.Body).Decode(&m)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		producer(m,p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}