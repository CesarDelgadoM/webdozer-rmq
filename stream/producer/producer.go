package producer

import (
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CesarDelgadoM/webdozer-rmq/database"
	"github.com/CesarDelgadoM/webdozer-rmq/stream"
)

type Producer struct {
	id        int
	procs     int
	keys      chan string
	keysQueue chan chan string
	quit      chan bool
	server    *stream.ServerMQ
}

func NewProducer(id int, keysQueue chan chan string, server *stream.ServerMQ) Producer {

	return Producer{
		id:        id,
		procs:     3,
		keys:      make(chan string),
		keysQueue: keysQueue,
		quit:      make(chan bool),
		server:    server,
	}
}

func (p *Producer) Start() {

	go func() {
		var m sync.Mutex
		for {
			p.keysQueue <- p.keys

			select {
			case topic := <-p.keys:

				topicSize := p.server.List.Length(topic)

				for i := 0; i < p.procs; i++ {

					go func(idgo int) {
						ch := p.server.Rmq.ExchangeDeclare(topic)
						defer ch.Close()

						queue := topic + "_queue_" + strconv.Itoa(idgo)
						bind := topic + "_bind_" + strconv.Itoa(idgo)

						p.server.Rmq.QueueBindDeclare(queue, topic, bind, ch)

						for {
							m.Lock()
							if topicSize == 0 {
								if p.server.Set.IsMember(stream.PRODUCER_CACHE, topic) {
									log.Printf("[producer-%d - goroutine-%d] - finished the process to topic: %s\n", p.id, idgo, topic)
									p.server.Set.Remove(stream.PRODUCER_CACHE, topic)
								}
								m.Unlock()
								break
							}

							msg := p.server.List.Pop(topic)
							if msg != "" {
								log.Printf("[producer-%d - gorutine-%d] - publish message to topic: %s\n", p.id, idgo, topic)
								p.server.Rmq.Publish(msg, topic, bind, ch)
							}
							atomic.AddInt64(&topicSize, -1)
							m.Unlock()

							time.Sleep(1 * time.Millisecond)
						}
					}(i + 1)
				}
			case <-p.quit:
				return
			}
		}
	}()
}

func StartProducers(server *stream.ServerMQ, redis *database.RedisPool, nproducers int) {

	keysQueue := make(chan chan string, nproducers)

	for i := 0; i < nproducers; i++ {
		p := NewProducer(i+1, keysQueue, server)
		p.Start()
	}

	go func() {

		for {
			rkeys := redis.GetKeys(stream.PATTERN_KEYS)
			if len(rkeys) != 0 {
				for _, key := range rkeys {
					if !server.Set.IsMember(stream.PRODUCER_CACHE, key) {
						log.Println("<-", key)
						server.Set.Add(stream.PRODUCER_CACHE, key)
						server.Set.Add(stream.PROCESSED_KEYS, key)
						keys := <-keysQueue
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
