package libutil

import (
	"strconv"

	"golang.org/x/sys/unix"
)

const IpRegex = `^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`
const CidrRegex = "^(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(?:\\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}\\/(?:3[0-2]|[0-2]?[0-9])$"
const MacRegex = "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
const PortRegex = "^([0-9]+)-([0-9]+)$"
const ProtoRegex = "^(tcp|udp|icmp|icmpv6)$"
const NumberRegex = "[0-9]+"
const NameRegex = "([a-zA-Z0-9_\\-\\.])+"
const TableRegex = "^[0-9]+$|^local$|^main$|^default$"
const FilePathRegex = ".+"
const UnitRegex = "[0-9]+[kKmMgGtT]?"

const (
	NET_PORT = 10000
	VM_PORT  = 10001
)

// get help string of regex
func GetRegexHelpString(regex string) string {
	switch regex {
	case IpRegex:
		return "IP(xxx.xxx.xxx.xxx)"
	case CidrRegex:
		return "CIDR(IP/MASK)"
	case MacRegex:
		return "MAC(XX:XX:XX:XX:XX:XX)"
	case PortRegex:
		return "PORT_RANGE(1-65535)"
	case ProtoRegex:
		return "PROTOCOL(tcp|udp|icmp|icmpv6)"
	case NumberRegex:
		return "NUMBER"
	case NameRegex:
		return "NAME(only number, letter, underscore, hyphen and dot)"
	case TableRegex:
		return "TABLE(0-255|local|main|default)"
	case FilePathRegex:
		return "FILE_PATH(only number, letter and /._-~)"
	case UnitRegex:
		return "UNIT(k|m|g|t)"
	default:
		return regex
	}
}

// convert table name to unix table id
func StringToUnixTableId(table string) int {
	if table == "" {
		return unix.RT_TABLE_MAIN
	}

	switch table {
	case "local":
		return unix.RT_TABLE_LOCAL
	case "main":
		return unix.RT_TABLE_MAIN
	case "default":
		return unix.RT_TABLE_DEFAULT
	default:
		num, _ := strconv.Atoi(table)
		return num
	}
}

// convert unix table id to table name
func UnixTableIdToString(table int) string {
	switch table {
	case unix.RT_TABLE_LOCAL:
		return "local"
	case unix.RT_TABLE_MAIN:
		return "main"
	case unix.RT_TABLE_DEFAULT:
		return "default"
	default:
		return strconv.Itoa(table)
	}
}
