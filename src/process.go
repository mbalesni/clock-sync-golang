package src

import (
	"container/list"
	"fmt"
	"strconv"
)

type Message struct {
	From int
	To   int
}

type Process struct {
	Id              int
	Name            string // e.g. A_0
	InitialTime     string
	Time            string // MVP without ticking and averaging
	HigherProcesses []*Process
	LowerProcesses  []*Process
	IsCoordinator   bool
	SendQueue       MessageQueue
	GetQueue        MessageQueue
}

func (p *Process) UpdateName() {
	fmt.Println(p.Name[2])
	electionCount, err := strconv.Atoi(string(p.Name[2]))
	if err != nil {
		panic(err)
	}
	electionCount = electionCount + 1
	p.Name = fmt.Sprintf("%s_%d", string(p.Name[:1]), electionCount)
}

func (p *Process) SyncTime() {

}

func (p *Process) StartElection() {
	p.UpdateName()
	// run logic of initiating election, i.e.
	// send messages to higher ids
	// wait for their replies etc etc
}

func (p *Process) Init() {
	p.SendQueue = MessageQueue{}
	p.GetQueue = MessageQueue{}
	p.SendQueue.Init()
	p.GetQueue.Init()
}

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
