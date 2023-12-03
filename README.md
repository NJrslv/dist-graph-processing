## Table of Contents
- [Introduction](#introduction)
- [Build](#build)
- [Details](#details)

## Introduction
`distgraphia` is an educational distributed system that is `not` fault-tolerant, with the aim of implementing the fundamental structure of a distributed system.


## Build
```bash
git clone https://github.com/NJrslv/distgraphia.git
cd distgraphia
docker build -t distgraphia .
docker run distgraphia
```

## Details
### Object Placement
- The distributed system is represented as a single process.
- Clients and the network are located on the stack.
- Each system node, embodied by a goroutine (method `node.Run()`) and its data, is isolated. The goroutine executes the node's functionality, and its stack holds a pointer to the data in the heap of the original process.
- Nodes communicate using channels

### Communication
- Clients communicate with the network via the channel `clientCh of requests`, and they receive responses from `request.replyCh`. Nodes communicate through internal channels, while the network selects a coordinator for the client.
![image](https://github.com/NJrslv/distgraphia/assets/108277031/2566403a-fc96-4a3c-b3bf-6a6e400fcfe3)



