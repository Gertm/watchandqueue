# watchandqueue
Watch for files being written and queue them when they are ready.

When a file is written in the folder of your choosing, a string containing the filename will be sent on a channel.

```go
// The channel 'c' will receive the full filenames of the files that have been written. 

	ctx := context.Background()
	c := make(chan string, 5)
	go func() {
		for {
			file := <-c
			fmt.Println("Got:", file)
		}
	}()
	fmt.Println("Starting file watcher")
	err := WatchForIncomingFiles(ctx, "/tmp", ".log", c)
	if err != nil {
		log.Println(err)
	}

```
