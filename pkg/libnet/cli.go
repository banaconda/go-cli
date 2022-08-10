package libnet

import (
	"context"
	"fmt"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libnet/networker"
	"go-cli/pkg/libutil"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var conn *grpc.ClientConn
var client networker.NetworkerClient
var nce = libcli.NewCommandElemWithoutFunc
var ncef = libcli.NewCommandElem

func InitCli(cli *libcli.GoCli) {
	var err error
	conn, err = grpc.Dial("localhost:10000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client = networker.NewNetworkerClient(conn)

	initCliLink(cli)
	initCliAddr(cli)
	initCliRule(cli)
	initCliRoute(cli)
}

type networkerQuery interface {
	// LINK
	*networker.NetLinkQuery | *networker.BridgeQuery |
		*networker.VethQuery | *networker.VlanQuery |
		// ADDR
		*networker.AddrQuery |
		// RULE
		*networker.RuleQuery |
		// ROUTE
		*networker.RouteQuery
}

type networkerReponse interface {
	// LINK
	*networker.NetLinkResponse |
		// ADDR
		*networker.AddrResponse |
		// RULE
		*networker.RuleResponse |
		// ROUTE
		*networker.RouteResponse
}

type queryInterface[Q networkerQuery, R networkerReponse] func(context.Context, Q, ...grpc.CallOption) (R, error)

func query[Q networkerQuery, R networkerReponse](f queryInterface[Q, R], queryElem Q) (R, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := f(ctx, queryElem)
	if err != nil {
		netLogger.Warn("%v", err)
	} else {
		netLogger.Info("%v", r)
	}

	return r, err
}

func initCliLink(cli *libcli.GoCli) {
	// show all links
	cli.AddCommandElem(
		nce("link", ""),
		ncef("show", "show all links", func(args []string) {
			resp, err := query(client.ShowNetLink, &networker.NetLinkQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// set link mac by link name
	cli.AddCommandElem(
		nce("link", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "link name"),
		nce("mac", ""),
		nce("set", "set link mac by link name"),
		ncef(libutil.MacRegex, "mac address", func(args []string) {
			resp, err := query(client.SetNetLinkMac, &networker.NetLinkQuery{
				Name: args[2],
				Mac:  args[5],
			})

			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// set link up by link name
	cli.AddCommandElem(
		nce("link", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "link name"),
		ncef("up", "", func(args []string) {
			resp, err := query(client.SetNetLinkUp, &networker.NetLinkQuery{
				Name: args[2],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// set link down by link name
	cli.AddCommandElem(
		nce("link", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "link name"),
		ncef("down", "", func(args []string) {
			resp, err := query(client.SetNetLinkDown, &networker.NetLinkQuery{
				Name: args[2],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// show all bridges
	cli.AddCommandElem(
		nce("bridge", ""),
		ncef("show", "show all bridges", func(args []string) {
			resp, err := query(client.ShowBridge, &networker.BridgeQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// show bridge slaves by bridge name
	cli.AddCommandElem(
		nce("bridge", ""),
		nce("show", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "bridge name"),
		ncef("slave", "show bridge slaves by bridge name", func(args []string) {
			resp, err := query(client.ShowBridgeSlave, &networker.BridgeQuery{
				Name: args[3],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// set bridge master by bridge name and slave name
	cli.AddCommandElem(
		nce("bridge", ""),
		nce("set", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "bridge name"),
		nce("slave", ""),
		ncef(libutil.StringRegex, "slave name", func(args []string) {
			resp, err := query(client.SetBridgeMaster, &networker.BridgeQuery{
				Name:      args[3],
				SlaveName: args[6],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// unset bridge master by bridge name and slave name
	cli.AddCommandElem(
		nce("bridge", ""),
		nce("unset", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "bridge name"),
		nce("slave", ""),
		ncef(libutil.StringRegex, "slave name", func(args []string) {
			resp, err := query(client.UnsetBridgeMaster, &networker.BridgeQuery{
				Name:      args[3],
				SlaveName: args[6],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// add bridge by bridge name
	cli.AddCommandElem(
		nce("bridge", ""),
		nce("add", ""),
		nce("name", ""),
		ncef(libutil.StringRegex, "bridge name", func(args []string) {
			resp, err := query(client.AddBridge, &networker.BridgeQuery{
				Name: args[2],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// del bridge by bridge name
	cli.AddCommandElem(
		nce("bridge", ""),
		nce("del", ""),
		nce("name", ""),
		ncef(libutil.StringRegex, "bridge name", func(args []string) {
			resp, err := query(client.DelBridge, &networker.BridgeQuery{
				Name: args[2],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// show all veths
	cli.AddCommandElem(
		nce("veth", ""),
		ncef("show", "show all veths", func(args []string) {
			resp, err := query(client.ShowVeth, &networker.VethQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// add veth by veth name and peer name
	cli.AddCommandElem(
		nce("veth", ""),
		nce("add", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "veth name"),
		nce("peer", ""),
		nce("name", ""),
		ncef(libutil.StringRegex, "peer name", func(args []string) {
			resp, err := query(client.AddVeth, &networker.VethQuery{
				Name:     args[2],
				PeerName: args[5],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))
	// del veth by veth name
	cli.AddCommandElem(
		nce("veth", ""),
		nce("del", ""),
		nce("name", ""),
		ncef(libutil.StringRegex, "veth name", func(args []string) {
			resp, err := query(client.DelVeth, &networker.VethQuery{
				Name: args[2],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// show all vlans
	cli.AddCommandElem(
		nce("vlan", ""),
		ncef("show", "show all vlans", func(args []string) {
			resp, err := query(client.ShowVlan, &networker.VlanQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// show vlan by vlan id
	cli.AddCommandElem(
		nce("vlan", ""),
		nce("show", ""),
		nce("id", ""),
		ncef(libutil.NumberRegex, "show vlan by vlan id", func(args []string) {
			vlanId, _ := strconv.Atoi(args[2])
			resp, err := query(client.ShowVlan, &networker.VlanQuery{
				VlanId: int32(vlanId),
			})
			if err != nil {
				fmt.Printf("%v", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// add vlan by vlan id and parent name
	cli.AddCommandElem(
		nce("vlan", ""),
		nce("add", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "vlan name"),
		nce("parent", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "parent name"),
		nce("id", ""),
		ncef(libutil.NumberRegex, "vlan id", func(args []string) {
			vlanId, _ := strconv.Atoi(args[8])
			resp, err := query(client.AddVlan, &networker.VlanQuery{
				Name:       args[3],
				ParentName: args[6],
				VlanId:     int32(vlanId),
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))

	// del vlan by vlan name
	cli.AddCommandElem(
		nce("vlan", ""),
		nce("del", ""),
		nce("name", ""),
		ncef(libutil.StringRegex, "vlan name", func(args []string) {
			resp, err := query(client.DelVlan, &networker.VlanQuery{
				Name: args[2],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.NetLinks)
		}))
}

func initCliAddr(cli *libcli.GoCli) {
	// show all addresses
	cli.AddCommandElem(
		nce("addr", ""),
		ncef("show", "show all addresses", func(args []string) {
			resp, err := query(client.ShowAddr, &networker.AddrQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Addrs)
		}))

	// show address by name
	cli.AddCommandElem(
		nce("addr", ""),
		nce("show", ""),
		nce("name", ""),
		ncef(libutil.StringRegex, "address name", func(args []string) {
			resp, err := query(client.ShowAddr, &networker.AddrQuery{
				Name: args[3],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Addrs)
		}))

	// add ip with mask by name
	cli.AddCommandElem(
		nce("addr", ""),
		nce("add", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "address name"),
		nce("ipWithMask", ""),
		ncef(libutil.CidrRegex, "ip with mask", func(args []string) {
			resp, err := query(client.AddAddr, &networker.AddrQuery{
				Name:       args[3],
				IpWithMask: args[5],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Addrs)
		}))

	// del ip with mask by name
	cli.AddCommandElem(
		nce("addr", ""),
		nce("del", ""),
		nce("name", ""),
		nce(libutil.StringRegex, "address name"),
		nce("ipWithMask", ""),
		ncef(libutil.CidrRegex, "ip with mask", func(args []string) {
			resp, err := query(client.DelAddr, &networker.AddrQuery{
				Name:       args[3],
				IpWithMask: args[5],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Addrs)
		}))
}

type nameRegex struct {
	Name  string
	Desc  string
	Regex string
}

// get combinations of name and regex
func getCombinations(args []nameRegex) [][]nameRegex {
	var combinations [][]nameRegex
	for i := 0; i < len(args); i++ {
		var newCombinations [][]nameRegex
		if len(combinations) == 0 {
			newCombinations = append(newCombinations, []nameRegex{args[i]})
		} else {
			for _, combination := range combinations {
				newCombinations = append(newCombinations, append(combination, args[i]))
			}
		}
		combinations = append(combinations, newCombinations...)
	}
	return combinations
}

func ruleAddCombination(cli *libcli.GoCli, combinationArgs []nameRegex) {
	//get combination of args

	allCombination := make([][]nameRegex, 0)

	for i := 0; i < len(combinationArgs); i++ {
		combination := getCombinations(combinationArgs[i:])

		allCombination = append(allCombination, combination...)
	}
	queryFunc := func(args []string) {
		src := "any"
		dst := "any"
		sPort := "any"
		dPort := "any"
		proto := "any"
		priority := 0

		for i := 4; i < len(args); i += 2 {
			switch args[i] {
			case "src":
				src = args[i+1]
			case "dst":
				dst = args[i+1]
			case "sPort":
				sPort = args[i+1]
			case "dPort":
				dPort = args[i+1]
			case "proto":
				proto = args[i+1]
			case "priority":
				priority, _ = strconv.Atoi(args[i+1])
			}
		}

		resp, err := query(client.AddRule, &networker.RuleQuery{
			Table:    args[3],
			Priority: int32(priority),
			Src:      src,
			Dst:      dst,
			SPort:    sPort,
			DPort:    dPort,
			IpProto:  proto,
		})
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		libutil.PrintStructAll(resp.Rules)
	}

	for _, combination := range allCombination {
		var args []string
		var ceArgs []*libcli.CommandElem

		ceArgs = append(ceArgs,
			nce("rule", ""),
			nce("add", "add rule by 5 tuple"),
			nce("table", ""),
			nce(libutil.TableRegex, "table name or number"))

		for _, eachNameRegex := range combination {
			args = append(args, eachNameRegex.Name, libutil.GetRegexHelpString(eachNameRegex.Regex))
			ceArgs = append(ceArgs,
				nce(eachNameRegex.Name, eachNameRegex.Desc),
				nce(eachNameRegex.Regex, ""))
		}
		ceArgs[len(ceArgs)-1].Func = queryFunc

		netLogger.Info("add rule %v", args)

		cli.AddCommandElem(
			ceArgs...)

	}
	netLogger.Info("len %d", len(allCombination))
}

func initCliRule(cli *libcli.GoCli) {
	// show all rules
	cli.AddCommandElem(
		nce("rule", ""),
		ncef("show", "show all rules", func(args []string) {
			resp, err := query(client.ShowRule, &networker.RuleQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Rules)
		}))

	// show rule by table
	cli.AddCommandElem(
		nce("rule", ""),
		nce("show", ""),
		nce("table", ""),
		ncef(libutil.TableRegex, "table name or num", func(args []string) {
			resp, err := query(client.ShowRule, &networker.RuleQuery{
				Table: args[3],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Rules)
		}))

	ruleAddCombination(cli, []nameRegex{
		{
			Name:  "src",
			Desc:  "source ip",
			Regex: libutil.CidrRegex,
		},
		{
			Name:  "dst",
			Desc:  "destination cidr",
			Regex: libutil.CidrRegex,
		},
		{
			Name:  "sPort",
			Desc:  "source port",
			Regex: libutil.PortRegex,
		},
		{
			Name:  "dPort",
			Desc:  "destination port",
			Regex: libutil.PortRegex,
		},
		{
			Name:  "proto",
			Desc:  "ip protocol",
			Regex: libutil.ProtoRegex,
		},
	})

	// del rule by table
	cli.AddCommandElem(
		nce("rule", ""),
		nce("del", "delete rule by table id or table id and priority"),
		nce("table", ""),
		ncef(libutil.TableRegex, "table name or number", func(args []string) {
			resp, err := query(client.DelRule, &networker.RuleQuery{
				Table:    args[3],
				Priority: 0,
				Src:      "any",
				Dst:      "any",
				SPort:    "any",
				DPort:    "any",
				IpProto:  "any",
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Rules)
		}))

	// del rule by table and priority
	cli.AddCommandElem(
		nce("rule", ""),
		nce("del", ""),
		nce("table", ""),
		nce(libutil.TableRegex, "table name or number"),
		nce("priority", "priority which whil be deleted"),
		ncef(libutil.NumberRegex, "priority", func(args []string) {
			priority, _ := strconv.Atoi(args[5])
			resp, err := query(client.DelRule, &networker.RuleQuery{
				Table:    args[3],
				Priority: int32(priority),
				Src:      "any",
				Dst:      "any",
				SPort:    "any",
				DPort:    "any",
				IpProto:  "any",
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Rules)
		}))
}

func initCliRoute(cli *libcli.GoCli) {
	// show route
	cli.AddCommandElem(
		nce("route", ""),
		ncef("show", "show all routes", func(args []string) {
			resp, err := query(client.ShowRoute, &networker.RouteQuery{})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// show route by table
	cli.AddCommandElem(
		nce("route", ""),
		nce("show", ""),
		nce("table", ""),
		ncef(libutil.TableRegex, "table name or num", func(args []string) {
			resp, err := query(client.ShowRoute, &networker.RouteQuery{
				Table: args[3],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// add route by destination and nexthop
	cli.AddCommandElem(
		nce("route", ""),
		nce("add", "add route"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.AddRoute, &networker.RouteQuery{
				Table:       "main",
				Protocol:    "static",
				Destination: args[3],
				NextHop:     args[5],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// add route by destination and nexthop and table
	cli.AddCommandElem(
		nce("route", ""),
		nce("add", "add route"),
		nce("table", ""),
		nce(libutil.TableRegex, "table number"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.AddRoute, &networker.RouteQuery{
				Table:       args[3],
				Protocol:    "static",
				Destination: args[5],
				NextHop:     args[7],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// add route by destination, source and nexthop
	cli.AddCommandElem(
		nce("route", ""),
		nce("add", "add route"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("src", "source ip"),
		nce(libutil.CidrRegex, "source ip"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.AddRoute, &networker.RouteQuery{
				Table:       "main",
				Protocol:    "static",
				Destination: args[3],
				Source:      args[5],
				NextHop:     args[7],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// add route by destination, source, nexthop and table
	cli.AddCommandElem(
		nce("route", ""),
		nce("add", "add route"),
		nce("table", ""),
		nce(libutil.TableRegex, "table number"),
		nce("src", "source ip"),
		nce(libutil.CidrRegex, "source ip"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("source", "source ip"),
		nce(libutil.IpRegex, "source ip"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.AddRoute, &networker.RouteQuery{
				Table:       args[3],
				Protocol:    "static",
				Destination: args[5],
				Source:      args[7],
				NextHop:     args[9],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// del route by destination
	cli.AddCommandElem(
		nce("route", ""),
		nce("del", "delete route"),
		nce("dst", "destination cidr"),
		ncef(libutil.CidrRegex, "destination cidr", func(args []string) {
			resp, err := query(client.DelRoute, &networker.RouteQuery{
				Table:       "main",
				Protocol:    "static",
				Destination: args[3],
				NextHop:     "any",
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// del route by destination and table
	cli.AddCommandElem(
		nce("route", ""),
		nce("del", "delete route"),
		nce("table", ""),
		nce(libutil.TableRegex, "table number"),
		nce("dst", "destination cidr"),
		ncef(libutil.CidrRegex, "destination cidr", func(args []string) {
			resp, err := query(client.DelRoute, &networker.RouteQuery{
				Table:       args[3],
				Protocol:    "static",
				Destination: args[5],
				NextHop:     "any",
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// del route by destination and nexthop
	cli.AddCommandElem(
		nce("route", ""),
		nce("del", "delete route"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.DelRoute, &networker.RouteQuery{
				Table:       "main",
				Protocol:    "static",
				Destination: args[3],
				NextHop:     args[5],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// del route by destination, nexthop and table
	cli.AddCommandElem(
		nce("route", ""),
		nce("del", "delete route"),
		nce("table", ""),
		nce(libutil.TableRegex, "table number"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.DelRoute, &networker.RouteQuery{
				Table:       args[3],
				Protocol:    "static",
				Destination: args[5],
				NextHop:     args[7],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// del route by destination, source and nexthop
	cli.AddCommandElem(
		nce("route", ""),
		nce("del", "delete route"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("src", "source ip"),
		nce(libutil.CidrRegex, "source ip"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.DelRoute, &networker.RouteQuery{
				Table:       "main",
				Protocol:    "static",
				Destination: args[3],
				Source:      args[5],
				NextHop:     args[7],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))

	// del route by destination, source, nexthop and table
	cli.AddCommandElem(
		nce("route", ""),
		nce("del", "delete route"),
		nce("table", ""),
		nce(libutil.TableRegex, "table number"),
		nce("dst", "destination cidr"),
		nce(libutil.CidrRegex, "destination cidr"),
		nce("src", "source ip"),
		nce(libutil.CidrRegex, "source ip"),
		nce("nexthop", "nexthop ip"),
		ncef(libutil.IpRegex, "nexthop ip", func(args []string) {
			resp, err := query(client.DelRoute, &networker.RouteQuery{
				Table:       args[3],
				Protocol:    "static",
				Destination: args[5],
				Source:      args[7],
				NextHop:     args[9],
			})
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			libutil.PrintStructAll(resp.Routes)
		}))
}
