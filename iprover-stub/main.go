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
	num := flag.Int("n", 5, "# of iterations")
	float := flag.String("float", "he", "string")
	flag.Parse()
	fmt.Println("launched with n", *num)
	fmt.Println("launched with float", *float)

	var loopNo uint64
	logs := make(chan Log)
	go handleMessages(logs)

	var input string
	var allInputs []string
	scanner := bufio.NewScanner((os.Stdin))
	for {
		scanner.Scan()
		input = scanner.Text()
		if len(input) == 0 {
			break
		}
		allInputs = append(allInputs, input)
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
