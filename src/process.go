package src

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	From        *Process
	To          *Process
	MessageType string // one of ELECTION|ELECTION_RESPONSE|COORDINATOR
	ElectionId  int
	Time        *Time
}

type Process struct {
	Id                 int
	InitialCount       string
	Name               string
	InitialTime        *Time
	Time               *Time
	HigherProcesses    map[int]*Process
	LowerProcesses     map[int]*Process
	Coordinator        *Process
	SendQueue          MessageQueue
	GetQueue           MessageQueue
	WaitingElection    int
	WaitingCoordinator int
	MaxElectionWait    int
	MaxCoordinatorWait int
	Verbose            bool
	Frozen             bool
	mu                 sync.Mutex
}

func NewProcess(id int, name string, initialTime *Time, initialCount string, verbose bool) *Process {

	initialTimeCopy := *initialTime

	process := Process{Id: id, Name: name, InitialTime: &initialTimeCopy,
		Time: initialTime, WaitingElection: -1, WaitingCoordinator: -1,
		MaxElectionWait: 1, Verbose: verbose, InitialCount: initialCount}
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

func (p *Process) SyncTime(time *Time) {
	if !p.Frozen {
		p.mu.Lock()
		p.Time.Add(time)
		p.mu.Unlock()
	}
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
		if target.Frozen != true {
			message := p.NewElectionMessage(target, "ELECTION", electionId)
			p.SendQueue.Add(message)
		}
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

func (p *Process) NewSyncMessage(to *Process, messageType string, time *Time) Message {

	return Message{Time: time, From: p, To: to, MessageType: messageType}

}

func (p *Process) NewMessage(to *Process, messageType string) Message {
	return Message{Time: CurrentTime(), From: p, To: to, MessageType: messageType}

}

func (p *Process) NewElectionMessage(to *Process, messageType string, electionId int) Message {
	return Message{Time: CurrentTime(), From: p, To: to, MessageType: messageType, ElectionId: electionId}
}

func (p *Process) ProcessMessages() {
	for p.GetQueue.queue.Len() > 0 {
		message := p.GetQueue.Pop()
		switch messageType := message.MessageType; messageType {

		case "ELECTION":
			{
				if p.Verbose {
					fmt.Println("P=", p.Id, "got an ELECTION from", message.From.Id)
				}
				// respond to election request
				electionResponseMessage := p.NewMessage(message.From, "ELECTION_RESPONSE")
				p.SendQueue.Add(electionResponseMessage)
				p.RunElection(message.ElectionId)
			}
		case "ELECTION_RESPONSE":
			{
				if p.Verbose {
					fmt.Println("P=", p.Id, "got an ELECTION_RESPONSE from", message.From.Id)
				}
				p.WaitingElection = -1   // stop waiting
				p.WaitingCoordinator = 0 // start waiting
			}
		case "COORDINATOR":
			{
				if p.Verbose {
					fmt.Println("P=", p.Id, "got a COORDINATOR from", message.From.Id)
				}
				p.UpdateElectionCount()
				p.Coordinator = message.From
				p.WaitingCoordinator = -1
			}
		case "CLOCK_SYNC":
			{
				//fmt.Println("Got a clock sync")
				if p.Verbose {
					fmt.Println("P=", p.Id, "got a CLOCK_SYNC from", message.From.Id)
				}
				p.SyncTime(message.Time)
			}
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
		p.Coordinator = p
		for _, process := range p.LowerProcesses {
			if !process.Frozen {
				electionMessage := p.NewMessage(process, "COORDINATOR")
				p.SendQueue.Add(electionMessage)

			}
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
	p.HigherProcesses = make(map[int]*Process)
	p.LowerProcesses = make(map[int]*Process)
	p.Frozen = false
	p.SendQueue.Init()
	p.GetQueue.Init()
}
