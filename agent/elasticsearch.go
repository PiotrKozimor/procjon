package agent

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var elasticsearchStatuses = []string{
	"green",
	"yellow",
	"red",
	"unknown",
}

type clusterHealth struct {
	Status string
}

// Elasticsearch holds URL of host to monitor and http client.
type Elasticsearch struct {
	Host   string
	Client httpClient
}

// httpClient is defined for testing purposes.
type httpClient interface {
	Get(url string) (resp *http.Response, err error)
}

// GetCurrentStatus fetches Elasticsearch cluster health from e.host
func (e *Elasticsearch) GetCurrentStatus() uint32 {
	var clusterStatus clusterHealth
	resp, err := e.Client.Get(fmt.Sprintf("%s/_cluster/health", e.Host))
	if err != nil {
		log.Print(err)
		return 3
	}
	if resp.StatusCode != 200 {
		log.Print(err)
		return 3
	}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&clusterStatus)
	if err != nil {
		log.Print(err)
		return 3
	}
	for code, status := range elasticsearchStatuses {
		if status == clusterStatus.Status {
			return uint32(code)
		}
	}
	log.Errorf("Could not find received status in statuses!")
	return 3
}

// GetStatuses defined for ElasticsearchMonitor.
func (e *Elasticsearch) GetStatuses() []string {
	return elasticsearchStatuses
}
