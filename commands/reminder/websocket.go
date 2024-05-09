package reminder

import (
	"MoneyGoblin4/db"
	"MoneyGoblin4/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func init() {
	// Connect to the WebSocket server
	url := "ws://127.0.0.1:8000/ws"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	// Start a goroutine to continuously read messages
	go func() {
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err.Error())
				continue
			}
			// Decode the JSON body into the struct
			var timer_data []structs.Timer_Model
			if err := json.NewDecoder(bytes.NewReader(message)).Decode(&timer_data); err != nil {
				log.Fatal(err.Error())
				continue
			}

			// Insert each timer model into the database
			for _, timer := range timer_data {
				written := false
				for _, w := range db.World_Statuses {
					for _, fc := range w.Free_Company_List {
						if fc.ID != timer.Fc_id {
							continue
						}
						fc.Submersible_List[timer.Sub_id].Name = timer.Name
						fc.Submersible_List[timer.Sub_id].Return_Time = timer.Return_time
						written = true
						break
					}
					if written {
						break
					}
				}
			}

			fmt.Printf("Received: %s\n", message)

		}
	}()
}
