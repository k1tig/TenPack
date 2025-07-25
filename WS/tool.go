package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
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
