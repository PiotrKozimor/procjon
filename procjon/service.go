package procjon

import (
	"errors"
	"fmt"
	"log"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	"google.golang.org/protobuf/proto"
)

// LoadService from badger. Service is uniquely identified by
// pb.ServiceStatus.ServiceIdentifier.
func LoadService(db *badger.DB, service *pb.Service, status *pb.ServiceStatus) error {
	err := db.View(func(txn *badger.Txn) error {
		marshService, err := txn.Get([]byte(status.ServiceIdentifier))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("Service not registered")
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

// SaveService to badger using protobuf.
func SaveService(db *badger.DB, service *pb.Service) error {
	err := db.Update(func(txn *badger.Txn) error {
		marshalled, err := proto.Marshal(service)
		log.Println(marshalled)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(service.ServiceIdentifier), marshalled)
		return err
	})
	return err
}
