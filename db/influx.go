package db

import (
	"context"
	"fmt"
	"os"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type InfluxDB struct {
	Client   influxdb2.Client
	WriteAPI api.WriteAPIBlocking
	QueryAPI api.QueryAPI
	Org      string
	Bucket   string
}

func NewInfluxDB() (*InfluxDB, error) {
	url := os.Getenv("INFLUXDB_URL")
	if url == "" {
		url = "http://localhost:8086"
	}
	token := os.Getenv("INFLUXDB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("INFLUXDB_TOKEN is required")
	}
	org := os.Getenv("INFLUXDB_ORG")
	if org == "" {
		org = "my-org"
	}
	bucket := os.Getenv("INFLUXDB_BUCKET")
	if bucket == "" {
		bucket = "my-bucket"
	}

	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)
	queryAPI := client.QueryAPI(org)

	// Check connection
	_, err := client.Ready(context.Background())
	if err != nil {
		return nil, err
	}

	return &InfluxDB{
		Client:   client,
		WriteAPI: writeAPI,
		QueryAPI: queryAPI,
		Org:      org,
		Bucket:   bucket,
	}, nil
}

func (i *InfluxDB) Close() {
	i.Client.Close()
}
