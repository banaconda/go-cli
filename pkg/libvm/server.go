package libvm

import (
	"fmt"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
	"log"
	"net"
	"os"

	nblogger "github.com/banaconda/nb-logger"
	"google.golang.org/grpc"
)

var logger nblogger.Logger
var vmerDB *VmerDB

type server struct {
	vmer.UnimplementedVmerServer
}

func handlerRequests(host string, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Error("failed to listen: %v", err)
		panic(err)
	}

	s := grpc.NewServer()

	vmer.RegisterVmerServer(s, &server{})
	logger.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve: %v", err)
		panic(err)
	}
}

func VmServer() {
	logPath := "log/vm.log"
	if _, err := os.Stat("log/"); os.IsNotExist(err) {
		os.MkdirAll("log", 0700)
	}

	var err error
	logger, err = nblogger.NewLogger(logPath, nblogger.Info, 1000, nblogger.LstdFlags|nblogger.Lshortfile|nblogger.Lmicroseconds)
	if err != nil {
		log.Fatalf("logger init fail: %v", err)
	}

	vmerDB, err = NewVmerDB("local.db")
	if err != nil {
		logger.Error("failed to open vmer db: %v", err)
		return
	}

	handlerRequests("", libutil.VM_PORT)
}
