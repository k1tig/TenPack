package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var localAddress = flag.String("ipAddress", "none", "static local address")
var urlStr string
var ipAddr string

// Add CSV export option
type raceInfo struct {
	Lap      int     `json:"lap"`
	Gate     int     `json:"gate"`
	Time     float64 `json:"time"`
	Finished bool    `json:"finished"`
	Uid      int     `json:"uid"`
}

type race struct {
	id                              time.Time
	uid                             int
	username                        string
	aborted                         bool
	lap1Gates, lap2Gates, lap3Gates []float64
	raceTimes                       struct {
		lap1, lap2, lap3, final, holeshot float64
	}
}

var raceRecords []race // check limit; if raceRecords[10]!=nil...
var raceCounter = 0    //packlimit

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	addr, ipExists := os.LookupEnv("IP_ADDR")
	if ipExists {
		ipAddr = addr
	} else {
		fmt.Println("No WS IP addr set")
	}

	flag.Parse()
	foundIpFlag := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "ipAddress" {
			foundIpFlag = true
		}
	})

	if foundIpFlag {
		ipAddr = *localAddress
	}
	urlStr = "ws://" + ipAddr + ":60003/velocidrone"

	done := make(chan struct{})
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(urlStr, nil) //check for static ip
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("Connected to VD")
		if ipExists {
			if addr != ipAddr {
				os.Setenv("IP_ADDR", ipAddr)
			}
		}
	}
	defer conn.Close()
	go msgHandler(done, conn)
	go pingGenerator(done, conn)
	select {
	case <-done:
		return
	}

}

func msgHandler(done chan struct{}, conn *websocket.Conn) {
	var raceData race

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
			if messageCheck[0] != rune(65533) {
				err = json.Unmarshal(message, &rawMsg)
				if err != nil {
					fmt.Println("error unmarshaling raw message", err)
					return
				}
				for key := range rawMsg {
					//fmt.Println(key)
					switch key {
					case "spectatorChange":
						var spectatorChange map[string]string
						err = json.Unmarshal(message, &spectatorChange)
						if err != nil {
							fmt.Println("Error unmarshaling 'spectator change'", err)
						}
						for key, value := range spectatorChange {
							fmt.Println(key, ":", value)
						}
					case "racestatus":
						var raceStatus map[string]map[string]string
						err = json.Unmarshal(message, &raceStatus)
						if err != nil {
							fmt.Println("Error unmarshaling 'race status'", err)
						}
						for key, value := range raceStatus {
							for key2, value2 := range value {
								fmt.Printf("%s: %s '%s'\n", key, key2, value2)
								switch value2 {
								case "start":
									raceData = race{id: time.Now()} // need to ensure this erases other fields when writing new timestamp
								case "race finished": // need to handle submitting empty structs
									r := &raceData
									ft := r.raceTimes.lap1 + r.raceTimes.lap2 + r.raceTimes.lap3 + r.raceTimes.holeshot
									raceData.raceTimes.final = ft
									raceRecords = append(raceRecords, raceData)
									raceCounter++
									// end program
								case "race aborted":
									raceRecords = append(raceRecords, race{
										aborted: true,
										id:      raceData.id})
									raceCounter++
									fmt.Println("Reached pack number:", raceCounter)
									// end program
								}
							}
						}
					case "countdown":
						var countdown map[string]map[string]string
						err = json.Unmarshal(message, &countdown)
						if err != nil {
							fmt.Println("Error unmarshaling 'countdown'", err)
						}
						for key, value := range countdown {
							for _, value2 := range value {
								fmt.Printf("%s: %s\n", key, value2)
							}
						}
					case "racedata":
						//nope
						var msgData map[string]map[string]raceInfo
						err = json.Unmarshal(message, &msgData)
						if err != nil {
							fmt.Println("Error unmarshaling 'racedata'", err)
						}

						for _, value := range msgData {
							for racerName, raceinfo := range value {
								r := &raceinfo
								switch r.Lap {
								case 1:
									raceData.lap1Gates = append(raceData.lap1Gates, r.Time)
									if r.Gate == 1 {
										raceData.uid = r.Uid
										raceData.username = racerName
										raceData.raceTimes.holeshot = r.Time
									}
								case 2:
									raceData.lap2Gates = append(raceData.lap2Gates, r.Time)
									if r.Gate == 1 {
										raceData.raceTimes.lap1 = r.Time
									}
								case 3:
									raceData.lap3Gates = append(raceData.lap3Gates, r.Time)
									if r.Gate == 1 {
										raceData.raceTimes.lap2 = r.Time
									}
									if r.Finished {
										raceData.raceTimes.lap3 = r.Time
									}
								}
							}
						}
					default:
						fmt.Println("Unknown message header")
						fmt.Printf("%s\n\n", string(message))
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
				log.Println("Ping sent.")
			}
		}
	}
}
