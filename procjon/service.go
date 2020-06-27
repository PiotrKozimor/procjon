package procjon

import (
	"errors"
	"fmt"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	"google.golang.org/protobuf/proto"
)

func LoadService(db *badger.DB, service *pb.Service, status *pb.ServiceStatus) error {
	err := db.View(func(txn *badger.Txn) error {
		marshService, err := txn.Get([]byte(status.ServiceIdentifier))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("Please register service")
		}
		if err != nil {
			return err
		}
		err = marshService.Value(func(val []byte) error {
			err = proto.Unmarshal(val, service)
			return err
		})
		return err
	})
	return err
}

func SaveService(db *badger.DB, service *pb.Service) error {
	err := db.Update(func(txn *badger.Txn) error {
		key, err := txn.Get([]byte(service.ServiceIdentifier))
		if badger.ErrKeyNotFound == err {
			marshalled, err := proto.Marshal(service)
			if err != nil {
				return err
			}
			txn.Set([]byte(service.ServiceIdentifier), marshalled)
		}
		if key != nil {
			return fmt.Errorf("Service %s already exists", service.ServiceIdentifier)
		}
		return nil
	})
	return err
}
