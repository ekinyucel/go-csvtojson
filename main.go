package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
)

var logger = log.New(os.Stdout, "converter: ", log.LstdFlags)
var fileList []File
var fileType = "csv"

func observeDirectory(c chan []File) {
	logger.Println("observing directory")

	cron := cron.New()
	cron.AddFunc("0 * * * *", func() { fileChanges(c) })
	cron.Start()
}

func fileChanges(c chan []File) {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if getInputFileFormat(file, fileType) {
			fileList = append(fileList, File{filename: file.Name(), processed: false})
		}
	}
	c <- fileList
}

func main() {
	fileChan := make(chan []File)

	go observeDirectory(fileChan)
	ctx := shutdown(context.Background())

	go func() {
		for i := range <-fileChan {
			go processFile(&fileList[i])
		}
	}()

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
