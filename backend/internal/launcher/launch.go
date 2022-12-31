package launcher

import (
	"backend/internal/conf"
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"syscall"
	"time"
)

const (
	terminator = "\n\n"
)

func Launch(ctx context.Context, inputParams string) (chan []byte, error) {
	fmt.Println("launching process with args ", inputParams)
	output := make(chan []byte)

	go func() {
		cmd := exec.Command(conf.BinPath)

		stdout, _ := cmd.StdoutPipe()

		stdin, err := cmd.StdinPipe()
		if err != nil {
			fmt.Println("got error while writing", err)
		}

		cmd.Start()
		fmt.Println(stdin)
		_, err = io.WriteString(stdin, inputParams)
		_, err = io.WriteString(stdin, terminator)
		if err != nil {
			fmt.Println("got error while writing", err)
		}
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			output <- NewAppLog(m)
			if ctx.Err() != nil {
				cmd.Process.Signal(syscall.SIGINT)
				output <- NewSysLog("sending SYGINT to the iProver")
			}
			go func() {
				time.Sleep(conf.KillTimeout)
				cmd.Process.Kill()
			}()
		}
		cmd.Wait()
		close(output)

	}()
	return output, nil
}
