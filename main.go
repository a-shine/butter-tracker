package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var todoList []string

func getCmd(input string) string {
	inputArr := strings.Split(input, " ")
	return inputArr[0]
}

func getMessage(input string) string {
	inputArr := strings.Split(input, " ")
	var result string
	for i := 1; i < len(inputArr); i++ {
		result += inputArr[i]
	}
	return result
}

func updateTodoList(input string) {
	tmpList := todoList
	todoList = []string{}
	for _, val := range tmpList {
		if val == input {
			continue
		}
		todoList = append(todoList, val)
	}
}

func main() {
	db := ConnectToDB()
	query := "SELECT * FROM nodes"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(rows)
	http.HandleFunc("/network", func(w http.ResponseWriter, r *http.Request) {
		// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}
		defer conn.Close()

		// Continuously read and write message
		for {
			output := "Current Todos: \n"
			for _, todo := range todoList {
				output += "\n - " + todo + "\n"
			}
			output += "\n----------------------------------------"
			message := []byte(output)
			err = conn.WriteMessage(1, message)
			if err != nil {
				log.Println("write failed:", err)
				break
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":8080", nil)
}
