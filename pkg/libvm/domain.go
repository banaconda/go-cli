package libvm

import (
	"context"
	"fmt"
	"go-cli/pkg/libcli"
	"go-cli/pkg/libutil"
	"go-cli/pkg/libvm/vmer"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/libvirt/libvirt-go"
)

// init domain cli
func initDomainCli(cli *libcli.GoCli) {
	// show domain
	cli.AddCommandElem(
		nce("domain", "domain"),
		ncef("show", "show domains", func(args []string) {
			stream, err := client.ShowDomains(context.Background(), &vmer.DomainMessage{})
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			var streamInterface StreamInterface[*vmer.DomainMessage] = stream

			messages, err := recvStream(streamInterface)
			if err != nil {
				logger.Warn("%v", err)
				return
			}

			libutil.PrintStructAll(messages)
		}))

	// show domain by name

	// create domain by name cpu memory disksize mac ip key volume network
	cli.AddCommandElem(
		nce("domain", "domain"),
		nce("create", "create domain"),
		nce("name", "domain name"),
		nce(libutil.NameRegex, "domain name"),
		nce("cpu", "cpu number"),
		nce(libutil.NumberRegex, "cpu number"),
		nce("memory", "memory size"),
		nce(libutil.UnitRegex, "memory size"),
		nce("disk-size", "disk size"),
		nce(libutil.UnitRegex, "disk size"),
		nce("mac", "mac address"),
		nce(libutil.MacRegex, "mac address"),
		nce("ip", "ip address with mask"),
		nce(libutil.CidrRegex, "ip address"),
		nce("key", "key name"),
		nce(libutil.NameRegex, "key name"),
		nce("image", "image name"),
		nce(libutil.NameRegex, "base image name"),
		nce("network", "network name"),
		nce(libutil.NameRegex, ""),
		nce("bridge", "bridge name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			vcpu, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				logger.Warn("failed to parse vcpu: %v", err)
				return
			}

			_, err = client.CreateDomain(context.Background(), &vmer.DomainMessage{
				Name:       args[3],
				Vcpu:       vcpu,
				Memory:     args[7],
				Mac:        args[11],
				Ip:         args[13],
				Key:        &vmer.KeyMessage{Name: args[15]},
				DiskSize:   args[9],
				Origin:     &vmer.BaseImageMessage{Name: args[17]},
				Network:    &vmer.NetworkMessage{Name: args[19]},
				BridgeName: args[21],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// create domain by name cpu memory distsize ip key volume network
	cli.AddCommandElem(
		nce("domain", "domain"),
		nce("create", "create domain"),
		nce("name", "domain name"),
		nce(libutil.NameRegex, "domain name"),
		nce("cpu", "cpu number"),
		nce(libutil.NumberRegex, "cpu number"),
		nce("memory", "memory size"),
		nce(libutil.UnitRegex, "memory size"),
		nce("disk-size", "disk size"),
		nce(libutil.UnitRegex, "disk size"),
		nce("ip", "ip address with mask"),
		nce(libutil.CidrRegex, "ip address"),
		nce("key", "key name"),
		nce(libutil.NameRegex, "key name"),
		nce("image", "image name"),
		nce(libutil.NameRegex, "base image name"),
		nce("network", "network name"),
		nce(libutil.NameRegex, ""),
		nce("bridge", "bridge name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			vcpu, err := strconv.ParseInt(args[5], 10, 64)
			if err != nil {
				logger.Warn("failed to parse vcpu: %v", err)
				return
			}

			_, err = client.CreateDomain(context.Background(), &vmer.DomainMessage{
				Name:       args[3],
				Vcpu:       vcpu,
				Memory:     args[7],
				DiskSize:   args[9],
				Ip:         args[11],
				Key:        &vmer.KeyMessage{Name: args[13]},
				Origin:     &vmer.BaseImageMessage{Name: args[15]},
				Network:    &vmer.NetworkMessage{Name: args[17]},
				BridgeName: args[19],
			})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// delete domain
	cli.AddCommandElem(
		nce("domain", "domain"),
		nce("delete", "delete domain"),
		nce("name", "domain name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			_, err := client.DeleteDomain(context.Background(), &vmer.DomainMessage{Name: args[3]})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// start domain
	cli.AddCommandElem(
		nce("domain", "domain"),
		nce("start", "start domain"),
		nce("name", "domain name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			_, err := client.StartDomain(context.Background(), &vmer.DomainMessage{Name: args[3]})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// stop domain
	cli.AddCommandElem(
		nce("domain", "domain"),
		nce("stop", "stop domain"),
		nce("name", "domain name"),
		ncef(libutil.NameRegex, "", func(args []string) {
			_, err := client.StopDomain(context.Background(), &vmer.DomainMessage{Name: args[3]})
			if err != nil {
				logger.Warn("%v", err)
				return
			}
		}))

	// reboot domain

	// attach volume

	// detach volume
}

// show domain
func (s *server) ShowDomains(in *vmer.DomainMessage, stream vmer.Vmer_ShowDomainsServer) error {
	var domains []Domain

	var err error
	if in.Name == "" {
		domains, err = vmerDB.GetAllDomains()
	} else {
		domain, err := vmerDB.GetDomainByName(in.Name)
		if err != nil {
			return err
		}
		domains = append(domains, *domain)
	}

	if err != nil {
		logger.Warn("failed to get base images: %v", err)
		return err
	}

	libvirtConn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		logger.Warn("failed to connect to libvirt: %v", err)
		return err
	}
	defer libvirtConn.Close()

	for _, domain := range domains {
		var state int
		// find libvirt domain by name
		dom, err := libvirtConn.LookupDomainByName(domain.Name)
		if err == nil {
			_, state, err = dom.GetState()
			if err != nil {
				logger.Warn("failed to get domain state: %v", err)
				state = 0
			}

		} else {
			logger.Warn("failed to lookup domain: %v", err)
		}

		if err := stream.Send(&vmer.DomainMessage{
			Name:     domain.Name,
			Vcpu:     domain.Cpu,
			Memory:   libutil.ConvertBytesToSize(int64(domain.Memory)),
			Mac:      domain.Mac,
			Ip:       domain.Ip,
			Key:      &vmer.KeyMessage{Name: domain.Key.Name},
			DiskSize: libutil.ConvertBytesToSize(int64(domain.Volume.Size)),
			Origin:   &vmer.BaseImageMessage{Name: domain.Volume.Origin.Name},
			Network:  &vmer.NetworkMessage{Name: domain.Network.Name},
			State:    vmer.DomainMessage_State(state),
		}); err != nil {
			return err
		}
	}

	return nil
}

// create domain
func (s *server) CreateDomain(ctx context.Context, in *vmer.DomainMessage) (*vmer.DomainMessage, error) {
	// check base image
	baseImage, err := vmerDB.GetBaseImageByName(in.Origin.Name)
	if err != nil {
		logger.Warn("failed to get base image: %v", err)
		return nil, err
	}

	// open base image file
	baseImageFile, err := os.Open(baseImage.Path)
	if err != nil {
		logger.Warn("failed to open base image file: %v", err)
		return in, err
	}

	domainBootVolumePath := fmt.Sprintf("%s/%s.qcow2", "/var/local/libvirt/volume", in.Name)

	// copy base image file to volume path
	_, err = libutil.CopyFile(baseImageFile, domainBootVolumePath)
	if err != nil {
		logger.Warn("failed to copy base image file to volume path: %v", err)
		return in, err
	}

	diskSize, err := libutil.ConvertSizeToBytes(in.DiskSize)
	if err != nil {
		logger.Warn("failed to convert size to bytes: %v", err)
		return in, err
	}

	// resize qcow2 image
	out, err := libutil.ResizeQcow2Image(domainBootVolumePath, diskSize)
	if err != nil {
		logger.Warn("failed to resize volume file: %v, out:%s", err, out)
		return in, err
	}

	bootVolumeName := fmt.Sprintf("%s-boot", in.Name)
	if err := vmerDB.InsertVolume(&Volume{
		Name:   bootVolumeName,
		Path:   domainBootVolumePath,
		Size:   diskSize,
		Format: "qcow2",
		Origin: *baseImage,
	}); err != nil {
		logger.Warn("failed to create volume: %v", err)
		return in, err
	}

	// get volume
	volume, err := vmerDB.GetVolumeByName(bootVolumeName)
	if err != nil {
		logger.Warn("failed to get volume: %v", err)
		return in, err
	}

	// check network
	network, err := vmerDB.GetNetworkByName(in.Network.Name)
	if err != nil {
		logger.Warn("failed to get network: %v", err)
		return nil, err
	}

	// check key
	key, err := vmerDB.GetKeyByName(in.Key.Name)
	if err != nil {
		logger.Warn("failed to get key: %v", err)
		return nil, err
	}

	// KiB
	memory, err := libutil.ConvertSizeToBytes(in.Memory)
	if err != nil {
		logger.Warn("failed to convert memory: %v", err)
		return nil, err
	}

	mac := ""
	if len(in.Mac) > 0 {
		mac = in.Mac
	} else {
		mac, err = libutil.GenerateMacAddressByIp(in.Ip)
		if err != nil {
			logger.Warn("failed to generate mac address: %v", err)
			return nil, err
		}
	}

	// make image
	out, err = libutil.RunExternalImageMaker(domainBootVolumePath, key.Username, key.Rsa, mac, network.Vlan,
		in.Ip, network.Gateway, network.Dns)
	if err != nil {
		logger.Warn("failed to make image: %v, out: %s", err, out)
		return nil, err
	}

	// domain xml
	domainXml, err := generateDomainXML(in, domainBootVolumePath)
	if err != nil {
		logger.Warn("failed to generate domain xml: %v", err)
		return nil, err
	}

	libvirtConn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		logger.Warn("failed to connect to libvirt: %v", err)
		return nil, err
	}
	defer libvirtConn.Close()

	libvirtDomain, err := libvirtConn.DomainDefineXML(domainXml)
	if err != nil {
		logger.Warn("failed to define domain: %v", err)
		return nil, err
	}

	libvirtDomain.Create()
	libvirtDomain.SetAutostart(true)

	// create domain
	domain := Domain{
		Name:    in.Name,
		Cpu:     in.Vcpu,
		Memory:  memory,
		Mac:     mac,
		Ip:      in.Ip,
		Key:     *key,
		Volume:  *volume,
		Network: *network,
	}

	// insert domain
	err = vmerDB.InsertDomain(&domain)
	if err != nil {
		logger.Warn("failed to insert domain: %v", err)
		return nil, err
	}

	return in, nil
}

// delete domain
func (s *server) DeleteDomain(ctx context.Context, in *vmer.DomainMessage) (*vmer.DomainMessage, error) {
	// get domain
	domain, err := vmerDB.GetDomainByName(in.Name)
	if err != nil {
		logger.Warn("failed to get domain: %v", err)
		return nil, err
	}

	libvirtConn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		logger.Warn("failed to connect to libvirt: %v", err)
		return nil, err
	}
	defer libvirtConn.Close()

	libvirtDomain, err := libvirtConn.LookupDomainByName(domain.Name)
	if err != nil && err.(libvirt.Error).Code != libvirt.ERR_NO_DOMAIN {
		logger.Warn("failed to lookup domain: %v", err)
		return nil, err
	}

	if err == nil {
		active, err := libvirtDomain.IsActive()
		if err != nil {
			logger.Warn("failed to get domain active: %v", err)
			return nil, err
		}

		if active {
			libvirtDomain.Destroy()
		}

		libvirtDomain.Undefine()
	}

	// delete domain
	err = vmerDB.DeleteDomain(domain)
	if err != nil {
		logger.Warn("failed to delete domain: %v", err)
		return nil, err
	}

	return in, nil
}

// start domain
func (s *server) StartDomain(ctx context.Context, in *vmer.DomainMessage) (*vmer.DomainMessage, error) {
	// get domain
	domain, err := vmerDB.GetDomainByName(in.Name)
	if err != nil {
		logger.Warn("failed to get domain: %v", err)
		return nil, err
	}

	libvirtConn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		logger.Warn("failed to connect to libvirt: %v", err)
		return nil, err
	}
	defer libvirtConn.Close()

	libvirtDomain, err := libvirtConn.LookupDomainByName(domain.Name)
	if err != nil {
		logger.Warn("failed to lookup domain: %v", err)
		return nil, err
	}

	active, err := libvirtDomain.IsActive()
	if err != nil {
		logger.Warn("failed to get domain active: %v", err)
		return nil, err
	}

	if !active {
		libvirtDomain.Create()
	}

	return in, nil
}

// stop domain
func (s *server) StopDomain(ctx context.Context, in *vmer.DomainMessage) (*vmer.DomainMessage, error) {
	// get domain
	domain, err := vmerDB.GetDomainByName(in.Name)
	if err != nil {
		logger.Warn("failed to get domain: %v", err)
		return nil, err
	}

	libvirtConn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		logger.Warn("failed to connect to libvirt: %v", err)
		return nil, err
	}
	defer libvirtConn.Close()

	libvirtDomain, err := libvirtConn.LookupDomainByName(domain.Name)
	if err != nil {
		logger.Warn("failed to lookup domain: %v", err)
		return nil, err
	}

	active, err := libvirtDomain.IsActive()
	if err != nil {
		logger.Warn("failed to get domain active: %v", err)
		return nil, err
	}

	if active {
		libvirtDomain.Destroy()
	}

	return in, nil
}

// generate domain xml
func generateDomainXML(domainMessage *vmer.DomainMessage, volumePath string) (string, error) {
	format := `<domain type='kvm'>
		<name>%s</name>
		<uuid>%s</uuid>
		<vcpu>%d</vcpu>
		<memory unit='KiB'>%d</memory>
		<os>
			<type arch='x86_64' machine='pc'>hvm</type>
			<boot dev='hd'/>
		</os>
		<devices>
			<emulator>/usr/bin/kvm</emulator>
			<disk type='file' device='disk'>
				<driver name='qemu' type='qcow2'/>
				<source file='%s'/>
				<target dev='vda' bus='virtio'/>
				<address type='pci' domain='0x0000' bus='0x00' slot='0x03' function='0x0'/>
			</disk>
			<interface type='bridge'>
				<target dev='%s-port'/>
				<mac address='%s'/>
				<source bridge='%s'/>
				<virtualport type='openvswitch'/>
				<address type='pci' domain='0x0000' bus='0x00' slot='0x04' function='0x0'/>
				<model type='virtio'/>
			</interface>
			<graphics type='vnc' port='-1'/>
			<serial type='pty'>
					<target type='isa-serial' port='0'>
							<model name='isa-serial'/>
					</target>
			</serial>
			<console type='pty'>
					<target type='serial' port='0'/>
			</console>
		</devices>
	</domain>`

	// generate uuid
	uuid := uuid.New()

	memory, err := libutil.ConvertSizeToBytes(domainMessage.Memory)
	memory /= 1024
	if err != nil {
		logger.Warn("failed to convert memory: %v", err)
		return "", err
	}

	mac := ""
	if len(domainMessage.Mac) > 0 {
		mac = domainMessage.Mac
	} else {
		mac, err = libutil.GenerateMacAddressByIp(domainMessage.Ip)
		if err != nil {
			logger.Warn("failed to generate mac address: %v", err)
			return "", err
		}
	}

	return fmt.Sprintf(format, domainMessage.Name, uuid, domainMessage.Vcpu, memory, volumePath, domainMessage.Name, mac, domainMessage.BridgeName), nil
}
