package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
)

var logger = log.New(os.Stdout, "main: ", log.LstdFlags)
var fileList []File
var fileChannel = make(chan []File)
var fileType = "csv"

//var fileChan = make(chan []File)

func observeDirectory() {
	logger.Println("observing directory")

	cron := cron.New()
	cron.AddFunc("0 * * * *", func() {
		go trackFiles()
		go func() {
			for i := range <-fileChannel {
				if !fileList[i].processed {
					go processFile(&fileList[i])
				}
			}
		}()
	})
	cron.Start()
}

func main() {
	// fileChan := make(chan []File)

	go observeDirectory()
	ctx := shutdown(context.Background())

	<-ctx.Done()
	for i := range fileList {
		logger.Println("after stopping the app ", &fileList[i])
	}
}

func shutdown(ctx context.Context) context.Context {
	ctx, done := context.WithCancel(ctx) // creating a context to handle graceful shutdown, cancelling requests

	quit := make(chan os.Signal, 1)                      // quit channel
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // when pressing CTRL+C, it passes value to quit channel

	go func() {
		defer done() // deferring the call done until the end of this function

		<-quit // start executing the shutdown logic
		signal.Stop(quit)
		close(quit) // closing the quit channel. it means that no more value will be sent through this channel

		logger.Printf("Application is shutting down \n")

		// it is not a server yet, so this block is omitted
		/*ctx, cancel := context.WithTimeout(ctx, 10*time.Second) // setting a deadline on requests to backend server
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the application: %s\n", err)
		}*/
	}()

	return ctx
}
