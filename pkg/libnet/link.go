package libnet

import (
	"context"
	"go-cli/pkg/libnet/networker"
	"net"
	"reflect"

	"github.com/vishvananda/netlink"
)

func getLinkTypeString(link netlink.Link) string {
	t := reflect.TypeOf(link).Elem()

	switch t {
	case reflect.TypeOf(netlink.Device{}):
		return "device"
	case reflect.TypeOf(netlink.Bridge{}):
		return "bridge"
	case reflect.TypeOf(netlink.Vlan{}):
		return "vlan"
	case reflect.TypeOf(netlink.Veth{}):
		return "veth"
	default:
		return "unknown"
	}
}

func (s *server) listLink() ([]*networker.NetLink, error) {
	linkSlice, err := netlink.LinkList()
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	linkList := make([]*networker.NetLink, 0)
	for _, link := range linkSlice {
		var parentName string
		var masterName string

		parentLink, err := netlink.LinkByIndex(link.Attrs().ParentIndex)
		if err == nil {
			parentName = parentLink.Attrs().Name
		}

		masterLink, err := netlink.LinkByIndex(link.Attrs().MasterIndex)
		if err == nil {
			masterName = masterLink.Attrs().Name
		}

		typeName := getLinkTypeString(link)

		vlanId := int32(0)
		vlanProtocol := ""
		if typeName == "vlan" {
			vlan, _ := link.(*netlink.Vlan)
			vlanId = int32(vlan.VlanId)
			vlanProtocol = netlink.VlanProtocolToString[vlan.VlanProtocol]
		}

		netLogger.Info("%s %s %v", typeName, reflect.TypeOf(link), link)
		linkList = append(linkList, &networker.NetLink{
			Name:         link.Attrs().Name,
			Type:         typeName,
			Mac:          link.Attrs().HardwareAddr.String(),
			Status:       link.Attrs().OperState.String(),
			Parent:       parentName,
			Master:       masterName,
			VlanId:       vlanId,
			VlanProtocol: vlanProtocol,
		})
	}

	return linkList, err
}

func (s *server) ShowNetLink(ctx context.Context, in *networker.NetLinkQuery) (*networker.NetLinkResponse, error) {
	linkList, err := s.listLink()
	return &networker.NetLinkResponse{NetLinks: linkList}, err
}

func (s *server) SetNetLinkMac(ctx context.Context, in *networker.NetLinkQuery) (*networker.NetLinkResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	err = netlink.LinkSetHardwareAddr(link, net.HardwareAddr(in.Mac))
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.NetLinkResponse{}, err
}

func (s *server) SetNetLinkUp(ctx context.Context, in *networker.NetLinkQuery) (*networker.NetLinkResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkSetUp(link)
	return &networker.NetLinkResponse{}, err
}

func (s *server) SetNetLinkDown(ctx context.Context, in *networker.NetLinkQuery) (*networker.NetLinkResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkSetDown(link)
	return &networker.NetLinkResponse{}, err
}

// bridge
func (s *server) ShowBridge(ctx context.Context, in *networker.BridgeQuery) (*networker.NetLinkResponse, error) {
	linkList, err := s.listLink()
	bridgeList := make([]*networker.NetLink, 0)
	for _, link := range linkList {
		if link.Type == "bridge" {
			bridgeList = append(bridgeList, link)
		}
	}

	return &networker.NetLinkResponse{NetLinks: bridgeList}, err
}

func (s *server) ShowBridgeSlave(ctx context.Context, in *networker.BridgeQuery) (*networker.NetLinkResponse, error) {
	master := in.Name

	linkList, err := s.listLink()
	slaveList := make([]*networker.NetLink, 0)
	for _, link := range linkList {
		if link.Master == master {
			slaveList = append(slaveList, link)
		}
	}

	return &networker.NetLinkResponse{NetLinks: slaveList}, err
}

func (s *server) AddBridge(ctx context.Context, in *networker.BridgeQuery) (*networker.NetLinkResponse, error) {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = in.Name
	newBridge := &netlink.Bridge{LinkAttrs: linkAttrs}
	err := netlink.LinkAdd(newBridge)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkSetUp(link)

	bridgeList := make([]*networker.NetLink, 0)
	bridgeList = append(bridgeList, &networker.NetLink{
		Name:   linkAttrs.Name,
		Type:   "bridge",
		Mac:    linkAttrs.HardwareAddr.String(),
		Status: linkAttrs.OperState.String(),
	})

	return &networker.NetLinkResponse{NetLinks: bridgeList}, err
}

func (s *server) DelBridge(ctx context.Context, in *networker.BridgeQuery) (*networker.NetLinkResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkDel(link)
	return &networker.NetLinkResponse{}, err
}

func (s *server) SetBridgeMaster(ctx context.Context, in *networker.BridgeQuery) (*networker.NetLinkResponse, error) {
	bridge, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	slave, err := netlink.LinkByName(in.SlaveName)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	err = netlink.LinkSetMaster(slave, bridge)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.NetLinkResponse{}, err
}

func (s *server) UnsetBridgeMaster(ctx context.Context, in *networker.BridgeQuery) (*networker.NetLinkResponse, error) {
	slave, err := netlink.LinkByName(in.SlaveName)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	err = netlink.LinkSetNoMaster(slave)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	return &networker.NetLinkResponse{}, err
}

// veth
func (s *server) ShowVeth(ctx context.Context, in *networker.VethQuery) (*networker.NetLinkResponse, error) {
	linkList, err := s.listLink()
	vethList := make([]*networker.NetLink, 0)
	for _, link := range linkList {
		if link.Type == "veth" {
			vethList = append(vethList, link)
		}
	}

	return &networker.NetLinkResponse{NetLinks: vethList}, err
}

func (s *server) AddVeth(ctx context.Context, in *networker.VethQuery) (*networker.NetLinkResponse, error) {
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = in.Name
	newVeth := &netlink.Veth{
		LinkAttrs: linkAttrs,
		PeerName:  in.PeerName,
	}
	err := netlink.LinkAdd(newVeth)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	link1, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	link2, err := netlink.LinkByName(in.PeerName)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkSetUp(link1)
	netlink.LinkSetUp(link2)

	vethList := make([]*networker.NetLink, 0)
	vethList = append(vethList, &networker.NetLink{
		Name:   linkAttrs.Name,
		Type:   "veth",
		Mac:    linkAttrs.HardwareAddr.String(),
		Status: linkAttrs.OperState.String(),
	})
	return &networker.NetLinkResponse{NetLinks: vethList}, err
}

func (s *server) DelVeth(ctx context.Context, in *networker.VethQuery) (*networker.NetLinkResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkDel(link)
	return &networker.NetLinkResponse{}, err
}

// vlan
func (s *server) ShowVlan(ctx context.Context, in *networker.VlanQuery) (*networker.NetLinkResponse, error) {
	linkList, err := s.listLink()
	vlanList := make([]*networker.NetLink, 0)
	for _, link := range linkList {
		if link.Type == "vlan" {
			if (in.Name == "" || in.Name == link.Name) && (in.VlanId == 0 || in.VlanId == link.VlanId) {
				vlanList = append(vlanList, link)
			}
		}
	}

	return &networker.NetLinkResponse{NetLinks: vlanList}, err
}

func (s *server) AddVlan(ctx context.Context, in *networker.VlanQuery) (*networker.NetLinkResponse, error) {
	parent, err := netlink.LinkByName(in.ParentName)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = in.Name
	linkAttrs.ParentIndex = parent.Attrs().Index
	newVlan := &netlink.Vlan{
		LinkAttrs:    linkAttrs,
		VlanId:       int(in.VlanId),
		VlanProtocol: netlink.VLAN_PROTOCOL_8021Q,
	}

	err = netlink.LinkAdd(newVlan)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}

	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkSetUp(link)

	return &networker.NetLinkResponse{}, err

}

func (s *server) DelVlan(ctx context.Context, in *networker.VlanQuery) (*networker.NetLinkResponse, error) {
	link, err := netlink.LinkByName(in.Name)
	if err != nil {
		netLogger.Warn("%v\n", err)
		return nil, err
	}
	netlink.LinkDel(link)
	return &networker.NetLinkResponse{}, err
}
