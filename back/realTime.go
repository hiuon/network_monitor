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
	timings := [4]int{5, 15, 60, 120}
	HurstParam := [4]float64{}
	HurstDisp := [4]float64{}
	index := 0
	flagWaiting := true
	device := getDeviceName()
	for true {
		if index == 240 {
			index = 0
			flagWaiting = false
		}
		readDataFromFile(&stats, getRealData(device), index)
		index += 5
		for i := 0; i < len(timings); i++ {
			getHRS(stats, timings[i], &HurstParam, &HurstDisp, i)
		}
		if flagWaiting {
			fmt.Println("I'm here:", index)
		} else {
			for i := 0; i < len(HurstParam); i++ {
				if HurstParam[i] > hurstTest[i] + 3 * hurstTestDisp[i] || HurstParam[i] < hurstTest[i] - 3 * hurstTestDisp[i] {
					fmt.Println("Warning! Something wrong with your network...")
					fmt.Println("Test data: ", hurstTest)
					fmt.Println("Test data disp: ", hurstTestDisp)

				}
			}
			fmt.Println(HurstParam)
			fmt.Println(HurstDisp)
		}
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
		err := w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
		if err != nil {
			return ""
		}
		if time.Since(start).Seconds() > 5.0 {
			break
		}
	}

	return fileName
}
