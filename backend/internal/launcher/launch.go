package launcher

import (
	"backend/internal"
	"bufio"
	"context"
	"io"
	"os/exec"
	"syscall"
	"time"
)

const (
	terminator = "\n\n"
)

func NewLauncher(logger internal.Logger) *Launcher {
	return &Launcher{logger: logger}
}

type Launcher struct {
	logger internal.Logger
}

func (l *Launcher) Launch(ctx context.Context, args LaunchArgs) (chan []byte, error) {
	l.logger.WithField("args", args).Info(ctx, "launching process with args")
	output := make(chan []byte)

	go func() {
		cmd := exec.Command(internal.BinPath) // TODO add args

		stdout, err := cmd.StdoutPipe()
		stdin, err := cmd.StdinPipe()
		if err != nil {
			l.logger.Error(ctx, "got error while setting up pipes", err)
		}

		err = cmd.Start()
		if err != nil {
			l.logger.Error(ctx, "couldn't start process", err)
		}
		_, err = io.WriteString(stdin, "hello world") // TODO read file
		if err != nil {
			l.logger.Error(ctx, "got error while writing to stdin", err)
		}
		_, err = io.WriteString(stdin, terminator)
		if err != nil {
			l.logger.Error(ctx, "got error while writing to stdin", err)
		}
		_, err = io.WriteString(stdin, terminator)
		if err != nil {
			l.logger.Error(ctx, "got error while writing to stdin", err)
		}
		if err != nil {
			l.logger.Error(ctx, "got error while writing to stdin", err)
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
				time.Sleep(internal.KillTimeout)
				cmd.Process.Kill()
			}()
		}
		cmd.Wait()
		close(output)

	}()
	return output, nil
}
