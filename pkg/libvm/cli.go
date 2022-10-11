package libvm

import (
	"fmt"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var client vmer.VmerClient

var nce = libcli.NewCommandElemWithoutFunc
var ncef = libcli.NewCommandElem

type StreamInterface[V ValueType] interface {
	Recv() (V, error)
	grpc.ClientStream
}

type ValueType interface {
	*vmer.BaseImageMessage | *vmer.KeyMessage | *vmer.NetworkMessage | *vmer.VolumeMessage | *vmer.DomainMessage
}

func recvStream[V ValueType](stream StreamInterface[V]) ([]V, error) {
	var messages []V
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			logger.Warn("%v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func InitCli(cli *libcli.GoCli) {
	var err error
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", "localhost", libutil.VM_PORT),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client = vmer.NewVmerClient(conn)

	initBaseImageCli(cli)
	initNetworkCli(cli)
	initKeyCli(cli)
	initVolumeCli(cli)
	initDomainCli(cli)
}
