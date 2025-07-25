package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity, adjust in production
	},
}
var mu sync.Mutex

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	} else {
		fmt.Println("Connected:", conn.LocalAddr())
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(2)
	sendData(conn, startData, &wg)
	go sendData(conn, racer1, &wg)
	//go sendData(conn, racer2, &wg)
	//go sendData(conn, racer3, &wg)
	wg.Wait()
	wg.Add(1)
	sendData(conn, endData, &wg)
	wg.Wait()
	fmt.Println("done")

	// Launch a goroutine to send timed messages

	// Keep the main handler alive to receive messages from the client if needed
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		// Process incoming messages here if necessary
	}
}

func main() {

	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Server started...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func makeJson(raw string) []byte {
	jsonData := []byte(raw)
	return jsonData
}

func sendData(conn *websocket.Conn, data []string, wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(4 * time.Second)
	for _, text := range data {
		msg := makeJson(text)
		mu.Lock()

		err := conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error writing message: %v", err)
			break // Exit the goroutine if writing fails
		}
		fmt.Println("Sent:", msg)
		mu.Unlock()

	}
}

var startData = []string{
	`{"spectatorChange":"k1tig"}`,
	`{"spectatorChange":"k1tig"}`,
	`{"racestatus":{"raceAction":"start"}}`,
	`{"racetype":{"raceMode":"THREE_LAP_SINGLE_CLASS","raceFormat":"NORMAL","raceLaps":"3"}}`,
	`{"countdown":{"countValue":"3"}}`,
	`{"countdown":{"countValue":"2"}}`,
	`{"countdown":{"countValue":"1"}}`,
	`{"countdown":{"countValue":"0"}}`,
	`{"FinishGate":{"StartFinishGate":"True"}}`,
}

var racer1 = []string{
	`{"racedata":{"k1tig":{"position":"1","lap":"1","gate":"1","time":"1.185","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"1","gate":"2","time":"1.278","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"1","gate":"3","time":"1.714","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"1","gate":"4","time":"2.207","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"2","gate":"1","time":"2.327","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"2","gate":"2","time":"2.477","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"2","gate":"3","time":"2.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"2","gate":"4","time":"3.153","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"3","gate":"1","time":"3.487","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"3","gate":"2","time":"4.722","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"3","gate":"3","time":"4.900","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"3","gate":"4","time":"5.403","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"k1tig":{"position":"1","lap":"3","gate":"5","time":"5.952","finished":"True","colour":"FF00FF","uid":304901}}}`,

	`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"1","time":"1.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"2","time":"1.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"3","time":"5.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"4","time":"6.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"1","time":"17.027","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"2","time":"18.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"3","time":"19.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"4","time":"21.153","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"1","time":"33.187","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"2","time":"34.722","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"3","time":"36.000","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"4","time":"37.403","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"5","time":"48.152","finished":"True","colour":"FF00FF","uid":304901}}}`,

	`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"1","time":"1.185","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"2","time":"1.278","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"3","time":"1.714","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"4","time":"2.207","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"1","time":"2.327","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"2","time":"2.477","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"3","time":"2.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"4","time":"3.153","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"1","time":"3.487","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"2","time":"4.722","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"3","time":"4.900","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"4","time":"5.403","finished":"False","colour":"FF00FF","uid":304901}}}`,
	`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"5","time":"5.952","finished":"True","colour":"FF00FF","uid":304901}}}`,
}

/*
	var racer2 = []string{
		`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"1","time":"1.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"2","time":"1.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"3","time":"5.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"1","gate":"4","time":"6.985","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"1","time":"17.027","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"2","time":"18.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"3","time":"19.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"2","gate":"4","time":"21.153","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"1","time":"33.187","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"2","time":"34.722","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"3","time":"36.000","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"4","time":"37.403","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"TestGuy":{"position":"1","lap":"3","gate":"5","time":"48.152","finished":"True","colour":"FF00FF","uid":304901}}}`,
	}

	var racer3 = []string{
		`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"1","time":"1.185","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"2","time":"1.278","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"3","time":"1.714","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"1","gate":"4","time":"2.207","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"1","time":"2.327","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"2","time":"2.477","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"3","time":"2.610","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"2","gate":"4","time":"3.153","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"1","time":"3.487","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"2","time":"4.722","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"3","time":"4.900","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"4","time":"5.403","finished":"False","colour":"FF00FF","uid":304901}}}`,
		`{"racedata":{"Venier":{"position":"1","lap":"3","gate":"5","time":"5.952","finished":"True","colour":"FF00FF","uid":304901}}}`,
	}
*/
var endData = []string{
	`{"racestatus":{"raceAction":"race finished"}}`,
}
