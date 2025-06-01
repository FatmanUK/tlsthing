package main

import (
	"os"
	"fmt"
)

const (
	LL_ERROR = uint16(1)
	LL_WARN  = uint16(2)
	LL_INFO  = uint16(3)
	LL_DEBUG = uint16(4)
)

type View struct {
	loglevel uint16
}

func (re View) begin() {
	re.log(LL_DEBUG, "Starting View object.")
}

func (re View) end() {
	re.log(LL_DEBUG, "Stopping View object.")
}

func (re View) log(i uint16, s string) {
	if re.loglevel >= i {
		fmt.Print("[")
		switch i {
			case LL_ERROR: {
				fmt.Print("XXX")
			}
			case LL_WARN: {
				fmt.Print("!!!")
			}
			case LL_INFO: {
				fmt.Print("   ")
			}
			case LL_DEBUG: {
				fmt.Print(">>>")
			}
		}
		fmt.Print("] ")
		fmt.Println(s)
		os.Stdout.Sync()
	}
}
