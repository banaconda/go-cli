package main

import (
	"fmt"
	"go-cli/pkg/cli"
	"os"

	nblogger "github.com/banaconda/nb-logger"
)

func main() {
	cliLogPath := "log/cli.log"
	if _, err := os.Stat("log/"); os.IsNotExist(err) {
		os.MkdirAll("log", 0700)
	}

	cliLogger, err := nblogger.NewLogger(cliLogPath, nblogger.Info, 1000, nblogger.LstdFlags)
	if err != nil {
		fmt.Printf("logger init fail: %v", err)
	}

	cli := cli.GoCli{}
	cli.Init(cliLogger)

	cli.Run()
}
