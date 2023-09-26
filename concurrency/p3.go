package main

import (
	"bufio"
	"log/slog"
	"os"
	"sync"
)

// read file line by line, send data to siteChan
func getSite3(file string, siteChan chan string) error {
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

// TODO: 协程的数量太多，CPU和内存可能会不足
// consume data from siteChan, do digManage2 and send result to resultChan
func digManage3(siteChan chan string, resultChan chan string) {
	var wg sync.WaitGroup
	for site := range siteChan {
		wg.Add(1)
		go func(site string) {
			defer wg.Done()

			if err := dig(site); err != nil {
				return
			}

			resultChan <- site
		}(site)
	}

	wg.Wait()
}
