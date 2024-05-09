package db

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"MoneyGoblin4/structs"
)

var World_Statuses []*structs.World_Status
var Status_Last_Updated int64 = 0

func init() {
	World_Statuses = make([]*structs.World_Status, 0)
	err := Update_World_Status()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Update_World_Status() error {
	resp, err := http.Get("http://127.0.0.1:8000/overview")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&World_Statuses)
	if err != nil {
		return err
	}

	Status_Last_Updated = time.Now().Unix()

	return nil
}

func Fetch_Report(days int) ([]*structs.Loot_Model, error) {

	report := make([]*structs.Loot_Model, 0)

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:8000/report?days=%d", days))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&report)
	if err != nil {
		return nil, err
	}

	return report, nil
}
