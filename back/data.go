package main

type data_stats struct {
	srcPort map[int]int
	dstPort map[int]int
	protocols map[string]int
	srcAddrIp map[int]int
	dstAddrIp map[int]int

}

func writeStats() {
	//myData := data_stats {}
	mymap := map[int]int{}
	mymap[1] = 0
	mymap[1] += 1
	println(mymap[1])
}
