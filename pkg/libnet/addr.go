package libnet

import (
	"context"
	"go-cli/pkg/libnet/networker"
	"net"

	"github.com/vishvananda/netlink"
)

func (s *server) ShowAddr(ctx context.Context, in *networker.AddrQuery) (*networker.AddrResponse, error) {
	linkList, err := netlink.LinkList()
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	addrList := make([]*networker.Addr, 0)
	for _, link := range linkList {
		addrs, err := netlink.AddrList(link, 0)
		if err != nil {
			logger.Warn("%v\n", err)
			return nil, err
		}

		for _, addr := range addrs {
			if len(addr.IP) == net.IPv6len {
				continue
			}
			logger.Info("%s %v", link.Attrs().Name, addr)

			if in.Name == "" || in.Name == link.Attrs().Name {
				addrList = append(addrList, &networker.Addr{
					Name:       link.Attrs().Name,
					IpWithMask: addr.IPNet.String(),
				})
			}
		}
	}

	return &networker.AddrResponse{Addrs: addrList}, err
}

func (s *server) AddAddr(ctx context.Context, in *networker.AddrQuery) (*networker.AddrResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	addr, err := netlink.ParseAddr(in.IpWithMask)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	err = netlink.AddrAdd(link, addr)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.AddrResponse{}, err
}

func (s *server) DelAddr(ctx context.Context, in *networker.AddrQuery) (*networker.AddrResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	addr, err := netlink.ParseAddr(in.IpWithMask)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	err = netlink.AddrDel(link, addr)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.AddrResponse{}, err
}
