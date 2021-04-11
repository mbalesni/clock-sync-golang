# Clock Synchronization by Bully


## Resources
- [Task Specification](https://courses.cs.ut.ee/LTAT.06.007/2021_spring/uploads/Main/Task2-2021.pdf)
- Bully Election Algorithm: [Seminar 7](https://courses.cs.ut.ee/2021/ds/spring/Main/Instructions3), [Wiki](https://en.wikipedia.org/wiki/Bully_algorithm) (better)
- Berkeley Time Sync: [Lecture, Slide 27](https://courses.cs.ut.ee/LTAT.06.007/2021_spring/uploads/Main/Lecture6-2021.pdf) (better read task specification, we don't need that complexity)
## Plan

Fix issues:

- [X] Killing doesn't notify who the new coordinator is and which process started the election (as per specification)
- [X] After reloading the coordinator doesn't change to the correct one, or several coordinators appear (is election happening?)
- [X] Freeze doesn't stop Time

- [X] Time is not rendered properly
- [X] Killing all processes results crashes the program
- [X] Attempting to Freeze/Unfreeze/Set-time to a non-existing process crashes the program
- [ ] ??? Freezing the coordinator isn't followed by time sync with new coordinator
- [ ] ??? Should update `currentProcess.MaxCoordinatorWait` on processes changes
- [ ] ??? Synchronizing Clock after Killing/Loading/ is not instant and takes several seconds (why?)
- [ ] ??? When a coordinator is killed, which process starts election, basically a random one (first in the list, but the list order changes within one running program)?
- [ ] ??? When sending out messages, processes check if target is not frozen. It might make more sense if the Network does that.

Implement:

- [X] Process Class (simulated, not actual multiprocessing)
- [X] Inter-process communication
- [X] Election through Bullying
- [x] Clock ticking
- [x] Clock synchronization
- [x] Network Class

- [x] Operations
  - [X] List
  - [X] Clock
  - [x] Kill
  - [x] Set-time
  - [x] Freeze
  - [x] Unfreeze
  - [x] Reload
- [x] Loading from file
- [x] CLI 

## Run tests

```bash
cd src
go test
```

## How to run

I have compiled several binaries, just for you. Choose the one `clock-sync-by-bully*` that matches your OS and at least one should work. If neither of the binaries works, check [How to compile](#how-to-compile)

```bash
clock-sync-by-bully processes_file
```

## How to compile

To compile the program yourself:

1. Install Golang 1.15.5 from [here](https://golang.org/dl/#go1.15.5)
2. Run in the project directory :
```bash
go build 
```

