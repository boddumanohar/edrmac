package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

func main() {
	// Initialize Elasticsearch client
	esCfg := elasticsearch7.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esClient, err := elasticsearch7.NewClient(esCfg)
	if err != nil {
		fmt.Println("Error creating Elasticsearch client:", err)
		return
	}

	fmt.Println("COnnected to ES")
	// Command to execute
	// cmd := exec.Command("log", "stream", "--style", "json").Output()
	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	fmt.Println("Error creating stdout pipe:", err)
	// 	return
	// }

	// // Start the command
	// err = cmd.Start()
	// if err != nil {
	// 	fmt.Println("Error starting command:", err)
	// 	return
	// }

	// // Read and process output
	// decoder := json.NewDecoder(stdout)
	for {
		fmt.Println("Started decoder")
		// var event map[string]interface{}
		// if err := decoder.Decode(&event); err != nil {
		// 	fmt.Println("Error decoding JSON:", err)
		// 	break
		// }

		// cmd, err := exec.Command("log", "stream", "--style", "json").Output()
		// if err != nil {
		// 	panic(err)
		// }

		args := "stream --style json"
		cmd := exec.Command("log", strings.Split(args, " ")...)

		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println("Received event:", m)
			index := "eclexys_edr_mac_index"
			// data, _ := json.Marshal(event)
			esClient.Index(index, bytes.NewReader([]byte(m)))
		}
		cmd.Wait()

		// Start the command
		// err = cmd.Start()
		// if err != nil {
		// 	fmt.Println("Error starting command:", err)
		// 	return
		// }

		// Send event to Elasticsearch

	}

	// err = cmd.Wait()
	// if err != nil {
	// 	fmt.Println("Error waiting for command to finish:", err)
	// }
}

// TODO: Forward the last 1 minute logs
