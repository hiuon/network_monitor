package main

import "fmt"

func main() {
	stats := [240]dataStats{}
	timings := [4]int{5, 15, 60, 120}
	HurstParam := [4]float64{}
	HurstDisp := [4]float64{}
	var testData []int
	for i := 0; i < 240; i++ {
		stats[i].srcPort = make(map[int]int)
		stats[i].dstPort = make(map[int]int)
		stats[i].protocols = make(map[string]int)
		stats[i].srcAddrIp = make(map[int]int)
		stats[i].dstAddrIp = make(map[int]int)
	}
	startGetStatistics(&stats)
	testData = make([]int, 240)
	for i := 0; i < 240; i++ {
		testData[i] = stats[i].protocols["IPv4"]
	}
	for i := 0; i < len(timings); i++ {
		getHRS(testData, timings[i], &HurstParam, &HurstDisp, i)
	}

	fmt.Println(HurstParam)
}






//func printStatistics(pcapfile string){
//	handle, err = pcap.OpenOffline(pcapfile)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for packet := range gopacket.PacketSource{}
//}
