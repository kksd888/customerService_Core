package handle

import (
	"os"
	"os/signal"
	"syscall"
)

var sigs = make(chan os.Signal, 1)

// 监听 优雅停服
func Listen() {

	// 优雅停机
	go func() {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		os.Exit(0)
	}()
}
