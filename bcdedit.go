package go_bcdedit

import (
	"github.com/gabriel-samfira/go-hivex"
	"github.com/jc-lab/go-bcdedit/internal/bcdtemplate"
	"github.com/jc-lab/go-bcdedit/model"
	"github.com/pkg/errors"
	"io"
	"os"
)

type Bcdedit interface {
	io.Closer
	Enumerate(objectId string) (map[string]BcdObject, error)
	UpsertObject(objectId string, description model.BcdDescription) (BcdObject, error)
	GetObject(objectId string) (BcdObject, error)
}

func CreateStore(store string) (Bcdedit, error) {
	err := os.WriteFile(store, bcdtemplate.EMPTY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file")
	}

	h, err := hivex.NewHivex(store, hivex.WRITE)
	if err != nil {
		return nil, errors.Wrap(err, "opening hive file")
	}

	return NewWithHive(h, true)
}

func OpenStore(store string, writable bool) (Bcdedit, error) {
	var flags = hivex.READ
	if writable {
		flags |= hivex.WRITE
	}
	h, err := hivex.NewHivex(store, flags)
	if err != nil {
		return nil, errors.Wrap(err, "opening hive file")
	}

	return NewWithHive(h, writable)
}
