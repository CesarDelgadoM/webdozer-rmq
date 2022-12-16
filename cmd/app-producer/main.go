package main

import (
	"runtime"

	"github.com/CesarDelgadoM/webdozer-rmq/database"
	"github.com/CesarDelgadoM/webdozer-rmq/stream"
	"github.com/CesarDelgadoM/webdozer-rmq/stream/producer"
)

func main() {
	signal := make(chan bool, 1)
	procs := runtime.GOMAXPROCS(runtime.NumCPU())
	nproducers := 3

	redis := database.NewRedisPool("localhost:6379", procs)
	defer redis.Close()

	rmq := stream.NewRabbitClient("amqp://guest:guest@localhost:5672/")

	server := stream.NewSeverMQ(rmq, redis)

	producer.StartProducers(server, redis, nproducers)

	<-signal
}
