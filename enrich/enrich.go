package enrich

import (
	"github.com/RayofLightz/ingestion/extract"
	"net"
)

func ReverseLookUp(data *[]map[string]string) error {
	//Iterate over the data for dest addresses
	for index, _ := range *data {
		if _, ok := (*data)[index]["src_ip"]; ok {
			ret, err := net.LookupAddr((*data)[index]["src_ip"])
			if err != nil {
				return err
			}
			(*data)[index]["reverse_lookup"] = ret[0]
		}
	}
	return nil
}

func CheckKnownMalware(data *[]map[string]string) error {
	hosts, err := extract.ExtractCsvRecord()
	if err != nil {
		return err
	}
	for index, _ := range *data {
		if _, ok := (*data)[index]["reverse_lookup"]; ok {
			for _, val := range hosts {
				if val == (*data)[index]["reverse_lookup"] {
					(*data)[index]["known_domain"] = val
				} else {
					(*data)[index]["known_domain"] = "false"
				}
			}
		}
	}
	return nil
}
