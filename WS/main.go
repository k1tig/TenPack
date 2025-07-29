package main

// Test Flags
// Clean up Timer Switch and functions
// Add CSV export option

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/gorilla/websocket"
)

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
	lastMsg                         raceInfo
	raceTimes                       struct {
		lap1, lap2, lap3, final, holeshot float64
	}
}
type race struct {
	id             time.Time
	aborted        bool
	pilots         []pilot
	finishedPilots int
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
		if messageCheck[0] != rune(65533) { // really need to figure out this message....
			err = json.Unmarshal(message, &rawMsg)
			if err != nil {
				fmt.Println("error unmarshaling raw message", err)
				return
			}
		}
		//fmt.Println(messageStr)
		go raceData.msgHandler(message, rawMsg)
	}
}

func (r *race) msgHandler(message []byte, rawMsg map[string]json.RawMessage) {
	mu.Lock()
	defer mu.Unlock()
	for key, rawJson := range rawMsg {
		//fmt.Println(key)
		switch key {
		case "spectatorChange":
		case "FinishGate":
		case "racetype":
		case "racestatus":
			var raceStatus map[string]string
			err := json.Unmarshal(rawJson, &raceStatus)
			if err != nil {
				fmt.Println("Error unmarshaling 'race status'", err)
			}
			for _, raceAction := range raceStatus {
				switch raceAction {
				case "start": // Can't use til Ash fixes API data for non-host
				//	started := "\n**New Race " + start.Render("Start") + "**\n"
				//fmt.Println(started)
				//	r = &race{id: time.Now()} // need to ensure this erases other fields when writing new timestamp
				case "race finished": // Can't use til Ash fixes API data for non-host
				/*	ended := "\n**Race " + end.Render("Ended") + "**\n"
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
						fmt.Println("Race-")
						for _, pilot := range i.pilots {
							fmt.Println(pilot.name)
						}
					}*/
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
		case "countdown":
		case "racedata":
			var msgData map[string]raceInfo
			err := json.Unmarshal(rawJson, &msgData)
			if err != nil {
				fmt.Println("Error unmarshaling 'racedata'", err)
			}
			for pilotName, pilotData := range msgData {
				if len(r.pilots) == 0 {
					started := "\n**New Race " + start.Render("Start") + "**\n"
					fmt.Println(started)
					r.id = time.Now() //check this time is writting right later
					newPilot := pilot{
						name: pilotName,
						uid:  pilotData.Uid,
					}
					newPilot.raceTimes.holeshot = pilotData.Time.float64
					newPilot.lap1Gates = append(newPilot.lap1Gates, pilotData.Time.float64)
					r.pilots = append(r.pilots, newPilot)
					fmt.Println(live.Render("New Pilot added:"), newPilot.name)
					fmt.Println(live.Render(pilotName, "Holeshot:", strconv.FormatFloat(pilotData.Time.float64, 'f', 3, 64)))
				} else {
					rMsg := &pilotData
					pilotFound := false
					for i, racer := range r.pilots {
						p := &r.pilots[i]

						if pilotName == racer.name {
							if p.lastMsg != pilotData {
								if pilotData.Finished.bool {
									raceData.finishedPilots++
								}

								switch rMsg.Lap.int {
								case 1:
									p.lap1Gates = append(p.lap1Gates, rMsg.Time.float64)
									p.lastMsg = pilotData
								case 2:
									p.lap2Gates = append(p.lap2Gates, rMsg.Time.float64)
									if rMsg.Gate.int == 1 {
										p.raceTimes.lap1 = rMsg.Time.float64 - p.raceTimes.holeshot
										fmt.Println(live.Render(pilotName, "Lap1:", strconv.FormatFloat(p.raceTimes.lap1, 'f', 3, 64)))
									}
									p.lastMsg = pilotData
								case 3:
									p.lap3Gates = append(p.lap3Gates, rMsg.Time.float64)
									if rMsg.Gate.int == 1 {
										p.raceTimes.lap2 = rMsg.Time.float64 - p.raceTimes.lap1 - p.raceTimes.holeshot
										fmt.Println(live.Render(pilotName, "Lap2:", strconv.FormatFloat(p.raceTimes.lap2, 'f', 3, 64)))
									}
									if rMsg.Finished.bool {
										p.raceTimes.lap3 = rMsg.Time.float64 - p.raceTimes.lap2 - p.raceTimes.lap1 - p.raceTimes.holeshot
										p.raceTimes.final = rMsg.Time.float64
										fmt.Println(live.Render(pilotName, "Lap3:", strconv.FormatFloat(p.raceTimes.lap3, 'f', 3, 64)))
									}
									p.lastMsg = pilotData
								default:
									fmt.Println("Unknown message header")
									fmt.Printf("%s\n\n", string(message))
								}
							}
							pilotFound = true
							if raceData.finishedPilots == len(raceData.pilots) {
								ended := "\n**Race " + end.Render("Ended") + "**\n"
								fmt.Println(ended)
								fmt.Printf("\n  ~Accumulated Race Times~\n\n")

								rows := [][]string{}
								var t1, t2, t3, fin float64
								sort.Slice(r.pilots, func(i, j int) bool { // yeet this is it's fucked

									return r.pilots[i].raceTimes.final < r.pilots[j].raceTimes.final
								})

								for _, r := range r.pilots {
									r1 := roundFloat(r.raceTimes.lap1, 3)
									r2 := roundFloat(r.raceTimes.lap2, 3)
									r3 := roundFloat(r.raceTimes.lap3, 3)
									rf := roundFloat(r.raceTimes.final, 3)

									lap1 := strconv.FormatFloat(r.raceTimes.lap1, 'f', 3, 64)
									lap2 := strconv.FormatFloat(r.raceTimes.lap2, 'f', 3, 64)
									lap3 := strconv.FormatFloat(r.raceTimes.lap3, 'f', 3, 64)
									final := strconv.FormatFloat(r.raceTimes.final, 'f', 3, 64)

									if t1 == 0 {
										t1 = r1
									} else {
										if r1 < t1 {
											t1 = r1
										}
									}
									if t2 == 0 {
										t2 = r2
									} else {
										if r2 < t2 {
											t2 = r2
										}
									}
									if t3 == 0 {
										t3 = r3
									} else {
										if r3 < t3 {
											t3 = r3
										}
									}
									if fin == 0 {
										fin = rf
									} else {
										if rf < fin {
											fin = rf
										}
									}

									pilotRow := []string{r.name, lap1, lap2, lap3, final}
									rows = append(rows, pilotRow)
								}

								t := table.New().
									Border(lipgloss.NormalBorder()).
									BorderStyle(lipgloss.NewStyle().Foreground(purple)).
									StyleFunc(func(row, col int) lipgloss.Style {
										if row == table.HeaderRow {
											return headerStyle
										}
										even := row%2 == 0

										switch col {
										case 1:
											cellValue, err := strconv.ParseFloat(rows[row][1], 64)
											if err != nil {
												fmt.Println("Error converting string to float:", err)
											}
											if cellValue <= t1 {
												return baseStyle.Foreground(lipgloss.Color("#5ede58ff"))
											} else {
												return baseStyle.Foreground(lipgloss.Color("#bb27b9ff"))
											}

										case 2:
											cellValue, err := strconv.ParseFloat(rows[row][2], 64)
											if err != nil {
												fmt.Println("Error converting string to float:", err)
											}
											if cellValue <= t2 {
												return baseStyle.Foreground(lipgloss.Color("#5ede58ff"))
											} else {
												return baseStyle.Foreground(lipgloss.Color("#bb27b9ff"))
											}
										case 3:
											cellValue, err := strconv.ParseFloat(rows[row][3], 64)
											if err != nil {
												fmt.Println("Error converting string to float:", err)
											}
											if cellValue <= t3 {
												return baseStyle.Foreground(lipgloss.Color("#5ede58ff"))
											} else {
												return baseStyle.Foreground(lipgloss.Color("#bb27b9ff"))
											}
										case 4:
											cellValue, err := strconv.ParseFloat(rows[row][4], 64)
											if err != nil {
												fmt.Println("Error converting string to float:", err)
											}
											if cellValue <= fin {
												return baseStyle.Foreground(lipgloss.Color("#5ede58ff"))
											} else {
												return baseStyle.Foreground(lipgloss.Color("#bb27b9ff"))
											}
										}

										if even {
											return evenRowStyle
										}
										return oddRowStyle
									}).
									Headers("Pilot", "Lap 1", "Lap 2", "Lap 3", "Final").
									Rows(rows...)
								fmt.Println(t)
								splitTotal := len(r.pilots[0].lap1Gates)
								var gateSplits []int
								for i := 1; i < splitTotal; i++ {
									gateSplits = append(gateSplits, i)
								}
								var trackSplits = gateSplits //[]int{} // neon: 2,9,15,19,20,25  //CAWFB: 2, 6, 16, 18, 22, 24, 26 //caw: 3,8
								leadTelem := r.leadSplits(trackSplits...)

								for _, pilot := range r.pilots {
									pSplit := pilotSplits(pilot, leadTelem, trackSplits...)
									title := "\nPilot: " + pilot.name
									fmt.Println(baseStyle.Render(title))
									fmt.Println(pSplit)
								}

								raceFinal := *r
								raceRecords = append(raceRecords, raceFinal)
								raceCounter++
								/*for _, i := range raceRecords {
									fmt.Println("Race-")
									for _, pilot := range i.pilots {
										fmt.Println(pilot.name)
									}
								}*/
								raceData = race{}
								break
							}
						}
					}
					if !pilotFound {
						newPilot := pilot{name: pilotName, uid: pilotData.Uid}
						if rMsg.Lap.int == 1 && rMsg.Gate.int == 1 {
							newPilot.raceTimes.holeshot = pilotData.Time.float64
							newPilot.lastMsg = pilotData
							r.pilots = append(r.pilots, newPilot)
							fmt.Println(live.Render("New Pilot added:") + newPilot.name)
							fmt.Println(live.Render(pilotName, "Holeshot:", strconv.FormatFloat(pilotData.Time.float64, 'f', 3, 64)))
						}
					}
				}
			}
		}
	}
}

// fix this stupid thing later. Git Gud.
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
