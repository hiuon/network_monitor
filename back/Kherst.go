package main

import "math"

func getHRS(ds []int, time int){
	// Get max and min values in each interval
	// Вычисляем размах для каждого интервала
	statR := make([]int, int(len(ds)/time))
	statS := make([]float64, int(len(ds)/time))
	statH := make([]float64, int(len(ds)/time))
	for i := 0; i < int(len(ds)/time); i++ {
		mean := 0.0
		max := 0
		min := 0
		// Математическое ожидание для данного интервала
		for j := i*time; j < (i + 1) * time; j++ {
			mean += float64(ds[j])
		}
		mean /= float64(time)
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
		statS[i] = math.Pow(float64(time), 0.5) * disp
		statH[i] = math.Log(float64(statR[i])/statS[i]) / math.Log(float64(time) * 0.5)
	}
}


