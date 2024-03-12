package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func main() {
	// Initialize Elasticsearch client
	esCfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esClient, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Println("Error creating Elasticsearch client:", err)
		return
	}

	// Command to execute
	cmd := exec.Command("log", "stream", "--style", "json")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Read and process output
	decoder := json.NewDecoder(stdout)
	for {
		var event map[string]interface{}
		if err := decoder.Decode(&event); err != nil {
			fmt.Println("Error decoding JSON:", err)
			break
		}

		fmt.Println("Received event:", event)

		// Send event to Elasticsearch
		index := "eclexys_edr_mac_index"
		req := esapi.IndexRequest{
			Index:      index,
			DocumentID: "",
			Body:       strings.NewReader(fmt.Sprintf("%v", event)),
			Refresh:    "true",
		}

		// Perform the request
		res, err := req.Do(context.Background(), esClient)
		if err != nil {
			fmt.Println("Error indexing document:", err)
		} else {
			defer res.Body.Close()
			if res.IsError() {
				fmt.Printf("Error indexing document: %s", res.Status())
			} else {
				fmt.Println("Document indexed successfully.")
			}
		}

		time.Sleep(1 * time.Second)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error waiting for command to finish:", err)
	}
}

// TODO: Forward the last 1 minute logs
