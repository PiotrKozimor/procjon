package procjon

import (
	"reflect"
	"testing"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
)

func TestServiceSaveLoad(t *testing.T) {
	db, err := badger.Open(badger.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Error(err)
	}
	service := pb.Service{
		Identifier: "foo",
		Timeout:    5,
		Statuses: []string{
			"green",
			"yellow",
			"red",
		},
	}
	err = SaveService(db, &service)
	if err != nil {
		t.Error(err)
	}
	var loadedService pb.Service
	err = LoadService(db, &loadedService, &pb.ServiceStatus{Identifier: "foo", StatusCode: 0})
	if err != nil {
		t.Error(err)
	}
	if service.Timeout != loadedService.Timeout && service.Identifier != loadedService.Identifier {
		t.Errorf("\nloadedService: %+v, \nservice: %+v", loadedService, service)
	}
	if !reflect.DeepEqual(loadedService.Statuses, service.Statuses) {
		t.Errorf("\nloadedService: %+v, \nservice: %+v", loadedService, service)
	}
	err = LoadService(db, &loadedService, &pb.ServiceStatus{Identifier: "bar"})
	if err == nil {
		t.Error("Wanted err, go nil")
	}
	// Corrupt database
	db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte("foo"), []byte("bar"))
	})
	err = LoadService(db, &loadedService, &pb.ServiceStatus{Identifier: "foo", StatusCode: 0})
	if err == nil {
		t.Error("Wanted err, go nil")
	}
	err = SaveService(db, &pb.Service{})
	if err == nil {
		t.Error("Wanted err, go nil")
	}
	// t.Errorf("Got error %v", err)
}
