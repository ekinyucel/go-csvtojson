package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "net/http/pprof"
)

const (
	// CSV type indicates the file format
	CSV string = "csv"
	// JSON type indicates the file format
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

	tick := time.Tick(1000 * time.Millisecond)

	for range tick {
		go trackFiles()
	}

	/*cron := cron.New()
	cron.AddFunc("0 * * * *", func() {
		go trackFiles()
	})
	cron.Start()*/
}

func main() {
	flag.StringVar(&fileType, "filetype", CSV, "input file format")
	flag.StringVar(&targetType, "targetType", JSON, "target file format")
	flag.StringVar(&folderName, "folder", "C:\\Users\\user\\Desktop\\", "folder name")
	flag.Parse()

	// for profiling
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	// for profiling

	go observeDirectory()

	go func() {
		for {
			select {
			case fileList := <-fileChannel:
				for i := range fileList {
					if !fileList[i].processed {
						processFile(&fileList[i]) // decide whether using go keyword here or not
					}
				}
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
