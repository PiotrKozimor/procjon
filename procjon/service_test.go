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
		ServiceIdentifier: "foo",
		Timeout:           5,
		Statuses:          make(map[int32]string),
	}
	service.Statuses[0] = "green"
	service.Statuses[1] = "yellow"
	service.Statuses[2] = "red"
	err = SaveService(db, &service)
	if err != nil {
		t.Error(err)
	}
	var loadedService pb.Service
	err = LoadService(db, &loadedService, &pb.ServiceStatus{ServiceIdentifier: "foo", StatusCode: 0})
	if err != nil {
		t.Error(err)
	}
	if service.Timeout != loadedService.Timeout && service.ServiceIdentifier != loadedService.ServiceIdentifier {
		t.Errorf("\nloadedService: %+v, \nservice: %+v", loadedService, service)
	}
	if !reflect.DeepEqual(loadedService.Statuses, service.Statuses) {
		t.Errorf("\nloadedService: %+v, \nservice: %+v", loadedService, service)
	}
}
