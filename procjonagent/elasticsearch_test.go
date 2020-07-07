package procjonagent

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"testing"
)

type ElasticMock struct {
	cnt       int
	responses []string
}

func (e *ElasticMock) Get(url string) (resp *http.Response, err error) {
	if e.cnt == 4 {
		log.Println("Returning error")
		e.cnt++
		return &http.Response{}, errors.New("foo")
	}
	if e.cnt == 5 {
		log.Println("Returning 400 status code")
		e.cnt++
		return &http.Response{StatusCode: 400}, nil
	}
	resp = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewReader([]byte(e.responses[e.cnt])))}
	e.cnt++
	return resp, nil
}

func TestElasticsearchStatus(t *testing.T) {
	responses := []string{
		`{
			"cluster_name" : "testcluster",
			"status" : "green",
			"timed_out" : false,
			"number_of_nodes" : 1,
			"number_of_data_nodes" : 1,
			"active_primary_shards" : 1,
			"active_shards" : 1,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 1,
			"delayed_unassigned_shards": 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch": 0,
			"task_max_waiting_in_queue_millis": 0,
			"active_shards_percent_as_number": 50.0
		  }`,
		`{
			"cluster_name" : "testcluster",
			"status" : "yellow",
			"timed_out" : false,
			"number_of_nodes" : 1,
			"number_of_data_nodes" : 1,
			"active_primary_shards" : 1,
			"active_shards" : 1,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 1,
			"delayed_unassigned_shards": 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch": 0,
			"task_max_waiting_in_queue_millis": 0,
			"active_shards_percent_as_number": 50.0
		  }`,
		`{
			"cluster_name" : "testcluster",
			"status" : "red",
			"timed_out" : false,
			"number_of_nodes" : 1,
			"number_of_data_nodes" : 1,
			"active_primary_shards" : 1,
			"active_shards" : 1,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 1,
			"delayed_unassigned_shards": 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch": 0,
			"task_max_waiting_in_queue_millis": 0,
			"active_shards_percent_as_number": 50.0
		  }`,
		`{
			"cluster_name" : "testcluster",
			"status" : "red",
			"timed_out" : false,
			"number_of_nodes" : 1,
			"number_of_data_nodes" : 1,
			"active_primary_shards" : 1,
			"active_shards" : 1,
			"relocating_shards" : 0,
			"initializing_shards" : 0,
			"unassigned_shards" : 1,
			"delayed_unassigned_shards": 0,
			"number_of_pending_tasks" : 0,
			"number_of_in_flight_fetch": 0,
			"task_max_waiting_in_queue_millis": 0,
			"active_shards_percent_as_number": 50.0
		  `,
	}
	eMock := ElasticMock{cnt: 0, responses: responses}
	e := ElasticsearchMonitor{
		Host:   "localhost",
		Client: &eMock,
	}
	for i := 0; i < 3; i++ {
		statusCode := e.GetCurrentStatus()
		if statusCode != int32(i) {
			t.Errorf("Got status code: %d, wanted: %d", statusCode, i)
		}
	}
	for i := 0; i < 3; i++ {
		statusCode := e.GetCurrentStatus()
		if statusCode != 3 {
			t.Errorf("Got status code: %d, wanted: %d", statusCode, 3)
		}
	}
}

func TestGetStatuses(t *testing.T) {
	e := ElasticsearchMonitor{}
	statuses := e.GetStatuses()
	if !reflect.DeepEqual(statuses, elasticsearchStatuses) {
		t.Errorf("Got: %+v, wanted: %+v", statuses, elasticsearchStatuses)
	}
}
