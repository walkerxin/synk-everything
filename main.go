package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
)

// [需求]关闭界面时结束主进程，杀死主进程时关闭界面
func main() {
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.Writer.Write([]byte("abcdefg"))
		})
		router.Run(":8080")
	}()

	ui, _ := lorca.New("http://localhost:8080", "", 600, 500)

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}
