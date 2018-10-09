package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"ws-socket/models"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan models.Message)    // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// HandleConnections handle newly connection
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	messages := readHistory()
	if messages != nil {
		ws.WriteJSON(messages)
	}

	for {
		msg := models.Message{}
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			// Not Log on closing socket
			if strings.Index(err.Error(), "close 1001") == -1 {
				log.Printf("error connection: %s", err.Error())
			}
			delete(clients, ws)
			break
		}
		msg.UpdatedAt = time.Now().Format("02 Jan 2006 15:04:05")
		writeHistory(msg)

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

// HandleMessages handle newly received message
func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error handle message: %s", err.Error())
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func readHistory() []models.Message {
	// open history file if there is none do nothing and return nil
	jsonFile, err := os.Open("./history.json")
	if err != nil {
		if strings.Index(err.Error(), "no such file or directory") == -1 {
			log.Println(err)
		}
		return nil
	}
	// read history file to []byte
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Println(err)
		return nil
	}

	newStr := "[" + string(byteValue) + "]" // make history to array

	messages := []models.Message{}
	//bind all history messages
	if err := json.Unmarshal([]byte(newStr), &messages); err != nil {
		log.Println(err)
		return nil
	}
	return messages
}

func writeHistory(msg models.Message) {
	// convert struct message to []byte
	dataJSON, _ := json.Marshal(msg)
	// open file history.json and append to file if there is none then create the file
	f, err := os.OpenFile("./history.json", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileStat, _ := f.Stat() // get stat of the file
	wAt := fileStat.Size()  // get lenght of byte in the file
	newData := string(dataJSON)
	if wAt > 0 {
		newData = ",\n" + newData // if there is a record append , and new line to it
	}

	// write at the last index of the file
	if _, err := f.Write([]byte(newData)); err != nil {
		panic(err)
	}
}
