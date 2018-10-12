package main

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"runtime"
	"time"
)

func runcmd(command string) (res string) {
	var shell, flag string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		flag = "/c"
	} else {
		shell = "/bin/sh"
		flag = "-c"
	}
	res, err := run(10, shell, flag, command)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

func run(timeout int, command string, args ...string) (res string, err error) {

	// instantiate new command
	cmd := exec.Command(command, args...)

	// get pipe to standard output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return res, errors.New("cmd.StdoutPipe() error: " + err.Error())
	}

	// start process via command
	if err = cmd.Start(); err != nil {
		return res, errors.New("cmd.Start() error: " + err.Error())
	}

	// setup a buffer to capture standard output
	var buf bytes.Buffer

	// create a channel to capture any errors from wait
	done := make(chan error)
	go func() {
		if _, err := buf.ReadFrom(stdout); err != nil {
			panic("buf.Read(stdout) error: " + err.Error())
		}
		done <- cmd.Wait()
	}()

	// block on select, and switch based on actions received
	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			return res, errors.New("failed to kill: " + err.Error())
		}
		return "", errors.New("command timed out")
	case err = <-done:
		if err != nil {
			close(done)
			return res, errors.New("process done, with error: " + err.Error())
		}
		return buf.String(), nil
	}
	return "", nil
}
