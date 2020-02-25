package exitdo

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

type sig int

const Signal sig = iota

type done struct {
	Done func(func(sig os.Signal))
}

func (s sig) ListenKill() *done {
	var signalList = []os.Signal{syscall.SIGINT, syscall.SIGTERM,syscall.SIGKILL}
	// 创建信号
	signalChan := make(chan os.Signal, 1)
	// 通知
	signal.Notify(signalChan, signalList...)
	return &done{Done: func(f func(signal os.Signal)) {
		// 阻塞
		f(<-signalChan)
		// 停止
		signal.Stop(signalChan)
	}}
}

func (s sig) Listen(sig ...os.Signal) *done {
	var signalList = sig
	// 创建信号
	signalChan := make(chan os.Signal, 1)
	// 通知
	signal.Notify(signalChan, signalList...)
	return &done{Done: func(f func(signal os.Signal)) {
		// 阻塞
		f(<-signalChan)
		// 停止
		signal.Stop(signalChan)
	}}
}

func (s sig) Kill(pid int) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout
	return kill.Run()
}
