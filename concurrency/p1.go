package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func main() {

}

// TODO: 大文件可能会导致内存不足
func getSite(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		slog.Error("os.ReadFile() error", "err", err)
		return nil, fmt.Errorf("os.ReadFile() error: %v", err)
	}

	return strings.Split(string(data), "\n"), nil
}

func digManage(sites []string) []string {
	var cloudflareSites []string
	for _, site := range sites {
		if err := dig(site); err != nil {
			continue
		}

		cloudflareSites = append(cloudflareSites, site)
	}

	return cloudflareSites
}
