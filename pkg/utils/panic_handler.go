package utils

import (
	"fmt"
	"runtime/debug"
)

func HandlePanic(handler func(interface{})) {
	if panicMsg := recover(); panicMsg != nil {
		panicStack := debug.Stack()

		logPanicStack(panicMsg, string(panicStack))
		if handler != nil {
			handler(panicMsg)
		}
	}
}

func logPanicStack(msg interface{}, stack string) {
	fmt.Printf("Execution Panic!\nMessage: %s\n%s", msg, stack)
}
