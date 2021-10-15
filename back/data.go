package main

import (
	"encoding/binary"
	"net"
)

type dataStats struct {
	srcPort map[int]int
	dstPort map[int]int
	protocols map[string]int
	srcAddrIp map[int]int
	dstAddrIp map[int]int
}

func (ds dataStats) init() {
	ds.srcPort = make(map[int]int)
	ds.dstPort = make(map[int]int)
	ds.protocols = make(map[string]int)
	ds.srcAddrIp = make(map[int]int)
	ds.dstAddrIp = make(map[int]int)
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func writeStatsPort(src int, dst int, data dataStats) {
	data.srcPort[src] += 1
	data.dstPort[dst] += 1
}

func writeStatsAddrIp(src net.IP, dst net.IP, data dataStats) {
	srcAddr := ip2int(src)
	dstAddr := ip2int(dst)
	data.srcAddrIp[int(srcAddr)] += 1
	data.dstAddrIp[int(dstAddr)] += 1
}

func writeStatsProtocol(protocol string, data dataStats) {
	data.protocols[protocol] += 1
}
