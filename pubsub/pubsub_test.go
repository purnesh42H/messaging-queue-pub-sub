package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testTopic1   = "topic1"
	testTopic2   = "topic2"
	testMessage1 = "test message 1"
	testMessage2 = "test message 2"
)

var (
	pubSub = NewConcurrentPubSub()
)

func TestPubSub(t *testing.T) {
	subscriber1 := pubSub.Subscribe(testTopic1)
	subscriber2 := pubSub.Subscribe(testTopic2)
	subscriber3 := pubSub.Subscribe(testTopic1)

	subscribers := map[chan string]string{}

	go func() {
		subscribers[subscriber1] = <-subscriber1
	}()
	go func() {
		subscribers[subscriber2] = <-subscriber2
	}()
	go func() {
		subscribers[subscriber3] = <-subscriber3
	}()

	pubSub.Publish(Message{Topic: testTopic1, Payload: testMessage1})
	pubSub.Publish(Message{Topic: testTopic2, Payload: testMessage2})

	time.Sleep(time.Second)

	assert.Equal(t, testMessage1, subscribers[subscriber1])
	assert.Equal(t, testMessage2, subscribers[subscriber2])
	assert.Equal(t, testMessage1, subscribers[subscriber3])
}
