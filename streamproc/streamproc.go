package streamproc

import (
	"encoding/json"
	"fmt"
	"github.com/RayofLightz/ingestion/enrich"
    "github.com/RayofLightz/ingestion/export"
	"github.com/RayofLightz/ingestion/extract"
	"net"
	"reflect"
	"sync"
    "strconv"
)

func listner(pipe chan string, conf map[string]bool) error {
	var addr string
	if conf["local"] == true {
		addr = "127.0.0.1:8080"
	} else {
		addr = "0.0.0.0:8080"
	}
	fmt.Println("Listening on ", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		buf := make([]byte, 1024*3)
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("error reading:", err)
			fmt.Println("request length:", reqLen)
		}
		bfpipe := string(buf)
		pipe <- bfpipe[:reqLen]
		conn.Write([]byte("Good"))
		conn.Close()
	}
}

func worker(wg *sync.WaitGroup, id int, pipe chan string, procp chan []map[string]string) {
	defer wg.Done()
	processData(pipe, procp)
}



func processData(pipe chan string, procp chan []map[string]string) {
    var lockedret struct{
            mux sync.Mutex
            valueArray []map[string]string
    }
	for {
		pdata := <-pipe
		data, err := extract.ExtractJsonRecord(pdata)
        defer lockedret.mux.Lock()
        lockedret.valueArray = recurs(data)
        defer lockedret.mux.Unlock()
		if err != nil {
			fmt.Println("error in extraction")
			fmt.Println(err)
		} else {
            defer lockedret.mux.Lock()
			procp <- lockedret.valueArray
            defer lockedret.mux.Unlock()
		}
	}
}

func recurs(pipe_var map[string]interface{}) []map[string]string {
    tmp_mp := make([]map[string]string,0,20)
	return_mp := make(map[string]string)
	for index, _ := range pipe_var {
		typ := reflect.TypeOf(pipe_var[index])
		if typ.String() == "[]interface {}" {
			np := pipe_var[index].([]interface{})
			for indexn, _ := range np {
				typn := reflect.TypeOf(np[indexn])
				if typn.String() == "map[string]interface {}" {
					mmp := np[indexn].(map[string]interface{})
					tmp := recurs(mmp)
					for indexb, _ := range tmp {
						for key, val := range tmp[indexb] {
							return_mp[key] = val
						}
					}
				} else {
					fmt.Printf("%s>", index)
					fmt.Println(np[indexn])
				}
			}
		} else if typ.String() == "json.Number" {
			var tmpjson json.Number
			tmpjson = pipe_var[index].(json.Number)
			add_val := tmpjson.String()
			return_mp[index] = add_val
		}else if typ.String() == "bool"{
            var tmpstr string
            tmpstr = strconv.FormatBool(pipe_var[index].(bool))
            return_mp[index] = tmpstr
        } else if typ.String() == "map[string]interface {}" {
			mp := pipe_var[index].(map[string]interface{})
			tmp := recurs(mp)
			for indexb, _ := range tmp {
				for key, val := range tmp[indexb] {
					return_mp[key] = val
				}
			}
		} else {
			return_mp[index] = pipe_var[index].(string)
		}
	}
	tmp_mp = append(tmp_mp, return_mp)
	return tmp_mp
}

func calculate(pipe chan []map[string]string, conf map[string]bool) {
	for {
		pipe_var := <-pipe
		if conf["rev_lookup"] == true {
			err := enrich.ReverseLookUp(&pipe_var)
			if err != nil {
				fmt.Println(err)
			} else if conf["check_known_malware"] == true {
				err := enrich.CheckKnownMalware(&pipe_var)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
        fmt.Println(pipe_var)
        err := export.ExportJson(pipe_var)
        if err != nil{
                fmt.Println(err)
        }
	}
}

func StartProcessor(conf map[string]bool) {
	var wg sync.WaitGroup
	pipe := make(chan string, 4)
	processing_pipe := make(chan []map[string]string, 4)
	go listner(pipe, conf)
	go calculate(processing_pipe, conf)
	for i := 0; i != 3; i++ {
		wg.Add(1)
		go worker(&wg, i, pipe, processing_pipe)
	}
	wg.Wait()
}
