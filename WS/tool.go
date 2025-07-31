package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type intConvert struct{ int }
type floatConvert struct{ float64 }
type boolConvert struct{ bool }

// helper funcs
func readJSONFromFile(filename string, v interface{}) error {
	jsonData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, v)
}
func writeJSONToFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ") // Use MarshalIndent for pretty-printed JSON
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}
func (intC *intConvert) UnmarshalJSON(data []byte) error {
	var intStr string
	if err := json.Unmarshal(data, &intStr); err != nil {
		return err
	}
	i, err := strconv.Atoi(intStr)
	if err != nil {
		return err
	}
	intC.int = i
	return nil
}
func (floatC *floatConvert) UnmarshalJSON(data []byte) error {
	var floatStr string
	if err := json.Unmarshal(data, &floatStr); err != nil {
		return err
	}
	f, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return err
	}
	floatC.float64 = f
	return nil
}
func (boolC *boolConvert) UnmarshalJSON(data []byte) error {
	var boolStr string
	if err := json.Unmarshal(data, &boolStr); err != nil {
		return err
	}
	switch boolStr {
	case "True":
		boolC.bool = true
	case "False":
		boolC.bool = false
	default:
		return fmt.Errorf("raceInfo bool cannot be converted")
	}
	return nil
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func pilotSplits(pilot pilot, leadSplit []float64, split ...int) *table.Table {
	var splitTimes [][]float64

	var (
		lap1Gates []float64
		lap2Gates []float64
		lap3Gates []float64
	)

	//this could be brought out of the function
	/*
		totalTrackGates := len(pilot.lap1Gates)
		var lastUserGate = 0
		for _, i := range split {
			if i > lastUserGate {
				lastUserGate = i
			}
		}
		if lastUserGate > totalTrackGates {
			split = []int{}
			for i := 1; i < totalTrackGates; i++ {
				split = append(split, i)
			}
		}*/

	//////////// Holeshot is gate 1!!!!!///////////
	lap1Gates = append(lap1Gates, pilot.lap1Gates[0])
	for _, i := range split {
		lap1Gates = append(lap1Gates, pilot.lap1Gates[i])
	}
	lap1Gates = append(lap1Gates, pilot.lap2Gates[0])
	splitTimes = append(splitTimes, findSplitTimes(lap1Gates))

	/////////
	lap2Gates = append(lap2Gates, pilot.lap2Gates[0])
	for _, i := range split {

		lap2Gates = append(lap2Gates, pilot.lap2Gates[i])
	}
	lap2Gates = append(lap2Gates, pilot.lap3Gates[0])
	splitTimes = append(splitTimes, findSplitTimes(lap2Gates))
	lap3Gates = append(lap3Gates, pilot.lap3Gates[0])
	for _, i := range split {
		lap3Gates = append(lap3Gates, pilot.lap3Gates[i])
	}
	lap3Gates = append(lap3Gates, pilot.lap3Gates[len(pilot.lap3Gates)-1])
	splitTimes = append(splitTimes, findSplitTimes(lap3Gates))

	var rows [][]string
	for _, i := range splitTimes {
		var cleanLap []float64
		for _, x := range i {
			num := roundFloat(x, 3)
			cleanLap = append(cleanLap, num)
		}
		t := slices.Equal(cleanLap, leadSplit)
		if !t {
			var row []string
			for index, pTime := range i {
				diff := pTime - leadSplit[index]
				num := roundFloat(diff, 3)
				x := strconv.FormatFloat(num, 'f', -1, 64)
				row = append(row, x)
			}
			rows = append(rows, row)
		} else {
			var row []string
			for _, pTime := range i {
				num := roundFloat(pTime, 3)
				x := strconv.FormatFloat(num, 'f', -1, 64)
				row = append(row, x)
			}
			rows = append(rows, row)
		}
	}
	headerRow := []string{}
	for i := range rows[0] {
		header := "Split " + strconv.Itoa(i+1)
		headerRow = append(headerRow, header)
	}
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		Headers(headerRow...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			} else {
				num := rows[row][col]
				val, err := strconv.ParseFloat(num, 64)
				if err != nil {
					fmt.Println(err)
				}
				ltStyle := false
				for _, num := range leadSplit {
					if val == num {
						ltStyle = true
					}
				}
				switch {
				case val > 0 && !ltStyle:
					return baseStyle.Foreground(lipgloss.Color("#de5f58ff"))
				case val < 0 && !ltStyle:
					return baseStyle.Foreground(lipgloss.Color("#2cee17ff"))
				default:
					return baseStyle.Foreground(lipgloss.Color("#35bff5ff"))
				}
			}
		}).
		Rows(rows...)
	return t
}

func (r race) leadSplits(split ...int) []float64 { // 2,6,16,18,22,24,26
	var splitTimes []float64
	var pilotGateList []float64
	var hotLap int
	topPilot := r.pilots[0]
	if topPilot.raceTimes.lap1 < topPilot.raceTimes.lap2 {
		hotLap = 1
	} else {
		if topPilot.raceTimes.lap2 < topPilot.raceTimes.lap3 {
			hotLap = 2
		} else {
			hotLap = 3
		}
	}
	/*
		totalTrackGates := len(topPilot.lap1Gates)
		var lastUserGate = 0
		for _, i := range split {
			if i > lastUserGate {
				lastUserGate = i
			}
		}
		if lastUserGate > totalTrackGates {
			fmt.Println("Alert: Specified Track Splits exceed track gate range")
			split = []int{}
			for i := 1; i < totalTrackGates; i++ {
				split = append(split, i)
			}
		}
	*/
	switch hotLap {
	case 1:
		pilotGateList = append(pilotGateList, topPilot.lap1Gates[0])
		for _, i := range split {
			pilotGateList = append(pilotGateList, topPilot.lap1Gates[i])
		}
		pilotGateList = append(pilotGateList, topPilot.lap2Gates[0])
		splitTimes = findSplitTimes(pilotGateList)

	case 2:
		pilotGateList = append(pilotGateList, topPilot.lap2Gates[0])
		for _, i := range split {
			pilotGateList = append(pilotGateList, topPilot.lap2Gates[i])
		}
		pilotGateList = append(pilotGateList, topPilot.lap3Gates[0])
		splitTimes = findSplitTimes(pilotGateList)
	case 3:
		pilotGateList = append(pilotGateList, topPilot.lap3Gates[0])
		for _, i := range split {
			pilotGateList = append(pilotGateList, topPilot.lap3Gates[i])
		}
		pilotGateList = append(pilotGateList, topPilot.lap3Gates[len(topPilot.lap3Gates)-1])
		splitTimes = findSplitTimes(pilotGateList)

	}
	var floatSplitTimes []float64
	for _, i := range splitTimes {
		i = roundFloat(i, 3)
		floatSplitTimes = append(floatSplitTimes, i)
	}
	return floatSplitTimes
}

func findSplitTimes(nums []float64) []float64 {
	var numList []float64
	numLen := len(nums)
	for i := 1; i < numLen; i++ {
		x := nums[i] - nums[i-1]
		numList = append(numList, x)
	}
	return numList
}
