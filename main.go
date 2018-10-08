package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"ws-socket/handler"
)

func main() {

	// Capture Interrupt ^C signal
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		select {
		// sig is a ^C, handle it
		case <-exit:
			os.Remove("./history.json")
			os.Exit(1)
		}
	}()

	// Create a simple file server
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handler.HandleConnections)
	http.HandleFunc("/join", handler.Join)

	// Start listening for incoming chat messages
	go handler.HandleMessages()
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
