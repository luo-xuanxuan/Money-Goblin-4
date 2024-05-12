package reminder

import (
	"MoneyGoblin4/db"
	"MoneyGoblin4/structs"
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
			var timer_data structs.Timer_Model
			if err := json.Unmarshal(message, &timer_data); err != nil {
				fmt.Printf("Error decoding JSON: %v\n", err)
				continue
			}

			// Insert each timer model into the database
			for _, w := range db.World_Statuses {
				for _, fc := range w.Free_Company_List {
					if fc.ID != timer_data.Fc_id {
						continue
					}
					fc.Submersible_List[timer_data.Sub_id].Name = timer_data.Name
					fc.Submersible_List[timer_data.Sub_id].Return_Time = timer_data.Return_time
					break
				}
			}

			fmt.Printf("Received: %s\n", message)

		}
	}()
}
