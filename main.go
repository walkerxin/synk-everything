package main

import (
	"os/exec"
)

// [需求]关闭界面时结束主进程，杀死主进程时关闭界面
func main() {
	chromePath := "/usr/bin/google-chrome-stable"
	cmd := exec.Command(chromePath, "--app=https://baidu.com")
	cmd.Start()

	// ui, _ := lorca.New("https://bilibili.com", "", 600, 400)
	// chSignal := make(chan os.Signal, 1)
	// signal.Notify(chSignal, syscall.SIGINT, syscall.SIGTERM)
	// select {
	// case <-chSignal:
	// case <-ui.Done():
	// }
	// ui.Close()
}