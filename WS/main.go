package main

import (
	"fmt"
	"log"

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
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Conn not available:", err)
				close(done)
				return
			}
			fmt.Println(string(message))
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
