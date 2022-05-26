package queue

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/adjust/rmq/v4"
)

const (
	prefetchLimit = 10
	pollDuration  = 1000 * time.Millisecond
	numConsumers  = 5
)

var (
	queue *Queue
	once  sync.Once
)

type Queue struct {
	conn  rmq.Connection
	queue rmq.Queue
}

func getConnection(role string) (rmq.Connection, error) {
	errChan := make(chan error, 10)
	go logErrors(errChan)

	return rmq.OpenConnection(role, "tcp", "localhost:6379", 2, errChan) // TODO: use env
}

func initQueue(role, queueName string) {
	conn, err := getConnection(role)
	if err != nil {
		log.Panic(err)
	}

	q, err := conn.OpenQueue(queueName)
	if err != nil {
		log.Panic(err)
	}

	queue = &Queue{
		conn:  conn,
		queue: q,
	}
}

func New(role, queueName string) *Queue {
	once.Do(func() {
		initQueue(role, queueName)
	})

	return queue
}

func (q *Queue) StartConsuming(consumer func(string) error) {
	err := q.queue.StartConsuming(prefetchLimit, pollDuration)
	if err != nil {
		log.Panic(err)
	}
	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("consumer %d", i)
		if _, err := q.queue.AddConsumer(name, NewConsumer(i, consumer)); err != nil {
			log.Panic(err)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals // wait for signal
	go func() {
		<-signals // hard exit on second signal (in case shutdown gets stuck)
		os.Exit(1)
	}()

	<-q.conn.StopAllConsuming() // wait for all Consume() calls to finish
}

type Consumer struct {
	name     string
	consumer func(string) error
}

func NewConsumer(tag int, consumer func(string) error) *Consumer {
	return &Consumer{
		name:     fmt.Sprintf("consumer%d", tag),
		consumer: consumer,
	}
}

func (consumer *Consumer) Consume(delivery rmq.Delivery) {
	payload := delivery.Payload()
	err := consumer.consumer(payload)
	if err != nil {
		if err := delivery.Reject(); err != nil {
			log.Printf("failed to reject %s: %s", payload, err)
		} else {
			log.Printf("rejected %s", payload)
		}
	} else {
		if err := delivery.Ack(); err != nil {
			log.Printf("failed to ack %s: %s", payload, err)
		} else {
			log.Printf("acked %s", payload)
		}
	}

}

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *rmq.ConsumeError:
			log.Print("consume error: ", err)
		case *rmq.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}
