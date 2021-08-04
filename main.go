package main

import (
	"github.com/bitrise-io/go-utils/log"
	"os"
)

func mainR() error {

	return nil
}

func main() {
	if err := mainR(); err != nil {
		log.Errorf("%s", err)
		os.Exit(1)
	}
}
