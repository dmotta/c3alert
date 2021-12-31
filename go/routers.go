package c3alert

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(producer *kafka.Producer) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.Methods("GET").Path("/").Name("Home").HandlerFunc(home)
	r.Methods("POST").Path("/services/{workspace}/{channel}/{channel_id}").Name("SlackHook").HandlerFunc(serviceHandler(producer))
	return r
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "home\n")
}

