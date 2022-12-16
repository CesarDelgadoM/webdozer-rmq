package stream

import (
	"context"
	"log"
	"time"

	"github.com/CesarDelgadoM/webdozer-rmq/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TOPIC_TYPE     string = "topic"
	PATTERN_KEYS   string = "urls-*"
	SEARCH_WORDS   string = "search-words"
	PRODUCER_CACHE string = "producer-cache"
	CONSUMER_CACHE string = "consumer-cache"
	PROCESSED_KEYS string = "processed-keys"
)

type ServerMQ struct {
	List database.IList
	Hash database.IHash
	Set  database.ISet
	Rmq  *RabbitClient
}

func NewSeverMQ(rmq *RabbitClient, redis *database.RedisPool) *ServerMQ {

	return &ServerMQ{
		List: database.NewList(redis),
		Set:  database.NewSet(redis),
		Hash: database.NewHash(SEARCH_WORDS, redis),
		Rmq:  rmq,
	}
}

type RabbitClient struct {
	conn *amqp.Connection
}

func NewRabbitClient(url string) *RabbitClient {

	log.Println("url-rabbitmq:", url)

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed connect to broker rabbitmq: ", err)
	}

	log.Println("Connect to rabbitmq")

	return &RabbitClient{
		conn: conn,
	}
}

func (rc *RabbitClient) ExchangeDeclare(name string) *amqp.Channel {

	ch, err := rc.conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
	}

	err = ch.ExchangeDeclare(name, TOPIC_TYPE, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare an exchange")
	}

	return ch
}

func (rc *RabbitClient) QueueBindDeclare(name, topic, key string, channel *amqp.Channel) amqp.Queue {

	q, err := channel.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare an queue:", err)
	}

	err = channel.QueueBind(q.Name, key, topic, false, nil)
	if err != nil {
		log.Fatal("Failed to create bind")
	}

	return q
}

func (rc *RabbitClient) Publish(msg, topic, key string, channel *amqp.Channel) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx, topic, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		log.Println("Failed to publish a message:", err)
	}
}

func (rc *RabbitClient) ConsumeDeclare(queue, consumer string, channel *amqp.Channel) <-chan amqp.Delivery {

	messages, err := channel.Consume(queue, consumer, true, false, false, false, nil)
	if err != nil {
		log.Println("Failed to declare consume:", err)
	}

	return messages
}
