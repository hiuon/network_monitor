package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"time"
)

var (
	promiscuous bool = false
	err         error
	timeout     time.Duration = 1 * time.Second
	handle      *pcap.Handle
	snapshotLen uint32 = 1024
	packetCount int    = 0
)

func startGetStatistics() {
	f, _ := os.Create("test.pcap")
	w := pcapgo.NewWriter(f)
	err := w.WriteFileHeader(snapshotLen, layers.LinkTypeEthernet)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//Choose device and get its name
	device := getDeviceName()

	fmt.Print("Enter test time duration (seconds): ")
	var timeForStatistics int
	//_, err = fmt.Scanln(&timeForStatistics)
	//if err != nil {
	//	log.Fatal(err)
	//}
	timeForStatistics = 240
	stats := [240]dataStats{}
	for i := 0; i < 240; i++ {
		stats[i].srcPort = make(map[int]int)
		stats[i].dstPort = make(map[int]int)
		stats[i].protocols = make(map[string]int)
		stats[i].srcAddrIp = make(map[int]int)
		stats[i].dstAddrIp = make(map[int]int)
	}

	// Open device
	handle, err = pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())


	//init stats structure


	//get current time for cycle
	start := time.Now()
	//start
	seconds := 0
	for packet := range packetSource.Packets() {
		if int(time.Since(start).Seconds()) > seconds + 1 {
			seconds++
		}
		printPacketInfo(packet, stats[seconds])
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++
		if time.Since(start).Seconds() > float64(timeForStatistics) {
			break
		}
	}

	fmt.Println("Protocol stats")
	for i := 0; i < 240; i++ {
		fmt.Println("Second: ", i + 1)
		for k, v := range stats[i].protocols {
			fmt.Println(k, " : ", v)
		}
	}
}

func getDeviceName() string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	for index, device := range devices {
		fmt.Println("\n", index+1, "-> Name: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}

	fmt.Print("\nEnter number of the device: ")
	var number int
	_, err = fmt.Scanln(&number)
	if err != nil {
		log.Fatal(err)
	}
	return devices[number-1].Name
}

func printPacketInfo(packet gopacket.Packet, data dataStats) {
	etherLayer := packet.Layer(layers.LayerTypeEthernet)
	if etherLayer != nil {
		ether, _ := etherLayer.(*layers.Ethernet)
		writeStatsProtocol(ether.EthernetType.String(), data)
	}
	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		writeStatsAddrIp(ip.SrcIP, ip.DstIP, data)
		writeStatsProtocol(ip.Protocol.String(), data)
	}
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		tcp, _ := tcpLayer.(*layers.TCP)
		writeStatsPort(int(tcp.SrcPort), int(tcp.DstPort), data)
	}
	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		writeStatsPort(int(udp.SrcPort), int(udp.DstPort), data)
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}

}


