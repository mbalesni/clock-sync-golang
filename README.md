# Clock Synchronization by Bully

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

