package main

import (
	"encoding/json"
	"fmt"
)

type raceData struct {
	Name     string `json:"name"`
	Position string `json:"position"`
	Lap      string `json:"lap"`
	Gate     string `json:"gate"`
	Time     string `json:"time"`
	Finished string `json:"finished"`
	Colour   string `json:"colour"`
	UID      int    `json:"uid"`
}

func main() {
	d1 := map[string]interface{}{
		"racedata": map[string]interface{}{
			"k1tig": map[string]interface{}{
				"position": "1",
				"lap":      "1",
				"gate":     "1",
				"time":     "1.625",
				"finished": "False",
				"colour":   "FF00FF",
				"uid":      304901,
			},
		},
	}

	//////this just makes the fake data from VD
	dBytes, err := json.Marshal(d1)
	if err != nil {
		fmt.Println("y")
	}
	///////////////////////////

	var msg = make(map[string]map[string]interface{})
	err = json.Unmarshal(dBytes, &msg)
	if err != nil {
		fmt.Println("Error Unmarshalling:", err)
	}

	//key is message type
	for msgType := range msg {
		switch msgType {
		case "racedata":
			var r1 raceData
			for _, value := range msg {
				for racerName, value2 := range value {
					tempMap, ok := value2.(map[string]interface{})
					if !ok {
						fmt.Println("assertion to temp map failed")
					}
					newBytes, err := json.Marshal(tempMap)
					if err != nil {
						fmt.Println("error:", err)
					}

					err = json.Unmarshal(newBytes, &r1)
					if err != nil {
						fmt.Println("Error unmarshaling:", err)
						return
					}
					r1.Name = racerName
					fmt.Printf("%s's data: %+v\n", r1.Name, r1)
				}
			}
		}
	}
}
