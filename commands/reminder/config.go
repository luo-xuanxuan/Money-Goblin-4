package reminder

import (
	"log"
	"os"
)

var channel_id string = "" //could support multiple channels, but dont want to deal with removing channels
var filename string = "./reminder.txt"

func init() {
	// Open file for reading
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer file.Close()

	// Get the file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Read the file content into a byte slice
	content := make([]byte, fileInfo.Size())
	_, err = file.Read(content)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	channel_id = string(content)
}

func update_reminder_channel(id string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write content to the file
	_, err = file.WriteString(id)
	if err != nil {
		return err
	}

	channel_id = id

	return nil
}
