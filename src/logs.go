package main

import (
	"log"
	"fmt"
	"runtime"
	"strings"
	"unicode"
)

/*
Usage:
...
	logs := makeLogChan()
	go threadMainLoop(args, logs)
	threadLoggingLoop(uint16(threshold), logs)
...

func threadMainLoop(args argMap, logs chan Log) {
	defer close(logs)
...
	logs <- Error("Error")
	logs <- Warn("Warn" )
	logs <- Info("Info" )
	logs <- Debug("Debug")
...
}
*/

type Log struct {
	Level uint
	Message string
	Prefix string
}

func Error(err string) Log {
	return Log{ 10, err, "ERR" }
}

func Warn(err string) Log {
	return Log{ 20, err, "WRN" }
}

func Info(err string) Log {
	return Log{ 30, err, "INF" }
}

func Debug(err string) Log {
	return Log{ 40, err, "DBG" }
}

func MakeLogChan() chan Log {
	return make(chan Log)
}

func ThreadLoggingLoop(threshold uint, logs chan Log) {
	runtime.LockOSThread()
	pre := map[uint]string{
		Error("").Level: Error("").Prefix,
		 Warn("").Level:  Warn("").Prefix,
		 Info("").Level:  Info("").Prefix,
		Debug("").Level: Debug("").Prefix,
	}
	levels := []uint{
		Error("").Level,
		Warn("").Level,
		Info("").Level,
		Debug("").Level,
	}
	for msg := range logs {
		if levels[threshold] >= msg.Level {
			outputMessage(pre[msg.Level], msg.Message)
		}
	}
}

func outputMessage(prefix string, message string) {
	message_strings := strings.Split(message, "\\n")
	for _, line := range message_strings {
		if line != "" {
			line = fmt.Sprintf("[%s] %s", prefix, line)
			log.Println(line)
		}
	}
}

func DebugOutputByteArray(buffer []byte, logs chan Log) {
	logs <- Debug(fmt.Sprintf("Buffer size: %d b", len(buffer)))
	msg := "/"
	for b := range buffer {
		msg += fmt.Sprintf(" %d", int(buffer[b]))
		if unicode.IsPrint(rune(buffer[b])) {
			msg += fmt.Sprintf(" (char %c)", buffer[b])
		}
		msg += " /"
	}
	logs <- Debug(msg)
}
