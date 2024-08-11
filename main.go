package main

import (
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {
	cassandraAddr := os.Getenv("CAS_ADDRESS")
	cassandraKeyspace := os.Getenv("CAS_KEYSPACE")
	appPort := os.Getenv("APP_PORT")

	//connect to Cassandra
	log.Println("Create Cassandra session")
	cluster := gocql.NewCluster(cassandraAddr)
	cluster.Keyspace = cassandraKeyspace
	cassandraSession, err := cluster.CreateSession()
	failOnError(err, "Error creating Cassandra session")
	defer cassandraSession.Close()

	// create Server
	server := NewServer(cassandraSession)
	router := mux.NewRouter()
	router.Use(server.GetXRequestIdMiddleware)
	// define authenticated route
	privateRouter := router.PathPrefix("/").Subrouter()
	privateRouter.Use(server.GetAuthMiddleware)
	privateRouter.HandleFunc("/dialog/{user_id}/send", server.GetSendMessageHandler).Methods("POST")
	privateRouter.HandleFunc("/dialog/{user_id}/list", server.GetDialogListHandler).Methods("GET")

	// start server
	log.Println("Start listening server on port " + appPort)
	http.ListenAndServe("0.0.0.0:"+appPort, router)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
