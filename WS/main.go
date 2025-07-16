package main

// Test Flags
// Clean up Timer Switch and functions
// Add CSV export option

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

var results = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("171"))
var live = lipgloss.NewStyle().
	Foreground(lipgloss.Color("8"))
var start = lipgloss.NewStyle().
	Foreground(lipgloss.Color("46"))
var end = lipgloss.NewStyle().
	Foreground(lipgloss.Color("124"))

var localAddress = flag.String("ipAddress", "none", "static local address")
var urlStr string

type settings struct {
	IpAddr string `json:"IP_ADDR"`
}
type intConvert struct{ int }
type floatConvert struct{ float64 }
type boolConvert struct{ bool }
type raceInfo struct {
	Lap      intConvert   `json:"lap"`
	Gate     intConvert   `json:"gate"`
	Time     floatConvert `json:"time"`
	Finished boolConvert  `json:"finished"`
	Uid      int          `json:"uid"`
}
type pilot struct {
	name                            string
	uid                             int
	lap1Gates, lap2Gates, lap3Gates []float64
	raceTimes                       struct {
		lap1, lap2, lap3, final, holeshot float64
	}
}
type race struct {
	id      time.Time
	aborted bool
	pilots  []pilot
}

var raceRecords []race // check limit; if raceRecords[10]!=nil...
var raceCounter = 0    //packlimit

func main() {
	var userSettings settings
	err := readJSONFromFile("settings.json", &userSettings)
	if err != nil {
		log.Fatalf("Error reading JSON: %v", err)
	}
	fmt.Println("Settings loaded successfully...")
	flag.Parse()
	foundIpFlag := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "ipAddress" {
			foundIpFlag = true
		}
	})
	if foundIpFlag {
		userSettings.IpAddr = *localAddress
	}
	urlStr = "ws://" + userSettings.IpAddr + ":60003/velocidrone"
	done := make(chan struct{})
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(urlStr, nil) //check for static ip
	if err != nil {
		log.Println("Check IP address")
		log.Panic(err)
	} else {
		log.Println("Connected to VD")
		if foundIpFlag {
			err := writeJSONToFile("settings.json", userSettings)
			if err != nil {
				log.Fatalf("Error writing JSON: %v", err)
			}
			fmt.Println("(Provided IP Saved)")
		}
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	go msgHandler(done, conn, &wg)
	go pingGenerator(done, conn)
	wg.Wait()

}

func msgHandler(done chan struct{}, conn *websocket.Conn, wg *sync.WaitGroup) {
	var raceData race
	defer wg.Done()
	for {
		select {
		case <-done:
			return
		default:
			var rawMsg map[string]json.RawMessage
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Conn not available:", err)
				close(done)
				return
			}
			messageStr := string(message)
			messageCheck := []rune(messageStr)
			if messageCheck[0] != rune(65533) { // really need to figure out this message....
				err = json.Unmarshal(message, &rawMsg)
				if err != nil {
					fmt.Println("error unmarshaling raw message", err)
					return
				}
				for key := range rawMsg {
					//fmt.Println(key)
					switch key {
					case "spectatorChange":
					case "FinishGate":
					case "racetype":
					case "racestatus":
						var raceStatus map[string]map[string]string
						err = json.Unmarshal(message, &raceStatus)
						if err != nil {
							fmt.Println("Error unmarshaling 'race status'", err)
						}
						for _, value := range raceStatus {
							for _, value2 := range value {
								switch value2 {
								case "start":
									started := "\n**New Race " + start.Render("Start") + "**\n"
									fmt.Println(started)
									raceData = race{id: time.Now()} // need to ensure this erases other fields when writing new timestamp
								case "race finished": // need to handle submitting empty structs
									ended := "\n**Race " + end.Render("Ended") + "**\n"
									fmt.Println(ended)
									fmt.Printf("\n~Accumulated Race Times~\n")
									for _, r := range raceData.pilots {
										fmt.Println(results.Render("Pilot - ", r.name))
										fmt.Println("HoleShot:", r.raceTimes.holeshot)
										fmt.Println("Lap1:", roundFloat((r.raceTimes.lap1), 3))
										fmt.Println("Lap2:", roundFloat((r.raceTimes.lap2), 3))
										fmt.Println("Lap3:", roundFloat((r.raceTimes.lap3), 3))
										fmt.Printf("Final: %v\n\n", roundFloat((r.raceTimes.final), 3))
										//fmt.Println("Reached pack number:", raceCounter)
									}
									raceRecords = append(raceRecords, raceData)
									raceCounter++
									//fmt.Println(raceRecords)
									// end program
								case "race aborted":
									raceData = race{
										aborted: true,
										id:      raceData.id}
									raceRecords = append(raceRecords, raceData)
									raceCounter++
									fmt.Println("Reached pack number:", raceCounter)
									fmt.Println("Race Aborted")
									// end program
								}
							}
						}
					case "countdown":
					case "racedata":
						//nope
						var msgData map[string]map[string]raceInfo
						err = json.Unmarshal(message, &msgData)
						if err != nil {
							fmt.Println("Error unmarshaling 'racedata'", err)
						}
						for _, value := range msgData {
							for msgName, raceinfo := range value { //msg name is pilots name in msg, rename later
								if len(raceData.pilots) == 0 {
									//fmt.Println("empty pilot list")
									newPilot := pilot{name: msgName, uid: raceinfo.Uid}
									raceData.pilots = append(raceData.pilots, newPilot)
									fmt.Println(live.Render("New Pilot added:"), newPilot.name)
								} else {
									for i, racer := range raceData.pilots {
										r := &raceinfo
										p := &raceData.pilots[i]
										if msgName == racer.name {

											switch r.Lap.int {
											case 1:
												p.lap1Gates = append(p.lap1Gates, r.Time.float64)
												if r.Gate.int == 1 {
													p.raceTimes.holeshot = r.Time.float64
													fmt.Println(live.Render(msgName, "Holeshot:", strconv.FormatFloat(r.Time.float64, 'f', 3, 64)))
												}
											/////// Had the wrong index for lap times. Fix the others later
											case 2:
												p.lap2Gates = append(p.lap2Gates, r.Time.float64)
												if r.Gate.int == 1 {

													p.raceTimes.lap1 = r.Time.float64 - p.raceTimes.holeshot
													fmt.Println(live.Render(msgName, "Lap1:", strconv.FormatFloat(p.raceTimes.lap1, 'f', 3, 64)))
												}
											case 3:
												p.lap3Gates = append(p.lap3Gates, r.Time.float64)
												if r.Gate.int == 1 {
													p.raceTimes.lap2 = r.Time.float64 - p.raceTimes.lap1 - p.raceTimes.holeshot
													fmt.Println(live.Render(msgName, "Lap2:", strconv.FormatFloat(p.raceTimes.lap2, 'f', 3, 64)))
												}
												if r.Finished.bool {
													p.raceTimes.lap3 = r.Time.float64 - p.raceTimes.lap2 - p.raceTimes.lap1 - p.raceTimes.holeshot
													p.raceTimes.final = r.Time.float64
													fmt.Println(live.Render(msgName, "Lap3:", strconv.FormatFloat(p.raceTimes.lap3, 'f', 3, 64)))
												}
											default:
												fmt.Println("Unknown message header")
												fmt.Printf("%s\n\n", string(message))
											}
										} else {
											newPilot := pilot{name: msgName, uid: raceinfo.Uid}
											if r.Lap.int == 1 && r.Gate.int == 1 {
												newPilot.raceTimes.holeshot = raceinfo.Time.float64
												fmt.Println(msgName, "Holeshot:", roundFloat(newPilot.raceTimes.holeshot, 3))
												raceData.pilots = append(raceData.pilots, newPilot)
												fmt.Println("New Pilot added", newPilot.name)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func pingGenerator(done chan struct{}, c *websocket.Conn) {
	for {
		select {
		case <-done:
			close(done)
			return
		default:
			ticker := time.NewTicker(time.Second * 30) // Send ping every 30 seconds
			defer ticker.Stop()
			for range ticker.C {
				err := c.WriteControl(websocket.PingMessage, []byte(""), time.Now().Add(time.Second*10)) // 10-second write deadline
				if err != nil {
					log.Println("write ping error:", err)
					return
				}
				//log.Println("Ping sent.")
			}
		}
	}
}

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
