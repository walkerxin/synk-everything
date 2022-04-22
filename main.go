package main

import (
	"os"
	"os/signal"
	"syscall"

	"walkerxin/synk-everything.git/server"

	"github.com/zserge/lorca"
)

func main() {
	port := "27149"
	go server.Run(port)
	ui := startBrowser(port)
	chSignal := listenToSignal()
	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}

func startBrowser(port string) lorca.UI {
	ui, _ := lorca.New("http://localhost:"+port+"/static", "", 1000, 600)
	return ui
}

func listenToSignal() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return chSignal
}
