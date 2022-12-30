package launcher

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"syscall"
)

const (
	path       = "/Users/eyukorovin/iprover-demo/iprover-stub/iprover"
	terminator = "\n"
)

func Launch(ctx context.Context, inputParams string) (chan string, error) {
	fmt.Println("launching process with args ", inputParams)
	output := make(chan string)

	go func() {
		cmd := exec.Command(path)

		stdout, _ := cmd.StdoutPipe()

		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println("got error while writing", err)
		}

		cmd.Start()
		fmt.Println(stdin)
		_, err = io.WriteString(stdin, inputParams)
		_, err = io.WriteString(stdin, terminator)
		_, err = io.WriteString(stdin, terminator)
		if err != nil {
			fmt.Println("got error while writing", err)
		}
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			output <- m

			if ctx.Err() != nil {
				cmd.Process.Signal(syscall.SIGINT)
			}
		}
		cmd.Wait()

	}()
	return output, nil
}
