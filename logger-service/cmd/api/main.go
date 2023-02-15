package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// Create a context in order to desconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// clsoe connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// Register the rpc server
	err = rpc.Register(new(RPCServer)) // Tell the application that we're gonna be accepting RPC requests
	go app.rpcListen()

	// listen for grpc connecitons
	go app.gRPCListen()

	// start server creating a goroutine
	// go app.serve()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic()
	}
}

// We're copying the content of this func and pasting it into main because we don't want it running in hte background, by this way it will run and stop.
func (app *Config) rpcListen() {
	log.Println("Starting RPC server on port ", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort)) // Listen in all ports
	if err != nil {
		log.Println("Error when trying to listen with rpc", err)
	}
	defer listen.Close() // avoid resource leak

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// Create conneciton options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// Connect
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting: ", err)
		return nil, err
	}

	return c, nil

}
