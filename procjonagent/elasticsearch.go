package procjonagent

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var elasticsearchStatuses = map[int32]string{
	0: "green",
	1: "yellow",
	2: "red",
	3: "unknown",
}

type clusterHealth struct {
	Status string
}

type ElasticsearchMonitor struct {
	host string
	http HttpClient
}

type HttpClient interface {
	Get(url string) (resp *http.Response, err error)
}

// GetCurrentStatus fetches Elasticsearch cluster health from e.host
func (e *ElasticsearchMonitor) GetCurrentStatus() int32 {
	var clusterStatus clusterHealth
	resp, err := e.http.Get(fmt.Sprintf("http://%s:9200/_cluster/health", e.host))
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
			return code
		}
	}
	log.Errorf("Could not find received status in statuses!")
	return 3
}

func (e *ElasticsearchMonitor) GetStatuses() map[int32]string {
	return elasticsearchStatuses
}
