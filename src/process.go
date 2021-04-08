package src

import (
	"fmt"
	"math/rand"
	"strconv"

	"time"
)

type Message struct {
	From        int
	To          int
	MessageType string // one of ELECTION|ELECTION_RESPONSE|COORDINATOR
	ElectionId  int
	Time        int64
}

type ProcessShallow struct {
	Id int
}

type Process struct {
	Id                 int
	Name               string
	InitialTime        string
	InitialSysTime     time.Time
	Time               string
	HigherProcesses    []*ProcessShallow
	LowerProcesses     []*ProcessShallow
	Coordinator        ProcessShallow
	SendQueue          MessageQueue
	GetQueue           MessageQueue
	WaitingElection    int
	WaitingCoordinator int
	MaxElectionWait    int
	MaxCoordinatorWait int
	Verbose            bool
	Frozen             bool
}

func NewProcess(id int, name string, initialTime string, verbose bool) *Process {
	process := Process{Id: id, Name: name, InitialTime: initialTime,
		Time: initialTime, WaitingElection: -1, WaitingCoordinator: -1,
		MaxElectionWait: 1, Verbose: verbose}
	process.Init()
	return &process
}

func (p *Process) UpdateElectionCount() {
	electionCount, err := strconv.Atoi(string(p.Name[2:]))
	if err != nil {
		panic(err)
	}
	electionCount = electionCount + 1
	p.Name = fmt.Sprintf("%s_%d", string(p.Name[:1]), electionCount)
}

func (p *Process) SyncTime() {

}

func (p *Process) RunElection(electionId int) {
	p.WaitingElection = 0

	if electionId < 0 { // new election
		rand.Seed(time.Now().UnixNano())
		electionId = rand.Intn(10000)
		if p.Verbose {
			fmt.Println("P=", p.Id, "Starting new Election:", electionId)
		}
	} else {
		if p.Verbose {
			fmt.Println("P=", p.Id, "Continuing Election:", electionId)
		}
	}

	for _, target := range p.HigherProcesses {
		message := p.NewElectionMessage(target.Id, "ELECTION", electionId)
		p.SendQueue.Add(message)
	}

	// run logic of initiating election, i.e.
	// send messages to higher ids
	// wait for their replies etc etc
}

func (p *Process) GetElectionCount() int {
	electionCount, err := strconv.Atoi(string(p.Name[2:]))
	if err != nil {
		panic(err)
	}
	return electionCount
}

func (p *Process) NewMessage(toId int, messageType string) Message {
	return Message{Time: time.Now().Unix(), From: p.Id, To: toId, MessageType: messageType}
}

func (p *Process) NewElectionMessage(toId int, messageType string, electionId int) Message {
	return Message{Time: time.Now().Unix(), From: p.Id, To: toId, MessageType: messageType, ElectionId: electionId}
}

func (p *Process) ProcessMessages() {
	for p.GetQueue.queue.Len() > 0 {
		message := p.GetQueue.Pop()
		if message.MessageType == "ELECTION" {
			if p.Verbose {
				fmt.Println("P=", p.Id, "got an ELECTION from", message.From)
			}
			// respond to election request
			electionResponseMessage := p.NewMessage(message.From, "ELECTION_RESPONSE")
			p.SendQueue.Add(electionResponseMessage)
			p.RunElection(message.ElectionId)
		} else if message.MessageType == "ELECTION_RESPONSE" {
			if p.Verbose {
				fmt.Println("P=", p.Id, "got an ELECTION_RESPONSE from", message.From)
			}
			p.WaitingElection = -1   // stop waiting
			p.WaitingCoordinator = 0 // start waiting
		} else if message.MessageType == "COORDINATOR" {
			if p.Verbose {
				fmt.Println("P=", p.Id, "got a COORDINATOR from", message.From)
			}
			p.UpdateElectionCount()
			p.Coordinator.Id = message.From
			p.WaitingCoordinator = -1 // stop waiting
		}
	}
}

func (p *Process) Cycle() {

	// increment waiting counters (if waiting)
	if p.WaitingCoordinator > -1 {
		p.WaitingCoordinator += 1
	}
	if p.WaitingElection > -1 {
		p.WaitingElection = p.WaitingElection + 1
	}

	p.ProcessMessages()

	// waited for Election reply for too long -> become coordinator
	if p.WaitingElection > p.MaxElectionWait {
		if p.Verbose {
			fmt.Println("P=", p.Id, "waited too long for an election responose. Becoming a coordinator!")
		}
		p.UpdateElectionCount()
		p.WaitingElection = -1 // stop waiting
		p.Coordinator.Id = p.Id
		for _, process := range p.LowerProcesses {
			electionMessage := p.NewMessage(process.Id, "COORDINATOR")
			p.SendQueue.Add(electionMessage)
		}
	}

	// waited for Coordinator message for too long -> re-init election
	if p.WaitingCoordinator > p.MaxCoordinatorWait {
		p.WaitingCoordinator = -1
		p.RunElection(-1) // new election
	}
}

func (p *Process) Init() {
	p.SendQueue = MessageQueue{}
	p.GetQueue = MessageQueue{}
	p.SendQueue.Init()
	p.GetQueue.Init()
}
