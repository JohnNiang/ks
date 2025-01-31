package common

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

// ExecCommand runs a command
func ExecCommand(name string, arg ...string) (err error) {
	command := exec.Command(name, arg...)

	//var stdout []byte
	//var errStdout error
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()
	err = command.Start()
	if err != nil {
		return err
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	_, _ = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = command.Wait()
	return
}

//ExecCommandGetOutput run command and get output
func ExecCommandGetOutput(name string, arg ...string) (string, error) {
	command := exec.Command(name, arg...)
	output, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(output), err
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}
