package libvm

import (
	"context"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
	"strconv"
)

func initNetworkCli(cli *libcli.GoCli) {
	// show network
	cli.AddCommandElem(
		nce("network", "network"),
		ncef("show", "show network", func(args []string) {
			stream, err := client.ShowNetworks(context.Background(), &vmer.NetworkMessage{})
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			var streamInterface StreamInterface[*vmer.NetworkMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// show network by name
	cli.AddCommandElem(
		nce("network", "network"),
		nce("show", "show network"),
		nce("name", "network name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			stream, err := client.ShowNetworks(context.Background(), &vmer.NetworkMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
			var streamInterface StreamInterface[*vmer.NetworkMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// create network
	cli.AddCommandElem(
		nce("network", "network"),
		nce("create", "create network"),
		nce("name", "network name"),
		nce(libutil.NameRegex, ""),
		nce("vlan", "vlan"),
		nce(libutil.NumberRegex, ""),
		nce("cidr", "cidr"),
		nce(libutil.CidrRegex, ""),
		nce("gateway", "gateway"),
		nce(libutil.IpRegex, ""),
		nce("dns", "dns"),
		ncef(libutil.IpRegex, "", func(args []string) {
			vlanId, err := strconv.ParseInt(args[5], 10, 32)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			_, err = client.CreateNetwork(context.Background(), &vmer.NetworkMessage{
				Name:    args[3],
				Vlan:    int32(vlanId),
				Cidr:    args[7],
				Gateway: args[9],
				Dns:     args[11],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// delete network
	cli.AddCommandElem(
		nce("network", "network"),
		nce("delete", "delete network"),
		nce("name", "network name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			_, err := client.DeleteNetwork(context.Background(), &vmer.NetworkMessage{
				Name: args[3],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))
}

// show network
func (s *server) ShowNetworks(in *vmer.NetworkMessage, stream vmer.Vmer_ShowNetworksServer) error {
	var networks []Network

	var err error
	if in.Name == "" {
		networks, err = vmerDB.GetAllNetworks()
	} else {
		network, err := vmerDB.GetNetworkByName(in.Name)
		if err != nil {
			return err
		}
		networks = append(networks, *network)
	}

	if err != nil {
		logger.Warn("failed to get base images: %v", err)
		return err
	}

	for _, network := range networks {
		if err := stream.Send(&vmer.NetworkMessage{
			Name:    network.Name,
			Vlan:    network.Vlan,
			Cidr:    network.Cidr,
			Gateway: network.Gateway,
			Dns:     network.Dns,
		}); err != nil {
			return err
		}
	}

	return nil

}

// create network
func (s *server) CreateNetwork(ctx context.Context, in *vmer.NetworkMessage) (*vmer.NetworkMessage, error) {
	network := Network{
		Name:    in.Name,
		Vlan:    in.Vlan,
		Cidr:    in.Cidr,
		Gateway: in.Gateway,
		Dns:     in.Dns,
	}

	if err := vmerDB.InsertNetwork(&network); err != nil {
		logger.Warn("failed to create network: %v", err)
		return nil, err
	}

	return in, nil
}

// delete network
func (s *server) DeleteNetwork(ctx context.Context, in *vmer.NetworkMessage) (*vmer.NetworkMessage, error) {
	// get network
	network, err := vmerDB.GetNetworkByName(in.Name)
	if err != nil {
		logger.Warn("failed to get network: %v", err)
		return nil, err
	}

	if err := vmerDB.DeleteNetwork(network); err != nil {
		logger.Warn("failed to delete network: %v", err)
		return nil, err
	}

	return in, nil
}
