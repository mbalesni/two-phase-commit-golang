# Two-Phase Commit Protocol

## Resources
- [Task Specification](https://courses.cs.ut.ee/LTAT.06.007/2021_spring/uploads/Main/Task3-2021.pdf)
- Two-Phase Commit protocol (2PC): [Lecture 13](https://courses.cs.ut.ee/LTAT.06.007/2021_spring/uploads/Main/Lecture12-2021.pdf)
## Plan

What's implemented:

- [X] Process Class (simulated, not actual multiprocessing)
- [X] Inter-process communication
- [X] Core Commit Process (see `network_test.go`)
- [ ] Synchronization after node failure
- [ ] Operations
  - [ ] Set-Value
  - [ ] Rollback
  - [ ] Add
  - [ ] Remove
  - [ ] Time-failure
  - [ ] Arbitraty-failure
  - [ ] Reload
- [ ] Error messages
- [ ] Loading from file
- [ ] CLI 

## Run tests

```bash
cd src
go test
```

## How to run

I have compiled several binaries, just for you. Choose the one `two-phase-program*` that matches your OS and at least one should work. If neither of the binaries works, check [How to compile](#how-to-compile)

```bash
two-phase-program 2PC.txt
```

## How to compile

To compile the program yourself:

1. Install Golang 1.15.5 from [here](https://golang.org/dl/#go1.15.5)
2. Run in the project directory :
```bash
go build 
```

