package export

import (
	"encoding/json"
	"os"
)

func appendFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	_, err = file.WriteString(data + "\n")
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func ExportJson(data interface{}) error {
	jsonstring, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = appendFile("logs/postprocessing.json", string(jsonstring))
	if err != nil {
		return err
	}
	return nil
}
