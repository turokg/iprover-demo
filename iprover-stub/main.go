package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	sleepPeriod    = time.Second
	defaultMessage = "going to sleep"
	terminator     = "\n"
)

type Log struct {
	RunTime uint64 `json:"runTime"`
	Message string `json:"message"`
}

func NewLog(runtime uint64, message string) Log {
	return Log{RunTime: runtime, Message: message}
}
func main() {
	num := flag.Bool("stdin", true, "# of iterations")
	flag.Parse()
	fmt.Println("launched with n", *num)

	var loopNo uint64
	logs := make(chan Log)
	go handleMessages(logs)

	var input string
	var allInputs []string
	scanner := bufio.NewScanner(os.Stdin)
	logs <- NewLog(0, "started scanning")
	for scanner.Scan() {
		input = scanner.Text()
		logs <- NewLog(0, fmt.Sprintf("got line from stdin: %s", input))

		allInputs = append(allInputs, input)
	}
	if err := scanner.Err(); err != nil {
		logs <- NewLog(0, fmt.Sprintf("error reading standard input: %s", err.Error()))
	}
	logs <- NewLog(0, fmt.Sprintf("Starting. Read from stdin: %s", strings.Join(allInputs, "\n")))

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGINT)
	go func() {
		sig := <-signalChannel
		logs <- NewLog(0, fmt.Sprintf("Recieved os signal [%s]. Exitng.", sig))
		time.Sleep(time.Second)
		os.Exit(0)
	}()

	for {
		logs <- NewLog(loopNo, defaultMessage)
		time.Sleep(sleepPeriod)
		loopNo += 1
	}
}

func handleMessages(logs chan Log) {
	for log := range logs {
		js, _ := json.Marshal(log)
		fmt.Println(string(js))
	}
}
