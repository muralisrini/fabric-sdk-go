/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package keyvaluestore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	utils "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/bccsp/utils"
	logging "github.com/op/go-logging"
)

var logger = logging.MustGetLogger("fabric_sdk_go")

// FileKeyValueStore ...
type FileKeyValueStore struct {
	path string
}

// CreateNewFileKeyValueStore ...
func CreateNewFileKeyValueStore(path string) (*FileKeyValueStore, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("FileKeyValueStore path is empty")
	}
	createDirIfNotExists(path)
	return &FileKeyValueStore{path: path}, nil
}

// Value ...
/**
 * Get the value associated with name.
 * @param {string} name
 * @returns []byte for the value
 */
func (fkvs *FileKeyValueStore) Value(key string) ([]byte, error) {
	file := path.Join(fkvs.path, key+".json")
	value, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// SetValue ...
/**
 * Set the value associated with name.
 * @param {string} name of the key to save
 * @param {[]byte} value to save
 */
func (fkvs *FileKeyValueStore) SetValue(key string, value []byte) error {
	file := path.Join(fkvs.path, key+".json")
	err := ioutil.WriteFile(file, value, 0600)
	if err != nil {
		return err
	}
	return nil
}

// createDirIfNotExists
func createDirIfNotExists(path string) error {
	missing, err := utils.DirMissingOrEmpty(path)
	logger.Debugf("KeyStore path [%s] missing [%t]: [%s]", path, missing, err)

	if missing {
		os.MkdirAll(path, 0755)
	}

	return nil
}