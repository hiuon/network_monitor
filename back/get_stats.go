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

func startGetStatistics(stats *[240]dataStats ) {
	fmt.Println("Do you want to create new test file? (yes/no)")
	var answer string
	_, err = fmt.Scanln(&answer)
	if err != nil {
		log.Fatal(err)
	}
	if answer != "yes" {
		writeDataFromFile(stats)
		return
	}
	f, _ := os.Create("test.pcap")
	w := pcapgo.NewWriter(f)
	err := w.WriteFileHeader(snapshotLen, layers.LinkTypeEthernet)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//Choose device and get its name
	device := getDeviceName()

	fmt.Println("Default test time duration 240 seconds...\nPlease wait...")
	fmt.Println("Start time:", time.Now())
	timeForStatistics := 240

	// Open device
	handle, err = pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	start := time.Now()
	for packet := range packetSource.Packets() {
		err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			return
		}
		if time.Since(start).Seconds() > float64(timeForStatistics) {
			break
		}
	}
	fmt.Println("End time: ", time.Now())

	fmt.Println("Data has been wrote to test.pcap file")

	writeDataFromFile(stats)

	fmt.Println("Program is starting...")



}

func writeDataFromFile(stats *[240]dataStats) {

	handle, err = pcap.OpenOffline("test.pcap")
	if err != nil { log.Fatal(err) }
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	flagTime := false
	seconds := 0
	numberPackets := 0
	var start time.Time
	for packet := range packetSource.Packets() {
		numberPackets++
		if !flagTime {
			start = packet.Metadata().Timestamp
			flagTime = true
		}

		if packet.Metadata().Timestamp.Sub(start).Microseconds() > int64((seconds + 1)*1000000) {
			seconds++
		}
		if seconds == 240 {
			seconds = 239
		}
		printPacketInfo(packet, stats[seconds])
	}
	//get current time for cycle
	//start := time.Now()
	//start
	//seconds := 0
	/*for packet := range packetSource.Packets() {
		if int(time.Since(start).Seconds()) > seconds + 1 {
			seconds++
		}
		printPacketInfo(packet, stats[seconds])
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++
		if time.Since(start).Seconds() > float64(timeForStatistics) {
			break
		}
	}*/

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


