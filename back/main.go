package main

func main() {
	stats := [240]dataStats{}
	for i := 0; i < 240; i++ {
		stats[i].srcPort = make(map[int]int)
		stats[i].dstPort = make(map[int]int)
		stats[i].protocols = make(map[string]int)
		stats[i].srcAddrIp = make(map[int]int)
		stats[i].dstAddrIp = make(map[int]int)
	}
	startGetStatistics(&stats)
}






//func printStatistics(pcapfile string){
//	handle, err = pcap.OpenOffline(pcapfile)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for packet := range gopacket.PacketSource{}
//}
