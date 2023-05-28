package main

import (
	"fmt"
	"messaging-queue-pub-sub/pubsub"
	"time"
)

func main() {
	pubSub := pubsub.NewConcurrentPubSub()

	subsciber1 := pubSub.Subscribe("topic1")
	go func() {
		for {
			payload := <-subsciber1
			fmt.Printf("Subscriber1 got payload %v on %v\n", payload, time.Now().UnixMilli())
		}
	}()

	subsciber2 := pubSub.Subscribe("topic2")
	go func() {
		for {
			payload := <-subsciber2
			fmt.Printf("Subscriber2 got payload %v on %v\n", payload, time.Now().UnixMilli())
		}
	}()

	for i := 1; i <= 8; i++ {
		if i <= 4 {
			go pubSub.Publish(pubsub.Message{Topic: "topic1", Payload: fmt.Sprintf("%d", i)})
		} else {
			go pubSub.Publish(pubsub.Message{Topic: "topic2", Payload: fmt.Sprintf("%d", i)})
		}
	}
	time.Sleep(time.Second)
	time.Sleep(time.Second)

	messages, _ := pubSub.GetSubscriberMessages(subsciber1)
	fmt.Printf("Subscriber1 all messages %v", messages)
	pubSub.ResetOffset(subsciber1, 1)
	time.Sleep(time.Second)
}
