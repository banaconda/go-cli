package main

import (
	"fmt"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libnet"
	"os"

	nblogger "github.com/banaconda/nb-logger"
)

func main() {
	go libnet.NetServer()

	cliLogPath := "log/cli.log"
	if _, err := os.Stat("log/"); os.IsNotExist(err) {
		os.MkdirAll("log", 0700)
	}

	cliLogger, err := nblogger.NewLogger(cliLogPath, nblogger.Info, 1000,
		nblogger.LstdFlags|nblogger.Lshortfile|nblogger.Lmicroseconds)
	if err != nil {
		fmt.Printf("logger init fail: %v", err)
	}

	cli := libcli.GoCli{}
	cli.Init(cliLogger)
	libnet.InitCli(&cli)

	cli.Run()
}
