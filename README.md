## Description
[LMAX-Disruptor](https://lmax-exchange.github.io/disruptor) is a ulta low latency inter thread communication as an alternative to bounded queueu, which makes use of padding to avoid memory false sharing, pre assign things to be cache friendly. 
This is the implemention of LMAX-Disruptor in Golang. 

It send data to consumers/threads at constant rate, it does this using below primary techniq:
1. It avoids locks/mutex at OS level.
2. It produces no garbage at runtime, by preallocating the ring buffer/array, By avoiding garbage, the need for a garbage collection and the stop-the-world application pauses  introduced can be almost entirely avoided.
3. Cache line padding

and many more.


## Features
On a MacBook Pro (Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz) using Go 1.20.5, it can pass many hundreds of millions of messages per second (yes, its true) from one goroutine to another goroutine, which is far fastest then go channels.

In Go, the current implementation of a channel (`chan`) maintains a lock around send, receive, and `len` operations and it maxes out around 25 million messages per second for uncontended access&mdash;more than an orders of magnitude slower when compared to the Disruptor.  The same channel, when contended between OS threads only pushes about 7 million messages per second.

## Prerequisites
1. Install latest go lang binary (1.18.3 and above) 
2. Install latest vs code
## Run Test Case
```Shell
    go test
```
## Perfomance
Running all unit test cases in benchmarks.
```Shell
    go test -v -bench=Benchmark -benchtime=1000000000x -count 3 ./v1/benchmarks
    go test -v -bench=Benchmark -benchtime=1000000000x -count 3 ./v2/benchmarks
```
Benchmarks
----------------------------
Each of the following benchmark tests sends an incrementing sequence message from one goroutine to another. The receiving goroutine asserts that the message is received is the expected incrementing sequence value. Any failures cause a panic.

* CPU: `Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz`
* Go Runtime: `Go 1.20.5`

Scenario | Per Operation Time | Bytes per Operation | Allocation per operation | Version
-------- | ------------------ | ------------------- | ------------------------ | --------
Disruptor: Reserve One & One Consumer | 8.973 ns/op  | 0 B/op | 0 allocs/op | v1
Disruptor: Reserve Many & One Consumer | 1.449 ns/op | 0 B/op | 0 allocs/op | v1
Disruptor: Reserve One Multiple Consumer | 9.656 ns/op | 0 B/op | 0 allocs/op | v1
Disruptor: Reserve Many & Multiple Consumer | 1.488 ns/op | 0 B/op | 0 allocs/op | v1
Disruptor: Reserve One & One Consumer | 10.403 ns/op  | 0 B/op | 0 allocs/op | v2
Disruptor: Reserve Many & One Consumer | 2.449 ns/op | 0 B/op | 0 allocs/op | v2
Disruptor: Reserve One Multiple Consumer | 10.656 ns/op | 0 B/op | 0 allocs/op | v2
Disruptor: Reserve Many & Multiple Consumer | 2.488 ns/op | 0 B/op | 0 allocs/op | v2

No new memory/bytes was assigned in each operation, this is the power of preallocated ring buffer, and considering cache line and Mechanical Sympathy.

## Example
Example is implemented in example/main.go
 
## Note
V1 version contains raw level of code in which user have to create ring buffer on his side and do the publish,commit, reserve, like operations.

V2 version contains generic level of disruptor package in which user will only initialise, publish, recieve the data.
    If you are new to LMAX Disuptor code will recomment to use v2 version.

Latency can seen in above Benchmarks for both the versions.

## Pending
1. Implementation for handling multiple data type for publish and consume.(use Golang Generic Function)
2. Support for multiple producers.
