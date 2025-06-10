// Package logs parses the log file and prints out statistics.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Logline struct {
	LogTime     string `json:"time"`
	Path        string `json:"path"`
	RequestData string `json:"request_data"`
	Took        int    `json:"took"`
	StatusCode  int    `json:"status_code"`
	OtherTime   string `json:"log_time"`
	Message     string `json:"message"`
}

func main() {
	var logs []Logline
	var tookCount map[string]int

	println("Loading logs...")
	filename := "logs/log.log"
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	logs = make([]Logline, 0)
	tookCount = make(map[string]int)
	for scanner.Scan() {
		line := scanner.Text()
		lineBytes := []byte(line)
		var log Logline
		err = json.Unmarshal(lineBytes, &log)
		logs = append(logs, log)

	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	var log_str []byte
	for _, log := range logs {
		log_str, _ = json.Marshal(log)
		fmt.Println(string(log_str))
		tookCount[log.Path] += log.Took
	}
	for k, v := range tookCount {
		fmt.Printf("%s took total of %d", k, v)
	}
}
