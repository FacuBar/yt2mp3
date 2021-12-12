package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/FacuBar/yt2mp3"
)

func main() {
	srv := yt2mp3.NewServer()

	go srv.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("Shutdown server ...")
}
