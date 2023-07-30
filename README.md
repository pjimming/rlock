# RLock

## Status
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Usage
```shell
go get -u github.com/pjimming/rlock
```

## [Lua Scripts](./lua.md)
> Hint: Your redis should support lua script.

## RLock Type
### [Distributed RLock](./rlock_distributed.go)
Redis distributed lock is a mechanism for implementing concurrency control in distributed systems. It is based on Redis, a high-performance in-memory database, and uses its single-threaded, atomic operations, and fast network access features to ensure data consistency and thread safety when multiple clients access shared resources at the same time.

The main purpose of distributed locks is to prevent multiple clients from performing sensitive operations on a shared resource at the same time, thereby preventing data corruption, dirty reads, or other concurrency issues. These sensitive operations may include updating database records, performing some critical calculations, or other tasks that require mutual exclusion.

### [Reentry RLock](./rlock_reentry.go)
Redis reentrant lock is an implementation of distributed lock, which allows the same client to acquire the lock multiple times without being blocked by the lock it holds after acquiring the lock. This locking mechanism enables the same client to safely use distributed locks in multi-level nested calls or recursive functions without causing deadlocks or concurrency issues due to repeated lock acquisitions.
