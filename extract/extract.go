package extract

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"strings"
)

func ExtractJsonRecord(json_blob string) (map[string]interface{}, error) {
	var jsonReturn map[string]interface{}
	jsonReturn = make(map[string]interface{})
	d := json.NewDecoder(strings.NewReader(json_blob))
	d.UseNumber()
	err := d.Decode(&jsonReturn)
	if err != nil {
		return nil, err
	}
	return jsonReturn, err
}

func ExtractCsvRecord() ([]string, error) {
	var retval []string
	data, err := ioutil.ReadFile("config/knownhosts.csv")
	if err != nil {
		return nil, err
	}
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	for index, _ := range records {
		retval = append(retval, records[index][1])
	}
	return retval, nil
}

func JsonConf() (map[string]bool, error) {
	//Extract config values
	data, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		return nil, err
	}
	var jsonReturn map[string]bool
	jsonReturn = make(map[string]bool)
	err = json.Unmarshal(data, &jsonReturn)
	if err != nil {
		return nil, err
	}
	return jsonReturn, nil
}
