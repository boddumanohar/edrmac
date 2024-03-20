package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/MaxSchaefer/macos-log-stream/pkg/mls"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
)

const (
	INDEX = "logs-eclexys_edr_mac_index"
)

func main() {
	// Initialize Elasticsearch client
	esCfg := elasticsearch8.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	esClient, err := elasticsearch8.NewClient(esCfg)
	if err != nil {
		fmt.Println("Error creating Elasticsearch client:", err)
		return
	}

	fmt.Println("Started decoder")

	logs := mls.NewLogs()
	// logs.Predicate = ""

	if err := logs.StartGathering(); err != nil {
		panic(err)
	}

	for log := range logs.Channel {
		fmt.Println(log)
		data, err := json.Marshal(log)
		if err != nil {
			panic(err)
		}

		res, err := esClient.Index(INDEX, bytes.NewReader(data))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Indexed document with ID: %s\n", res.String())
	}
}
