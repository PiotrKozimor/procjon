package procjon

import (
	"errors"
	"fmt"

	"github.com/PiotrKozimor/procjon"
	"github.com/dgraph-io/badger/v2"
	"google.golang.org/protobuf/proto"
)

func loadService(db *badger.DB, service *procjon.Service, status *procjon.ServiceStatus) error {
	err := db.View(func(txn *badger.Txn) error {
		marshService, err := txn.Get([]byte(status.ServiceIdentifier))
		if errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("Please register service")
		}
		if err != nil {
			return err
		}
		err = proto.Unmarshal(marshService.Key(), service)
		return err
	})
	return err
}

func saveService(db *badger.DB, service *procjon.Service) error {
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

	})
	return err
}
