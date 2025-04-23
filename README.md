# README

Author: `William Cruz`

## How to run

The `Dockerfile` defines a self-contained Go reference environment.
Build and run the program using [Docker](https://docs.docker.com/get-started/get-docker/):
```
$ docker build -t challenge .
$ docker run --rm -it challenge --auth=<token>
```
Feel free to modify the `Dockerfile` as you see fit.

If go `1.23` or later is locally installed, run the program directly for convenience:
```
$ go run main.go --auth=<token>
```

To run the application tests:
```
$ go run test ./...
```

## Concurrency

The concurrency treatment on this solution is based on Golang implementation of _Mutual exclusion_.
Defined by `sync.Mutex`, every time this variable is locked it prevents other threads from accessing the same block of code simultaneously, only when unlocked can a new thread continue the execution.

At first, I tried sharing the same `sync.Mutex` pointer on all the `Storage` implementations. That didn't turn out well, because the _Move_ action would require locking two Storages at the same time, which wouldn't work as it is mutual exclusion.

Therefore, the final solution was to use the lock on the _Place_ and _Pickup_ functions of the kitchen, for a guaranteed consistent state of every action. Futhermore, each `Storage` uses the hash map implementation called `sync.Map`, which is for thread safe read and write to the map in Golang.

## Discard criteria

For starters, I decided on a metric to calculate for how long a given order stays fresh, based on the `freshness (sec)` value.
Using the bellow formula, we can know the exact timestamp a food _stops being fresh_.

```
time to spoil = order timestamp (unix micro) + freshness by storage in seconds (unix micro)
```

This value is assigned to the `TTL` variable in the `Order` struct the first time its stored anywhere, and does not change even when moved.

Then, by sorting the orders from least to greatest based on the `TTL`, the first one will always be the closest to spoiling (If not already).

For that purpose, we can use a _Priority Queue_ (a.k.a MinHeap) to keep track of orders as they are stored. This data structure can track new orders and retrieve the minimal very efficiently, since we have a defined value of `TTL` for comparisons of lesser/greater.
