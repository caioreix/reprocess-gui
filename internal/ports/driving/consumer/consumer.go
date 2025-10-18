package consumer

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Consumer struct {
}

func (sc *Consumer) NewConsumer() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "",
		"group.id":          "",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return
	}

	eventCount := 0
	messages := make(map[string][]*kafka.Message) //Storage for messages to be batch processed. Our case requires a map for messages of keys.
	for {
		msg, err := sc.Instance.ReadMessage(time.Second)
		if err != nil {
			if k, ok := err.(kafka.Error); ok && k.Code() == kafka.ErrTimedOut {
				//In case of a timeout, do not wait reaching the BATCH_SIZE. Process stored messages.
				if len(messages) > 0 {
					sc.process(&messages, &eventCount)
					messages = make(map[string][]*kafka.Message)
				}
			}
			continue
		}
		put(&messages, msg) //Store messages to batch process later.
		eventCount++

		if eventCount%BATCH_SIZE == 0 {
			sc.process(&messages, &eventCount)
			messages = make(map[string][]*kafka.Message)
		}
	}
}

func (sc *Consumer) process(messages *map[string][]*kafka.Message, eventCount *int) {
	sc.processMessages(messages) //Async processing
	*eventCount = 0
	sc.commit()
}
