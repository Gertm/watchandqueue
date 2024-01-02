package watchandqueue

import (
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	count           int
	monitored_files sync.Map
	pollInterval    = 3
)

func SetPollInterval(interval int) {
	pollInterval = interval
}

func WatchForIncomingFiles(ctx context.Context, watchDirectory, extension string, c chan<- string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var dir string

	if watchDirectory == "" {
		dir, _ = os.Getwd()
	} else {
		dir = watchDirectory
	}

	done := make(chan bool)
	go func(ctx context.Context) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					if _, ok := monitored_files.Load(event.Name); ok {
						// already watching this file
						log.Printf("Already watching %v\n", event.Name)
						continue
					} else {
						monitored_files.Store(event.Name, true)
						count++
					}
					go func() {
						if strings.HasSuffix(strings.ToLower(event.Name), ".mkv") {
							log.Printf("A new file is being written: %v\n", event.Name)
							err := waitForUploadToFinish(event.Name)
							if err != nil {
								log.Println(err)
							} else {
								c <- event.Name
							}
						}
					}()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				log.Println("Being cancelled by context.")
				return
			}

		}
	}(ctx)

	log.Printf("Watching directory %v for incoming .mkv files.\n", dir)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return nil
}

// This is a bit of a naive way of checking if the file is done writing.
// Yet it works quite well in practise for me. Then again, I have quite
// reliable internet, so that helps. So this can certainly be improved.
func waitForUploadToFinish(file string) error {
	var size int64
	size = 0
	sameSizeCount := 0
	log.Printf("Waiting for write operations to stop on %v\n", file)
	defer func() {
		count--
		monitored_files.Delete(file)
	}()
	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		fi, err := os.Stat(file)
		if err != nil {
			return err
		}
		currentSize := fi.Size()
		if currentSize > size {
			size = currentSize
			sameSizeCount = 0
			continue
		}
		if currentSize == fi.Size() {
			sameSizeCount += 1
		}
		if sameSizeCount == 3 {
			return nil
		}
	}
}
