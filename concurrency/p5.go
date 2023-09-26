package main

import (
	"bufio"
	"context"
	"log/slog"
	"os"
	"runtime"

	"golang.org/x/sync/semaphore"
)

// read file line by line, send data to siteChan
func getSite5(file string, siteChan chan string) error {
	f, err := os.Open(file)
	if err != nil {
		slog.Error("os.Open() error", "err", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		siteChan <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		slog.Error("scanner.Err() error", "err", err)
		os.Exit(1)
	}

	return nil
}

// consume data from siteChan, do digManage2 and send result to resultChan
func digManage5(siteChan chan string, resultChan chan string) {
	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	for site := range siteChan {
		if err := sem.Acquire(context.TODO(), 1); err != nil {
			slog.Error("sem.Acquire() error", "err", err)
			return
		}

		go func(site string) {
			defer sem.Release(1)

			if err := dig(site); err != nil {
				return
			}

			resultChan <- site
		}(site)
	}

	if err := sem.Acquire(context.TODO(), int64(runtime.GOMAXPROCS(0))); err != nil {
		slog.Error("sem.Acquire() error", "err", err)
		return
	}
}
