# Clock Synchronization by Bully


## Resources
- [Task Specification](https://courses.cs.ut.ee/LTAT.06.007/2021_spring/uploads/Main/Task2-2021.pdf)
- Bully Election Algorithm: [Seminar 7](https://courses.cs.ut.ee/2021/ds/spring/Main/Instructions3), [Wiki](https://en.wikipedia.org/wiki/Bully_algorithm) (better)
- Berkeley Time Sync: [Lecture, Slide 27](https://courses.cs.ut.ee/LTAT.06.007/2021_spring/uploads/Main/Lecture6-2021.pdf) (better read task specification, we don't need that complexity)
## Plan

What's implemented:

- [X] Process Class (simulated, not actual multiprocessing)
- [X] Inter-process communication
- [X] Election through Bullying
- [ ] Clock ticking
- [ ] Clock synchronization
- [ ] Network Class (now inter-process communication is done by a test script `network_test.go`)

- [ ] Operations
  - [X] List
  - [X] Clock
  - [ ] Kill
  - [ ] Set-time
  - [ ] Freeze
  - [ ] Unfreeze
  - [ ] Reload
- [ ] Loading from file
- [ ] CLI 

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

