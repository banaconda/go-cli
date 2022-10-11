package libnet

import (
	"fmt"
	"go-cli/pkg/libnet/networker"
	"go-cli/pkg/libutil"
	"log"
	"net"
	"os"

	nblogger "github.com/banaconda/nb-logger"
	"google.golang.org/grpc"
)

var logger nblogger.Logger

type server struct {
	networker.UnimplementedNetworkerServer
}

func handlerRequests(host string, port int) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		logger.Error("failed to listen: %v", err)
		panic(err)
	}

	s := grpc.NewServer()

	networker.RegisterNetworkerServer(s, &server{})
	logger.Info("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve: %v", err)
		panic(err)
	}
}

func NetServer() {
	netLogPath := "log/net.log"
	if _, err := os.Stat("log/"); os.IsNotExist(err) {
		os.MkdirAll("log", 0700)
	}

	var err error
	logger, err = nblogger.NewLogger(netLogPath, nblogger.Info, 1000, nblogger.LstdFlags|nblogger.Lshortfile|nblogger.Lmicroseconds)
	if err != nil {
		log.Fatalf("logger init fail: %v", err)
	}
	handlerRequests("", libutil.NET_PORT)
}
