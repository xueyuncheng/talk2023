package main

import (
	"bufio"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

// read file line by line, send data to siteChan
func getSite4(file string, siteChan chan string) error {
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

// TODO: 同时使用channel和WaitGroup,感觉代码优点多，能不能更简单一点。 控制同时运行的协程数量，和semaphore类似
// consume data from siteChan, do digManage2 and send result to resultChan
func digManage4(siteChan chan string, resultChan chan string) {
	ch := make(chan struct{}, runtime.GOMAXPROCS(0))

	var wg sync.WaitGroup
	for site := range siteChan {
		ch <- struct{}{}
		wg.Add(1)

		go func(site string) {
			defer func() { <-ch }()
			defer wg.Done()

			if err := dig(site); err != nil {
				return
			}

			resultChan <- site
		}(site)
	}

	wg.Wait()
}
