package main

import (
	"os/exec"

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

	chromePath := "/usr/bin/google-chrome-stable"
	cmd := exec.Command(chromePath, "--app=http://localhost:8080")
	cmd.Start()
	select {}
	// 主线程结束，其中的协程也结束
	/*
		for {}
		go程 无限 1k～1w个
		ErrServerClose
		
		阻塞 就像等待戈多（同步读） 什么都不能做 挂起(快照，进程消失？)
		死循环并不是阻塞
	*/

	// ui, _ := lorca.New("https://bilibili.com", "", 600, 400)
	// chSignal := make(chan os.Signal, 1)
	// signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	// select {
	// case <-chSignal:
	// case <-ui.Done():
	// }
	// ui.Close()
}