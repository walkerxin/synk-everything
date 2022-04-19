package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/gin-gonic/gin"
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

	var chromePath string
	switch runtime.GOOS {
	case "windows":
		chromePath = os.Getenv("ProgramFiles(x86)") + "/Google/Chrome/Application/chrome.exe"
	default:
		chromePath = "/usr/bin/google-chrome-stable"
	}
	cmd := exec.Command(chromePath, "--app=http://localhost:8080")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-chSignal
	err := cmd.Process.Kill()
	if err != nil {
		log.Println(cmd.Process.Pid, err)
	}

	// ui, _ := lorca.New("https://bilibili.com", "", 600, 400)
	// chSignal := make(chan os.Signal, 1)
	// signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	// select {
	// case <-chSignal:
	// case <-ui.Done():
	// }
	// ui.Close()
}
