package main

import (
	"fmt"
	"os"

	"golang.org/x/mod/semver"
)

func main() {
	fmt.Println(semver.IsValid(os.Args[1]))
}
