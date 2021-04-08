package src

import "container/list"

type MessageQueue struct {
	queue *list.List
}

func (mq *MessageQueue) Init() {
	mq.queue = list.New()
}

func (mq *MessageQueue) Add(message Message) {
	mq.queue.PushBack(message)
}

func (mq *MessageQueue) Pop() Message {
	elem := mq.queue.Front()        // First element
	message := elem.Value.(Message) // Cast abstract Queue Element to Message
	mq.queue.Remove(elem)           // Dequeue
	return message
}
