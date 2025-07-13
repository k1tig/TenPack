package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var localAddress = flag.String("ipAddress", "none", "static local address")
var urlStr string
var ipAddr string

// Add CSV export option

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

	err := godotenv.Load("config.env")
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
						/*
							var spectatorChange map[string]string
							err = json.Unmarshal(message, &spectatorChange)
							if err != nil {
								fmt.Println("Error unmarshaling 'spectator change'", err)
							}
							for key, value := range spectatorChange {
								fmt.Println(key, ":", value)
							}
						*/
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
									fmt.Println("New Race Start")
									raceData = race{id: time.Now()} // need to ensure this erases other fields when writing new timestamp
								case "race finished": // need to handle submitting empty structs
									r := &raceData
									ft := r.raceTimes.lap1 + r.raceTimes.lap2 + r.raceTimes.lap3 + r.raceTimes.holeshot
									raceData.raceTimes.final = ft
									raceRecords = append(raceRecords, raceData)
									raceCounter++
									fmt.Println("Race Ended")
									fmt.Println("Reached pack number:", raceCounter)
									fmt.Println(" ~Race Summary~ ")
									fmt.Println("HoleShot:", r.raceTimes.holeshot)
									fmt.Println("Lap1:", (r.raceTimes.lap1 + r.raceTimes.holeshot))
									fmt.Println("Lap2:", (r.raceTimes.lap2 + r.raceTimes.lap1 + r.raceTimes.holeshot))
									fmt.Println("Lap3:", (r.raceTimes.lap3 + r.raceTimes.lap2 + r.raceTimes.lap1 + r.raceTimes.holeshot))
									fmt.Println("Final Time:", ft)

									// end program
								case "race aborted":
									raceRecords = append(raceRecords, race{
										aborted: true,
										id:      raceData.id})
									raceCounter++
									fmt.Println("Race Aborted")
									fmt.Println("Reached pack number:", raceCounter)
									// end program
								}
							}
						}
					case "countdown":
						/*
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
						*/
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
								switch r.Lap.int {
								case 1:
									raceData.lap1Gates = append(raceData.lap1Gates, r.Time.float64)
									if r.Gate.int == 1 {
										raceData.uid = r.Uid
										raceData.username = racerName
										raceData.raceTimes.holeshot = r.Time.float64
										fmt.Println(racerName, "Holeshot:", r.Time)
									}

								/////// Had the wrong index for lap times. Fix the others later
								case 2:
									raceData.lap2Gates = append(raceData.lap2Gates, r.Time.float64)
									if r.Gate.int == 1 {
										lapLen := len(raceData.lap1Gates)
										lap1 := raceData.lap1Gates[lapLen-1]
										raceData.raceTimes.lap1 = lap1 - raceData.raceTimes.holeshot
										fmt.Println(racerName, "Lap1:", raceData.raceTimes.lap1)
										fmt.Println(racerName, "Lap1(raw):", r.Time.float64)
									}
								/////broke past here
								case 3:
									raceData.lap3Gates = append(raceData.lap3Gates, r.Time.float64)
									if r.Gate.int == 1 {
										lapLen := len(raceData.lap2Gates)
										lap2 := raceData.lap2Gates[lapLen-1]
										raceData.raceTimes.lap2 = lap2 - raceData.raceTimes.lap1 - (2 * raceData.raceTimes.holeshot)
										fmt.Println(racerName, "Lap2:", raceData.raceTimes.lap2)
									}
									if r.Finished.bool {
										raceData.raceTimes.lap3 = r.Time.float64 - raceData.raceTimes.lap2 - raceData.raceTimes.lap1 - (3 * raceData.raceTimes.holeshot)
										fmt.Println(racerName, "Lap3:", raceData.raceTimes.lap3)
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
