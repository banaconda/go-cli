package libnet

import (
	"context"
	"go-cli/pkg/libnet/networker"
	"go-cli/pkg/libutil"
	"net"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// show route
func (s *server) ShowRoute(ctx context.Context, in *networker.RouteQuery) (*networker.RouteResponse, error) {
	linkList, err := netlink.LinkList()
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	routeList := make([]*networker.Route, 0)
	for _, link := range linkList {
		routes, err := netlink.RouteListFiltered(unix.AF_INET,
			&netlink.Route{LinkIndex: link.Attrs().Index}, netlink.RT_FILTER_OIF|netlink.RT_FILTER_TABLE)
		if err != nil {
			logger.Warn("%v\n", err)
			return nil, err
		}

		for _, route := range routes {
			logger.Info("%s %v %v %v", link.Attrs().Name, route, route.Protocol, route.Scope)

			if (route.Gw != nil && len(route.Gw) == net.IPv6len) || (route.Src != nil &&
				len(route.Src) == net.IPv6len) || (route.Dst != nil && len(route.Dst.IP) == net.IPv6len) ||
				route.Protocol == unix.RTPROT_KERNEL || in.Table != "" && route.Table != libutil.StringToUnixTableId(in.Table) {
				continue
			}

			destination := "default"
			source := "any"
			gateway := "any"

			if route.Dst != nil {
				destination = route.Dst.String()
			}
			if route.Src != nil {
				source = route.Src.String()
			}
			if route.Gw != nil {
				gateway = route.Gw.String()
			}

			table := libutil.UnixTableIdToString(route.Table)

			if in.Device == "" || in.Device == link.Attrs().Name {

				routeList = append(routeList, &networker.Route{
					Protocol:    route.Protocol.String(),
					Table:       table,
					Destination: destination,
					Source:      source,
					NextHop:     gateway,
					Device:      link.Attrs().Name,
				})
			}
		}

	}

	return &networker.RouteResponse{Routes: routeList}, err
}

// add route
func (s *server) AddRoute(ctx context.Context, in *networker.RouteQuery) (*networker.RouteResponse, error) {
	var protocol = unix.RTPROT_STATIC
	var dst *net.IPNet = nil
	var nextHop net.IP = nil
	var src net.IP = nil

	table := libutil.StringToUnixTableId(in.Table)
	if in.Destination != "" {
		_, dst, _ = net.ParseCIDR(in.Destination)
	}
	if in.Source != "" {
		src = net.ParseIP(in.Source)
	}
	if in.NextHop != "" {
		nextHop = net.ParseIP(in.NextHop)
	}

	route := &netlink.Route{
		Protocol: netlink.RouteProtocol(protocol),
		Table:    table,
		Dst:      dst,
		Src:      src,
		Gw:       nextHop,
	}

	err := netlink.RouteAdd(route)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.RouteResponse{}, err
}

// del route
func (s *server) DelRoute(ctx context.Context, in *networker.RouteQuery) (*networker.RouteResponse, error) {
	var protocol = unix.RTPROT_STATIC
	var dst *net.IPNet = nil
	var nextHop net.IP = nil
	var src net.IP = nil

	table := libutil.StringToUnixTableId(in.Table)
	if in.Destination != "" {
		_, dst, _ = net.ParseCIDR(in.Destination)
	}
	if in.Source != "" {
		src = net.ParseIP(in.Source)
	}
	if in.NextHop != "" {
		nextHop = net.ParseIP(in.NextHop)
	}

	route := &netlink.Route{
		Protocol: netlink.RouteProtocol(protocol),
		Table:    table,
		Dst:      dst,
		Src:      src,
		Gw:       nextHop,
	}

	err := netlink.RouteDel(route)
	if err != nil {
		logger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.RouteResponse{}, err
}
