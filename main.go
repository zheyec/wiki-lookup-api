package main

import (
	"context"
	"flag"
	"lxm-ency/services"
	"lxm-ency/utils"
	"os"
	"os/signal"
	"syscall"

	"fmt"
)

func main() {

	var port int
	flag.IntVar(&port, "p", 8000, "端口号")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		services.Run(ctx, port)
	}()

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		cancel()
		done <- true
	}()

	clean()

	<-done
	fmt.Println("[IAM]=>stop service.")
	clean()
}

// clean temp files
func clean() {
	utils.RemoveAllTemp(services.SynpPath)
	utils.RemoveAllTemp(utils.CardPath)
}
