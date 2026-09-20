package main

import (
	"hash"
	"sync"
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{typ: "string", str: "PONG"}
	}

	return Value{typ: "string", str: args[0].bulk}
}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET": set,
	"GET": get,
	"HSET": hset,
	"HGET": hget,

}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutux{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "err wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "err wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	value, ok := SETs[key]

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "err wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash][key]; !ok {
		HSETs[hash] = map[string]string{}
	}

	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "err wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[0].bulk

	HSETsMu.Lock()
	value, ok := HSETs[hash][key]
	HSETsMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bluk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "err wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	HSETsMu.Lock()
	key, ok := SETs[hash]
	value := SETs[key]
	HSETsMu.Unlock()
	
	
}