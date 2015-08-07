package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

const (
	TASK_COUNT = 256
)

var (
	ch   = make(chan bool)
	stop = make(chan bool)
)

func main() {
	spawnListeners()
	spawnSender()
forever:
	for {
		select {
		case _ = <-stop:
			break forever
		}
	}
}

func handle(i int) {
	var (
		stderr bytes.Buffer
		stdout bytes.Buffer
	)

	cmd := &exec.Cmd{
		Path:   "/bin/ls",
		Args:   []string{},
		Env:    []string{"LC_ALL=C", "LANG=C"},
		Stdout: &stdout,
		Stderr: &stderr,
	}
	defer func() {
		if eerr := recover(); eerr != nil {
			fmt.Printf("panic: %s ", eerr)
		}
	}()
	err := cmd.Run()

	fmt.Printf("%d %s: %s %s", i, err, stdout.String(), stderr.String())
}

func listener(i int) {
	for {
		select {
		case _ = <-ch:
			handle(i)
		}
	}
}

func sender() {
	for i := 0; i < TASK_COUNT; i++ {
		ch <- true
	}
	time.Sleep(time.Second)
	sender()
}

func spawnListeners() {
	c := runtime.NumCPU() * 4
	for i := 0; i < c; i++ {
		go listener(i)
	}
}

func spawnSender() {
	sender()
}
