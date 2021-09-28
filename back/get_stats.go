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
	_, err = fmt.Scanln(&timeForStatistics)
	if err != nil {
		log.Fatal(err)
	}

	// Open device
	handle, err = pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	//get current time for cycle
	start := time.Now()
	for packet := range packetSource.Packets() {
		printPacketInfo(packet)
		w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		packetCount++
		// Only capture 100 and then stop
		if time.Since(start).Seconds() > float64(timeForStatistics) {
			break
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

func printPacketInfo(packet gopacket.Packet) {
	etherLayer := packet.Layer(layers.LayerTypeEthernet)
	if etherLayer != nil {
		fmt.Println("Ethernet layer detected.")
		ether, _ := etherLayer.(*layers.Ethernet)
		fmt.Println("Ethernet type: ", ether.EthernetType)
		fmt.Println()
	}
	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Println("IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, ID, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Printf("From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Println("Protocol: ", ip.Protocol)
		tcpLayer, _ := ipLayer.(*layers.TCP)
		if tcpLayer != nil {
			fmt.Println("Source Port: ", tcpLayer.SrcPort)
			fmt.Println("Destination Port: ", tcpLayer.DstPort)
		}
		udpLayer, _ := ipLayer.(*layers.UDP)
		if udpLayer != nil {
			fmt.Println("Source Port: ", udpLayer)
			fmt.Println("Destination Port: ", udpLayer)
		}
		fmt.Println()
	}

	if err := packet.ErrorLayer(); err != nil {
		fmt.Println("Error decoding some part of the packet:", err)
	}

	writeStats()
}


