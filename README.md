# BMO-Redis

A lightweight Redis clone written in Go. It speaks the Redis wire protocol (RESP), so you can talk to it with the standard `redis-cli` or any Redis client library, and it persists writes to disk using an append-only file (AOF).

Built as a learning project to understand how Redis works under the hood: TCP servers, protocol parsing, command dispatch, and durability.

## Features

- **RESP parser and writer**: reads and writes the Redis Serialization Protocol over TCP
- **Command handlers**: commands are dispatched through a handler map, so adding a new command is one function
- **AOF persistence**: write commands (`SET`, `HSET`) are appended to `database.aof` and replayed on startup, so data survives restarts
- **Redis-compatible**: works with `redis-cli` out of the box on the default port `6379`

## Project structure

```
bmo-redis/
├── main.go      # TCP listener, startup AOF replay, request/response loop
├── resp.go      # RESP protocol reader and writer
├── handler.go   # Command handlers and the command → handler map
├── aof.go       # Append-only file: write, read, replay
└── README.md
```

## How it works

1. **Startup**: the server listens on TCP port `6379` and opens (or creates) `database.aof`.
2. **Recovery**: every command stored in the AOF is read back and executed through the same handlers, rebuilding the in-memory state.
3. **Serving**: each incoming request is parsed as a RESP array. The first element is the command name (case-insensitive) and the rest are its arguments.
4. **Dispatch**: the command is looked up in the handler map and executed, and the result is encoded back to the client as RESP.
5. **Durability**: `SET` and `HSET` requests are written to the AOF before they are executed.

## Getting started

### Prerequisites

- [Go](https://go.dev/dl/) 1.18 or newer
- [`redis-cli`](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/) (optional, for testing)

### Installation

```bash
git clone https://github.com/student-ankitpandit/bmo-redis.git
cd bmo-redis
```

If the repo does not contain a `go.mod` yet, initialise the module first:

```bash
go mod init bmo-redis
```

### Run the server

```bash
go run .
```

You should see:

```
Listening on port :6379
```

### Connect

In another terminal:

```bash
redis-cli -p 6379
```

## Usage

```
127.0.0.1:6379> SET name bmo
OK
127.0.0.1:6379> GET name
"bmo"
127.0.0.1:6379> HSET user1 name ankit
OK
127.0.0.1:6379> HGET user1 name
"ankit"
127.0.0.1:6379> HGETALL user1
1) "name"
2) "ankit"
```

## Supported commands

| Command   | Description                         | Persisted to AOF |
| --------- | ----------------------------------- | :--------------: |
| `PING`    | Health check                        |        No        |
| `SET`     | Set a string value for a key        |        Yes       |
| `GET`     | Get the string value of a key       |        No        |
| `HSET`    | Set a field in a hash               |        Yes       |
| `HGET`    | Get a field from a hash             |        No        |
| `HGETALL` | Get all fields and values in a hash |        No        |

## Persistence

Writes are stored in `database.aof` in the directory where you run the server. The file contains the raw RESP-encoded commands. To reset the database, stop the server and delete `database.aof`.

## Current limitations

- Accepts a **single client connection** at a time
- Only the string and hash data types are supported
- The AOF is never compacted or rewritten, so it grows with every write
- No key expiry, transactions, or pub/sub

## Roadmap

- [ ] Handle multiple concurrent clients with goroutines
- [ ] Guard shared state with a mutex (`sync.RWMutex`)
- [ ] Configurable port and AOF path via flags or env vars
- [ ] AOF fsync policy and compaction
- [ ] More commands (`DEL`, `EXISTS`, `EXPIRE`, `KEYS`, lists, sets)
- [ ] Unit tests for the RESP parser and handlers

## Contributing

Issues and pull requests are welcome. For larger changes, please open an issue first to discuss what you would like to change.

## Author

**Ankit Pandit**: [@student-ankitpandit](https://github.com/student-ankitpandit)
