package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"log"
	"os"
	"strconv"
	"time"
)

func startSniffer(hurstTest [4]float64, hurstTestDisp [4]float64) {
	stats := [240]dataStats{}
	for i := 0; i < 240; i++ {
		stats[i].srcPort = make(map[int]int)
		stats[i].dstPort = make(map[int]int)
		stats[i].protocols = make(map[string]int)
		stats[i].srcAddrIp = make(map[int]int)
		stats[i].dstAddrIp = make(map[int]int)
	}
	HurstParam := hurstTest
	index := 0
	device := getDeviceName()
	for true {
		if index == 240 {
			index = 0
		}
		index += 5
		name := getRealData(device)
		readDataFromFile(&stats, name, index - 5)
		getHRSReal(stats, index, &HurstParam, 0, 5)

		//if index % 15 == 0 {
			getHRSReal(stats, index, &HurstParam, 1, 15)
		//}
		//if index % 60 == 0 {
			getHRSReal(stats, index, &HurstParam, 2, 60)
		//}
		//if index % 120 == 0 {
			getHRSReal(stats, index, &HurstParam, 3, 120)
		//}

		
		for i := 0; i < len(HurstParam); i++ {
			if HurstParam[i] > hurstTest[i] + 2 * hurstTestDisp[i] || HurstParam[i] < hurstTest[i] - 2 * hurstTestDisp[i] {
				fmt.Println("Warning! Something wrong with your network...")
			}
		}
		fmt.Println(index - 5)
		fmt.Printf("Test data: %.2f\n", hurstTest)
		fmt.Printf("Test data disp: %.2f\n", hurstTestDisp)
		fmt.Printf("Real data: %.2f\n", HurstParam)
	}

}

func readDataFromFile(stats *[240]dataStats, fileName string, index int) {
	handle, err = pcap.OpenOffline(fileName + ".pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	flagTime := false
	seconds := 0
	var start time.Time
	for packet := range packetSource.Packets() {
		if !flagTime {
			start = packet.Metadata().Timestamp
			flagTime = true
		}

		if packet.Metadata().Timestamp.Sub(start).Microseconds() > int64((seconds+1)*1000000) {
			seconds++
		}
		if seconds + index == 240 {
			seconds -= 1
		}
		printPacketInfo(packet, stats[seconds+index])
	}
	handle.Close()
	e := os.Remove(fileName + ".pcap")
	if e != nil {
		log.Fatal(e)
	}
}

func getRealData(device string) string {
	date := time.Now()
	fileName := date.Month().String() + strconv.Itoa(date.Day()) + strconv.Itoa(date.Year()) + "-" + strconv.Itoa(date.Hour()) + strconv.Itoa(date.Minute()) + strconv.Itoa(date.Second())
	f, _ := os.Create(fileName + ".pcap")
	w := pcapgo.NewWriter(f)
	err := w.WriteFileHeader(snapshotLen, layers.LinkTypeEthernet)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	handle, err = pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	start := time.Now()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		if time.Since(start).Seconds() > 6.0 {
			break
		}
		err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			return ""
		}

	}

	return fileName
}
