package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	done := make(chan struct{})
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://192.168.68.85:60003/velocidrone", nil) //check for static ip
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("Connected to VD")
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
					case "FinishGate":
					case "racedata":
						type raceInfo struct {
							position string
							lap      string
							gate     string
							time     string
							finished string
							color    string
							uid      int
						}

						var racedata map[string]map[string]raceInfo
						err = json.Unmarshal(message, &racedata)
						if err != nil {
							fmt.Println("Error unmarshaling 'racedata'", err)
						}

						for _, value := range racedata {
							for key, raceinfo := range value {
								r := &raceinfo
								fmt.Printf("%s\nPosition: %s\nLap: %s\nGate: %s\nTime: %s\nFinished: %s\nUID: %d\n", key, r.position, r.lap, r.gate, r.time, r.finished, r.uid)
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

/*
		if err := json.Unmarshal(message, &rxMsg); err != nil {
			log.Fatal(err)
		}
		topKey := maps.Keys(rxMsg)
		header := topKey[0]

		switch {
		case header == "racedata":
			if err := json.Unmarshal(rxMsg[header], &racedata); err != nil {
				log.Fatal(err)
			}

			x := maps.Keys(racedata)
			racerName := x[0]

			if err := json.Unmarshal(racedata[racerName], &person); err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Racer's Name: %s\n", racerName)
			for k, v := range person {
				fmt.Printf("%s: %s\n", k, v)
			}
			println()
		case header == "racestatus":

		case header == "racetype":

		case header == "countdown":

		}

		//x := maps.Keys(data["racedata"])
		clear(message)
	}
}
*/
//clear()
