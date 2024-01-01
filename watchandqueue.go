package watchandqueue

import (
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func watchForFiles(watchDirectory string, f func() error) {
	count = 0
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
	go func() {
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
						return
					} else {
						monitored_files.Store(event.Name, true)
						count++
					}
					go func() {
						if strings.HasSuffix(strings.ToLower(event.Name), ".mkv") {
							log.Printf("A new file is being written: %v\n", event.Name)
							err := waitForUploadToFinish(event.Name)
							if err != nil {
								// Could not wait for the file correctly. Something must have gone awry.
								log.Println(err)
							} else {
								count -= 1
								if count <= 0 {
									log.Println("No more files being written to.")
									if err := f(); err != nil {
										log.Println(err)
									}
								} else {
									log.Println("Some files are still being written to...", count)
								}
							}
						}
					}()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	log.Printf("Watching directory %v for incoming .mkv files.\n", dir)
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
