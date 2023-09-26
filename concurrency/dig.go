package main

import (
	"fmt"
	"net"
	"strings"

	"log/slog"
)

var ErrNotCloudflare = fmt.Errorf("not cloudflare nameserver")

// 查询域名的DNS记录，并判断ns是否为cloudflare
func dig(site string) error {
	nss, err := net.LookupNS(site)
	if err != nil {
		slog.Error("net.LookupNS() error", "err", err)
		return fmt.Errorf("net.LookupNS() error: %v", err)
	}

	for _, ns := range nss {
		if strings.HasSuffix(ns.Host, "cloudflare.com") {
			slog.Info("cloudflare", "site", site)
			return nil
		}
	}

	return ErrNotCloudflare
}
