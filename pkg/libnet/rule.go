package libnet

import (
	"context"
	"fmt"
	"go-cli/pkg/libnet/networker"
	"go-cli/pkg/libutil"
	"net"
	"strconv"

	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

// RulePortRange to string
func rulePortRangeToString(r *netlink.RulePortRange) string {
	if r == nil {
		return "any"
	}

	if r.Start == r.End {
		return strconv.Itoa(int(r.Start))
	} else {
		return fmt.Sprintf("%d-%d", r.Start, r.End)
	}
}

// ipNetToString convert IPNet to string
func ipNetToString(ipNet *net.IPNet) string {
	if ipNet == nil {
		return "any"
	}

	return ipNet.String()
}

// ip proto convert to string
func ipProtoToString(proto int) string {
	switch proto {
	case unix.IPPROTO_TCP:
		return "tcp"
	case unix.IPPROTO_UDP:
		return "udp"
	case unix.IPPROTO_ICMP:
		return "icmp"
	case unix.IPPROTO_ICMPV6:
		return "icmpv6"
	default:
		return "any"
	}
}

// parse port range
func parsePortRange(portRange string) (start, end uint16) {
	if portRange == "any" {
		return 0, 0
	}

	if _, err := fmt.Sscanf(portRange, "%d-%d", &start, &end); err != nil {
		if _, err := fmt.Sscanf(portRange, "%d", &start); err != nil {
			return 0, 0
		}
		end = start
	}

	return start, end
}

// show rule
func (s *server) ShowRule(ctx context.Context, in *networker.RuleQuery) (*networker.RuleResponse, error) {
	ruleListV4, err := netlink.RuleList(unix.AF_INET)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	ruleList := make([]*networker.Rule, 0)
	for _, rule := range ruleListV4 {
		table := libutil.UnixTableIdToString(rule.Table)

		netLogger.Info("%v", rule)
		if in.Table != "" && in.Table != table {
			continue
		}

		ruleList = append(ruleList, &networker.Rule{
			Priority: int32(rule.Priority),
			Table:    table,
			Src:      ipNetToString(rule.Src),
			Dst:      ipNetToString(rule.Dst),
			SPort:    rulePortRangeToString(rule.Sport),
			DPort:    rulePortRangeToString(rule.Dport),
			IpProto:  ipProtoToString(rule.IPProto),
			IIfName:  rule.IifName,
			OIfName:  rule.OifName,
		})

	}

	return &networker.RuleResponse{Rules: ruleList}, err
}

// add rule by query
func (s *server) AddRule(ctx context.Context, in *networker.RuleQuery) (*networker.RuleResponse, error) {
	rule := netlink.NewRule()
	rule.Table = libutil.StringToUnixTableId(in.Table)
	netLogger.Info("in %v", in)

	if in.Src != "any" {
		_, rule.Src, _ = net.ParseCIDR(in.Src)
	}

	if in.Dst != "any" {
		_, rule.Dst, _ = net.ParseCIDR(in.Dst)
	}

	if in.SPort != "any" {
		start, end := parsePortRange(in.SPort)
		rule.Sport = &netlink.RulePortRange{
			Start: start,
			End:   end,
		}
	}

	if in.DPort != "any" {
		start, end := parsePortRange(in.DPort)
		rule.Dport = &netlink.RulePortRange{
			Start: start,
			End:   end,
		}
	}

	if in.IpProto != "any" {
		switch in.IpProto {
		case "tcp":
			rule.IPProto = unix.IPPROTO_TCP
		case "udp":
			rule.IPProto = unix.IPPROTO_UDP
		case "icmp":
			rule.IPProto = unix.IPPROTO_ICMP
		case "icmpv6":
			rule.IPProto = unix.IPPROTO_ICMPV6
		default:
			rule.IPProto = 0
		}
	}

	err := netlink.RuleAdd(rule)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.RuleResponse{}, err
}

// del rule by query
func (s *server) DelRule(ctx context.Context, in *networker.RuleQuery) (*networker.RuleResponse, error) {
	ruleListV4, err := netlink.RuleList(netlink.FAMILY_V4)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	for _, rule := range ruleListV4 {
		table := libutil.StringToUnixTableId(in.Table)
		if table != rule.Table {
			continue
		}

		if in.Priority != 0 && in.Priority != int32(rule.Priority) {
			continue
		}

		if in.Src != "any" && in.Src != ipNetToString(rule.Src) {
			continue
		}

		if in.Dst != "any" && in.Dst != ipNetToString(rule.Dst) {
			continue
		}

		if in.SPort != "any" && in.SPort != rulePortRangeToString(rule.Sport) {
			continue
		}

		if in.DPort != "any" && in.DPort != rulePortRangeToString(rule.Dport) {
			continue
		}

		if in.IpProto != "any" && in.IpProto != ipProtoToString(rule.IPProto) {
			continue
		}

		err = netlink.RuleDel(&rule)
		if err != nil {
			netLogger.Warn("%v\n", err)
			return nil, err
		}
	}

	return &networker.RuleResponse{}, err
}
