package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron"
)

const (
	// CSV type indicates file format
	CSV string = "csv"
	// JSON type indicates file format
	JSON string = "json"
)

var logger = log.New(os.Stdout, "main: ", log.LstdFlags)
var fileList []File
var fileChannel = make(chan []File)
var fileType string
var targetType string
var folderName string

func observeDirectory() {
	logger.Printf("observing this directory %s", folderName)

	cron := cron.New()
	cron.AddFunc("0 * * * *", func() {
		go trackFiles()
	})
	cron.Start()
}

func main() {
	flag.StringVar(&fileType, "filetype", CSV, "input file format")
	flag.StringVar(&targetType, "targetType", JSON, "target file format")
	flag.StringVar(&folderName, "folder", "C:\\Users\\user\\Desktop\\", "folder name")
	flag.Parse()

	go observeDirectory()

	go func() {
		for i := range <-fileChannel {
			if !fileList[i].processed {
				go processFile(&fileList[i])
			}
		}
	}()

	ctx := shutdown(context.Background())

	<-ctx.Done()
}

func shutdown(ctx context.Context) context.Context {
	ctx, done := context.WithCancel(ctx)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer done()

		<-quit
		signal.Stop(quit)
		close(quit)

		logger.Println("Application is shutting down")
	}()

	return ctx
}
