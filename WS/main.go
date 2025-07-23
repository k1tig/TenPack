package main

// Test Flags
// Clean up Timer Switch and functions
// Add CSV export option

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
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
var testServerBool = flag.Bool("tServer", true, "activates test server")
var testFlag bool

type settings struct {
	IpAddr string `json:"IP_ADDR"`
}

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
	lastMsg                         []byte
	raceTimes                       struct {
		lap1, lap2, lap3, final, holeshot float64
	}
}
type race struct {
	id      time.Time
	aborted bool
	pilots  []pilot
	//lastMsg []byte
}

var raceData race
var raceRecords []race // check limit; if raceRecords[10]!=nil...
var raceCounter = 0    //packlimit
var mu sync.Mutex

func main() {
	var userSettings settings
	done := make(chan struct{})

	err := readJSONFromFile("settings.json", &userSettings)
	if err != nil {
		log.Fatalf("Error reading JSON: %v", err)
	}
	fmt.Println("Settings loaded successfully...")
	flag.Parse()
	foundIpFlag := false
	foundTServerFlag := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "ipAddress" {
			foundIpFlag = true
		}
	})
	if foundIpFlag {
		userSettings.IpAddr = *localAddress
	}
	flag.Visit(func(f *flag.Flag) { //this seems fucked up. learn about flags
		if f.Name == "tServer" {
			foundTServerFlag = true
		}
	})
	if foundTServerFlag {
		testFlag = *testServerBool
	}

	var urlStr string
	if testFlag {
		urlStr = "ws://" + userSettings.IpAddr + ":8080/ws"
		//u = url.URL{Scheme: "ws", Host: userSettings.IpAddr, Path: wsPath}
	} else {

		urlStr = "ws://" + userSettings.IpAddr + ":60003/velocidrone"

	}

	log.Printf("connecting to %s", urlStr)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(urlStr, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	if err != nil {
		log.Println("Check IP address")
		log.Panic(err)
	} else {
		log.Println("Connected to VD")
		if foundIpFlag && !foundTServerFlag {
			err := writeJSONToFile("settings.json", userSettings)
			if err != nil {
				log.Fatalf("Error writing JSON: %v", err)
			}
			fmt.Println("(Provided IP Saved)")
		}
	}
	defer conn.Close()
	go pingGenerator(done, conn)

	for {
		var rawMsg map[string]json.RawMessage
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Conn not available:", err)
			close(done)
			return
		}
		messageStr := string(message)
		messageCheck := []rune(messageStr)
		//	fmt.Println(messageStr)
		if messageCheck[0] != rune(65533) { // really need to figure out this message....
			err = json.Unmarshal(message, &rawMsg)
			if err != nil {
				fmt.Println("error unmarshaling raw message", err)
				return
			}
		}
		fmt.Println(messageStr)
		go raceData.msgHandler(message, rawMsg)
	}
}

func (r *race) msgHandler(message []byte, rawMsg map[string]json.RawMessage) {
	mu.Lock()
	defer mu.Unlock()
	for key := range rawMsg {
		//fmt.Println(key)
		switch key {
		case "spectatorChange":
		case "FinishGate":
		case "racetype":
		case "racestatus":
			var raceStatus map[string]map[string]string
			err := json.Unmarshal(message, &raceStatus)
			if err != nil {
				fmt.Println("Error unmarshaling 'race status'", err)
			}
			for _, value := range raceStatus {
				for _, value2 := range value {
					switch value2 {
					case "start":
						started := "\n**New Race " + start.Render("Start") + "**\n"
						fmt.Println(started)
						r = &race{id: time.Now()} // need to ensure this erases other fields when writing new timestamp
					case "race finished": // need to handle submitting empty structs
						ended := "\n**Race " + end.Render("Ended") + "**\n"
						fmt.Println(ended)
						fmt.Printf("\n~Accumulated Race Times~\n")
						for _, r := range r.pilots {
							fmt.Println(results.Render(r.name))
							fmt.Println("HoleShot:", r.raceTimes.holeshot)
							fmt.Println("Lap1:", roundFloat((r.raceTimes.lap1), 3))
							fmt.Println("Lap2:", roundFloat((r.raceTimes.lap2), 3))
							fmt.Println("Lap3:", roundFloat((r.raceTimes.lap3), 3))
							fmt.Printf("Final: %v\n\n", roundFloat((r.raceTimes.final), 3))
						}
						raceFinal := *r
						raceRecords = append(raceRecords, raceFinal)
						raceCounter++
						for _, i := range raceRecords {
							for _, pilot := range i.pilots {
								fmt.Println(pilot.name)
							}
						}
					case "race aborted":
						r = &race{
							aborted: true,
							id:      r.id}
						raceRecords = append(raceRecords, *r)
						raceCounter++
						fmt.Println("Reached pack number:", raceCounter)
						fmt.Println("Race Aborted")
					}
				}
			}
		case "countdown":
		case "racedata":
			//nope
			var msgData map[string]map[string]raceInfo
			err := json.Unmarshal(message, &msgData)
			if err != nil {
				fmt.Println("Error unmarshaling 'racedata'", err)
			}
			for _, value := range msgData {
				for msgName, raceinfo := range value { //msg name is pilots name in msg, rename later
					if len(r.pilots) == 0 {
						//fmt.Println("empty pilot list")
						newPilot := pilot{
							name: msgName,
							uid:  raceinfo.Uid,
						}
						newPilot.raceTimes.holeshot = raceinfo.Time.float64
						r.pilots = append(r.pilots, newPilot)
						fmt.Println(live.Render("New Pilot added:"), newPilot.name)
						fmt.Println(live.Render(msgName, "Holeshot:", strconv.FormatFloat(raceinfo.Time.float64, 'f', 3, 64)))
						break
					} else {
						for i, racer := range r.pilots {
							rMsg := &raceinfo
							p := &r.pilots[i]
							if msgName == racer.name {
								rawRacerMsg, err := json.Marshal(value)
								if err != nil {
									fmt.Println(err)
								}
								if !bytes.Equal(p.lastMsg, rawRacerMsg) {
									switch rMsg.Lap.int {
									case 1:
										p.lap1Gates = append(p.lap1Gates, rMsg.Time.float64)
									/////// Had the wrong index for lap times. Fix the others later
									case 2:
										p.lap2Gates = append(p.lap2Gates, rMsg.Time.float64)
										if rMsg.Gate.int == 1 {

											p.raceTimes.lap1 = rMsg.Time.float64 - p.raceTimes.holeshot
											fmt.Println(live.Render(msgName, "Lap1:", strconv.FormatFloat(p.raceTimes.lap1, 'f', 3, 64)))
										}
									case 3:
										p.lap3Gates = append(p.lap3Gates, rMsg.Time.float64)
										if rMsg.Gate.int == 1 {
											p.raceTimes.lap2 = rMsg.Time.float64 - p.raceTimes.lap1 - p.raceTimes.holeshot
											fmt.Println(live.Render(msgName, "Lap2:", strconv.FormatFloat(p.raceTimes.lap2, 'f', 3, 64)))
										}
										if rMsg.Finished.bool {
											p.raceTimes.lap3 = rMsg.Time.float64 - p.raceTimes.lap2 - p.raceTimes.lap1 - p.raceTimes.holeshot
											p.raceTimes.final = rMsg.Time.float64
											fmt.Println(live.Render(msgName, "Lap3:", strconv.FormatFloat(p.raceTimes.lap3, 'f', 3, 64)))
										}
									default:
										fmt.Println("Unknown message header")
										fmt.Printf("%s\n\n", string(message))
									}
								}
							} else {
								newPilot := pilot{name: msgName, uid: raceinfo.Uid}
								if rMsg.Lap.int == 1 && rMsg.Gate.int == 1 {
									newPilot.raceTimes.holeshot = raceinfo.Time.float64
									r.pilots = append(r.pilots, newPilot)
									fmt.Println(live.Render("New Pilot added:") + newPilot.name)
									fmt.Println(live.Render(msgName, "Holeshot:", strconv.FormatFloat(raceinfo.Time.float64, 'f', 3, 64)))
									break
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
			}
		}
	}
}
