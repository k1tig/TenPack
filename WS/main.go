package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

var rxMsg map[string]json.RawMessage
var racedata map[string]json.RawMessage
var person map[string]string

type rawMsg json.RawMessage

func main() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://192.168.68.85:60003/velocidrone", nil) //check for static ip
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("Connected to VD")
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		fmt.Println(string(message))
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
