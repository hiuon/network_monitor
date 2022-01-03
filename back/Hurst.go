package main

import (
	"fmt"
	"math"
)

func getHRS(stats [240]dataStats, time int, data *[4]float64, hurstdisp *[4]float64, item int) {

	ds := make([]int, 240)
	for i := 0; i < 240; i++ {
		ds[i] = stats[i].protocols["IPv4"]
	}
	// Get max and min values in each interval
	// Вычисляем размах для каждого интервала
	statR := make([]float64, len(ds)/time)
	statS := make([]float64, len(ds)/time)
	statMean := make([]float64, len(ds)/time)
	statH := make([]float64, len(ds)/time)

	for i := 0; i < len(ds)/time; i++ {
		//max := 0
		//min := ds[i*time]
		for j := i * time; j < (i+1)*time; j++ {
			statMean[i] += float64(ds[j])
			//if max < ds[j] {
			//	max = ds[j]
			//}
			//if min > ds[j] {
			//	min = ds[j]
			//}
		}
		statMean[i] /= float64(time)
		//statR[i] = max - min
		min := float64(ds[i*time]) - statMean[i]
		max := float64(ds[i*time]) - statMean[i]
		temp := 0.0
		for j := i * time; j < (i+1)*time; j++ {
			temp += float64(ds[j]) - statMean[i]
			if temp > max {
				max = temp
			}
			if temp < min {
				min = temp
			}
		}
		statR[i] = max - min
 	}

	for i := 0; i < len(ds)/time; i++ {
		for j := i * time; j < (i+1)*time; j++ {
			statS[i] += math.Pow(float64(ds[j])-statMean[i], 2)
		}
		statS[i] = math.Sqrt(statS[i])
		statS[i] *= math.Sqrt(1.0 / float64(time))


		if statS[i] == 0 || statR[i] == 0 {
			if i == 0 {
				statH[i] = 0.5
			} else {
				statH[i] = statH[i - 1]
			}
		} else {
			statH[i] = math.Log10(statR[i]/statS[i]) / math.Log10(float64(time)*0.5)
		}
	}

	mean := 0.0
	disp := 0.0
	for i := 0; i < len(statH); i++ {
		mean += statH[i]
	}
	fmt.Println("time:", time, statH)
	mean /= float64(len(statH))
	for i := 0; i < len(statH); i++ {
		disp += math.Pow(mean-statH[i], 2)
	}
	disp = math.Sqrt(disp*1.0/(float64(len(statH)) - 1))

	data[item] = mean
	hurstdisp[item] = disp

	/*
		for i := 0; i < len(ds)/time; i++ {
			mean := 0.0
			max := 0
			min := ds[i]
			// Математическое ожидание для данного интервала
			for j := i*time; j < (i + 1) * time; j++ {
				mean += float64(ds[j])
			}
			mean /= float64(len(ds)/time)
			// Размах накопленного отклонения
			for j := i*time; j < (i + 1) * time; j++ {
				if max < ds[j] {
					max = ds[j]
				}
				if min > ds[j] {
					min = ds[j]
				}
			}
			statR[i] = max - min
			// Среднеквадратичное отклонение
			disp := 0.0
			for j := i*time; j < (i + 1) * time; j++ {
				disp += math.Pow(mean - float64(ds[j]), 2)
			}
			statS[i] = math.Pow(1/float64(len(ds)/time), 0.5) * disp
			statH[i] = math.Log(float64(statR[i])/statS[i]) / math.Log(float64(len(ds)/time) * 0.5)
		}

		mean := 0.0
		disp := 0.0
		for i := 0; i < len(statH); i++ {
			mean += statH[i]
		}
		mean /= float64(len(statH))

		for i := 0; i < len(statH); i++ {
			disp += math.Pow(mean - statH[i], 2)
		}
		disp /= float64(len(statH))

		data[item] = mean
		hurstdisp[item] = disp
	*/
}

func getHRSReal(stats [240]dataStats, index int, data *[4]float64, hurstIndex int, length int) {

	ds := make([]int, length)
	for i := 0; i < length; i++ {
		if index - length < 0 {
			return
		}
		ds[i] = stats[index - length + i].protocols["IPv4"]
	}
	// Get max and min values in each interval
	// Вычисляем размах для каждого интервала
	statR := 0.0
	statS := 0.0
	statMean := 0.0

	for i := 0; i < length; i++ {
		statMean += float64(ds[i])
	}
	statMean /= float64(length)
	min := float64(ds[0]) - statMean
	max := float64(ds[0]) - statMean
	temp := 0.0
	for i := 0; i < length; i++ {
		temp += float64(ds[i]) - statMean
		if temp > max {
			max = temp
		}
		if temp < min {
			min = temp
		}
	}
	statR = max - min


	for i := 0; i < length; i++ {
		statS += math.Pow(float64(ds[i])-statMean, 2)
	}
	statS = math.Sqrt(statS)
	statS *= math.Sqrt(1.0 / float64(length))

	statH := 0.0
	if statS == 0 || statR == 0 {
		statH = 0.5
	} else {
		statH = math.Log10(statR/statS) / math.Log10(float64(length)*0.5)
	}
	data[hurstIndex] = statH
}




