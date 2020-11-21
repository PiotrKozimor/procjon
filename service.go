package procjon

import (
	"errors"
	"fmt"

	"github.com/PiotrKozimor/procjon/pb"
	"github.com/dgraph-io/badger/v2"
	"google.golang.org/protobuf/proto"
)

// LoadService from badger KV store. Service is uniquely identified by
// pb.ServiceStatus.Identifier.
func LoadService(db *badger.DB, service *pb.Service, status *pb.ServiceStatus) error {
	return db.View(func(txn *badger.Txn) error {
		marshService, err := txn.Get([]byte(status.Identifier))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("service not registered")
		}
		if err != nil {
			return err
		}
		return marshService.Value(func(val []byte) error {
			return proto.Unmarshal(val, service)
		})
	})
}

// SaveService to badger KV store using protobuf.
// If service with given pb.Service.Identifier exists it will be overwritten.
func SaveService(db *badger.DB, service *pb.Service) error {
	err := db.Update(func(txn *badger.Txn) error {
		marshalled, err := proto.Marshal(service)
		if err != nil {
			return err
		}
		err = txn.Set([]byte(service.Identifier), marshalled)
		return err
	})
	return err
}
