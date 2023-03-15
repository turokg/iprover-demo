package launcher

import (
	"backend/internal"
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

func NewLauncher(logger internal.Logger) *Launcher {
	return &Launcher{logger: logger}
}

type Launcher struct {
	logger internal.Logger
}

func (l *Launcher) Launch(ctx context.Context, wg *sync.WaitGroup, args LaunchArgs, output chan []byte) {
	wg.Add(1)
	defer func() { fmt.Println("FINISHED") }()
	defer wg.Done()

	l.logger.WithField("problemID", args.ProblemID).Info(ctx, "launching process")

	cmd := exec.Command(internal.BinPath, "--stdin", "true", "--fof", "true") // TODO add args

	stdout, err := cmd.StdoutPipe()
	stdin, err := cmd.StdinPipe()
	if err != nil {
		l.logger.Error(ctx, "got error while setting up pipes", err)
	}

	err = cmd.Start()

	if err != nil {
		l.logger.Error(ctx, "couldn't start process", err)
	}

	_, err = io.WriteString(stdin, args.ProblemText) // TODO read file
	if err != nil {
		l.logger.Error(ctx, "got error while writing to stdin", err)
	}
	err = stdin.Close()
	if err != nil {
		l.logger.Error(ctx, "got error while closing to stdin", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = cmd.Process.Signal(syscall.SIGINT)
				if err != nil {
					l.logger.Error(ctx, "error killing the process", err)
				} else {
					l.logger.Info(ctx, "sent SYGINT to the iProver")
					output <- NewSysLog("sending SYGINT to the iProver")
				}
				return
			}
		}
	}()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		//l.logger.WithField("text", m).Info(ctx, "got message from stdout") // TODO debug
		output <- NewAppLog(m)
	}
	l.logger.Info(ctx, "finished reading from stdin")
	err = scanner.Err()
	if err != nil {
		l.logger.Error(ctx, "error while scanning", err)
	}
	err = cmd.Wait()
	if err != nil {
		l.logger.Error(ctx, "error while waiting for cmd", err)
	}
	close(output)
	time.Sleep(5 * time.Second)
	l.logger.Info(ctx, "finished launch")
}
