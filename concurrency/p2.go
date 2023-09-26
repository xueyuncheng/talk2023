package main

import (
	"bufio"
	"log/slog"
	"os"
)

// read file line by line, send data to siteChan
func getSite2(file string, siteChan chan string) error {
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

// TODO: 缺少并发，可能会导致很耗时
// consume data from siteChan, do digManage2 and send result to resultChan
func digManage2(siteChan chan string, resultChan chan string) {
	for site := range siteChan {
		if err := dig(site); err != nil {
			continue
		}

		resultChan <- site
	}
}
