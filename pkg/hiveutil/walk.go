package hiveutil

import (
	"errors"
	"github.com/gabriel-samfira/go-hivex"
	"io/fs"
)

var SkipAll = fs.SkipAll

type ReadFunc = func(node int64, name string, err error) error

func GetObjectsNode(hive *hivex.Hivex) (int64, error) {
	var foundNode int64
	root, err := hive.Root()
	if err != nil {
		return 0, err
	}
	err = ReadNode(hive, root, func(node int64, name string, err error) error {
		if name == "Objects" {
			foundNode = node
			return SkipAll
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if foundNode == 0 {
		return foundNode, errors.New("could not find Root\\Objects")
	}
	return foundNode, nil
}

func ReadNode(hive *hivex.Hivex, parentNode int64, fn ReadFunc) error {
	children, err := hive.NodeChildren(parentNode)
	if err != nil {
		return err
	}
	for _, childNode := range children {
		childName, err := hive.NodeName(childNode)
		err = fn(childNode, childName, err)
		if err != nil {
			if errors.Is(err, SkipAll) {
				return nil
			}
			return err
		}
	}
	return nil
}

func FindChild(hive *hivex.Hivex, parentNode int64, targetKey string) (int64, error) {
	var targetNode int64
	err := ReadNode(hive, parentNode, func(childNode int64, name string, err error) error {
		if err != nil {
			return err
		}
		if name == targetKey {
			targetNode = childNode
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return targetNode, nil
}

func FindValue(hive *hivex.Hivex, parentNode int64, targetKey string) (int64, error) {
	values, err := hive.NodeValues(parentNode)
	if err != nil {
		return 0, err
	}
	for _, value := range values {
		key, err := hive.NodeValueKey(value)
		if err != nil {
			return 0, err
		}
		if key == targetKey {
			return value, nil
		}
	}
	return 0, nil
}

func UpsertNode(hive *hivex.Hivex, parent int64, targetName string) (int64, error) {
	var targetNode int64
	err := ReadNode(hive, parent, func(node int64, name string, err error) error {
		if name == targetName {
			targetNode = node
			return SkipAll
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if targetNode != 0 {
		return targetNode, nil
	}
	return hive.NodeAddChild(parent, targetName)
}
