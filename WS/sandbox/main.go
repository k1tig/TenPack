package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/exp/maps"
)

type Racer struct {
	race int
	name string
	time float64
}

type Split struct {
	split int
	heat  Racer
}

func main() {
	m := map[string]interface{}{
		"racedata": map[string]interface{}{
			"k1tig": map[string]interface{}{
				"time":     "69.420",
				"position": "1",
				"lap":      "3",
				"gate":     "4",
				"finished": "true",
				"color":    " #semen",
			},
		},
	}
	jsonStr, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	//x := maps.Keys(m["racedata"])
	//fmt.Println(x[0])
	//fmt.Println("Want: map[racedata:map[k1tig:map[colour:FF0000 finished:False gate:4 lap:1 position:1 time:4.492]]]")
	var data map[string]json.RawMessage
	var nextData map[string]json.RawMessage
	var racer map[string]string

	if err := json.Unmarshal(jsonStr, &data); err != nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(data["racedata"], &nextData); err != nil {
		fmt.Println(err)
	}
	x := maps.Keys(nextData)
	racerName := x[0]
	//fmt.Println("Key:  ", x[0])
	if err := json.Unmarshal(nextData[racerName], &racer); err != nil {
		fmt.Println(err)
	}
	//	fmt.Println(racer["time"])
	/*for x, i := range racer {
		fmt.Printf("\nKey: %s\nValue: %s\n", x, i)
	}*/

	var races Split
	races.split = 1
	if racer["finished"] == "true" {
		var r Racer
		t, err := strconv.ParseFloat(racer["time"], 32)
		if err != nil {
			fmt.Printf("racer time error: %s", err)
		}
		r.name = racerName
		r.race = 1
		r.time = t
		races.heat = r

	}

	fmt.Printf("\nCurrent Split:%v\nRacer: %s, Race: %v\nTime: %.2f", races.split, races.heat.name, races.heat.race, races.heat.time)
}
