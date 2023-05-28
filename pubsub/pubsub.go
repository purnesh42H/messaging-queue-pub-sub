package pubsub

import (
	"errors"
	"fmt"
	"sync"
)

type PubSub interface {
	Subscribe(topic string) chan string
	Publish(message Message) error
	ResetOffset(subscriber chan string, offset int) error
	GetSubscriberMessages(subscriber chan string) ([]string, error)
}

type concurrentPubsub struct {
	subscribers     map[string][]chan string
	messages        map[chan string][]string
	mutex           sync.RWMutex
	subscriberMutex sync.Mutex
}

func NewConcurrentPubSub() PubSub {
	return &concurrentPubsub{
		subscribers: make(map[string][]chan string),
		messages:    make(map[chan string][]string),
	}
}

func (pbs *concurrentPubsub) Subscribe(topic string) chan string {
	pbs.mutex.RLock()
	defer pbs.mutex.RUnlock()

	if _, exist := pbs.subscribers[topic]; !exist {
		pbs.subscribers[topic] = []chan string{}
	}

	ch := make(chan string)
	pbs.subscribers[topic] = append(pbs.subscribers[topic], ch)

	fmt.Printf("Created channel for topic %s\n", topic)

	pbs.messages[ch] = []string{}
	return ch
}

func (pbs *concurrentPubsub) Publish(message Message) error {
	pbs.mutex.Lock()
	defer pbs.mutex.Unlock()

	topic := message.Topic
	payload := message.Payload

	subscribers, exist := pbs.subscribers[topic]
	if !exist {
		return errors.New(fmt.Sprintf("topic %s does not exist\n", topic))
	}

	for _, ch := range subscribers {
		go func(ch chan string) {
			pbs.messages[ch] = append(pbs.messages[ch], payload)
			ch <- payload
		}(ch)
	}

	return nil
}

func (pbs *concurrentPubsub) ResetOffset(subscriber chan string, offset int) error {
	pbs.mutex.RLock()
	defer pbs.mutex.RUnlock()

	if _, exist := pbs.messages[subscriber]; !exist {
		return errors.New("subscriber does not exist")
	}

	fmt.Printf("%d len of pbs.messages[subscriber]\n", len(pbs.messages[subscriber]))

	go func() {
		for i := offset; i < len(pbs.messages[subscriber]); i++ {
			subscriber <- pbs.messages[subscriber][i]
		}
	}()

	return nil
}

func (pbs *concurrentPubsub) GetSubscriberMessages(subscriber chan string) ([]string, error) {
	pbs.mutex.RLock()
	defer pbs.mutex.RUnlock()

	if _, exist := pbs.messages[subscriber]; !exist {
		return []string{}, errors.New("subscriber does not exist")
	}

	fmt.Printf("%d len of pbs.messages[subscriber]\n", len(pbs.messages[subscriber]))

	return pbs.messages[subscriber], nil
}
