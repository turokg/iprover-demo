package launcher

import (
	"backend/internal"
	"bufio"
	"context"
	"io"
	"os/exec"
	"syscall"
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
	output := make(chan []byte, internal.LaunchBuffer)

	go func() {
		cmd := exec.Command(internal.BinPath, "--stdin", "true") // TODO add args

		stdout, err := cmd.StdoutPipe()
		stdin, err := cmd.StdinPipe()
		if err != nil {
			l.logger.Error(ctx, "got error while setting up pipes", err)
		}

		err = cmd.Start()

		if err != nil {
			l.logger.Error(ctx, "couldn't start process", err)
		}
		//go func() {
		//	time.Sleep(internal.KillTimeout)
		//	cmd.Process.Kill()
		//}()

		_, err = io.WriteString(stdin, args.ProblemText) // TODO read file
		if err != nil {
			l.logger.Error(ctx, "got error while writing to stdin", err)
		}
		err = stdin.Close()
		if err != nil {
			l.logger.Error(ctx, "got error while closing to stdin", err)
		}
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			l.logger.Info(ctx, "scanning cycler")
			m := scanner.Text()
			l.logger.WithField("text", m).Info(ctx, "got message from stdout")
			output <- NewAppLog(m)
			if ctx.Err() != nil {
				cmd.Process.Signal(syscall.SIGINT)
				output <- NewSysLog("sending SYGINT to the iProver")
			}

		}
		l.logger.Info(ctx, "got out of cycle")
		err = scanner.Err()
		if err != nil {
			l.logger.Error(ctx, "error while scanning", err)
		}
		err = cmd.Wait()
		if err != nil {
			l.logger.Error(ctx, "error while waiting for cmd", err)
		}
		close(output)

	}()
	return output, nil
}
