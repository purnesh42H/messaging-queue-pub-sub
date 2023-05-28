
## Requirements

- Pubsub can have multiple topics
- Multiple subscribers should be able to subscribe to multiple topics
- When a message is published to a topic by a publisher, all subscribers of that topic should get the message
- Publisher should be separate from subscriber. Subscriber should only get messages
- Subscribers should be able to run in parallel

### Additional Requirements
- Have a way to get all messages of a subscriber again
- Have a way to get all messages of a subscriber from a particular offset. Order of publish is not needed to be maintained