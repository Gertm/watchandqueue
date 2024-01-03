package watchandqueue

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestSetPollInterval(t *testing.T) {
	type args struct {
		interval int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetPollInterval(tt.args.interval)
		})
	}
}

func TestWatchForIncomingFiles(t *testing.T) {
	type args struct {
		ctx            context.Context
		watchDirectory string
		extension      string
		c              chan<- string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WatchForIncomingFiles(tt.args.ctx, tt.args.watchDirectory, tt.args.extension, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("WatchForIncomingFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_waitForUploadToFinish(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := waitForUploadToFinish(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("waitForUploadToFinish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_RunWatchForIncomingFiles(t *testing.T) {
	ctx := context.Background()
	c := make(chan string, 5)
	go func() {
		for {
			file := <-c
			fmt.Println("Got:", file)
		}
	}()
	fmt.Println("Starting file watcher")
	err := WatchForIncomingFiles(ctx, "/tmp", ".mkv", c)
	if err != nil {
		log.Println(err)
	}

}
