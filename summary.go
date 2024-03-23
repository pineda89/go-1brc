package main

import (
	"fmt"
)

type stationData struct {
	Name  string
	Min   int
	Max   int
	Sum   int
	Count int
}

type summarizedStationData struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func sumarize(workers []*worker) {
	var summaryMap = make(map[uint64]*summarizedStationData)
	for i := range workers {
		workers[i].swdata.Iter(func(k uint64, v *stationData) (stop bool) {
			if oldV, ok := summaryMap[k]; ok {
				oldV.Count += v.Count
				oldV.Sum += toFloat(v.Sum)
				if m := toFloat(v.Max); m > oldV.Max {
					oldV.Max = m
				}
				if m := toFloat(v.Min); m < oldV.Min {
					oldV.Min = m
				}
			} else {
				summaryMap[k] = &summarizedStationData{
					Name:  v.Name,
					Min:   toFloat(v.Min),
					Max:   toFloat(v.Max),
					Sum:   toFloat(v.Sum),
					Count: v.Count,
				}
			}
			return
		})
	}

	for _, v := range summaryMap {
		fmt.Printf("%v %v %v %v %v\n", v.Name, v.Min, v.Max, v.Sum, v.Count)
	}
	fmt.Printf("%v %v\n", len(summaryMap), len(workers))
}

func toFloat(sum int) float64 {
	return float64(sum) / 10
}
