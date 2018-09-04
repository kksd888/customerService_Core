package handle

import (
	"os"
	"os/signal"
	"syscall"
)

var (
	sigs = make(chan os.Signal, 1)
)

// 监听 优雅停服、客服下线，用户消息过期
func Listen() {

	// 优雅停机
	go func() {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		// todo 所有客服强制下线 ？
		os.Exit(0)
	}()
}
