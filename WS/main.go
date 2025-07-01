package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var timeoutTimer time.Duration

func main() {
	chanMsg := make(chan struct{})
	chanExit := make(chan bool)
	timeoutTimer = 10 * time.Minute

	defer close(chanExit)
	defer close(chanMsg)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://192.168.68.85:60003/velocidrone", nil) //check for static ip
	if err != nil {
		log.Panic(err)
	} else {
		log.Println("Connected to VD")
	}
	defer conn.Close()
	keepAlive(conn, timeoutTimer, chanMsg, chanExit)
outterloop:
	for {
		switch {
		case <-chanExit:
			conn.Close()
			fmt.Print("closed")
			break outterloop
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Println(string(message))
			chanMsg <- struct{}{}
		}
	}
	fmt.Println("\nExiting Program")
	time.Sleep(3 * time.Second)
}
func keepAlive(c *websocket.Conn, timeout time.Duration, msg chan struct{}, chExit chan bool) {
	LastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		LastResponse = time.Now()
		return nil
	})
	go func() {
		for {
			select {
			case <-msg:
				LastResponse = time.Now()
			default:
				err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
				if err != nil {
					fmt.Println("Error Keepalive:", err)
					return
				}
				time.Sleep(timeout / 2)
				fmt.Println(LastResponse)
				if time.Since(LastResponse) > timeout {
					fmt.Print("\nInnactive Client timeout: Closing connection....")
					chExit <- true
					return
				}
			}
		}
	}()
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
