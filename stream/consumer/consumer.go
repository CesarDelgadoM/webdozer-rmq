package consumer

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/CesarDelgadoM/webdozer-rmq/stream"
)

type Consumer struct {
	id         int
	procs      int
	pathFiles  string
	key        chan string
	keysQueues chan chan string
	quit       chan bool
	server     *stream.ServerMQ
}

func NewConsumer(id int, pathFiles string, keysQueues chan chan string, server *stream.ServerMQ) Consumer {

	return Consumer{
		id:         id,
		procs:      3,
		pathFiles:  pathFiles,
		key:        make(chan string),
		keysQueues: keysQueues,
		quit:       make(chan bool, 3),
		server:     server,
	}
}

func (c *Consumer) Consume() {

	go func() {
		var m sync.Mutex
		for {
			c.keysQueues <- c.key

			select {
			case topic := <-c.key:
				file, err := os.OpenFile(c.pathFiles+topic+".txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
				if err != nil {
					log.Println("Failed to create file:", err)
					continue
				}

				for i := 0; i < c.procs; i++ {

					go func(idgo int) {
						ch := c.server.Rmq.ExchangeDeclare(topic)
						defer ch.Close()

						searchWords := c.server.Hash.GetValues(topic)

						queue := topic + "_queue_" + strconv.Itoa(idgo)
						consumer := topic + "_consumer_" + strconv.Itoa(idgo)

						msgs := c.server.Rmq.ConsumeDeclare(queue, consumer, ch)
						for {
							select {
							case msg := <-msgs:
								log.Printf("[consumer-%d - goroutine-%d] - consume messages from topic: %s\n", c.id, idgo, topic)
								m.Lock()
								url := string(msg.Body)
								for _, word := range searchWords {
									if strings.Contains(url, word) {
										log.Printf("Writing message from topic %s\n", topic)
										file.WriteString(url + "\n")
									}
								}
								m.Unlock()
							case <-time.After(30 * time.Second):
								if c.server.Set.IsMember(stream.CONSUMER_CACHE, topic) {
									log.Printf("[consumer-%d - goroutine-%d] - finished the process to topic: %s\n", c.id, idgo, topic)
									c.server.Set.Remove(stream.CONSUMER_CACHE, topic)
									c.server.Set.Remove(stream.PROCESSED_KEYS, topic)
									file.Close()
								}
								return
							}
							time.Sleep(1 * time.Millisecond)
						}
					}(i + 1)
				}
			case <-c.quit:
				return
			}
		}
	}()
}

func StartConsumers(server *stream.ServerMQ, nconsumers int) {

	queueKeys := make(chan chan string, nconsumers)

	for i := 0; i < nconsumers; i++ {
		c := NewConsumer(i+1, "../../urls_files/", queueKeys, server)
		c.Consume()
	}

	go func() {
		for {
			keysr := server.Set.GetValues(stream.PROCESSED_KEYS)
			if len(keysr) != 0 {
				for _, key := range keysr {
					if !server.Set.IsMember(stream.CONSUMER_CACHE, key) {
						log.Println("<-", key)
						server.Set.Add(stream.CONSUMER_CACHE, key)
						keys := <-queueKeys
						keys <- key
					}
				}
			} else {
				log.Println("No keys to process...")
			}
			time.Sleep(30 * time.Second)
		}
	}()
}
